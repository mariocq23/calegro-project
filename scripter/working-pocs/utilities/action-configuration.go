package utilities

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"scripter/entities"
	"strings"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image" // New import for image-specific types
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/moby/go-archive"
	mobyterm "github.com/moby/term"
	xterm "golang.org/x/term"
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
	if configuration.HostOs == "darwin" && !stringHandler.ContainsString(feasiblePackInstallersInMac, configuration.PackageInstaller) && !configuration.Containerize {
		return false
	}
	if configuration.HostOs == "linux" && !stringHandler.ContainsString(feasiblePackInstallersInLinux, configuration.PackageInstaller) {
		return false
	}
	if configuration.HostOs == "linux" && configuration.SignalOs == "windows" && configuration.Containerize {
		return false
	}
	if configuration.HostOs == "windows" && !stringHandler.ContainsString(feasiblePackInstallersInWindows, configuration.PackageInstaller) && !configuration.Containerize {
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
	if containerize {
		installDependenciesOnLinuxContainer()
		return
	}

	for _, dep := range dependencies {
		if !isInstalledOnLinux(dep) {
			log.Printf("%s is not installed. Proceeding with installation...\n", dep)
			err := installDependencyOnLinux(dep)
			if err != nil {
				log.Printf("Failed to install %s.\n", dep)
			}
		}
	}
}

