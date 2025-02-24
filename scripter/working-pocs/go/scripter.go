package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	filePath := os.Args[1]
	yamls := readAllYamls(filePath)
	combinedYaml := combineYamls(reverseYamlArray(yamls))
}

func reverseYamlArray(yamls []YamlFile) []YamlFile {
	reversed := make([]YamlFile, len(yamls))
	for i, j := 0, len(yamls)-1; i < len(yamls); i, j = i+1, j-1 {
		reversed[i] = yamls[j]
	}
	return reversed
}

func combineYamls(yamls []YamlFile) YamlFile {

	yamlProperties := new([]YamlProperty)

	finalYaml := new(YamlFile)

	for index, yaml := range yamls {

		setPropertyValue("AgentOrLabel", yaml.Configuration.AgentOrLabel, finalYaml, yamlProperties)
		setPropertyValue("ContextName", yaml.Configuration.ContextName, finalYaml, yamlProperties)
		setPropertyValue("ExecutionMode", yaml.Configuration.ExecutionMode, finalYaml, yamlProperties)
		setPropertyValue("BypassSecurity", strconv.FormatBool(yaml.Configuration.BypassSecurity), finalYaml, yamlProperties)
		setPropertyValue("CertificateLocation", yaml.Configuration.Security.CertificateLocation, finalYaml, yamlProperties)
		setPropertyValue("PrivatePasswordLocation", yaml.Configuration.Security.PrivatePasswordLocation, finalYaml, yamlProperties)
		setPropertyValue("PublicPassword", yaml.Configuration.Security.PublicPassword, finalYaml, yamlProperties)
		setPropertyValue("User", yaml.Configuration.Security.User, finalYaml, yamlProperties)
		setPropertyValue("TemplateOrSource", yaml.Configuration.Security.TemplateOrSource, finalYaml, yamlProperties)
		setPropertyValue("Api", yaml.Action.Api, finalYaml, yamlProperties)
		setPropertyValue("NameOrFullPath", yaml.Action.NameOrFullPath, finalYaml, yamlProperties)
		setPropertyValue("AgentOrLabel", yaml.Configuration.AgentOrLabel, finalYaml, yamlProperties)
		setPropertyValue("AgentOrLabel", yaml.Configuration.AgentOrLabel, finalYaml, yamlProperties)
		setPropertyValue("AgentOrLabel", yaml.Configuration.AgentOrLabel, finalYaml, yamlProperties)
		setPropertyValue("AgentOrLabel", yaml.Configuration.AgentOrLabel, finalYaml, yamlProperties)
		setPropertyValue("AgentOrLabel", yaml.Configuration.AgentOrLabel, finalYaml, yamlProperties)
		setPropertyValue("AgentOrLabel", yaml.Configuration.AgentOrLabel, finalYaml, yamlProperties)
		setPropertyValue("AgentOrLabel", yaml.Configuration.AgentOrLabel, finalYaml, yamlProperties)
		setPropertyValue("AgentOrLabel", yaml.Configuration.AgentOrLabel, finalYaml, yamlProperties)
		setPropertyValue("AgentOrLabel", yaml.Configuration.AgentOrLabel, finalYaml, yamlProperties)
		setPropertyValue("AgentOrLabel", yaml.Configuration.AgentOrLabel, finalYaml, yamlProperties)
		setPropertyValue("AgentOrLabel", yaml.Configuration.AgentOrLabel, finalYaml, yamlProperties)
		setPropertyValue("AgentOrLabel", yaml.Configuration.AgentOrLabel, finalYaml, yamlProperties)
		setPropertyValue("AgentOrLabel", yaml.Configuration.AgentOrLabel, finalYaml, yamlProperties)
		setPropertyValue("AgentOrLabel", yaml.Configuration.AgentOrLabel, finalYaml, yamlProperties)
		setPropertyValue("AgentOrLabel", yaml.Configuration.AgentOrLabel, finalYaml, yamlProperties)
	}

}

func setPropertyValue(name string, value string, finalYaml *YamlFile, yamlProperties *[]YamlProperty) {
	yamlProperty := YamlProperty{Name: name, Value: value}
	if !strings.Contains(value, "$(overridable)") {
		yamlProperty.Sealed = true
	}
	if strings.Contains(value, "default") {
		yamlProperty.Default = true
	}
}

func readAllYamls(path string) []YamlFile {

	yamlsArray := make([]YamlFile, 0)

	yaml := readYaml(path)

	yamlsArray = append(yamlsArray, yaml)

	if strings.TrimSpace(yaml.Header.Inherits) != "" {
		parentPath, parentName := extractBeforeAndAfterValues(yaml.Header.Import)
		importInherit := ImportInherit{ParentPath: parentPath, ParentName: parentName}
		newYamlArray := readAllYamls(importInherit.ParentPath)
		yamlsArray = append(yamlsArray, newYamlArray...)
	}
	return yamlsArray
}

func extractBeforeAndAfterValues(input string) (string, string) {
	parts := strings.Split(input, "=>")
	if len(parts) == 2 {
		before := strings.TrimSpace(parts[0])
		after := strings.TrimSpace(parts[1])
		return before, after
	}
	return "", "" // Return empty strings if the split doesn't produce two parts
}

func readYaml(filePath string) YamlFile {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var yamlFile YamlFile
	err = yaml.Unmarshal(data, &yamlFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", yamlFile)

	return yamlFile

}

// func interpretYaml(yaml YamlFile) {
// 	//Parent

// 	if strings.TrimSpace(yaml.Header.Inherits) != "" {
// 		var parentYamlFile = readYaml(yaml.Header.Inherits)
// 		interpretYaml(parentYamlFile)
// 	}

// 	checkConfiguration(yaml)

// 	//Security
// 	//Chosen Context
// 	//Steps if any, otherwise the solely Action
// }

// func checkConfiguration(yaml YamlFile) {

// }

// func setLabels(labels []string) {

// }

// func execute(command string, args []string) {
// 	// Example with a config file:
// 	//cmd := exec.Command("dosbox", "-conf", "my_dosbox.conf")

// 	// Example with commands (using -c):

// 	// Capture output (optional)
// 	cmd := exec.Command(command, args...)
// 	out, err := cmd.CombinedOutput()

// 	if err != nil {
// 		log.Fatalf("cmd.Run() failed with %s\n", err)
// 	}
// 	fmt.Printf("combined out:\n%s\n", string(out))

// 	err = cmd.Start()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(command + " started with config/commands!")
// }
