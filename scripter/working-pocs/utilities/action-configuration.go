package utilities

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"scripter/entities"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image" // New import for image-specific types
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type ActionConfiguration struct {
	HostOs       string
	SignalOs     string
	Containerize bool
	Executor     struct {
		Self          bool
		NameOrAddress string
		UserName      string
		Os            string
		Password      string
	}
	Type             string
	PackageInstaller string
}

var feasiblePackInstallersInMac = []string{
	"homebrew",
}

var feasiblePackInstallersInLinux = []string{
	"apt",
	"yum",
}

var feasiblePackInstallersInWindows = []string{
	"chocolatey",
}

func (configuration ActionConfiguration) SetConfigurationFromSignal(signal entities.Signal) ActionConfiguration {
	configuration.HostOs = signal.HostOs
	configuration.SignalOs = signal.SignalOs
	configuration.Containerize = signal.Containerize
	configuration.PackageInstaller = signal.PackageInstaller

	return configuration
}

func setSignalAction(signal entities.Signal) {
	installDependencies(signal)
}

func (configuration ActionConfiguration) SetGeneralConfiguration(signal entities.Signal) {
	if !checkConfigurationPlatformCombination(configuration) {
		panic("Host/signal OS combination, signal type, or package installer not feasible (Currenltly only handling one virtualization level per signal)!")
	}
	setSignalAction(signal)
}

func checkConfigurationPlatformCombination(configuration ActionConfiguration) bool {
	if configuration.HostOs == "darwin" && configuration.SignalOs == "windows" && configuration.Containerize {
		return false
	}
	if configuration.HostOs == "darwin" && !stringHandler.ContainsString(feasiblePackInstallersInMac, configuration.PackageInstaller) {
		return false
	}
	if configuration.HostOs == "linux" && !stringHandler.ContainsString(feasiblePackInstallersInLinux, configuration.PackageInstaller) {
		return false
	}
	if configuration.HostOs == "linux" && configuration.SignalOs == "windows" && configuration.Containerize {
		return false
	}
	if configuration.HostOs == "windows" && !stringHandler.ContainsString(feasiblePackInstallersInWindows, configuration.PackageInstaller) {
		return false
	}
	if configuration.Type == "ui-window-program" && configuration.Containerize {
		return false
	}
	return true
}

func installDependencies(signal entities.Signal) {
	logFileName := "package_installer.log"
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetPrefix("INSTALL: ")
	log.SetFlags(log.Ldate | log.Ltime)

	installOsDependencies(signal)

	log.Println("Dependency check and installation complete.")
	fmt.Println("Dependency check and installation complete. See brew_install.log for details.")
}

func installOsDependencies(signal entities.Signal) {
	if signal.HostOs == "darwin" {
		installDependenciesOnMac(signal.InstallationDependencies, signal.Containerize)
		return
	}
	if signal.HostOs == "windows" {
		installDependenciesOnWindows(signal.InstallationDependencies, signal.Containerize, signal.SignalOs)
		return
	}
	if signal.HostOs == "linux" {
		installDependenciesOnLinux(signal.InstallationDependencies, signal.Containerize)
		return

	}
}

func installDependenciesOnLinux(dependencies []string, containerize bool) {
	panic("unimplemented")
}

func installDependenciesOnWindows(dependencies []string, containerize bool, signalOs string) {
	if containerize && signalOs == "linux" {
		installDependenciesOnLinuxContainer(dependencies)
		return
	}
	if containerize && signalOs == "windows" {
		installDependenciesOnWindowsContainer(dependencies)
		return
	}

	for _, dep := range dependencies {
		if !isInstalledOnWindows(dep) {
			log.Printf("%s is not installed. Proceeding with installation...\n", dep)
			err :=
				installDependencyOnWindows(dep, true, false)
			if err != nil {
				log.Printf("Failed to install %s.\n", dep)
			}
		}
	}
}

func installDependenciesOnWindowsContainer(dependencies []string) {
	panic("unimplemented")
}

