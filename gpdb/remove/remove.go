package remove

import (
	"../core"
	"../install"
	"github.com/op/go-logging"
	"strings"
)

var (
	log = logging.MustGetLogger("gpdb")
)

// Run gpdeletesystem
func deleteGPDBEnvUsingGpDeleteSystem() error {

	log.Info("Calling gpdeletesystem to remove the environment")

	// Script arguments
	var deleteSystemArgs []string
	deleteSystemArgs = append(deleteSystemArgs, "source "+install.EnvFileName)
	deleteSystemArgs = append(deleteSystemArgs, "gpdeletesystem -d $MASTER_DATA_DIRECTORY -f << EOF")
	deleteSystemArgs = append(deleteSystemArgs, "y")
	deleteSystemArgs = append(deleteSystemArgs, "y")
	deleteSystemArgs = append(deleteSystemArgs, "EOF")

	// Write it to the file.
	file := core.TempDir + "run_deletesystem.sh"
	err := install.ExecuteBash(file, deleteSystemArgs)
	if err != nil {
		return err
	}

	return nil
}

// Cleaning up all the environments files etc
func deleteGPDBEnvUsingManualMethod(version string, timestamp string) error {

	log.Info("Running cleanup script to clean the environment listing")

	// Script arguments
	var removeAllargs []string
	removeAllargs = append(removeAllargs, "/bin/sh "+core.UninstallDir+"uninstall_"+version+"_"+timestamp)

	// Write it to the file.
	file := core.TempDir + "run_cleanup_script.sh"
	err := install.ExecuteBash(file, removeAllargs)
	if err != nil {
		return err
	}

	return nil
}

// Main Remove method
func Remove(version string) error {

	log.Info("Starting program to uninstall the version: " + version)

	// store the variable
	core.RequestedInstallVersion = version

	// Check if the envfile for that version exists
	chosenEnvFile, err := install.PrevEnvFile("choose")
	if err != nil {
		return err
	}

	// If we receive none, then display the error to user
	var timestamp string
	if core.IsValueEmpty(chosenEnvFile) {
		log.Fatal("Cannot find any environment with the version: " + version)
	} else { // Else store the value
		install.EnvFileName = core.EnvFileDir + chosenEnvFile
		timestamp = strings.Split(chosenEnvFile, "_")[2]
	}

	log.Info("The choosen enviornment file to remove is: " + install.EnvFileName)

	// store the this database port and the GPHOME location
	err = install.ExtractEnvVariables(install.EnvFileName)
	if err != nil {
		return err
	}

	// Check if the database is running, if not then start the database
	err = install.StartDBifNotStarted()
	if err != nil {
		return err
	}

	// Cleanup GPCC if installed.
	if !core.IsValueEmpty(install.GPPERFMONHOME) {
		err = install.UninstallGPCC(timestamp, install.EnvFileName)
		if err != nil {
			return err
		}
	}

	// run delete command
	err = deleteGPDBEnvUsingGpDeleteSystem()
	if err != nil {
		log.Warning("gpdeletesystem failed, running manual cleanup")
		err = deleteGPDBEnvUsingManualMethod(version, timestamp)
		if err != nil {
			return err
		}
	} else {
		log.Info("Deleted GPDB via gpdeletesystem was a success.")
		err = deleteGPDBEnvUsingManualMethod(version, timestamp)
		if err != nil {
			return err
		}
	}

	log.Info("Uninstallation of environment \"" + chosenEnvFile + "\" was a success")
	log.Info("exiting ....")

	return nil
}
