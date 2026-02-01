package utilities

import (
	"reflect"
	"runtime"
	"scripter/entities"
	"scripter/entities/versions"
	"strings"
)

type ObjectHandler struct{}

var fileReader = FileReader{}

//Object Array Generator - More context logic related

func (objectHandler ObjectHandler) GenerateYamlProperties(yamls []*versions.YamlFile_Generic_01) ([]entities.YamlProperty, []entities.YamlContextProperty, []entities.SignalStep, []entities.Label) {

	yamlProperties := []entities.YamlProperty{}
	yamlContextProperties := []entities.YamlContextProperty{}
	signalSteps := []entities.SignalStep{}
	finalSignalSteps := []entities.SignalStep{}
	overridableSteps := []string{}
	labels := []entities.Label{}

	for _, yaml := range yamls {
		if yaml.Configuration.Containerize != nil {
			yamlProperties = append(yamlProperties, generateBoolProperty("Configuration.Containerize", yaml.Configuration.Containerize, yaml.Header.Name, yaml.Configuration.CanOverwrite))
		}
		yamlProperties = append(yamlProperties, generateProperty("Configuration.AgentOrLabel", yaml.Configuration.AgentOrLabel, yaml.Header.Name, yaml.Configuration.CanOverwrite))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.ContextName", yaml.Configuration.ContextName, yaml.Header.Name, yaml.Configuration.CanOverwrite))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.ExecutionMode", yaml.Configuration.ExecutionMode, yaml.Header.Name, yaml.Configuration.CanOverwrite))
		if yaml.Configuration.BypassSecurity != nil {
			yamlProperties = append(yamlProperties, generateBoolProperty("Configuration.BypassSecurity", yaml.Configuration.BypassSecurity, yaml.Header.Name, yaml.Configuration.CanOverwrite))
		}
		yamlProperties = append(yamlProperties, generateProperty("Configuration.Security.CertificationHub", yaml.Configuration.Security.CertificationHub, yaml.Header.Name, yaml.Configuration.CanOverwrite))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.Security.AuthenticationHub", yaml.Configuration.Security.AuthenticationHub, yaml.Header.Name, yaml.Configuration.CanOverwrite))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.Security.AuthorizationHub", yaml.Configuration.Security.AuthorizationHub, yaml.Header.Name, yaml.Configuration.CanOverwrite))

		yamlProperties = append(yamlProperties, generateProperty("Action.Api", yaml.Action.Api, yaml.Header.Name, yaml.Action.CanOverwrite))
		yamlProperties = append(yamlProperties, generateProperty("Action.NameOrFullPath", yaml.Action.NameOrFullPath, yaml.Header.Name, yaml.Action.CanOverwrite))
		yamlProperties = append(yamlProperties, generateProperty("Action.Type", yaml.Action.Type, yaml.Header.Name, yaml.Action.CanOverwrite))
		yamlProperties = append(yamlProperties, generateProperty("Action.ShutdownSignal", yaml.Action.ShutdownSignal, yaml.Header.Name, yaml.Action.CanOverwrite))
		yamlProperties = append(yamlProperties, generateProperty("Action.Platform.OsFamily", yaml.Action.Platform.OsFamily, yaml.Header.Name, yaml.Action.CanOverwrite))
		yamlProperties = append(yamlProperties, generateProperty("Action.Platform.PackageInstaller", yaml.Action.Platform.PackageInstaller, yaml.Header.Name, yaml.Action.CanOverwrite))
		yamlProperties = append(yamlProperties, generateArrayProperty("Action.Platform.InstallationDependencies", yaml.Action.InstallationDependencies, yaml.Header.Name, yaml.Action.CanOverwrite))
		yamlProperties = append(yamlProperties, generateArrayProperty("Action.InitialInputs", yaml.Action.InitialInputs, yaml.Header.Name, yaml.Action.CanOverwrite))
		labels = append(labels, generateLabels(yaml.Header.Labels, yaml.Header.Name)...)
		yamlProperties = append(yamlProperties, generateDictionaryProperty("Action.EnvironmentVariables", stringHandler.StringListToMap(stringHandler.RemoveUnnecessaryStringInArray(yaml.Action.EnvironmentVariables)), yaml.Header.Name, yaml.Action.CanOverwrite))

		for index, context := range yaml.Environment.Contexts {
			yamlContextProperties = append(yamlContextProperties, generateContextProperty("Environment.Context.Context", context.Context, yaml.Header.Name, index, yaml.Environment.CanOverwrite))
			yamlContextProperties = append(yamlContextProperties, generateContextArrayProperty("Environment.Context.Dependencies", context.Dependencies, yaml.Header.Name, index, yaml.Environment.CanOverwrite))
			yamlContextProperties = append(yamlContextProperties, generateContextArrayProperty("Environment.Context.ContextInitialInputs", context.ContextInitialInputs, yaml.Header.Name, index, yaml.Environment.CanOverwrite))
			yamlContextProperties = append(yamlContextProperties, generateContextDictionaryProperty("Environment.Context.EnvironmentVariables", stringHandler.StringListToMap(stringHandler.RemoveUnnecessaryStringInArray(context.EnvironmentVariables)), yaml.Header.Name, index, yaml.Environment.CanOverwrite))
		}

		if yaml.Steps.CanOverwrite != nil && (*yaml.Steps.CanOverwrite || len(yaml.Steps.List) == 0) {
			overridableSteps = append(overridableSteps, yaml.Header.Name)
		} else {
			signalSteps = append(signalSteps, generateSignalSteps(yaml)...)
		}

	}

	for _, step := range signalSteps {
		ancestors := fileReader.ObtainAllSourceAncestors(step.Source, yamls)

		if stringHandler.ContainsString(overridableSteps, step.Source) {
			finalSignalSteps = append(finalSignalSteps, step)
		} else {
			for _, ancestor := range ancestors {
				if stringHandler.ContainsString(overridableSteps, ancestor) {
					finalSignalSteps = appendNew(finalSignalSteps, step)
				}
			}
		}
	}

	return yamlProperties, yamlContextProperties, finalSignalSteps, labels
}

