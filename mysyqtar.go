package mysyqtar

type Cat struct {
	Name        string
	Description string
	Breed       string
	Gender      Gender
}

type Gender int

const (
	Male Gender = iota
	Female
	Other
)
