package entities

type Signal struct {
	Sender                   string
	Executor                 string
	ExecutionMode            string
	Type                     string
	BypassSecurity           bool
	User                     string
	Certificate              string
	Password                 string
	Token                    string
	AuthenticationHub        string
	AuthorizationHub         string
	CertificationHub         string
	Api                      string
	ExecutablePath           string
	OutputMode               string
	ShutdownSignal           string
	Arguments                []string
	Os                       string
	PackageInstaller         string
	InstallationDependencies []string
	ExecutionDependencies    []string
	Environment              string
	EnvironmentVariables     []string
	Steps                    []SignalStep
}
