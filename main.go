package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

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

	http.HandleFunc("POST /mysyq", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		fmt.Println("Name to add:", name)
		tag, err := conn.Exec(context.Background(), "INSERT INTO cats (name) VALUES ($1)", name)
		if err != nil && err != pgx.ErrNoRows {
			fmt.Printf("%v\n", err)
			w.Write([]byte("Could not add a new mysyq"))
			return
		}
		fmt.Println(tag.String())
		w.Write([]byte("Added new mysyq"))
	})

	http.HandleFunc("GET /mysyqtar", func(w http.ResponseWriter, r *http.Request) {
		var names []string
		rows, err := conn.Query(context.Background(), "SELECT name FROM  cats")
		for rows.Next() {
			var name string
			rows.Scan(&name)
			names = append(names, name)
		}
		if err != nil {
			fmt.Printf("%v\n", err)
			fmt.Println("QueryRow failed")
			return
		}
		for _, name := range names {
			w.Write([]byte(name))
			w.Write([]byte("\n"))
		}
	})

	http.HandleFunc("GET /mysyqtar/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			fmt.Println("No id")
			return
		}

		var name string
		err = conn.QueryRow(context.Background(), "SELECT name FROM cats WHERE id = $1", id).Scan(&name)
		if err != nil {
			fmt.Printf("%v\n", err)
			w.Write([]byte("Mysyq not found"))
			return
		}

		w.Write([]byte(name))
	})

	// default serve mux not recommended?
	http.ListenAndServe("", nil)
}
