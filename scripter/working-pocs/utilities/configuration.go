package utilities

import (
	"scripter/entities"
)

type Configuration struct {
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

func (configuration Configuration) SetConfigurationFromSignal(signal entities.Signal) Configuration {
	configuration.HostOs = signal.HostOs
	configuration.SignalOs = signal.SignalOs
	configuration.Containerize = signal.Containerize
	configuration.Vmize = signal.Vmize
	configuration.PackageInstaller = signal.PackageInstaller

	return configuration
}

func (configuration Configuration) SetGeneralConfiguration(signal entities.Signal) {
	if !checkConfigurationPlatformCombination(configuration) {
		panic("Host/signal OS combination, signal type, or package installer not feasible (Currenltly only handling one virtualization level per signal)!")
	}
}

func checkConfigurationPlatformCombination(configuration Configuration) bool {
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
