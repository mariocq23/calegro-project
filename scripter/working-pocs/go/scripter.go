package main

import (
	"fmt"
	"log"
	"os"
	"scripter/entities"
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

	signal := generateFinalYaml(generatedYamlProperties)

	fmt.Printf("%+v\n", signal)
}

func generateFinalYaml(generatedYamlProperties []entities.YamlProperty) entities.Signal {
	sealedProperties := []string{}
	signal := entities.Signal{}
	signal.Sender = generatedYamlProperties[len(generatedYamlProperties)-1].TemplateName
	for _, prop := range generatedYamlProperties {
		if containsString(sealedProperties, prop.Name) {
			continue
		}
		signal = updateSignal(signal, prop)
		if prop.Sealed {
			sealedProperties = append(sealedProperties, prop.Name)
		}
	}
	return signal
}

func stringToBoolStrict(s string) bool {
	b, err := strconv.ParseBool(strings.ToLower(s))
	if err != nil {
		return false // Return false on error
	}
	return b
}

func removeUnnecessaryString(rawString string) string {
	finalValue := strings.ReplaceAll(rawString, "$(overridable)", "")

	finalValue = strings.TrimSpace(finalValue)

	return finalValue
}

func removeUnnecessaryStringInArray(values []string) []string {
	var result []string
	for _, s := range values {
		if s != "$(overridable)" {
			cleanedString := strings.TrimSpace(s)
			result = append(result, cleanedString)
		}
	}
	return result
}

func updateSignal(signal entities.Signal, prop entities.YamlProperty) entities.Signal {

	if prop.Name == "Configuration.AgentOrLabel" && prop.Value != "" {
		signal.Executor = removeUnnecessaryString(prop.Value)
	}
	if prop.Name == "Configuration.ContextName" && prop.Value != "" {
		signal.Environment = removeUnnecessaryString(prop.Value)
	}
	if prop.Name == "Configuration.ExecutionMode" && prop.Value != "" {
		signal.ExecutionMode = removeUnnecessaryString(prop.Value)
	}
	if prop.Name == "Configuration.BypassSecurity" && prop.Value != "" {
		signal.BypassSecurity = stringToBoolStrict(removeUnnecessaryString(prop.Value))
	}
	if prop.Name == "Configuration.Security.CertificationHub" && prop.Value != "" {
		signal.CertificationHub = removeUnnecessaryString(prop.Value)
	}
	if prop.Name == "Configuration.Security.AuthenticationHub" && prop.Value != "" {
		signal.AuthenticationHub = removeUnnecessaryString(prop.Value)
	}
	if prop.Name == "Configuration.Security.AuthorizationHub" && prop.Value != "" {
		signal.AuthorizationHub = removeUnnecessaryString(prop.Value)
	}
	if prop.Name == "Action.Api" && prop.Value != "" {
		signal.Api = removeUnnecessaryString(prop.Value)
	}
	if prop.Name == "Action.NameOrFullPath" && prop.Value != "" {
		signal.ExecutablePath = removeUnnecessaryString(prop.Value)
	}
	if prop.Name == "Action.Type" && prop.Value != "" {
		signal.Type = removeUnnecessaryString(prop.Value)
	}
	if prop.Name == "Action.OutputMode" && prop.Value != "" {
		signal.OutputMode = removeUnnecessaryString(prop.Value)
	}
	if prop.Name == "Action.ShutdownSignal" && prop.Value != "" {
		signal.ShutdownSignal = removeUnnecessaryString(prop.Value)
	}
	if prop.Name == "Action.Platform.OsFamily" && prop.Value != "" {
		signal.Os = removeUnnecessaryString(prop.Value)
	}
	if prop.Name == "Action.Platform.PackageInstaller" && prop.Value != "" {
		signal.PackageInstaller = removeUnnecessaryString(prop.Value)
	}
	if prop.Name == "Action.Platform.ExecutionDependencies" && prop.Values != nil && len(prop.Values) > 0 {
		signal.InstallationDependencies = removeUnnecessaryStringInArray(prop.Values)
	}
	if prop.Name == "Action.InitialInputs" && prop.Values != nil && len(prop.Values) > 0 {
		signal.Arguments = removeUnnecessaryStringInArray(prop.Values)
	}

	/*contextContextPattern := `^Context\[\d]\.Context$`

	re := regexp.MustCompile(contextContextPattern)

	if re.MatchString(prop.Name) && prop.Value != "" {
		startIndex := strings.Index(prop.Name, "[")
		endIndex := strings.Index(prop.Name, "]")
		substring := prop.Name[startIndex+1 : endIndex]
		_, err := strconv.Atoi(substring)

		if err == nil {
			signal.Contexts = append(signal.Contexts, struct {
				Context      string `yaml:"context"`
				Dependencies struct {
					Location string   `yaml:"location"`
					List     []string `yaml:"list"`
				} `yaml:"dependencies"`
				ContextInitialInputs []string `yaml:"context-initial-inputs"`
				EnvironmentVariables []string `yaml:"environment-variables"`
			}{
				Context: prop.Name,
				Dependencies: struct {
					Location string   `yaml:"location"`
					List     []string `yaml:"list"`
				}{
					Location: "remote",
					List:     []string{"junit"},
				},
				ContextInitialInputs: []string{"./tests"},
				EnvironmentVariables: []string{"ENV=test"},
			})
		}
	}*/

	return signal
}

