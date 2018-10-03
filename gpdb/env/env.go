package env

import (
	"errors"
	"github.com/ielizaga/piv-go-gpdb/core"
	"github.com/op/go-logging"
	"github.com/ielizaga/piv-go-gpdb/install"
)

var (
	log = logging.MustGetLogger("gpdb")
)

// Set the environment
func SettheChoosenEnv(chosenEnvFile string, version string) error {
	// If we receive none, then display the error to user
	if core.IsValueEmpty(chosenEnvFile) {
		log.Fatal("Cannot find any environment with the version: " + version )
	} else { // Else store the value
		install.EnvFileName = core.EnvFileDir + chosenEnvFile
	}

	log.Info("The choosen enviornment file is: " + install.EnvFileName)

	// store the this database port and the GPHOME location
	err := install.ExtractEnvVariables(install.EnvFileName)
	if err != nil { return err }

	// Check if the database is running, if not then start the database
	err = install.StartDBifNotStarted()
	if err != nil { return err }

	log.Info("Environment has been setup and ready to use")

	// Open terminal after source the environment file
	err = install.SetVersionEnv(install.EnvFileName)
	if err != nil { return err }

	log.Info("Environment setup is complete")

	return nil
}


// Function to setup environment
func Environment(version string) error {

	// If no version is provided then the user is requesting to list all the environment installed
	if version == "" {
		log.Info("listing all the environment version installed on this cluster")
		chosenEnvFile, err := install.PrevEnvFile("listandchoose")
		if err != nil { return err }

		// set the chosen environment
		if core.IsValueEmpty(chosenEnvFile) {
			return errors.New("There is no any installation of GPDB, please install the product to list the environment here")
		} else {
			err = SettheChoosenEnv(chosenEnvFile, version)
			if err != nil { return err }
		}


	} else { // he is checking for a specific version

		log.Info("listing all the environment that has been installed with version: " + version)

		// store the variable
		core.RequestedInstallVersion = version

		// Get the env files that we know about
		chosenEnvFile, err := install.PrevEnvFile("choose")
		if err != nil { return err }

		// set the chosen environment
		err = SettheChoosenEnv(chosenEnvFile, version)
		if err != nil { return err }
	}

	log.Info("Exiting ..... ")
	return nil
}