package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

type YamlFile struct {
	Header struct {
		Import   []string `yaml:"import"`
		Inherits []string `yaml:"inherits"`
		Name     string   `yaml:"name"`
	} `yaml:"header"`
	Configuration struct {
		AgentOrLabel   string `yaml:"agent-or-label"`
		ExecutionMode  string `yaml:"execution-mode"`
		BypassSecurity bool   `yaml:"bypass-security"`
		Security       struct {
			User                    string `yaml:"user"`
			PublicPassword          string `yaml:"public-password"`
			PrivatePasswordLocation string `yaml:"private-password-location"`
			CertificateLocation     string `yaml:"certificate-location"`
			TemplateOrSource        string `yaml:"template-or-source"`
		} `yaml:"security"`
		ContextName string `yaml:"context-name"`
	} `yaml:"configuration"`
	Action struct {
		NameOrFullPath string   `yaml:"name-or-full-path"`
		Type           string   `yaml:"type"`
		Api            string   `yaml:"api"`
		OutputMode     string   `yaml:"output-mode"`
		ShutdownSignal string   `yaml:"shutdown-signal"`
		InitialInputs  []string `yaml:"shutdown-signal"` //If a specific context is defined and used, inputs should
		//be defined there instead
		Platform struct {
			OsFamily              string   `yaml:"os-family"`
			PackageInstaller      string   `yaml:"package-installer"`
			ExecutionDependencies []string `yaml:"execution-dependencies"`
		}
	} `yaml:"action"`
	Contexts []struct {
		Context      string `yaml:"context"`
		Dependencies struct {
			Location string   `yaml:"location"`
			List     []string `yaml:"list"`
		} `yaml:"dependencies"`
		ContextInitialInputs []string `yaml:"context-initial-inputs"`
		EnvironmentVariables []string `yaml:"environment-variables"`
	} `yaml:"contexts"`
	Steps []struct {
		Name    string `yaml:"name"`
		Pointer string `yaml:"pointer"`
	} `yaml:"steps"`
}

func main() {
	var filePath string = os.Args[1]
	//"C:/git/calegro-project/yaml-library/window-program.yaml"
	readYaml(filePath)
	//interpretYaml(result)
}

func readYaml(filePath string) {
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

}
func interpretYaml(yaml YamlFile) {
	//Parent
	//Security
	//Chosen Context
	//Steps if any, otherwise the solely Action
}

func setLabels(labels []string) {

}

func execute(command string, args []string) {
	// Example with a config file:
	//cmd := exec.Command("dosbox", "-conf", "my_dosbox.conf")

	// Example with commands (using -c):

	// Capture output (optional)
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(command + " started with config/commands!")
}
