package types

type PersonalityTrait struct {
	Trait string
}

type Persona struct {
	FirstName         string `db:first_name`
	LastName          string `db:last_name`
	PersonalityTraits []PersonalityTrait
}
