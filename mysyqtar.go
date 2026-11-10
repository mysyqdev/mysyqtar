package mysyqtar

type Cat struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Breed       string `json:"breed"`
	Gender      Gender `json:"gender"`
}

type Gender int

const (
	Male Gender = iota
	Female
	Other
)

func (c *Cat) SetGender(s string) {
	switch s {
	case "male":
		c.Gender = Male
	case "female":
		c.Gender = Female
	default:
		c.Gender = Other
	}
}

func (g Gender) String() string {
	switch g {
	case Male:
		return "male"
	case Female:
		return "female"
	default:
		return "other"
	}
}
