package versions

type YamlFile_Generic_01 struct {
	Header struct {
		Inherits string   `yaml:"inherits"`
		Name     string   `yaml:"name"`
		Labels   []string `yaml:"labels"`
	} `yaml:"header"`
	Configuration struct {
		Containerize          *bool  `yaml:"containerize"`
		ContainerOrchestrator string `yaml:"container-orchestrator"`
		IdempotentEngine      string `yaml:"idempotent-engine"`
		AgentOrLabel          string `yaml:"agent-or-label"`
		ExecutionMode         string `yaml:"execution-mode"`
		BypassSecurity        *bool  `yaml:"bypass-security"`
		Security              struct {
			AuthenticationHub string `yaml:"authentication-hub"`
			AuthorizationHub  string `yaml:"authorization-hub"`
			CertificationHub  string `yaml:"certification-hub"`
		} `yaml:"security"`
		ExecuteLocally *bool `yaml:"execute-locally"`
		Executor       struct {
			Os   string `yaml:"os"`
			User string `yaml:"user"`
			Host string `yaml:"host"`
		} `yaml:"executor"`
		ContextName  string `yaml:"context-name"`
		CanOverwrite *bool  `yaml:"can-overwrite"`
	} `yaml:"configuration"`
	Action struct {
		NameOrFullPath string   `yaml:"name-or-full-path"`
		Type           string   `yaml:"type"`
		Api            string   `yaml:"api"`
		ShutdownSignal string   `yaml:"shutdown-signal"`
		InitialInputs  []string `yaml:"initial-inputs"` //If a specific context is defined and used, inputs should
		//be defined there instead
		EnvironmentVariables []string `yaml:"environment-variables"` //If a specific context is defined and used, inputs should
		//be defined there instead
		Platform struct {
			OsFamily         string `yaml:"os-family"`
			PackageInstaller string `yaml:"package-installer"`
		}
		InstallationDependencies []string `yaml:"installation-dependencies"`
		ExecutionDependencies    []string `yaml:"execution-dependencies"`
		CanOverwrite             *bool    `yaml:"can-overwrite"`
	} `yaml:"action"`
	Environment struct {
		Contexts []struct {
			Context              string   `yaml:"context"`
			Dependencies         []string `yaml:"dependencies"`
			ContextInitialInputs []string `yaml:"context-initial-inputs"`
			EnvironmentVariables []string `yaml:"environment-variables"`
		} `yaml:"contexts"`
		CanOverwrite *bool `yaml:"can-overwrite"`
	}

	Steps struct {
		List []struct {
			Step    string `yaml:"step"`
			Pointer string `yaml:"pointer"`
		} `yaml:"list"`
		CanOverwrite *bool `yaml:"can-overwrite"`
	} `yaml:"steps"`

	Parent *YamlFile_Generic_01
}