func filterBy(values []string, filter string) []string {
	var filtered []string
	for _, s := range values {
		if strings.HasPrefix(s, filter) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func reverseYamlArray(yamls []entities.YamlFile) []entities.YamlFile {
	reversed := make([]entities.YamlFile, len(yamls))
	for i, j := 0, len(yamls)-1; i < len(yamls); i, j = i+1, j-1 {
		reversed[i] = yamls[j]
	}
	return reversed
}

func generateYamlProperties(yamls []entities.YamlFile) []entities.YamlProperty {

	yamlProperties := []entities.YamlProperty{}

	for _, yaml := range yamls {

		yamlProperties = append(yamlProperties, generateProperty("Configuration.AgentOrLabel", yaml.Configuration.AgentOrLabel, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.ContextName", yaml.Configuration.ContextName, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.ExecutionMode", yaml.Configuration.ExecutionMode, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.BypassSecurity", strconv.FormatBool(yaml.Configuration.BypassSecurity), yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.Security.CertificationHub", yaml.Configuration.Security.CertificationHub, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.Security.AuthenticationHub", yaml.Configuration.Security.AuthenticationHub, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.Security.AuthorizationHub", yaml.Configuration.Security.AuthorizationHub, yaml.Header.Name))

		yamlProperties = append(yamlProperties, generateProperty("Action.Api", yaml.Action.Api, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Action.NameOrFullPath", yaml.Action.NameOrFullPath, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Action.Type", yaml.Action.Type, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Action.OutputMode", yaml.Action.OutputMode, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Action.ShutdownSignal", yaml.Action.ShutdownSignal, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Action.Platform.OsFamily", yaml.Action.Platform.OsFamily, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Action.Platform.PackageInstaller", yaml.Action.Platform.PackageInstaller, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateArrayProperty("Action.Platform.ExecutionDependencies", yaml.Action.Platform.InstallationDependencies, yaml.Header.Name))
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

			yamlProperties = append(yamlProperties, generateProperty(stepName, step.Step, yaml.Header.Name))
			yamlProperties = append(yamlProperties, generateProperty(stepPointer, step.Pointer, yaml.Header.Name))
		}

	}
	return yamlProperties
}

func generateProperty(name string, value string, templateName string) entities.YamlProperty {
	yamlProperty := entities.YamlProperty{Name: name, Value: value, TemplateName: templateName}
	if !strings.Contains(value, "$(overridable)") {
		yamlProperty.Sealed = true
	}
	if strings.Contains(value, "default") {
		yamlProperty.Default = true
	}
	return yamlProperty
}

func generateArrayProperty(name string, values []string, templateName string) entities.YamlProperty {
	yamlProperty := entities.YamlProperty{Name: name, Values: values, TemplateName: templateName}
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

func readAllYamls(path string) []entities.YamlFile {

	yamlsArray := make([]entities.YamlFile, 0)

	yaml := readYaml(path)

	yamlsArray = append(yamlsArray, yaml)

	if strings.TrimSpace(yaml.Header.Inherits) != "" {
		parentPath, parentName := extractBeforeAndAfterValues(yaml.Header.Import)
		importInherit := entities.ImportInherit{ParentPath: parentPath, ParentName: parentName}
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

func readYaml(filePath string) entities.YamlFile {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var yamlFile entities.YamlFile
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