func installDependenciesOnWindows(dependencies []string, containerize bool, signalOs string) {
	if containerize && signalOs == "linux" {
		installDependenciesOnLinuxContainer()
		return
	}
	if containerize && signalOs == "windows" {
		installDependenciesOnWindowsContainer(dependencies, true)
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

func installDependenciesOnWindowsContainer(dependencies []string, preinstall bool) {
	if !preinstall {
		runWindowsContainer(dependencies)
		return
	}
	// Generate the Dockerfile content
	dockerfileContent, err := generateDockerfile("windows", dependencies)
	if err != nil {
		log.Fatalf("Failed to generate Dockerfile: %v", err)
	}
	fmt.Println("\nGenerated Dockerfile:\n", dockerfileContent)

	// Define a consistent container name
	// containerName := fmt.Sprintf("my-dev-container-%s", "windows")

	// // Build the image and run the container
	// if err := buildAndRunDockerfile(dockerfileContent, containerName); err != nil {
	// 	log.Fatalf("Failed to build and run container: %v", err)
	// }

	fmt.Println("Go program finished.")
}

func runWindowsContainer(dependencies []string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	defer cli.Close()

	// Use a common Windows Server Core image. Nano Server is smaller but might lack some tools.
	// You might need to pick a specific tag like :ltsc2022 or :20H2 depending on your host OS.
	imageName := "mcr.microsoft.com/windows/servercore:ltsc2022"
	containerName := "my-windows-dev-interactive"

	fmt.Printf("Attempting to provision Windows container '%s' with dependencies: %s\n", containerName, strings.Join(dependencies, ", "))
	fmt.Printf("!!! Ensure you are running on a Windows host with Docker Desktop in Windows Container mode. !!!\n")

	// --- Step 1: Pre-check and Clean Up Existing Container ---
	_, err = cli.ContainerInspect(ctx, containerName)
	if err == nil {
		fmt.Printf("Container '%s' already exists. Stopping and removing...\n", containerName)
		timeout := 5
		if err := cli.ContainerStop(ctx, containerName, container.StopOptions{Timeout: &timeout}); err != nil {
			log.Printf("Warning: Failed to stop existing container '%s': %v. Attempting to remove anyway.", containerName, err)
		}
		if err := cli.ContainerRemove(ctx, containerName, container.RemoveOptions{}); err != nil {
			log.Fatalf("Failed to remove existing container '%s': %v", containerName, err)
		}
		fmt.Printf("Existing container '%s' removed.\n", containerName)
	}

	// --- Step 2: Pull Image ---
	fmt.Printf("Attempting to pull image '%s'...\n", imageName)
	reader, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		log.Fatalf("Failed to pull image '%s': %v", imageName, err)
	}
	defer reader.Close()
	// Windows container pull output might not use multiplexing, so stdcopy might not be strictly necessary
	// but it's generally safe to use.
	if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, reader); err != nil {
		log.Printf("Warning: Error streaming image pull output: %v", err)
	}
	fmt.Printf("Successfully pulled image '%s'\n", imageName)

	// --- Step 3: Container Configuration for Interactive Session ---
	config := &container.Config{
		Image:        imageName,
		Tty:          true, // Allocate a TTY (behavior may vary compared to Linux)
		OpenStdin:    true, // Keep STDIN open
		AttachStdout: true, // Attach stdout
		AttachStderr: true, // Attach stderr
		// For Windows, default shell is cmd.exe. We start powershell.exe for more flexibility.
		Cmd: []string{"powershell.exe"},
	}

	hostConfig := &container.HostConfig{
		// No specific host config needed for simple interactive use
	}

	// --- Step 4: Create and Start Container ---
	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
	if err != nil {
		log.Fatalf("Failed to create container '%s': %v", containerName, err)
	}
	fmt.Printf("Successfully created container with ID: %s\n", resp.ID)

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Fatalf("Failed to start container '%s': %v", containerName, err)
	}
	fmt.Printf("Successfully started container '%s'\n", containerName)

	// --- Step 5: Install Dependencies via Exec ---
	// For Windows, installing dependencies often involves different commands.
	// We'll use a simple PowerShell command. `winget` is often not in base images.
	// For example, to install a feature or a specific package.
	// For this example, we'll simulate a simple "installation" command.
	// Real-world Windows package management (like chocolatey or direct MSIs) is more complex.
	fmt.Println("\n--- Installing Dependencies ---")
	// Example: Create a directory and a file.
	// Real dependencies might involve `dism` for Windows features or specific installers.
	commands := []string{
		"mkdir C:\\Dependencies",
		fmt.Sprintf("echo Dependencies_installed_GoLang > C:\\Dependencies\\%s.txt", strings.Join(dependencies, "_")),
		// Example of installing a Windows feature using Dism.exe
		// "dism.exe /online /enable-feature /featurename:IIS-WebServerRole /all /NoRestart",
	}

	for _, cmd := range commands {
		fmt.Printf("Executing command in container: %s\n", cmd)
		execConfig := container.ExecOptions{
			User:         "ContainerAdministrator", // Typical user for Windows containers
			Privileged:   false,
			Tty:          true, // Use TTY for clearer output
			AttachStdout: true,
			AttachStderr: true,
			Cmd:          []string{"powershell.exe", "-Command", cmd}, // Use PowerShell to execute
		}

		execResp, err := cli.ContainerExecCreate(ctx, resp.ID, execConfig)
		if err != nil {
			log.Fatalf("Failed to create exec command for '%s': %v", cmd, err)
		}

		attachResp, err := cli.ContainerExecAttach(ctx, execResp.ID, container.ExecAttachOptions{})
		if err != nil {
			log.Fatalf("Failed to attach to exec command for '%s': %v", cmd, err)
		}
		// Windows container output might not be multiplexed, so io.Copy might be more direct.
		// However, stdcopy.StdCopy is resilient.
		if _, err := io.Copy(os.Stdout, attachResp.Reader); err != nil { // Directly copy to stdout
			log.Printf("Warning: Error streaming output for command '%s': %v", cmd, err)
		}
		attachResp.Close()

		exitCodeResp, err := cli.ContainerExecInspect(ctx, execResp.ID)
		if err != nil {
			log.Fatalf("Error inspecting exec command for '%s': %v", cmd, err)
		}
		if exitCodeResp.ExitCode != 0 {
			log.Fatalf("Command '%s' failed with exit code: %d", cmd, exitCodeResp.ExitCode)
		}
		fmt.Printf("Successfully executed command: %s\n", cmd)
	}

	fmt.Println("\n--- Dependencies Provisioned. Attaching to Windows Container. ---")
	fmt.Println("Type 'exit' to detach from the container (container will remain running).")

	// --- Step 6: Attach to the Container's Primary Process (Interactive Shell) ---
	oldState, err := xterm.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Failed to set terminal to raw mode: %v", err)
	}
	defer func() {
		fmt.Println("\n--- Detaching from container. Restoring terminal state. ---")
		xterm.Restore(int(os.Stdin.Fd()), oldState)
	}()

	width, height, err := xterm.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Printf("Warning: Failed to get terminal size: %v. Using default.", err)
		width, height = 80, 24
	}

	// Windows containers might have specific attach options or behaviors for TTY/Stdin/Stdout.
	// The `AttachOptions` are generally the same.
	attachOptions := container.AttachOptions{
		Stream:     true,
		Stdin:      true,
		Stdout:     true,
		Stderr:     true,
		Logs:       false,
		DetachKeys: "ctrl-p,ctrl-q", // Standard Docker detach keys
	}

	hijackedResp, err := cli.ContainerAttach(ctx, resp.ID, attachOptions)
	if err != nil {
		log.Fatalf("Failed to attach to container: %v", err)
	}
	defer hijackedResp.Close()

	// Goroutine to copy container stdout/stderr to os.Stdout/Stderr
	// For Windows, often direct io.Copy is used as streams might not be multiplexed by default.
	go func() {
		_, err := io.Copy(os.Stdout, hijackedResp.Reader)
		if err != nil {
			// This error can often be an EOF when the container exits, which is normal.
			// log.Printf("Error copying stdout from container: %v", err)
		}
	}()

	// Goroutine to copy os.Stdin to container stdin
	go func() {
		_, err := io.Copy(hijackedResp.Conn, os.Stdin)
		if err != nil {
			// This error can often be an EOF when your program terminates or connection closes.
			// log.Printf("Error copying stdin to container: %v", err)
		}
	}()

	// Initial resize call
	cli.ContainerResize(ctx, resp.ID, container.ResizeOptions{Width: uint(width), Height: uint(height)})

	// Block indefinitely to keep the interactive session open.
	select {}
}

