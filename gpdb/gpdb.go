package main

import (
	"github.com/op/go-logging"
	"os"
	"github.com/ielizaga/piv-go-gpdb/core"
	"github.com/ielizaga/piv-go-gpdb/argParser"
	"github.com/ielizaga/piv-go-gpdb/download"
	"github.com/ielizaga/piv-go-gpdb/install"
	"github.com/ielizaga/piv-go-gpdb/env"
	"github.com/ielizaga/piv-go-gpdb/remove"
)

// Define the logging format, used in the project
var (
	log    = logging.MustGetLogger("gpdb")
	format = logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05.000}:%{level:s} > %{color:reset}%{message}`,
	)
)

func main() {

	// Logger for go-logging package
	// create backend for os.Stderr, set the format and update the logger to what logger to be used
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)

	// Extract all the OS command line arguments
	if !core.InstallAfterDownload {
		argParser.ArgParser()
	}

	// Read the configuration file and set the variable.
	log.Info("Extracting all the configuration parameters")
	err := core.Config()
	core.Fatal_handler(err)

	// Run Program based on what option is specified
	switch core.ArgOption {
	case "download":
		// Run Download
		err := download.Download()
		core.Fatal_handler(err)

		// If user has asked to run the installer after download then
		// set the argument and rerun the main to pick the right choice.
		if core.InstallAfterDownload {
			core.ArgOption = "install"
			core.RequestedInstallProduct = core.RequestedDownloadProduct
			main()
		}

	case "install":
		// If the product to install is GPDB then call the GPDB installation program
		if core.RequestedInstallProduct == "gpdb" { // Run Install GPDB
			err := install.InstallSingleNodeGPDB()
			core.Fatal_handler(err)
		} else if core.RequestedInstallProduct == "gpcc" { // else GPCC installer
			err := install.InstalSingleNodeGPCC()
			core.Fatal_handler(err)
		}
	case "remove":
		// Run remove
		err := remove.Remove(core.RequestedRemoveVersion)
		core.Fatal_handler(err)
	case "env":
		//Run env
		err := env.Environment(core.RequestedVersionEnv)
		core.Fatal_handler(err)
	default:
		// Error if command is invalid
		log.Fatal("Command not recognized")
	}

}