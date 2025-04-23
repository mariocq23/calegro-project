package utilities

import (
	"reflect"
	"scripter/entities"
	"strings"
)

type ObjectHandler struct{}

var fileReader = FileReader{}

//Object Array Generator - More context logic related

func (objectHandler ObjectHandler) GenerateYamlProperties(yamls []*entities.YamlFile) ([]entities.YamlProperty, []entities.YamlContextProperty, []entities.SignalStep) {

	yamlProperties := []entities.YamlProperty{}
	yamlContextProperties := []entities.YamlContextProperty{}
	signalSteps := []entities.SignalStep{}
	finalSignalSteps := []entities.SignalStep{}
	overridableSteps := []string{}

	for _, yaml := range yamls {
		if yaml.Configuration.Containerize != nil {
			yamlProperties = append(yamlProperties, generatePositiveBoolProperty("Configuration.Containerize", yaml.Configuration.Containerize, yaml.Header.Name))
		}
		yamlProperties = append(yamlProperties, generateProperty("Configuration.AgentOrLabel", yaml.Configuration.AgentOrLabel, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.ContextName", yaml.Configuration.ContextName, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateProperty("Configuration.ExecutionMode", yaml.Configuration.ExecutionMode, yaml.Header.Name))
		if yaml.Configuration.BypassSecurity != nil {
			yamlProperties = append(yamlProperties, generateNegativeBoolProperty("Configuration.BypassSecurity", yaml.Configuration.BypassSecurity, yaml.Header.Name))
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
		yamlProperties = append(yamlProperties, generateArrayProperty("Header.Labels", yaml.Header.Labels, yaml.Header.Name))
		yamlProperties = append(yamlProperties, generateDictionaryProperty("Action.EnvironmentVariables", stringHandler.StringListToMap(stringHandler.RemoveUnnecessaryStringInArray(yaml.Action.EnvironmentVariables)), yaml.Header.Name))

		for index, context := range yaml.Contexts {
			yamlContextProperties = append(yamlContextProperties, generateContextProperty("Context.Context", context.Context, yaml.Header.Name, index))
			yamlContextProperties = append(yamlContextProperties, generateContextArrayProperty("Context.Dependencies", context.Dependencies, yaml.Header.Name, index))
			yamlContextProperties = append(yamlContextProperties, generateContextArrayProperty("Context.ContextInitialInputs", context.ContextInitialInputs, yaml.Header.Name, index))
			yamlContextProperties = append(yamlContextProperties, generateContextDictionaryProperty("Context.EnvironmentVariables", stringHandler.StringListToMap(stringHandler.RemoveUnnecessaryStringInArray(context.EnvironmentVariables)), yaml.Header.Name, index))
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

	return yamlProperties, yamlContextProperties, finalSignalSteps
}

func generateContextArrayProperty(name string, values []string, templateName string, index int) entities.YamlContextProperty {
	yamlProperty := entities.YamlContextProperty{Name: name, Values: values, TemplateName: templateName, Position: index}
	if !stringHandler.ContainsString(values, "$(overridable)") && values != nil && len(values) > 0 {
		yamlProperty.Sealed = true
	}
	if stringHandler.ContainsString(values, "default") {
		yamlProperty.Default = true
	}
	return yamlProperty
}

func generateSignalStep(step string, pointer string, index int, templateName string) entities.SignalStep {
	signalStep := entities.SignalStep{
		Name:    step,
		Pointer: pointer,
		Source:  templateName,
	}

	return signalStep
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

//Objectj generator - more context logic related

func generateNegativeBoolProperty(name string, value *bool, templateName string) entities.YamlProperty {
	yamlProperty := entities.YamlProperty{Name: name, BoolValue: value, TemplateName: templateName}
	if value != nil && !*value && name == "Configuration.BypassSecurity" {
		yamlProperty.Sealed = true
		yamlProperty.Default = true
	}
	return yamlProperty
}

func generatePositiveBoolProperty(name string, value *bool, templateName string) entities.YamlProperty {
	yamlProperty := entities.YamlProperty{Name: name, BoolValue: value, TemplateName: templateName}
	if value != nil && *value && name == "Configuration.Containerize" {
		yamlProperty.Sealed = true
		yamlProperty.Default = true
	}
	return yamlProperty
}

//Objectj generator - more context logic related

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

//Objectj generator - more context logic related

func generateArrayProperty(name string, values []string, templateName string) entities.YamlProperty {
	yamlProperty := entities.YamlProperty{Name: name, Values: values, TemplateName: templateName}
	if !stringHandler.ContainsString(values, "$(overridable)") && values != nil && len(values) > 0 {
		yamlProperty.Sealed = true
	}
	if !stringHandler.ContainsString(values, "default") {
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

//Objectj generator - more context logic related

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

func (objectHandler ObjectHandler) GenerateSignal(generalProperties []entities.YamlProperty, contextProperties []entities.YamlContextProperty, steps []entities.SignalStep, originatorPath string, nickname string, requireAcknowledge string) entities.Signal {
	sealedProperties := []string{}
	signal := entities.Signal{}
	signal.Sender = generalProperties[len(generalProperties)-1].TemplateName
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
		emitQuay := entities.EmitQuay{Name: stringHandler.GetFilenameWithoutExtension(executionDependency), Path: executionDependency, Relationship: entities.Dependency, Priority: index}
		emitQuays = append(emitQuays, emitQuay)
	}

	for index, step := range steps {
		emitQuay := entities.EmitQuay{Name: stringHandler.GetFilenameWithoutExtension(step.Pointer), Path: step.Pointer, Relationship: entities.Step, Priority: index}
		emitQuays = append(emitQuays, emitQuay)
	}

	return emitQuays
}

func updateSignalGeneralProperties(signal entities.Signal, prop entities.YamlProperty) entities.Signal {

	if prop.Name == "Header.Labels" {
		signal.Labels = prop.Values
	}
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
	if prop.Name == "Action.OutputMode" && prop.Value != "" {
		signal.OutputMode = stringHandler.RemoveUnnecessaryString(prop.Value)
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
