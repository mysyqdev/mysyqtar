package mysyqtar

type Cat struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Breed       string `json:"breed"`
	Gender      Gender `json:"gender"`
}

func (c *Cat) ParseGender(s string) {
	switch s {
	case "Male":
		c.Gender = Male
	case "Female":
		c.Gender = Female
	default:
		c.Gender = Other
	}
}

type Gender int

const (
	Male Gender = iota
	Female
	Other
)
