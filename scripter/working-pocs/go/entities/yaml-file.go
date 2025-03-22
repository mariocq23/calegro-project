package entities

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
			AuthenticationHub string `yaml:"authentication-hub"`
			AuthorizationHub  string `yaml:"authorization-hub"`
			CertificationHub  string `yaml:"certification-hub"`
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
			OsFamily                 string   `yaml:"os-family"`
			PackageInstaller         string   `yaml:"package-installer"`
			InstallationDependencies []string `yaml:"installation-dependencies"`
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
		Step    string `yaml:"step"`
		Pointer string `yaml:"pointer"`
	} `yaml:"steps"`
}
