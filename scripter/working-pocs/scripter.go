package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"scripter/entities"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	filePath := os.Args[1]
	originatorPath := os.Args[2]
	yamls := readAllYamls(filePath)
	for _, yaml := range yamls {
		fmt.Printf("%+v\n", yaml)
	}
	generalProperties, contextProperties, signalSteps := generateYamlProperties(yamls)
	for _, prop := range generalProperties {
		fmt.Printf("%+v\n", prop)
	}
	for _, prop := range contextProperties {
		fmt.Printf("%+v\n", prop)
	}
	for _, prop := range signalSteps {
		fmt.Printf("%+v\n", prop)
	}

	signal := generateSignal(generalProperties, contextProperties, signalSteps, originatorPath)

	fmt.Printf("%+v\n", signal)
}

func generateSignal(generalProperties []entities.YamlProperty, contextProperties []entities.YamlContextProperty, steps []entities.SignalStep, originatorPath string) entities.Signal {
	sealedProperties := []string{}
	signal := entities.Signal{}
	signal.Sender = generalProperties[len(generalProperties)-1].TemplateName
	for _, prop := range generalProperties {
		if containsString(sealedProperties, prop.Name) {
			continue
		}
		signal = updateSignalGeneralProperties(signal, prop)
		if prop.Sealed {
			sealedProperties = append(sealedProperties, prop.Name)
		}
	}

	signal.Steps = steps

	if originatorPath != "" {
		signal.OriginatorQuay.SourceOrPath = originatorPath
		signal.OriginatorQuay.Name = getFilenameWithoutExtension(originatorPath)
	}

	if signal.Environment == "default" {
		return signal
	}

	signal = updateSignalContextProperties(signal, contextProperties)

	return signal
}

