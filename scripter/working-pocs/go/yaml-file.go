package main

type YamlFile struct {
	Header struct {
		Import   string `yaml:"import"`
		Inherits string `yaml:"inherits"`
		Name     string `yaml:"name"`
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
		InitialInputs  []string `yaml:"initial-inputs"` //If a specific context is defined and used, inputs should
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
