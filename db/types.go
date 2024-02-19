package db

type PersonalityTrait struct {
	Trait string
}

type Persona struct {
	FirstName         string `db:first_name`
	LastName          string `db:last_name`
	PersonalityTraits []PersonalityTrait
}

type User struct {
	Username string
}

type Message struct {
	Author  string `json:"role"`
	Content string `json:"content"`
}