func installDependenciesOnMac(dependencies []string, containerize bool) {
	if containerize {
		installDependenciesOnLinuxContainer(dependencies)
		return
	}

	for _, dep := range dependencies {
		if !isInstalledOnMac(dep) {
			log.Printf("%s is not installed. Proceeding with installation...\n", dep)
			err := installDependencyOnMac(dep)
			if err != nil {
				log.Printf("Failed to install %s.\n", dep)
			}
		}
	}
}

func installDependencyOnLinux(dep string) any {
	fmt.Printf("Installing %s with real-time verbose output...\n", dep)
	cmd := exec.Command("sudo", "apt", "install", dep, "-y")

	cmd.Stdout = os.Stdout // Connect brew's stdout to Go's stdout (your terminal)
	cmd.Stderr = os.Stderr // Connect brew's stderr to Go's stderr (your terminal)

	err := cmd.Run() // Use Run() instead of CombinedOutput()
	if err != nil {
		fmt.Printf("Error installing %s: %v\n", dep, err)
		return err
	}
	fmt.Printf("%s installation finished.\n", dep)
	return nil
}

// installDependenciesOnLinuxContainer sets up a Docker container, installs multiple dependencies,
// and leaves the container running.
func installDependenciesOnLinuxContainer(dependencies []string) { // Changed signature back to []string
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	defer cli.Close()

	imageName := "ubuntu:latest"
	containerName := "my-ubuntu-dev"

	if len(dependencies) == 0 {
		fmt.Println("No dependencies specified for container installation.")
		return
	}
	fmt.Printf("Dependencies to install in container: %s\n", strings.Join(dependencies, ", "))

	// Check if container already exists and stop/remove it to ensure a clean run
	// This is optional but good for repeated runs of the script
	_, err = cli.ContainerInspect(ctx, containerName)
	if err == nil { // Container exists
		fmt.Printf("Container '%s' already exists. Stopping and removing...\n", containerName)
		timeout := 0 // Stop immediately
		// Corrected: container.StopOptions{}
		if err := cli.ContainerStop(ctx, containerName, container.StopOptions{Timeout: &timeout}); err != nil {
			log.Printf("Warning: Failed to stop existing container '%s': %v. Attempting to remove anyway.", containerName, err)
		}
		// Corrected: container.RemoveOptions{}
		if err := cli.ContainerRemove(ctx, containerName, container.RemoveOptions{}); err != nil {
			log.Fatalf("Failed to remove existing container '%s': %v", containerName, err)
		}
		fmt.Printf("Existing container '%s' removed.\n", containerName)
	}

	// Pull the Ubuntu image if it doesn't exist
	// Corrected: image.PullOptions{}
	reader, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		log.Fatalf("Failed to pull image '%s': %v", imageName, err)
	}
	defer reader.Close()
	// Read from the reader to wait for the pull to complete and see progress
	_, _ = stdcopy.StdCopy(os.Stdout, os.Stderr, reader) // Copy pull output to stdout/stderr
	fmt.Printf("Successfully pulled image '%s'\n", imageName)

	// Container configuration
	// Corrected: container.Config{} (already correct, but clarifying)
	config := &container.Config{
		Image:        imageName,
		Tty:          true, // Keep TTY for interactive shell if needed
		OpenStdin:    true, // Keep OpenStdin for interactive shell if needed
		AttachStdout: true, // Important for `docker logs` and exec command output
		AttachStderr: true, // Important for `docker logs` and exec command output
	}

	// Host configuration (for interactive terminal)
	// Corrected: container.HostConfig{} (already correct, but clarifying)
	hostConfig := &container.HostConfig{}

	// Create the container
	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
	if err != nil {
		log.Fatalf("Failed to create container '%s': %v", containerName, err)
	}
	fmt.Printf("Successfully created container with ID: %s\n", resp.ID)

	// Start the container
	// Corrected: container.StartOptions{}
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Fatalf("Failed to start container '%s': %v", containerName, err)
	}
	fmt.Printf("Successfully started container '%s'\n", containerName)

	// Execute commands to install software
	commands := []string{
		"apt-get update",
		fmt.Sprintf("DEBIAN_FRONTEND=noninteractive apt-get install -y %s", strings.Join(dependencies, " ")),
	}

	for _, cmd := range commands {
		// Corrected: container.ExecOptions{}
		execConfig := container.ExecOptions{
			User:         "root",
			Privileged:   false, // Exec inside container typically doesn't need privileged
			Tty:          true,  // Attach TTY for clearer output (optional, but good for interactive)
			AttachStdout: true,
			AttachStderr: true,
			Cmd:          []string{"bash", "-c", cmd},
		}

		execResp, err := cli.ContainerExecCreate(ctx, resp.ID, execConfig)
		if err != nil {
			log.Fatalf("Failed to create exec command for '%s': %v", cmd, err)
		}

		// Corrected: container.ExecAttachOptions{}
		attachResp, err := cli.ContainerExecAttach(ctx, execResp.ID, container.ExecAttachOptions{})
		if err != nil {
			log.Fatalf("Failed to attach to exec command for '%s': %v", cmd, err)
		}
		_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, attachResp.Reader)
		if err != nil {
			log.Printf("Warning: Error streaming output for command '%s': %v", cmd, err)
		}
		attachResp.Close() // Close the attach stream immediately after copying to ensure resources are released

		exitCodeResp, err := cli.ContainerExecInspect(ctx, execResp.ID)
		if err != nil {
			log.Fatalf("Error inspecting exec command for '%s': %v", cmd, err)
		}
		if exitCodeResp.ExitCode != 0 {
			// Corrected the order of arguments in Fatalf to match the format string
			log.Fatalf("Command '%s' failed with exit code: %d", cmd, exitCodeResp.ExitCode)
		}
		fmt.Printf("Successfully executed command: %s\n", cmd)
	}

	fmt.Println("\nContainer is running with the specified dependencies installed.")
	fmt.Printf("You can attach to it using: docker attach %s\n", containerName)
	fmt.Printf("To stop and remove it: docker stop %s && docker rm %s\n", containerName, containerName)
}

