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
	AuthorId    int
	Content     string
	OrderNumber int
}