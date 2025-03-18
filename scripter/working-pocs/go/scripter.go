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
	generatedYamlProperties := generateYamlProperties(reverseYamlArray(yamls))

	for _, prop := range generatedYamlProperties {
		fmt.Printf("%+v\n", prop)
	}
}

func reverseYamlArray(yamls []YamlFile) []YamlFile {
	reversed := make([]YamlFile, len(yamls))
	for i, j := 0, len(yamls)-1; i < len(yamls); i, j = i+1, j-1 {
		reversed[i] = yamls[j]
	}
	return reversed
}

func generateYamlProperties(yamls []YamlFile) []YamlProperty {

	yamlProperties := []YamlProperty{}

	for _, yaml := range yamls {

		yamlProperties = append(yamlProperties, generateProperty("Configuration.AgentOrLabel", yaml.Configuration.AgentOrLabel, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.ContextName", yaml.Configuration.ContextName, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.ExecutionMode", yaml.Configuration.ExecutionMode, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.BypassSecurity", strconv.FormatBool(yaml.Configuration.BypassSecurity), yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.Security.CertificateLocation", yaml.Configuration.Security.CertificateLocation, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.Security.PrivatePasswordLocation", yaml.Configuration.Security.PrivatePasswordLocation, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.Security.PublicPassword", yaml.Configuration.Security.PublicPassword, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.Security.User", yaml.Configuration.Security.User, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.Security.TemplateOrSource", yaml.Configuration.Security.TemplateOrSource, yaml.Header.Name))

		yamlProperties = append(yamlProperties, generateProperty("Action.Api", yaml.Action.Api, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Action.NameOrFullPath", yaml.Action.NameOrFullPath, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Action.Type", yaml.Action.Type, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Action.OutputMode", yaml.Action.OutputMode, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Action.ShutdownSignal", yaml.Action.ShutdownSignal, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Action.Platform.OsFamily", yaml.Action.Platform.OsFamily, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Action.Platform.PackageInstaller", yaml.Action.Platform.PackageInstaller, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateArrayProperty("Action.Platform.ExecutionDependencies", yaml.Action.Platform.ExecutionDependencies, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateArrayProperty("Action.InitialInputs", yaml.Action.InitialInputs, yaml.Header.Name))

		for index, context := range yaml.Contexts {

			contextName := fmt.Sprintf("Context[%d].Context", index)
			location := fmt.Sprintf("Context[%d].Dependencies.Location", index)
			list := fmt.Sprintf("Context[%d].Dependencies.List", index)
			contextInitialInputs := fmt.Sprintf("Context[%d].ContextInitialInputs", index)
			environmentVariables := fmt.Sprintf("Context[%d].EnvironmentVariables", index)

			yamlProperties = append(yamlProperties, generateProperty(contextName, context.Context, yaml.Header.Name))
			yamlProperties = append(yamlProperties, generateProperty(location, context.Dependencies.Location, yaml.Header.Name))
			yamlProperties = append(yamlProperties, generateArrayProperty(list, context.Dependencies.List, yaml.Header.Name))
			yamlProperties = append(yamlProperties, generateArrayProperty(contextInitialInputs, context.ContextInitialInputs, yaml.Header.Name))
			yamlProperties = append(yamlProperties, generateArrayProperty(environmentVariables, context.EnvironmentVariables, yaml.Header.Name))
		}
		for index, step := range yaml.Steps {

			stepName := fmt.Sprintf("Step[%d].Name", index)
			stepPointer := fmt.Sprintf("Step[%d].Pointer", index)

			generateProperty(stepName, step.Name, yaml.Header.Name)
			generateProperty(stepPointer, step.Pointer, yaml.Header.Name)
		}

	}
	return yamlProperties
}

func generateProperty(name string, value string, templateName string) YamlProperty {
	yamlProperty := YamlProperty{Name: name, Value: value, TemplateName: templateName}
	if !strings.Contains(value, "$(overridable)") {
		yamlProperty.Sealed = true
	}
	if strings.Contains(value, "default") {
		yamlProperty.Default = true
	}
	return yamlProperty
}

func generateArrayProperty(name string, values []string, templateName string) YamlProperty {
	yamlProperty := YamlProperty{Name: name, Values: values, TemplateName: templateName}
	if !containsString(values, "$(overridable)") {
		yamlProperty.Sealed = true
	}
	if containsString(values, "default") {
		yamlProperty.Default = true
	}
	return yamlProperty
}

func containsString(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
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
