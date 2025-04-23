package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"scripter/entities"
	"scripter/utilities"
)

var osWrapper = utilities.OsWrapper{}
var security = utilities.Security{}
var fileReader = utilities.FileReader{}
var objectHandler = utilities.ObjectHandler{}

func main() {

	filePath := os.Args[1]
	originatorPath := os.Args[2]
	nickname := os.Args[3]
	requireAcknowledge := os.Args[4]

	yamls := fileReader.ReadAllYamls(filePath)

	generalProperties, contextProperties, signalSteps := objectHandler.GenerateYamlProperties(yamls)

	signal := objectHandler.GenerateSignal(generalProperties, contextProperties, signalSteps, originatorPath, nickname, requireAcknowledge)

	signal.HostOs = runtime.GOOS

	fmt.Printf("%+v\n", signal)

	osWrapper.SetHostOs(signal.HostOs)
	osWrapper.SetSignalOs(signal.SignalOs)
	osWrapper.ToggleContainerization(signal.Containerize)
}

// Executor Lobby
func interpretSignal(signal entities.Signal) {

	if !security.ValidateSecurity(signal) {
		return
	}

	setGeneralConfiguration(signal)

	setSignalAction(signal)

	queueQuaySignals(signal)

	//execute("algo", signal.Arguments)
}

// Prepare signals to be emitted
func queueQuaySignals(signal entities.Signal) {
	panic("unimplemented")
}

// Prepare environment for execution
func setSignalAction(signal entities.Signal) {
	installDependencies(signal)
}

// Install prior required dependencies in environment for execution
func installDependencies(signal entities.Signal) {
	panic("unimplemented")
}

// Set Execution Enviornment
func setGeneralConfiguration(signal entities.Signal) {
	panic("unimplemented")
}

// Validate security of the signal (or bypass it in dev environments)

// Set labels for the signal so a runner can pick it
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
