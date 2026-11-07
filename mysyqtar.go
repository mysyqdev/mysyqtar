package mysyqtar

type Cat struct {
	Name        string
	Description string
	Breed       string
	gender      Gender
}

type Gender int

const (
	Male Gender = iota
	Female
	Other
)