func isInstalledOnLinux(dep string) bool {
	cmd := exec.Command(dep)
	if strings.Contains(cmd.Path, "/") {
		return true
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false // Assume not installed if there's an error checking
	}
	return strings.Contains(string(output), dep)
}

func installDependencyOnWindows(dep string, containairize bool, vmize bool) any {
	fmt.Printf("Installing %s with real-time verbose output...\n", dep)
	cmd := exec.Command("choco", "install", dep)

	cmd.Stdout = os.Stdout // Connect brew's stdout to Go's stdout (your terminal)
	cmd.Stderr = os.Stderr // Connect brew's stderr to Go's stderr (your terminal)

	err := cmd.Run() // Use Run() instead of CombinedOutput()
	if err != nil {
		fmt.Printf("Error installing %s: %v\n", dep, err)
		return err
	}
	fmt.Printf("%s installation finished.\n", dep)
	return nil
}

func isInstalledOnWindows(dep string) bool {
	cmd := exec.Command("choco", "list", "--versions", dep)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false // Assume not installed if there's an error checking
	}
	return strings.Contains(string(output), dep)
}

func isInstalledOnMac(dep string) bool {
	cmd := exec.Command("brew", "list", "--versions", dep)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false // Assume not installed if there's an error checking
	}
	return strings.Contains(string(output), dep)
}

func installDependencyOnMac(dependency string) error {

	fmt.Printf("Installing %s with real-time verbose output...\n", dependency)
	cmd := exec.Command("brew", "install", dependency)

	cmd.Stdout = os.Stdout // Connect brew's stdout to Go's stdout (your terminal)
	cmd.Stderr = os.Stderr // Connect brew's stderr to Go's stderr (your terminal)

	err := cmd.Run() // Use Run() instead of CombinedOutput()
	if err != nil {
		fmt.Printf("Error installing %s: %v\n", dependency, err)
		return err
	}
	fmt.Printf("%s installation finished.\n", dependency)
	return nil
}

// stdCopy enhanced to handle string directly.
func stdCopy(dst io.Writer, dstErr io.Writer, src io.Reader) error {
	buf := make([]byte, 32*1024)
	var err error
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			var nw int
			var ew error
			s := string(buf[0:nr])
			if strings.Contains(s, "ERROR:") {
				nw, ew = dstErr.Write(buf[0:nr])
			} else {
				nw, ew = dst.Write(buf[0:nr])
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return err
}
