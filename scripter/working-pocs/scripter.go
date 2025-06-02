package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"scripter/entities"
	"scripter/utilities"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

var security = utilities.Security{}
var fileReader = utilities.FileReader{}
var objectHandler = utilities.ObjectHandler{}
var configuration = utilities.ActionConfiguration{}
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

func main2() {
	// 1. Set up the context.
	ctx := context.Background()

	// 2. Create a new Docker client.
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	defer cli.Close()

	// 3. Pull the Docker image.  We'll use a Windows image that has Chocolatey.
	//   Make sure the base image you select has chocolatey.
	reader, err := cli.ImagePull(ctx, "mcr.microsoft.com/windows/servercore:ltsc2022", image.PullOptions{})
	if err != nil {
		log.Fatalf("Failed to pull image: %v", err)
	}
	io.Copy(os.Stdout, reader)

	// 4. Configure the container.  Now, set the command to install
	//    a list of dependencies using Chocolatey, and then start a shell.
	deps := []string{"git", "curl", "nodejs"}                                      // Example: Install Git, curl, and Node.js.
	chocoInstallCmd := fmt.Sprintf("choco install %s -y", strings.Join(deps, ",")) // Join the dependencies with commas.
	cmdStr := fmt.Sprintf("%s; powershell.exe", chocoInstallCmd)                   // Install and then start powershell
	containerConfig := &container.Config{
		Image: "mcr.microsoft.com/windows/servercore:ltsc2022", // Use the Windows image.
		Cmd:   []string{"powershell", "-Command", cmdStr},      // Run the choco install command in PowerShell.
		Tty:   true,                                            // IMPORTANT: Set Tty to true to keep the container running and interactive
	}

	// 5. Configure the host.
	hostConfig := &container.HostConfig{}

	// 6. Create the container.
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "my-choco-container")
	if err != nil {
		log.Fatalf("Failed to create container: %v", err)
	}

	// 7. Start the container.
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Fatalf("Failed to start container: %v", err)
	}

	// 8. Print the container ID.
	fmt.Printf("Container ID: %s\n", resp.ID)

	// 9.  No longer wait, just print a message that it is running
	fmt.Println("Container is running.  You can connect to it using 'docker exec -it my-choco-container powershell'")
}

// Executor Lobby
func interpretSignal(signal entities.Signal) {

	if !security.ValidateSecurity(signal) {
		return
	}

	configuration.SetGeneralConfiguration(signal)

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
