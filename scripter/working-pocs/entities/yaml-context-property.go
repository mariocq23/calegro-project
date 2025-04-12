package entities

type YamlContextProperty struct {
	Sealed       bool
	Default      bool
	Name         string
	Value        string
	Values       []string
	TemplateName string
	Position     int
}
