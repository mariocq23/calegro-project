package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type YamlFile struct {
	Header struct {
		Import      []string `yaml:"import"`
		Id          string   `yaml:"id"`
		Name        string   `yaml:"name"`
		IsChildNode bool     `yaml:"is-child-node"`
		Inherits    string   `yaml:"inherits"`
		Implements  []string `yaml:"implements"`
	} `yaml:"header"`
	Configuration struct {
		User           string `yaml:"user"`
		Agent          string `yaml:"agent"`
		ExecutionMode  string `yaml:"execution-mode"`
		BypassSecurity bool   `yaml:"bypass-security"`
		Security       struct {
			PublicPassword          string `yaml:"public-password"`
			PrivatePasswordLocation string `yaml:"private-password-location"`
			CertificateLocation     string `yaml:"certificate-location"`
			TemplateOrSource        string `yaml:"template-or-source"`
		} `yaml:"security"`
		Location    string `yaml:"location"`
		ContextName string `yaml:"context-name"`
		Encoding    string `yaml:"encoding"`
	} `yaml:"configuration"`
	Action struct {
		Name                 string `yaml:"name"`
		Type                 string `yaml:"type"`
		MainFunction         string `yaml:"main-function"`
		Api                  string `yaml:"api"`
		UseApi               bool   `yaml:"use-api"`
		DisplayOutputConsole bool   `yaml:"display-output-console"`
		Platform             struct {
			ExecutionDependencties []string `yaml:"execution-dependencies"`
		}
		Location             string   `yaml:"location"`
		InitialInputs        []string `yaml:"initial-inputs"`
		EnvironmentVariables []string `yaml:"environment-variables"`
	} `yaml:"action"`
	Contexts []struct {
		Context       string `yaml:"context"`
		ContextAction struct {
			Dependencies struct {
				Location string   `yaml:"location"`
				List     []string `yaml:"list"`
			} `yaml:"dependencies"`
		} `yaml:"context-action"`
		CustomProperties []string `yaml:"custom-properties"`
	} `yaml:"contexts"`
	Steps []struct {
		Name    string `yaml:"name"`
		Pointer string `yaml:"pointer"`
	} `yaml:"steps"`
}

func main() {
	data, err := os.ReadFile("C:/git/calegro-project/examples/vikings/vikings.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var yamlFile YamlFile
	err = yaml.Unmarshal(data, &yamlFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", yamlFile)
}
