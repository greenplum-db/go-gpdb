package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Environment struct {
	GpHome           string
	MasterDir        string
	PgPort           string
	PgDatabase       string
	SingleOrMulti    string
	GpPerfmonHome    string
	GpccInstanceName string
	GpccPort         string
	GpccVersion      string
	GpccUninstallLoc string
}

// Extract the data from the envFile
func dataExtractor(envFile, search string) string {
	content := readFile(envFile)
	c := contentExtractor(content, fmt.Sprintf("/%s/ {print $2}", search), []string{"FS", "="})
	return removeBlanks(c.String())
}

// Load all the environment information
func environment(envFile string) Environment {
	e := new(Environment)
	e.GpHome = dataExtractor(envFile, "export GPHOME=")
	e.MasterDir = dataExtractor(envFile, "export MASTER_DATA_DIRECTORY=")
	e.PgPort = dataExtractor(envFile, "export PGPORT=")
	e.PgDatabase = dataExtractor(envFile, "export PGDATABASE=")
	e.SingleOrMulti = dataExtractor(envFile, "export singleOrMulti=")
	e.GpPerfmonHome = dataExtractor(envFile, "export GPPERFMONHOME=")
	e.GpccInstanceName = dataExtractor(envFile, "export GPCC_INSTANCE_NAME=")
	e.GpccPort = dataExtractor(envFile, "export GPCCPORT=")
	e.GpccVersion = dataExtractor(envFile, "export GPCCVersion=")
	e.GpccUninstallLoc = dataExtractor(envFile, "export GPCC_UNINSTALL_LOC=")
	return *e
}

// Create environment file of this installation
func (i *Installation) createEnvFile() {
	// Environment file fully qualified path
	i.EnvFile = Config.INSTALL.ENVDIR + "env_" + cmdOptions.Version + "_" + i.Timestamp
	Infof("Creating environment file for this installation at: " + i.EnvFile)

	// Create the file
	createFile(i.EnvFile)
	writeFile(i.EnvFile, []string{
		"export GPHOME=" + i.BinaryInstallationLocation,
		"export PYTHONPATH=$GPHOME/lib/python",
		"export PYTHONHOME=$GPHOME/ext/python",
		"export PATH=$GPHOME/bin:$PYTHONHOME/bin:$PATH",
		"export LD_LIBRARY_PATH=$GPHOME/lib:$PYTHONHOME/lib:$LD_LIBRARY_PATH",
		"export OPENSSL_CONF=$GPHOME/etc/openssl.cnf",
		"export MASTER_DATA_DIRECTORY=" + i.GPInitSystem.MasterDir + "/" + i.GPInitSystem.ArrayName + "-1",
		"export PGPORT=" + i.GPInitSystem.MasterPort,
		"export PGDATABASE=" + i.GPInitSystem.DBName,
		"export singleOrMulti=" + i.SingleORMulti,
	})
}

// Update Environment file
func (i *Installation) updateEnvFile() error {
	Infof("Updating the environment file \"%s\" with the GPCC environment", i.EnvFile)

	// Append to file
	appendFile(i.EnvFile, []string{
		"export GPPERFMONHOME=" + i.GPCC.GpPerfmonHome,
		"export PATH=$GPPERFMONHOME/bin:$PATH",
		"export LD_LIBRARY_PATH=$GPPERFMONHOME/lib:$LD_LIBRARY_PATH",
		"export GPCC_INSTANCE_NAME=" + i.GPCC.InstanceName,
		"export GPCCPORT=" + i.GPCC.InstancePort,
		"export GPCCVersion=" + cmdOptions.CCVersion,
		"export GPCC_UNINSTALL_LOC=" + i.GPCC.UninstallFile,
	})

	return nil
}

// Read file and find the content that we are interested
func readFileAndGatherInformation(file string) string {
	Debugf("Reading and gathering information for the file: %s", file)

	// Obtain the below information from env file
	var output string
	envs := environment(file)
	output = output + "|" + envs.PgPort

	// Is DB running
	if isDbHealthy(file, "") {
		output = output + "| RUNNING"
	} else {
		output = output + "| STOPPED"
	}

	output = output + "|" + envs.GpccInstanceName
	if !IsValueEmpty(envs.GpccPort) {
		output = output + "|" + fmt.Sprintf("http://%s:%s", GetLocalIP(), envs.GpccPort)
	}

	return output
}

