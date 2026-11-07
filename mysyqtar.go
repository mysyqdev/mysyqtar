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
