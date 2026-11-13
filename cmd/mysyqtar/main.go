package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/mysyqdev/mysyqtar"
)

func parseForm(r *http.Request, c *mysyqtar.Cat) error {
	name := r.FormValue("name")
	description := r.FormValue("description")
	breed := r.FormValue("breed")
	gender := r.FormValue("gender")

	if name == "" {
		return errors.New("Invalid form")
	}

	c.Name = name
	c.Description = description
	c.Breed = breed
	c.SetGender(gender)
	return nil
}

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("Unable to connect to db")
		return
	}
	defer conn.Close(context.Background())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("meow not found :(\n"))
	})

	http.HandleFunc("POST /mysyq/", func(w http.ResponseWriter, r *http.Request) {
		var cat mysyqtar.Cat
		// this or json?
		err := parseForm(r, &cat)
		if err != nil {
			fmt.Println("Could not parse the form")
			// write an error to w
			return
		}

		tag, err := conn.Exec(
			context.Background(),
			`INSERT INTO cats (name, description, breed, gender)
			VALUES ($1, $2, $3, $4)`,
			cat.Name, cat.Description, cat.Breed, cat.Gender.String(),
		)
		if err != nil && err != pgx.ErrNoRows {
			fmt.Printf("%v\n", err)
			w.Write([]byte("Could not add a new mysyq"))
			return
		}
		fmt.Println(tag.String())
		w.Write([]byte("Added new mysyq"))
	})

	http.HandleFunc("GET /mysyqtar/", func(w http.ResponseWriter, r *http.Request) {
		var cats []mysyqtar.Cat

		rows, err := conn.Query(
			context.Background(),
			"SELECT * FROM  cats",
		)
		if err != nil {
			fmt.Printf("%v\n", err)
			fmt.Println("QueryRow failed")
			return
		}

		for rows.Next() {
			var cat mysyqtar.Cat
			var gender string
			rows.Scan(&cat.ID, &cat.Name, &cat.Description, &gender, &cat.Breed)
			cat.SetGender(gender)
			cats = append(cats, cat)
			if rows.Err() != nil {
				fmt.Println(rows.Err().Error())
			}
		}

		re := make(map[string][]mysyqtar.Cat)
		re["cats"] = cats

		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			json.NewEncoder(w).Encode(re)
		} else {
			w.Write([]byte("Could not get mysyqs\n"))
		}
	})

	http.HandleFunc("GET /mysyqtar/{id}/", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid request\n"))
			fmt.Println("No id")
			return
		}

		var cat mysyqtar.Cat
		var gender string
		err = conn.QueryRow(
			context.Background(),
			"SELECT * FROM cats WHERE id = $1", id,
		).Scan(&cat.ID, &cat.Name, &cat.Description, &gender, &cat.Breed)
		cat.SetGender(gender)
		if err != nil {
			fmt.Printf("%v\n", err)
			w.Write([]byte("Mysyq not found"))
			return
		}

		response, err := json.Marshal(cat)
		if err == nil {
			w.Write(response)
		} else {
			w.Write([]byte("Could not get mysyqs\n"))
		}
	})

	// default serve mux not recommended?
	http.ListenAndServe("", nil)
}