// List all the environments available
func ListEnvironmentsInstalled(search string) []string {
	allEnv, _ := FilterDirsGlob(Config.INSTALL.ENVDIR, search)
	return allEnv
}

// List all the installed environment files.
func installedEnvFiles(search, confirmation string, ignoreErr bool) string {
	Debugf("Searching for installed env file using the search string: %s", search)

	var output = []string{`Index | Environment File | Master Port | Status | GPCC Instance Name | GPCC Instance URL`,
		`------|-----------------------|----------------- |------------------|----------------------------------------|------------------------------------------`,
	}

	// Search for environment files for that version
	allEnv := ListEnvironmentsInstalled(search)
	if (len(allEnv) == 0 && !IsValueEmpty(cmdOptions.Version)) && !ignoreErr {
		Fatalf("No installation found for the version: %s, try downloading and installing it", cmdOptions.Version)
	} else if len(allEnv) == 0 && IsValueEmpty(cmdOptions.Version) {
		Fatalf("No installation found, try downloading and installing it")
	}

	// All the environments
	for k, v := range allEnv {
		envName := strings.Replace(v, Config.INSTALL.ENVDIR, "", -1)
		output = append(output, fmt.Sprintf("%s|%s%s", strconv.Itoa(k+1), envName, readFileAndGatherInformation(v)))
	}

	if len(allEnv) > 0 && cmdOptions.ListEnv { // if the option is to list all the env, without prompt

		printOnScreen("Here is a list of env installed", output)
		os.Exit(0)

	} else if len(allEnv) > 0 && confirmation == "confirm" { // Found matching environment file of this installation, now ask for confirmation

		printOnScreen("Here is a list of env installed, confirm if you want to continue", output)

		// Now ask for the confirmation
		confirm := YesOrNoConfirmation()

		// What was the confirmation
		if confirm == "y" { // yes
			Infof("Continuing with the installation of version: %s", cmdOptions.Version)
			return allEnv[0]
		} else { // no
			Infof("Cancelling the installation...")
			os.Exit(0)
		}

	} else if len(allEnv) == 1 && confirmation == "choose" { // if there is only one , then there is no choose just provide the only one
		if !cmdOptions.Vars {
			startDBifNotStarted(allEnv[0])
		}
		return allEnv[0]

	} else if (len(allEnv) > 0 && confirmation == "choose") || confirmation == "list&choose" { // Choose from the given list

		printOnScreen(fmt.Sprintf("Found %d installation, choose from the list", len(allEnv)), output)

		// What is users choice
		choice := PromptChoice(len(allEnv))

		// return the enviornment file to the main function
		choosenEnv := allEnv[choice-1]
		if !cmdOptions.Vars {
			startDBifNotStarted(choosenEnv)
		}
		return choosenEnv
	}

	return ""
}

// List all the environment installed on this box
func env() {
	var envFile string
	// No version provided, show everything
	if cmdOptions.Version == "" {
		Infof("Listing all the environment installed")
		envFile = installedEnvFiles("*", "list&choose", false)
	} else { // Version given, search for env file
		// Don't display and info message when vars called, keep the screen clean
		if !cmdOptions.Vars {
			Infof("Listing all the environment installed with version: %s", cmdOptions.Version)
		}
		envFile = installedEnvFiles("*"+cmdOptions.Version+"*", "choose", false)
	}

	// User asked to print all variables for this environment
	if cmdOptions.Vars {
		cmdOut, err := executeOsCommandOutput("cat", envFile)
		if err != nil {
			Fatalf("Error when trying to read the contents of env file %v: %v", envFile, err)
		}
		fmt.Print(string(cmdOut))
	} else { // Guide user on how to set the environment up
		displayEnvFileToSource(envFile)
	}
}

// Display the env content on the screen
func displayEnvFileToSource(file string) {
	printOnScreen("Source the environment file to set the environment", []string{"source " + file})
}
