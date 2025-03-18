package main

type YamlProperty struct {
	Sealed       bool
	Default      bool
	Name         string
	Value        string
	Values       []string
	TemplateName string
}