func getFilenameWithoutExtension(path string) string {
	filename := filepath.Base(path)
	extension := filepath.Ext(filename)

	if extension != "" {
		return strings.TrimSuffix(filename, extension)
	}
	return filename
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

func updateSignalGeneralProperties(signal entities.Signal, prop entities.YamlProperty) entities.Signal {

	if prop.Name == "Configuration.AgentOrLabel" && prop.Value != "" {
		signal.Executor = removeUnnecessaryString(prop.Value)
	}
	if prop.Name == "Configuration.ContextName" && prop.Value != "" {
		signal.Environment = removeUnnecessaryString(prop.Value)
	}
	if prop.Name == "Configuration.ExecutionMode" && prop.Value != "" {
		signal.ExecutionMode = removeUnnecessaryString(prop.Value)
	}
	if prop.Name == "Configuration.BypassSecurity" {
		signal.BypassSecurity = prop.BoolValue
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
	if prop.Name == "Action.Platform.InstallationDependencies" && prop.Values != nil && len(prop.Values) > 0 {
		signal.InstallationDependencies = removeUnnecessaryStringInArray(prop.Values)
	}
	if prop.Name == "Action.InitialInputs" && prop.Values != nil && len(prop.Values) > 0 {
		signal.Arguments = removeUnnecessaryStringInArray(prop.Values)
	}
	if prop.Name == "Action.EnvironmentVariables" && prop.DictValues != nil && len(prop.DictValues) > 0 {
		signal.EnvironmentVariables = prop.DictValues
	}
	if prop.Name == "Action.ExecutionDependencies" && prop.Values != nil && len(prop.Values) > 0 {
		signal.ExecutionDependencies = prop.Values
	}

	return signal
}

// Contexts []struct {
// 	Context      string `yaml:"context"`
// 	Dependencies struct {
// 		Location string   `yaml:"location"`
// 		List     []string `yaml:"list"`
// 	} `yaml:"dependencies"`
// 	ContextInitialInputs []string `yaml:"context-initial-inputs"`
// 	EnvironmentVariables []string `yaml:"environment-variables"`
// } `yaml:"contexts"`

func updateSignalContextProperties(signal entities.Signal, props []entities.YamlContextProperty) entities.Signal {
	chosenContextIndex := 0

	for _, item := range props {
		if item.Name == "Context.Context" && item.Value == signal.Environment {
			chosenContextIndex = item.Position
		}
	}

	for index, item := range props {
		if index != chosenContextIndex {
			continue
		}
		if item.Name == "Context.Dependencies" {
			signal.ExecutionDependencies = item.Values
		}
		if item.Name == "Context.EnvironmentVariables" {
			signal.EnvironmentVariables = item.DictValues
		}
		if item.Name == "Context.ContextInitialInputs" {
			signal.Arguments = item.Values
		}
	}

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

func reverseYamlArray(yamls []*entities.YamlFile) []*entities.YamlFile {
	reversed := make([]*entities.YamlFile, len(yamls))
	for i, j := 0, len(yamls)-1; i < len(yamls); i, j = i+1, j-1 {
		reversed[i] = yamls[j]
	}
	return reversed
}

func generateYamlProperties(yamls []*entities.YamlFile) ([]entities.YamlProperty, []entities.YamlContextProperty, []entities.SignalStep) {

	yamlProperties := []entities.YamlProperty{}
	yamlContextProperties := []entities.YamlContextProperty{}
	signalSteps := []entities.SignalStep{}
	finalSignalSteps := []entities.SignalStep{}
	overridableSteps := []string{}

	for _, yaml := range yamls {

		yamlProperties = append(yamlProperties, generateProperty("Configuration.AgentOrLabel", yaml.Configuration.AgentOrLabel, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.ContextName", yaml.Configuration.ContextName, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.ExecutionMode", yaml.Configuration.ExecutionMode, yaml.Header.Name))
		if yaml.Configuration.BypassSecurity != nil {
			yamlProperties = append(yamlProperties, generateBoolProperty("Configuration.BypassSecurity", *yaml.Configuration.BypassSecurity, yaml.Header.Name))
		}
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
		yamlProperties = append(yamlProperties, generateArrayProperty("Action.Platform.InstallationDependencies", yaml.Action.Platform.InstallationDependencies, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateArrayProperty("Action.InitialInputs", yaml.Action.InitialInputs, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateDictionaryProperty("Action.EnvironmentVariables", stringListToMap(removeUnnecessaryStringInArray(yaml.Action.EnvironmentVariables)), yaml.Header.Name))

		for index, context := range yaml.Contexts {
			yamlContextProperties = append(yamlContextProperties, generateContextProperty("Context.Context", context.Context, yaml.Header.Name, index))
			yamlContextProperties = append(yamlContextProperties, generateContextArrayProperty("Context.Dependencies", context.Dependencies, yaml.Header.Name, index))
			yamlContextProperties = append(yamlContextProperties, generateContextArrayProperty("Context.ContextInitialInputs", context.ContextInitialInputs, yaml.Header.Name, index))
			yamlContextProperties = append(yamlContextProperties, generateContextDictionaryProperty("Context.EnvironmentVariables", stringListToMap(removeUnnecessaryStringInArray(context.EnvironmentVariables)), yaml.Header.Name, index))
		}
		for index, step := range yaml.Steps {
			if !strings.Contains(step.Step, "$(overridable)") && strings.Trim(step.Step, " ") != "" && !strings.Contains(step.Pointer, "$(overridable)") && strings.Trim(step.Pointer, " ") != "" {
				signalSteps = append(signalSteps, generateSignalStep(step.Step, step.Pointer, index, yaml.Header.Name))
			} else {
				overridableSteps = append(overridableSteps, yaml.Header.Name)
			}
		}

		if len(yaml.Steps) == 0 {
			overridableSteps = append(overridableSteps, yaml.Header.Name)
		}
	}

	//"$(append)"
	//for _, step := range overridableSteps {
	//if   {
	//	finalSignalSteps = append(finalSignalSteps, step)
	//}
	//}

	//"$(append)"

	for _, step := range signalSteps {
		ancestors := obtainAllSourceAncestors(step.Source, yamls)

		if containsString(overridableSteps, step.Source) {
			finalSignalSteps = append(finalSignalSteps, step)
		} else {
			for _, ancestor := range ancestors {
				if containsString(overridableSteps, ancestor) {
					finalSignalSteps = appendNew(finalSignalSteps, step)
				}
			}
		}
	}

	return yamlProperties, yamlContextProperties, finalSignalSteps
}

func stringListToMap(input []string) map[string]string {
	resultMap := make(map[string]string)

	for _, item := range input {
		parts := strings.SplitN(item, " ", 2)
		if len(parts) == 2 {
			key := strings.Trim(parts[0], "()") // Remove parentheses
			resultMap[key] = parts[1]
		}
	}

	return resultMap
}

func appendNew(steps []entities.SignalStep, currentStep entities.SignalStep) []entities.SignalStep {
	if !containsObject(steps, currentStep) {
		steps = append(steps, currentStep)
	}
	return steps
}

func readAllYamls(path string) []*entities.YamlFile {

	yamlsArray := make([]*entities.YamlFile, 0)

	yaml := readYaml(path)

	if strings.TrimSpace(yaml.Header.Inherits) != "" {
		parentPath, parentName := extractBeforeAndAfterValues(yaml.Header.Inherits)
		importInherit := entities.ImportInherit{ParentPath: parentPath, ParentName: parentName}
		yaml.Parent = readYaml(importInherit.ParentPath)
		newYamlArray := readAllYamls(importInherit.ParentPath)
		yamlsArray = append(yamlsArray, newYamlArray...)
	}

	yamlsArray = append(yamlsArray, yaml)

	return yamlsArray
}

func obtainAllSourceAncestors(templateName string, yamls []*entities.YamlFile) []string {
	yamlsArray := make([]string, 0)

	template := findYamlByName(yamls, templateName)

	if template != nil && strings.TrimSpace(template.Header.Inherits) != "" {
		parentPath, parentName := extractBeforeAndAfterValues(template.Header.Inherits)
		importInherit := entities.ImportInherit{ParentPath: parentPath, ParentName: parentName}
		newYamlArray := obtainAllSourceAncestors(importInherit.ParentName, yamls)
		yamlsArray = append(yamlsArray, newYamlArray...)
	}

	yamlsArray = append(yamlsArray, template.Header.Name)

	return yamlsArray

}

func findYamlByName(yamls []*entities.YamlFile, templateName string) *entities.YamlFile {
	for _, yaml := range yamls {
		if yaml.Header.Name == templateName {
			return yaml
		}
	}
	return nil
}

func generateSignalStep(step string, pointer string, index int, templateName string) entities.SignalStep {
	signalStep := entities.SignalStep{
		Name:    step,
		Pointer: pointer,
		Source:  templateName,
	}

	return signalStep
}

func generateProperty(name string, value string, templateName string) entities.YamlProperty {
	yamlProperty := entities.YamlProperty{Name: name, Value: value, TemplateName: templateName}
	if !strings.Contains(value, "$(overridable)") && value != "" {
		yamlProperty.Sealed = true
	}
	if strings.Contains(value, "default") {
		yamlProperty.Default = true
	}
	return yamlProperty
}

func generateBoolProperty(name string, value bool, templateName string) entities.YamlProperty {
	yamlProperty := entities.YamlProperty{Name: name, BoolValue: value, TemplateName: templateName}
	if !value && name == "Configuration.BypassSecurity" {
		yamlProperty.Sealed = true
		yamlProperty.Default = true
	}
	return yamlProperty
}

func generateContextProperty(name string, value string, templateName string, index int) entities.YamlContextProperty {
	yamlProperty := entities.YamlContextProperty{Name: name, Value: value, TemplateName: templateName, Position: index}
	if !strings.Contains(value, "$(overridable)") && value != "" {
		yamlProperty.Sealed = true
	}
	if strings.Contains(value, "default") {
		yamlProperty.Default = true
	}
	return yamlProperty
}

func generateArrayProperty(name string, values []string, templateName string) entities.YamlProperty {
	yamlProperty := entities.YamlProperty{Name: name, Values: values, TemplateName: templateName}
	if !containsString(values, "$(overridable)") && values != nil && len(values) > 0 {
		yamlProperty.Sealed = true
	}
	if containsString(values, "default") {
		yamlProperty.Default = true
	}
	return yamlProperty
}

func generateContextArrayProperty(name string, values []string, templateName string, index int) entities.YamlContextProperty {
	yamlProperty := entities.YamlContextProperty{Name: name, Values: values, TemplateName: templateName, Position: index}
	if !containsString(values, "$(overridable)") && values != nil && len(values) > 0 {
		yamlProperty.Sealed = true
	}
	if containsString(values, "default") {
		yamlProperty.Default = true
	}
	return yamlProperty
}

func generateDictionaryProperty(name string, values map[string]string, templateName string) entities.YamlProperty {
	yamlProperty := entities.YamlProperty{Name: name, DictValues: values, TemplateName: templateName}
	if !containsKeyValuePair(values, "$(overridable)") && values != nil && len(values) > 0 {
		yamlProperty.Sealed = true
	}
	if containsKeyValuePair(values, "default") {
		yamlProperty.Default = true
	}
	return yamlProperty
}

func generateContextDictionaryProperty(name string, values map[string]string, templateName string, index int) entities.YamlContextProperty {
	yamlProperty := entities.YamlContextProperty{Name: name, DictValues: values, TemplateName: templateName, Position: index}
	if !containsKeyValuePair(values, "$(overridable)") && values != nil && len(values) > 0 {
		yamlProperty.Sealed = true
	}
	if containsKeyValuePair(values, "default") {
		yamlProperty.Default = true
	}
	return yamlProperty
}

func containsKeyValuePair(values map[string]string, target string) bool {
	for key, value := range values {
		if key == target || value == target {
			return true
		}
	}
	return false
}

func containsObject(steps []entities.SignalStep, step entities.SignalStep) bool {
	arrVal := reflect.ValueOf(steps)
	if arrVal.Kind() != reflect.Slice {
		panic("arr must be a slice")
	}

	for i := 0; i < arrVal.Len(); i++ {
		if reflect.DeepEqual(arrVal.Index(i).Interface(), step) {
			return true
		}
	}
	return false
}

func containsString(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
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

func readYaml(filePath string) *entities.YamlFile {

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var yamlFile entities.YamlFile
	err = yaml.Unmarshal(data, &yamlFile)
	if err != nil {
		log.Fatal(err)
	}

	return &yamlFile
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
