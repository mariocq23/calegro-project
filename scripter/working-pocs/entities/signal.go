package entities

type Signal struct {
	Labels                   []string
	Containerize             bool
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
	ShutdownSignal           string
	Arguments                []string
	HostOs                   string
	SignalOs                 string
	ExecutorOs               string
	PackageInstaller         string
	InstallationDependencies []string
	ExecutionDependencies    []string
	Environment              string
	EnvironmentVariables     map[string]string
	OriginatorQuay           OriginatorQuay
	EmitQuays                []EmitQuay
}