func installDependenciesOnMac(dependencies []string, containerize bool) {
	if containerize {
		installDependenciesOnLinuxContainer()
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

func buildAndRunContainer(containerName, dockerfilePath string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()

	imageTag := strings.ToLower(containerName + ":latest")

	// --- Step 1: Read the Dockerfile content from the file system ---
	fmt.Printf("Reading Dockerfile from: %s\n", dockerfilePath)
	dockerfileBytes, err := os.ReadFile(dockerfilePath)
	if err != nil {
		return fmt.Errorf("failed to read Dockerfile: %w", err)
	}

	// --- Step 2: Create an in-memory tar archive with the Dockerfile ---
	fmt.Println("Creating in-memory Docker build context...")
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	header := &tar.Header{
		Name: "Dockerfile",
		Size: int64(len(dockerfileBytes)),
		Mode: 0600,
	}

	if err := tw.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write tar header: %w", err)
	}
	if _, err := tw.Write(dockerfileBytes); err != nil {
		return fmt.Errorf("failed to write Dockerfile content to tar: %w", err)
	}
	if err := tw.Close(); err != nil {
		return fmt.Errorf("failed to close tar writer: %w", err)
	}

	tarReader := bytes.NewReader(buf.Bytes())

	// --- Step 3: Build the Docker Image from the in-memory tarball ---
	fmt.Printf("Building Docker image '%s' from in-memory context...\n", imageTag)
	buildOptions := build.ImageBuildOptions{
		Tags:       []string{imageTag},
		Dockerfile: "Dockerfile",
		Remove:     true, // Clean up intermediate containers
	}

	buildResponse, err := cli.ImageBuild(ctx, tarReader, buildOptions)
	if err != nil {
		return fmt.Errorf("failed to start image build: %w", err)
	}
	defer buildResponse.Body.Close()

	// Stream build output to standard output
	if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, buildResponse.Body); err != nil {
		log.Printf("Warning: Error streaming build output: %v", err)
	}
	fmt.Printf("Successfully built image '%s'.\n", imageTag)

	// --- Step 4: Clean Up Old Container ---
	_, err = cli.ContainerInspect(ctx, containerName)
	if err == nil {
		fmt.Printf("Container '%s' already exists. Stopping and removing...\n", containerName)
		timeout := 5
		if err := cli.ContainerStop(ctx, containerName, container.StopOptions{Timeout: &timeout}); err != nil {
			log.Printf("Warning: Failed to stop existing container '%s': %v. Attempting to remove anyway.", containerName, err)
		}
		if err := cli.ContainerRemove(ctx, containerName, container.RemoveOptions{}); err != nil {
			return fmt.Errorf("failed to remove existing container '%s': %w", containerName, err)
		}
		fmt.Printf("Existing container '%s' removed.\n", containerName)
	}

	// --- Step 5: Run the Custom Image ---
	fmt.Printf("Creating and starting container '%s' from image '%s'...\n", containerName, imageTag)
	config := &container.Config{
		Image:        imageTag,
		Tty:          true,
		OpenStdin:    true,
		AttachStdout: true,
		AttachStderr: true,
	}

	resp, err := cli.ContainerCreate(ctx, config, &container.HostConfig{}, nil, nil, containerName)
	if err != nil {
		return fmt.Errorf("failed to create container '%s': %w", containerName, err)
	}
	fmt.Printf("Successfully created container with ID: %s\n", resp.ID)

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container '%s': %w", containerName, err)

	}
	fmt.Printf("Successfully started container '%s'\n", containerName)

	// --- Step 6: Attach to the Container ---
	fmt.Println("\n--- Attaching to Container. Services are now running. ---")
	fmt.Println("Type 'exit' to detach from the container (container will remain running).")

	// Set the terminal to raw mode to handle user input correctly
	oldState, err := xterm.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to set terminal to raw mode: %w", err)
	}
	defer func() {
		fmt.Println("\n--- Detaching from container. Restoring terminal state. ---")
		xterm.Restore(int(os.Stdin.Fd()), oldState)
	}()

	// Get the terminal size for resizing
	width, height, err := xterm.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Printf("Warning: Failed to get terminal size: %v. Using default.", err)
		width, height = 80, 24
	}

	attachOptions := container.AttachOptions{
		Stream:     true,
		Stdin:      true,
		Stdout:     true,
		Stderr:     true,
		Logs:       false,
		DetachKeys: "ctrl-p,ctrl-q",
	}

	hijackedResp, err := cli.ContainerAttach(ctx, resp.ID, attachOptions)
	if err != nil {
		return fmt.Errorf("failed to attach to container: %w", err)
	}
	defer hijackedResp.Close()

	// Use goroutines to handle stdin/stdout streams
	go func() {
		_, err := io.Copy(os.Stdout, hijackedResp.Reader)
		if err != nil && err != io.EOF {
			log.Printf("Error copying stdout from container: %v", err)
		}
	}()

	go func() {
		_, err := io.Copy(hijackedResp.Conn, os.Stdin)
		if err != nil && err != io.EOF {
			log.Printf("Error copying stdin to container: %v", err)
		}
	}()

	// Resize the container's TTY to match the terminal
	cli.ContainerResize(ctx, resp.ID, container.ResizeOptions{Width: uint(width), Height: uint(uint16(height))})

	// Block the main goroutine to keep the program running
	select {}
}

