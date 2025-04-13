package entities

type YamlProperty struct {
	Sealed       bool
	Default      bool
	Name         string
	BoolValue    bool
	Value        string
	Values       []string
	DictValues   map[string]string
	TemplateName string
}
