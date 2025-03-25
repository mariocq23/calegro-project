package entities

type YamlStepProperty struct {
	Sealed       bool
	Default      bool
	Name         string
	Value        string
	Values       []string
	TemplateName string
	Position     int
}
