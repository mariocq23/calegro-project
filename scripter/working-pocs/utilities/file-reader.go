package utilities

import (
	"log"
	"os"
	"scripter/entities"
	"scripter/entities/versions"

	"strings"

	"gopkg.in/yaml.v3"
)

type FileReader struct {
}

var stringHandler = StringHandler{}

//File reader

func (fileReader FileReader) ReadAllYamls(path string) []*versions.YamlFile_Generic_01 {

	yamlsArray := make([]*versions.YamlFile_Generic_01, 0)

	yaml := fileReader.ReadYaml(path)

	if strings.TrimSpace(yaml.Header.Inherits) != "" {
		parentPath, parentName := stringHandler.ExtractBeforeAndAfterValues(yaml.Header.Inherits)
		importInherit := entities.ImportInherit{ParentPath: parentPath, ParentName: parentName}
		yaml.Parent = fileReader.ReadYaml(importInherit.ParentPath)
		newYamlArray := fileReader.ReadAllYamls(importInherit.ParentPath)
		yamlsArray = append(yamlsArray, newYamlArray...)
	}

	yamlsArray = append(yamlsArray, yaml)

	return yamlsArray
}

//File reader

func (fileReader FileReader) ReadYaml(filePath string) *versions.YamlFile_Generic_01 {

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var yamlFile versions.YamlFile_Generic_01
	err = yaml.Unmarshal(data, &yamlFile)
	if err != nil {
		log.Fatal(err)
	}

	return &yamlFile
}

func (fileReader FileReader) ObtainAllSourceAncestors(templateName string, yamls []*versions.YamlFile_Generic_01) []string {
	yamlsArray := make([]string, 0)

	template := findYamlByName(yamls, templateName)

	if template != nil && strings.TrimSpace(template.Header.Inherits) != "" {
		parentPath, parentName := stringHandler.ExtractBeforeAndAfterValues(template.Header.Inherits)
		importInherit := entities.ImportInherit{ParentPath: parentPath, ParentName: parentName}
		newYamlArray := fileReader.ObtainAllSourceAncestors(importInherit.ParentName, yamls)
		yamlsArray = append(yamlsArray, newYamlArray...)
	}

	yamlsArray = append(yamlsArray, template.Header.Name)

	return yamlsArray

}

func findYamlByName(yamls []*versions.YamlFile_Generic_01, templateName string) *versions.YamlFile_Generic_01 {
	for _, yaml := range yamls {
		if yaml.Header.Name == templateName {
			return yaml
		}
	}
	return nil
}
