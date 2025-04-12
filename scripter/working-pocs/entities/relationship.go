package entities

type Relationship int

const (
	Dependency Relationship = iota
	Step
)