// installDependenciesOnLinuxContainer sets up a Docker container, installs multiple dependencies,
// and leaves the container running.
func installDependenciesOnLinuxContainer() {

	// Build the image and run the container

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	// Print the directory
	fmt.Println("Current directory is:", dir)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}
	dockerfilePath := filepath.Join(wd, "Dockerfile")

	err = buildImage("suck-it", dockerfilePath)
	if err != nil {
		log.Fatalf("Provisioning failed: %v", err)
	}

	fmt.Println("Go program finished.")
}

func buildImage(imageTag, contextPath string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()

	// --- Step 1: Create a tar archive of the entire build context directory ---
	fmt.Printf("Creating tar archive of build context: %s\n", contextPath)
	tarContext, err := archive.TarWithOptions(contextPath, &archive.TarOptions{})
	if err != nil {
		return fmt.Errorf("failed to create build context tar: %w", err)
	}
	defer tarContext.Close()

	// --- Step 2: Build the Docker Image from the tar archive ---
	fmt.Printf("Building Docker image '%s'...\n", imageTag)
	buildOptions := build.ImageBuildOptions{
		Tags:       []string{imageTag},
		Dockerfile: "Dockerfile", // This should be the name of the Dockerfile within the context
		Remove:     true,         // Clean up intermediate containers
	}

	buildResponse, err := cli.ImageBuild(ctx, tarContext, buildOptions)
	if err != nil {
		return fmt.Errorf("failed to start image build: %w", err)
	}
	defer buildResponse.Body.Close()

	// --- Step 3: Stream build output ---
	// This part of your original code is correct for streaming
	termFd, isTerm := mobyterm.GetFdInfo(os.Stderr)
	if err := jsonmessage.DisplayJSONMessagesStream(buildResponse.Body, os.Stderr, termFd, isTerm, nil); err != nil {
		// Use stdcopy if the output is not JSON-formatted (e.g., plain text logs)
		if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, buildResponse.Body); err != nil {
			log.Printf("Warning: Error streaming build output: %v", err)
		}
	}

	return nil
}

