package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"scripter/entities"
	"scripter/utilities"
)

var security = utilities.Security{}
var fileReader = utilities.FileReader{}
var objectHandler = utilities.ObjectHandler{}
var configuration = utilities.Configuration{}
var actionHandler = utilities.ActiondHandler{}
var queueHandler = utilities.QueueHandler{}

func main() {

	filePath := os.Args[1]
	originatorPath := os.Args[2]
	nickname := os.Args[3]
	requireAcknowledge := os.Args[4]

	yamls := fileReader.ReadAllYamls(filePath)

	generalProperties, contextProperties, signalSteps := objectHandler.GenerateYamlProperties(yamls)

	signal := objectHandler.GenerateSignal(generalProperties, contextProperties, signalSteps, originatorPath, nickname, requireAcknowledge)

	configuration = configuration.SetConfigurationFromSignal(signal)

	fmt.Printf("%+v\n", signal)

	interpretSignal(signal)
}

// Executor Lobby
func interpretSignal(signal entities.Signal) {

	if !security.ValidateSecurity(signal) {
		return
	}

	configuration.SetGeneralConfiguration(signal)

	actionHandler.SetSignalAction(signal)

	//queueHandler.QueueQuaySignals(signal)

	//execute("algo", signal.Arguments)
}

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