func generateLabels(rawlabels []string, templateName string) []entities.Label {
	labels := []entities.Label{}
	for _, rawLabel := range rawlabels {
		label := entities.Label{
			Label:    rawLabel,
			Template: templateName,
		}
		labels = append(labels, label)
	}
	return labels
}

func generateSignalSteps(yaml *versions.YamlFile_Generic_01) []entities.SignalStep {
	signalSteps := []entities.SignalStep{}
	for _, step := range yaml.Steps.List {
		signalStep := entities.SignalStep{
			Name:    step.Step,
			Pointer: step.Pointer,
			Source:  yaml.Header.Name,
		}
		signalSteps = append(signalSteps, signalStep)
	}
	return signalSteps
}

func generateContextArrayProperty(name string, values []string, templateName string, index int, override *bool) entities.YamlContextProperty {
	yamlProperty := entities.YamlContextProperty{Name: name, Values: values, TemplateName: templateName, Position: index}
	if override != nil {
		yamlProperty.Sealed = !*override
	}
	if stringHandler.ContainsString(values, "default") {
		yamlProperty.Default = true
	}
	return yamlProperty
}

func appendNew(steps []entities.SignalStep, currentStep entities.SignalStep) []entities.SignalStep {
	if !containsObject(steps, currentStep) {
		steps = append(steps, currentStep)
	}
	return steps
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

func generateContextDictionaryProperty(name string, values map[string]string, templateName string, index int, override *bool) entities.YamlContextProperty {
	yamlProperty := entities.YamlContextProperty{Name: name, DictValues: values, TemplateName: templateName, Position: index}
	if override != nil {
		yamlProperty.Sealed = !*override
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

//Objectj generator - more context logic related

func generateBoolProperty(name string, value *bool, templateName string, override *bool) entities.YamlProperty {
	yamlProperty := entities.YamlProperty{Name: name, BoolValue: value, TemplateName: templateName}
	if override != nil && (name == "Configuration.BypassSecurity" || name == "Configuration.Containerize") {
		yamlProperty.Sealed = !*override
	}
	return yamlProperty
}

//Objectj generator - more context logic related

func generateProperty(name string, value string, templateName string, override *bool) entities.YamlProperty {
	yamlProperty := entities.YamlProperty{Name: name, Value: value, TemplateName: templateName}
	if override != nil {
		yamlProperty.Sealed = !*override
	}
	if strings.Contains(value, "default") {
		yamlProperty.Default = true
	}
	return yamlProperty
}

//Objectj generator - more context logic related

func generateArrayProperty(name string, values []string, templateName string, override *bool) entities.YamlProperty {
	yamlProperty := entities.YamlProperty{Name: name, Values: values, TemplateName: templateName}

	if override != nil {
		yamlProperty.Sealed = !*override
	}

	if !stringHandler.ContainsString(values, "default") {
		yamlProperty.Default = true
	}
	return yamlProperty
}

func generateContextProperty(name string, value string, templateName string, index int, override *bool) entities.YamlContextProperty {
	yamlProperty := entities.YamlContextProperty{Name: name, Value: value, TemplateName: templateName, Position: index}

	if override != nil {
		yamlProperty.Sealed = !*override
	}

	if strings.Contains(value, "default") {
		yamlProperty.Default = true
	}
	return yamlProperty
}

func getDistinctLabels(labels []entities.Label) []string {
	// Use a map to store unique labels encountered.
	// The value (empty struct{}) is a common idiom for a set in Go,
	// as it takes up no memory.
	seenLabels := make(map[string]struct{})
	var distinctLabels []string // This slice will store the distinct labels

	for _, l := range labels {
		// Check if the label has already been added to our set
		if _, exists := seenLabels[l.Label]; !exists {
			// If not, add it to the map
			seenLabels[l.Label] = struct{}{}
			// And append it to our result slice
			distinctLabels = append(distinctLabels, l.Label)
		}
	}

	return distinctLabels
}

//Objectj generator - more context logic related

func generateDictionaryProperty(name string, values map[string]string, templateName string, override *bool) entities.YamlProperty {
	yamlProperty := entities.YamlProperty{Name: name, DictValues: values, TemplateName: templateName}
	if override != nil {
		yamlProperty.Sealed = !*override
	}

	if containsKeyValuePair(values, "default") {
		yamlProperty.Default = true
	}
	return yamlProperty
}

func (objectHandler ObjectHandler) GenerateSignal(generalProperties []entities.YamlProperty, contextProperties []entities.YamlContextProperty, steps []entities.SignalStep, labels []entities.Label, originatorPath string, nickname string, requireAcknowledge string) entities.Signal {
	sealedProperties := []string{}
	signal := entities.Signal{}
	signal.Sender = generalProperties[len(generalProperties)-1].TemplateName
	signal.HostOs = runtime.GOOS
	for _, prop := range generalProperties {
		if stringHandler.ContainsString(sealedProperties, prop.Name) {
			continue
		}
		signal = updateSignalGeneralProperties(signal, prop)
		if prop.Sealed {
			sealedProperties = append(sealedProperties, prop.Name)
		}
	}

	if originatorPath != "" {
		signal.OriginatorQuay.SourceOrPath = originatorPath
		signal.OriginatorQuay.Name = stringHandler.GetFilenameWithoutExtension(originatorPath)
		signal.OriginatorQuay.ProcessName = nickname
		signal.OriginatorQuay.RequireAcknowledge = stringHandler.InterpretStringAsBool(requireAcknowledge)
	}

	signal.Labels = getDistinctLabels(labels)

	if signal.Environment == "default" {
		return signal
	}

	signal = updateSignalContextProperties(signal, contextProperties)

	signal.EmitQuays = generateEmitQuays(signal, steps)

	return signal
}

func generateEmitQuays(signal entities.Signal, steps []entities.SignalStep) []entities.EmitQuay {
	emitQuays := make([]entities.EmitQuay, 0)

	for index, executionDependency := range signal.ExecutionDependencies {
		emitQuay := entities.EmitQuay{Name: stringHandler.GetFilenameWithoutExtension(executionDependency), Path: executionDependency, Relationship: entities.FlowDependency, Priority: index}
		emitQuays = append(emitQuays, emitQuay)
	}

	for index, step := range steps {
		emitQuay := entities.EmitQuay{Name: stringHandler.GetFilenameWithoutExtension(step.Pointer), Path: step.Pointer, Relationship: entities.StepDependency, Priority: index}
		emitQuays = append(emitQuays, emitQuay)
	}

	return emitQuays
}

//	if prop.Name == "Header.Labels" {
//		signal.Labels = prop.Values
//	}
func updateSignalGeneralProperties(signal entities.Signal, prop entities.YamlProperty) entities.Signal {
	if prop.Name == "Configuration.Containerize" {
		signal.Containerize = *prop.BoolValue
	}
	if prop.Name == "Configuration.AgentOrLabel" && prop.Value != "" {
		signal.Executor = stringHandler.RemoveUnnecessaryString(prop.Value)
	}
	if prop.Name == "Configuration.ContextName" && prop.Value != "" {
		signal.Environment = stringHandler.RemoveUnnecessaryString(prop.Value)
	}
	if prop.Name == "Configuration.ExecutionMode" && prop.Value != "" {
		signal.ExecutionMode = stringHandler.RemoveUnnecessaryString(prop.Value)
	}
	if prop.Name == "Configuration.BypassSecurity" {
		signal.BypassSecurity = *prop.BoolValue
	}
	if prop.Name == "Configuration.Security.CertificationHub" && prop.Value != "" {
		signal.CertificationHub = stringHandler.RemoveUnnecessaryString(prop.Value)
	}
	if prop.Name == "Configuration.Security.AuthenticationHub" && prop.Value != "" {
		signal.AuthenticationHub = stringHandler.RemoveUnnecessaryString(prop.Value)
	}
	if prop.Name == "Configuration.Security.AuthorizationHub" && prop.Value != "" {
		signal.AuthorizationHub = stringHandler.RemoveUnnecessaryString(prop.Value)
	}
	if prop.Name == "Action.Api" && prop.Value != "" {
		signal.Api = stringHandler.RemoveUnnecessaryString(prop.Value)
	}
	if prop.Name == "Action.NameOrFullPath" && prop.Value != "" {
		signal.ExecutablePath = stringHandler.RemoveUnnecessaryString(prop.Value)
	}
	if prop.Name == "Action.Type" && prop.Value != "" {
		signal.Type = stringHandler.RemoveUnnecessaryString(prop.Value)
	}
	if prop.Name == "Action.ShutdownSignal" && prop.Value != "" {
		signal.ShutdownSignal = stringHandler.RemoveUnnecessaryString(prop.Value)
	}
	if prop.Name == "Action.Platform.OsFamily" && prop.Value != "" {
		signal.SignalOs = stringHandler.RemoveUnnecessaryString(prop.Value)
	}
	if prop.Name == "Action.Platform.PackageInstaller" && prop.Value != "" {
		signal.PackageInstaller = stringHandler.RemoveUnnecessaryString(prop.Value)
	}
	if prop.Name == "Action.Platform.InstallationDependencies" && prop.Values != nil && len(prop.Values) > 0 {
		signal.InstallationDependencies = stringHandler.RemoveUnnecessaryStringInArray(prop.Values)
	}
	if prop.Name == "Action.InitialInputs" && prop.Values != nil && len(prop.Values) > 0 {
		signal.Arguments = stringHandler.RemoveUnnecessaryStringInArray(prop.Values)
	}
	if prop.Name == "Action.EnvironmentVariables" && prop.DictValues != nil && len(prop.DictValues) > 0 {
		signal.EnvironmentVariables = prop.DictValues
	}
	if prop.Name == "Action.ExecutionDependencies" && prop.Values != nil && len(prop.Values) > 0 {
		signal.ExecutionDependencies = prop.Values
	}

	return signal
}

func updateSignalContextProperties(signal entities.Signal, props []entities.YamlContextProperty) entities.Signal {
	chosenContextIndex := 0

	for _, item := range props {
		if item.Name == "Context.Context" && item.Value == signal.Environment {
			chosenContextIndex = item.Position
			break
		}
	}

	for _, item := range props {
		if item.Position != chosenContextIndex {
			continue
		}
		if item.Name == "Context.Dependencies" {
			signal.ExecutionDependencies = item.Values
			continue
		}
		if item.Name == "Context.EnvironmentVariables" {
			signal.EnvironmentVariables = item.DictValues
			continue
		}
		if item.Name == "Context.ContextInitialInputs" {
			signal.Arguments = item.Values
			continue
		}
	}

	return signal
}
