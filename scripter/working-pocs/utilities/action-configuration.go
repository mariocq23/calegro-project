package utilities

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"scripter/entities"
	"strings"
)

type ActionConfiguration struct {
	HostOs           string
	SignalOs         string
	Containerize     bool
	Vmize            bool
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
	configuration.Vmize = signal.Vmize
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
	if configuration.Vmize && configuration.Containerize {
		return false
	}
	if configuration.HostOs == "darwin" && configuration.SignalOs == "windows" && configuration.Containerize {
		return false
	}
	if configuration.HostOs == "darwin" && configuration.SignalOs == "linux" && !configuration.Containerize &&
		!configuration.Vmize {
		return false
	}
	if configuration.HostOs == "darwin" && !stringHandler.ContainsString(feasiblePackInstallersInMac, configuration.PackageInstaller) {
		return false
	}
	if configuration.HostOs == "linux" && (configuration.SignalOs == "windows" || configuration.SignalOs == "darwin") && (configuration.Containerize ||
		!configuration.Vmize) {
		return false
	}
	if configuration.HostOs == "linux" && !stringHandler.ContainsString(feasiblePackInstallersInLinux, configuration.PackageInstaller) {
		return false
	}
	if configuration.HostOs == "linux" && configuration.SignalOs == "windows" && configuration.Containerize {
		return false
	}
	if configuration.HostOs == "windows" && configuration.SignalOs == "linux" && !configuration.Containerize &&
		!configuration.Vmize {
		return false
	}
	if configuration.HostOs == "windows" && !stringHandler.ContainsString(feasiblePackInstallersInWindows, configuration.PackageInstaller) {
		return false
	}
	if configuration.HostOs == "windows" && configuration.SignalOs == "darwin" && (configuration.Containerize || !configuration.Vmize) {
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
		for _, dep := range signal.InstallationDependencies {
			if !isInstalledOnMac(dep) {
				log.Printf("%s is not installed. Proceeding with installation...\n", dep)
				err :=
					installDependencyOnMac(dep)
				if err != nil {
					log.Printf("Failed to install %s.\n", dep)
				}
			}
		}
		return
	}
	if signal.HostOs == "windows" {
		for _, dep := range signal.InstallationDependencies {
			if !isInstalledOnWindows(dep) {
				log.Printf("%s is not installed. Proceeding with installation...\n", dep)
				err :=
					installDependencyOnWindows(dep)
				if err != nil {
					log.Printf("Failed to install %s.\n", dep)
				}
			}
		}
		return
	}
	if signal.HostOs == "linux" {
		for _, dep := range signal.InstallationDependencies {
			if !isInstalledOnLinux(dep) {
				log.Printf("%s is not installed. Proceeding with installation...\n", dep)
				err :=
					installDependencyOnLinux(dep)
				if err != nil {
					log.Printf("Failed to install %s.\n", dep)
				}
			}
		}
		return
	}

}

func installDependencyOnLinux(dep string) any {
	fmt.Printf("Installing %s with real-time verbose output...\n", dep)
	cmd := exec.Command("sudo", "apt", "install", dep)

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

func isInstalledOnLinux(dep string) bool {
	if dep == "docker" {
		cmd := exec.Command(dep, "version")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return false // Assume not installed if there's an error checking
		}
		return strings.Contains(string(output), dep)
	}
	cmd := exec.Command("sudo", "apt", "list", "--versions", dep)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false // Assume not installed if there's an error checking
	}
	return strings.Contains(string(output), dep)
}

func installDependencyOnWindows(dep string) any {
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
