package utilities

import (
	"fmt"
	"os/exec"
	"scripter/entities"
)

type ActiondHandler struct{}

func (actionHandler ActiondHandler) SetSignalAction(signal entities.Signal) {
	installDependencies(signal.InstallationDependencies)
}

func installDependencies(dependencies []string) {
	command := fmt.Sprintf("brew install %s", dependencies)
	cmd := exec.Command("bash", "-c", command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error installing %s: %v\n", dependencies, err)
		fmt.Printf("Output:\n%s\n", string(output))
		return
	}

	fmt.Printf("%s installed successfully!\n", dependencies)
	fmt.Printf("Output:\n%s\n", string(output))
}