func executeCommand(ctx context.Context, cli *client.Client, containerID string, cmd string) {
	fmt.Printf("Executing command in container: %s\n", cmd)
	execConfig := container.ExecOptions{
		User:         "root",
		Privileged:   false,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"bash", "-c", cmd},
	}

	execResp, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		log.Fatalf("Failed to create exec command for '%s': %v", cmd, err)
	}

	attachResp, err := cli.ContainerExecAttach(ctx, execResp.ID, container.ExecAttachOptions{})
	if err != nil {
		log.Fatalf("Failed to attach to exec command for '%s': %v", cmd, err)
	}
	defer attachResp.Close()

	if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, attachResp.Reader); err != nil {
		log.Printf("Warning: Error streaming output for command '%s': %v", cmd, err)
	}

	exitCodeResp, err := cli.ContainerExecInspect(ctx, execResp.ID)
	if err != nil {
		log.Fatalf("Error inspecting exec command for '%s': %v", cmd, err)
	}
	if exitCodeResp.ExitCode != 0 {
		log.Fatalf("Command '%s' failed with exit code: %d", cmd, exitCodeResp.ExitCode)
	}
	fmt.Printf("Successfully executed command: %s\n", cmd)
}

func generateDockerfile(osType string, dependencies []string) (string, error) {
	var dockerfileContent strings.Builder

	switch strings.ToLower(osType) {
	case "linux":
		baseImage := "ubuntu:latest"
		dockerfileContent.WriteString(fmt.Sprintf("FROM %s\n", baseImage))
		dockerfileContent.WriteString("ENV DEBIAN_FRONTEND=noninteractive\n")
		dockerfileContent.WriteString("RUN apt-get update && apt-get install -y \\\n")
		for i, dep := range dependencies {
			dockerfileContent.WriteString(fmt.Sprintf("    %s%s", dep, func() string {
				if i < len(dependencies)-1 {
					return " \\"
				}
				return ""
			}()))
			if (i+1)%5 == 0 && i < len(dependencies)-1 { // Newline every 5 deps for readability
				dockerfileContent.WriteString("\n")
			}

		}
		dockerfileContent.WriteString("\n")

		dockerfileContent.WriteString("CMD [\"bash\"]\n") // Keep container open with bash
	case "windows":
		// Choose a Windows Server Core image. Ensure it's compatible with your Windows Docker host.
		// Use a specific tag like :ltsc2022 to avoid issues with host OS mismatch.
		baseImage := "mcr.microsoft.com/windows/servercore:ltsc2022"
		dockerfileContent.WriteString(fmt.Sprintf("FROM %s\n", baseImage))
		dockerfileContent.WriteString("SHELL [\"powershell.exe\", \"-Command\", \"$ErrorActionPreference = 'Stop'; $ProgressPreference = 'SilentlyContinue';\"]\n")

		// Windows package installation is more complex.
		// For simplicity, we'll demonstrate creating files/dirs.
		// Real installations would involve MSIEXEC, DISM, winget, or Chocolatey.
		dockerfileContent.WriteString("RUN mkdir C:\\Dependencies\n")
		dockerfileContent.WriteString(fmt.Sprintf("RUN echo Dependencies_installed_GoLang > C:\\Dependencies\\%s.txt\n", strings.Join(dependencies, "_")))

		// Example of installing a Windows feature (uncomment if needed and ensure compatibility)
		// dockerfileContent.WriteString("RUN dism.exe /online /enable-feature /featurename:IIS-WebServerRole /all /NoRestart\n")

		dockerfileContent.WriteString("CMD [\"powershell.exe\"]\n") // Keep container open with PowerShell
	default:
		return "", fmt.Errorf("unsupported OS type: %s. Choose 'linux' or 'windows'.", osType)
	}

	return dockerfileContent.String(), nil
}

