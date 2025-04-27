package entities

type Relationship int

const (
	FlowDependency Relationship = iota
	StepDependency
)
