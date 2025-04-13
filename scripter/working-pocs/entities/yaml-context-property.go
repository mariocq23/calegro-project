package entities

type YamlContextProperty struct {
	Sealed       bool
	Default      bool
	Name         string
	Value        string
	Values       []string
	DictValues   map[string]string // Dictionary as a property
	TemplateName string
	Position     int
}