// buildAndRunContainer builds a custom Docker image and runs a container from it.
func readAndBuildContainer(containerName, dockerfilePath string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()

	imageTag := strings.ToLower(containerName + ":latest")

	// --- Step 1: Read the Dockerfile content from the file system ---
	dockerfileBytes, err := os.ReadFile(dockerfilePath)
	if err != nil {
		return fmt.Errorf("failed to read Dockerfile at %s: %w", dockerfilePath, err)
	}

	// --- Step 2: Create an in-memory tar archive with the Dockerfile ---
	fmt.Println("Creating in-memory Docker build context...")
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	header := &tar.Header{
		Name: "Dockerfile",
		Size: int64(len(dockerfileBytes)),
		Mode: 0600,
	}

	if err := tw.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write tar header: %w", err)
	}
	if _, err := tw.Write(dockerfileBytes); err != nil {
		return fmt.Errorf("failed to write Dockerfile content to tar: %w", err)
	}
	if err := tw.Close(); err != nil {
		return fmt.Errorf("failed to close tar writer: %w", err)
	}

	tarReader := bytes.NewReader(buf.Bytes())

	// --- Step 3: Build the Docker Image from the in-memory tarball ---
	fmt.Printf("Building Docker image '%s' from in-memory context...\n", imageTag)
	buildOptions := build.ImageBuildOptions{
		Tags:       []string{imageTag},
		Dockerfile: "Dockerfile",
		Remove:     true, // Clean up intermediate containers
	}

	buildResponse, err := cli.ImageBuild(ctx, tarReader, buildOptions)
	if err != nil {
		return fmt.Errorf("failed to start image build: %w", err)
	}
	defer buildResponse.Body.Close()

	// Stream build output to standard output
	if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, buildResponse.Body); err != nil {
		log.Printf("Warning: Error streaming build output: %v", err)
	}
	fmt.Printf("Successfully built image '%s'.\n", imageTag)

	return nil
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
