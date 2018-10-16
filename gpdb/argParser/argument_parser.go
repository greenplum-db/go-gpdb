package argParser

import (
	"../core"
	"flag"
	"fmt"
	"os"
	"strings"
)

// OS Argument Parser
func ArgParser() {

	// Download Command Parser
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	DownloadProductFlag := downloadCmd.String("p", "gpdb", "What product do you want to Install? [OPTIONS: gpdb, gpcc, gpextras]")
	DownloadVersionFlag := downloadCmd.String("v", "", "OPTIONAL: Which GPDB version software do you want to list ?")
	DownloadInstallFlag := downloadCmd.Bool("install", false, "OPTIONAL: Install after downloaded (Only works with \"gpdb\")?")

	// Install Command Parser
	InstallCmd := flag.NewFlagSet("install", flag.ExitOnError)
	InstallProductFlag := InstallCmd.String("p", "gpdb", "What product do you want to Install [OPTIONS: gpdb, gpcc]?")
	InstallVersionFlag := InstallCmd.String("v", "", "Which version do you want to Install ?")
	InstallCCVersionFlag := InstallCmd.String("c", "", "What is the version of GPCC that you can to install (for only -p gpcc) ?")

	// Remove Command Parser
	RemoveCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	RemoveVersionFlag := RemoveCmd.String("v", "", "Provide the version from the installed list to remove")

	// Environment Command Parser
	EnvCmd := flag.NewFlagSet("env", flag.ExitOnError)
	EnvVersionFlag := EnvCmd.String("v", "", "Provide the version from the installed list to remove")

	// If no COMMAND keyword provided then show the help menu.
	if len(os.Args) == 1 {
		ShowHelp()
	}

	// If there is a command keyword provided then check to what is it and then parse the appropriate options
	switch os.Args[1] {
	case "download": // Download command parser
		core.ArgOption = "download"
		downloadCmd.Parse(os.Args[2:])
	case "install": // Install command parser
		core.ArgOption = "install"
		InstallCmd.Parse(os.Args[2:])
	case "env": // Env command parser
		core.ArgOption = "env"
		EnvCmd.Parse(os.Args[2:])
	case "remove": // Remove command parser
		core.ArgOption = "remove"
		RemoveCmd.Parse(os.Args[2:])
	case "version": // Version of the software
		fmt.Printf("Version: %.1f\n", core.Version)
		os.Exit(0)
	case "help": // Help Menu
		ShowHelp()
	default: // Error if command is invalid
		fmt.Printf("ERROR: %q is not valid command.\n", os.Args[1])
		ShowHelp()
	}

	// If the command send is download, then parse the commandline arguments
	if downloadCmd.Parsed() {

		// If the product parameter is passed, then check if its valid value.
		if *DownloadProductFlag != "" {

			// If its valid then we are going to store it
			if core.Contains(core.AcceptedDownloadProduct, *DownloadProductFlag) {
				core.RequestedDownloadProduct = strings.ToLower(*DownloadProductFlag)
			} else { // Else print error to choose the right value
				fmt.Println("ERROR: Invalid options provided to the argument -p")
				fmt.Print("Usage of download: \n")
				downloadCmd.PrintDefaults()
				os.Exit(2)
			}
		}

		// If the version parameter is passed, then store the value
		if *DownloadVersionFlag != "" {
			core.RequestedDownloadVersion = *DownloadVersionFlag
		}

		// If the user request to install the product after download
		if *DownloadInstallFlag && core.RequestedDownloadProduct == "gpdb" {
			core.InstallAfterDownload = true
		} else if *DownloadInstallFlag {
			fmt.Println("ERROR: -install only works with -p \"gpdb\"")
			fmt.Print("Usage of download: \n")
			downloadCmd.PrintDefaults()
			os.Exit(2)
		}
	}

	// Install command argument parser
	if InstallCmd.Parsed() {
		if *InstallProductFlag != "" {
			// If its valid then we are going to store it
			if core.Contains(core.AcceptedInstallProduct, *InstallProductFlag) {
				core.RequestedInstallProduct = strings.ToLower(*InstallProductFlag)
			} else { // Else print error to choose the right value
				fmt.Println("ERROR: Invalid options provided to the argument -p")
				fmt.Print("Usage of install: \n")
				InstallCmd.PrintDefaults()
				os.Exit(2)
			}
		}

		// If the version parameter is passed, then store the value
		if *InstallVersionFlag != "" {
			core.RequestedInstallVersion = *InstallVersionFlag
		} else { // On install this is a required parameter, if not provided then terminate and ask to enter the version
			fmt.Println("ERROR: GPDB Version missing, Please provide the version that needs to be used")
			fmt.Print("Usage of install: \n")
			InstallCmd.PrintDefaults()
			os.Exit(2)
		}

		// If the request to install is GPCC then check if the cc version is provided
		if core.RequestedInstallProduct == "gpcc" && *InstallCCVersionFlag == "" {
			fmt.Println("ERROR: GPCC Version missing, Please provide the version that need to be installed")
			fmt.Print("Usage of install: \n")
			InstallCmd.PrintDefaults()
			os.Exit(2)
		} else {
			core.RequestedCCInstallVersion = *InstallCCVersionFlag
		}
	}

	// If the command is to extract the env listing then
	if EnvCmd.Parsed() {

		if *EnvVersionFlag != "" {
			core.RequestedVersionEnv = *EnvVersionFlag
		}

	}

	// If the command is to remove the environment then remove it
	if RemoveCmd.Parsed() {

		if *RemoveVersionFlag == "" {
			fmt.Println("ERROR: Missing version information for remove command")
			fmt.Print("Usage of remove: \n")
			RemoveCmd.PrintDefaults()
			os.Exit(2)
		} else {
			core.RequestedRemoveVersion = *RemoveVersionFlag
		}

	}

}
