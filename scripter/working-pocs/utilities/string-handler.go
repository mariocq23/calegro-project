package utilities

import (
	"path/filepath"
	"strings"
)

type StringHandler struct{}

func (stringHandler StringHandler) ExtractBeforeAndAfterValues(input string) (string, string) {
	parts := strings.Split(input, "=>")
	if len(parts) == 2 {
		before := strings.TrimSpace(parts[0])
		after := strings.TrimSpace(parts[1])
		return before, after
	}
	return "", "" // Return empty strings if the split doesn't produce two parts
}

func (stringHandler StringHandler) ContainsString(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

func (stringHandler StringHandler) GetFilenameWithoutExtension(path string) string {
	filename := filepath.Base(path)
	extension := filepath.Ext(filename)

	if extension != "" {
		return strings.TrimSuffix(filename, extension)
	}
	return filename
}

func (stringHandler StringHandler) RemoveUnnecessaryString(rawString string) string {
	finalValue := strings.ReplaceAll(rawString, "$(overridable)", "")

	finalValue = strings.TrimSpace(finalValue)

	return finalValue
}

func (stringHandler StringHandler) RemoveUnnecessaryStringInArray(values []string) []string {
	var result []string
	for _, s := range values {
		if s != "$(overridable)" {
			cleanedString := strings.TrimSpace(s)
			result = append(result, cleanedString)
		}
	}
	return result
}

func (stringHandler StringHandler) FilterBy(values []string, filter string) []string {
	var filtered []string
	for _, s := range values {
		if strings.HasPrefix(s, filter) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func (stringHandler StringHandler) StringListToMap(input []string) map[string]string {
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

func (stringHandler StringHandler) InterpretStringAsBool(requireAcknowledge string) bool {
	if requireAcknowledge == "" || requireAcknowledge == "0" || requireAcknowledge == "false" || requireAcknowledge == "no" {
		return false
	}
	if requireAcknowledge == "1" || requireAcknowledge == "true" || requireAcknowledge == "yes" {
		return true
	}
	return false
}
