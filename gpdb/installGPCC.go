package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func (i *Installation) preGPCCChecks() {
	Infof("Running the pre checks to install GPCC version \"%s\" on the GPDB Version \"%s\"", cmdOptions.CCVersion, cmdOptions.Version)

	// If this os is ubuntu, currently we dont have any version available to install GPCC on ubuntu
	if getSystemInfoAndCheckIfItsUbuntu() {
		Fatalf("Currently install of GPCC via %s tool is not supported", programName)
	}

	// Check if there is already a version of GPDB installed
	i.EnvFile = installedEnvFiles(fmt.Sprintf("*%s*", cmdOptions.Version), "choose", false)

	// Extract the environment information
	i.extractEnvValues()

	// Extract current date and time
	i.Timestamp = time.Now().Format("20060102150405")

	// Check if this env have GPCC Installed
	i.doesThisEnvHasGPCCInstalled()

	// Check if the database is running, if not then start the database
	startDBifNotStarted(i.EnvFile)
}

// Extract environment values from the env file
func (i *Installation) extractEnvValues() {
	Infof("Extracting the environment information from the file: %s", i.EnvFile)
	envs := environment(i.EnvFile)
	i.GPInitSystem.MasterPort = envs.PgPort
	i.GPInitSystem.MasterDir = envs.MasterDir
	i.GPCC.InstanceName = envs.GpccInstanceName
	i.GPCC.InstancePort = envs.GpccPort
	i.GPCC.GpPerfmonHome = envs.GpPerfmonHome
}

// Check if this version of database has GPCC installed.
func (i *Installation) doesThisEnvHasGPCCInstalled() {
	Debugf("Checking if there is already a version of gpcc installed on this database version")
	// If exists then is there a GPCC already installed, if yes then ask for confirmation
	if i.GPCC.InstanceName != "" {
		Warnf("Found a instance of GPCC already installed on this environment, please confirm if the existing GPCC can be uninstalled")
		// Now ask for the confirmation
		confirm := YesOrNoConfirmation()
		// What was the confirmation
		if confirm == "y" { // yes, then uninstall the old GPCC installation
			removeGPCC(i.EnvFile)
		} else { // no then exit
			Infof("Cancelling the installation...")
			os.Exit(0)
		}
	}
}

// Install the product that is requested
func (i *Installation) installGPCCProduct() {

	Infof("Installing GPCC version \"%s\" on the GPDB Version \"%s\"", cmdOptions.CCVersion, cmdOptions.Version)

	// Hostfile is create by the install gpdb, since we can't run the
	// install gpcc without gpdb, so we use the file created by the install gpdb
	i.WorkingHostFileLocation = os.Getenv("HOME") + "/hostfile"
	exists, _ := doesFileOrDirExists(i.WorkingHostFileLocation)
	if !exists {
		Fatalf("No host file found on the home location of the user")
	}

	// The installation b/w GPCC 3 & GPCC 4 is completely different so
	// we need to create two working libraries to make it work.
	if isThis4x() {
		i.installGPCC4xAndAbove()
	} else {
		i.installGPCCBelow4x()
	}

}

// Is it a 4.x to <4.x version
func isThis4x() bool {
	v := extractVersion(cmdOptions.CCVersion)
	if v >= 4.0 { // newer GPCC code
		return true
	} else { // legacy GPCC code
		return false
	}
}

// Install GPCC 4.x
func (i *Installation) installGPCC4xAndAbove() {
	// unzip the binaries
	i.GPCC.GPCCBinaryLoc, _ = locateAndExtractPackage(fmt.Sprintf("*%s*.zip", cmdOptions.CCVersion))

	// Install the Gpperfmon database
	i.installGpperfmon()

	// Install the gpcc binaries.
	i.installGPCCBinaries4x()
}

// Install GPCC less than 4.x
func (i *Installation) installGPCCBelow4x() {

	// Install the product
	i.installGPCCBinariesBelow4x()

	// If its a multi host deployment then install the gpcc on all the host
	i.InstallGPCCBinariesIfMultiHost()

	// Install the Gpperfmon database
	i.installGpperfmon()

	// Install the GPCC web Interface
	i.InstallWebUIBelow4x()
}

// Install the product that is requested
func (i *Installation) postGPCCInstall() {

	Infof("Running the post steps for the installation of GPCC version \"%s\" on the GPDB Version \"%s\"", cmdOptions.CCVersion, cmdOptions.Version)

	// Store the last used port for future use
	i.savePort()

	// Create the uninstall script
	i.uninstallGPCCScript()

	// Update the envfile
	i.updateEnvFile()

	// start the GPCC
	startGPCC(i.EnvFile)

	// Installation complete, print on the screen the env file to source and cleanup temp files
	displayEnvFileToSource(i.EnvFile)

	// get the ip of the host and display on the screen
	Infof("GPCC Web URL: http://%s:%s, username / password: gpmon / %s", GetLocalIP(), i.GPCC.InstancePort, Config.INSTALL.GPMONPASS)
}

// Installed GPPERFMON software to the database
func (i *Installation) installGpperfmon() error {

	Infof("Installing GPPERFMON Database on the database with environment file: %s", i.EnvFile)

	// Installation script
	GpPerfMonFile := Config.CORE.TEMPDIR + "install_gpperfmon.sh"
	generateBashFileAndExecuteTheBashFile(GpPerfMonFile, "/bin/sh", []string{
		"source " + i.EnvFile,
		"gpperfmon_install --enable --password " + Config.INSTALL.GPMONPASS + " --port $PGPORT",
		"echo \"local    gpperfmon   gpmon         md5\" >> $MASTER_DATA_DIRECTORY/pg_hba.conf",
		"echo \"host      all     gpmon    0.0.0.0/0    md5\" >> $MASTER_DATA_DIRECTORY/pg_hba.conf",
		"echo \"host      all     gpmon    127.0.0.1/28    md5\" >> $MASTER_DATA_DIRECTORY/pg_hba.conf",
		"echo \"host     all         gpmon         ::1/128       md5\" >> $MASTER_DATA_DIRECTORY/pg_hba.conf",
	})

	// Verify if the gpperfmon process are running
	i.verifyGpperfmon()
	return nil
}

// Verify if database gpperfmon is installed correctly
func (i *Installation) verifyGpperfmon() {

	// Restart
	Infof("Restarting the database so that the command center agents can be started")
	stopDB(i.EnvFile)
	startDB(i.EnvFile)

	// Verify GPmon process
	Infof("Checking if GPMMON process is running")
	_, err := executeOsCommandOutput("pgrep", "gpmmon")
	if err != nil {
		Fatalf("GPPMON process is not running, the installation of gpperfmon was a failure, err: %$v", err)
	} else {
		Infof("GPMMON process is running, check if the gpperfmon database has all the required tables.")
	}

	// Verify gpperfmon database
	Infof("Checking if GPPERFMON database can be accessed")
	queryString := "select * from system_now"
	gpperfmonHealthFile := Config.CORE.TEMPDIR + "gpperfmon_health.sh"
	deleteFile(gpperfmonHealthFile)
	createFile(gpperfmonHealthFile)
	writeFile(gpperfmonHealthFile, []string{
		fmt.Sprintf("source %s", i.EnvFile),
		fmt.Sprintf("psql -d gpperfmon -U gpmon -c '%s'", queryString),
	})
	_, err = executeOsCommandOutput("/bin/sh", gpperfmonHealthFile)
	if err != nil {
		Fatalf("GPPERFMON database is not accessible, err: %v", err)
	} else {
		Infof("GPPERFMON Database can be accessed, continuing the script...")
	}
}

// Generate command center instance name
func commandCenterInstanceName() string {
	return fmt.Sprintf("ccversion_%s_%s", strings.Replace(cmdOptions.Version, ".", "", -1), strings.Replace(cmdOptions.CCVersion, ".", "", -1))
}

// Create configuration file
func (i *Installation) createGPCCCOnfigurationFile() string {
	i.GPCC.InstancePort = i.validatePort("GPCC_PORT", defaultGpccPort)
	i.GPCC.WebSocketPort = i.validatePort("WEBSOCKET_PORT", defaultWebSocket)
	i.GPCC.InstanceName = commandCenterInstanceName()
	gpccInstallConfig := Config.CORE.TEMPDIR + fmt.Sprintf("gpcc_config_%s_%s", cmdOptions.CCVersion, i.Timestamp)
	deleteFile(gpccInstallConfig)
	createFile(gpccInstallConfig)
	writeFile(gpccInstallConfig, []string{
		"path = /usr/local",
		"display_name = " + i.GPCC.InstanceName,
		"master_port = " + i.GPInitSystem.MasterPort,
		"web_port = " + i.GPCC.InstancePort,
		"rpc_port = " + i.GPCC.WebSocketPort,
		"enable_ssl = false",
		"enable_kerberos = false",
	})
	return gpccInstallConfig
}

// Installing gpcc instance >= 4.x
func (i *Installation) installGPCCBinaries4x() {
	Infof("Installing gpcc binaries for the cc version: %s", cmdOptions.CCVersion)
	executeGPPCFile := Config.CORE.TEMPDIR + "execute_gpcc.sh"
	generateBashFileAndExecuteTheBashFile(executeGPPCFile, "/bin/sh", []string{
		fmt.Sprintf("source %s", i.EnvFile),
		fmt.Sprintf("%s -c %s &>/dev/null << EOF", i.GPCC.GPCCBinaryLoc, i.createGPCCCOnfigurationFile()),
		"y",
		"EOF",
	})
	i.extractGPPERFMON()
}

// Installing gpcc instance < 4.x
func (i *Installation) installGPCCBinariesBelow4x() {
	Infof("Installing gpcc binaries for the cc version: %s", cmdOptions.CCVersion)
	binFile, _ := getBinaryFile(cmdOptions.CCVersion)
	i.GPCC.GpPerfmonHome = fmt.Sprintf("/usr/local/greenplum-cc-web-dbv-%s-ccv-%s", cmdOptions.Version, cmdOptions.CCVersion)
	var scriptOption = []string{"yes", i.GPCC.GpPerfmonHome, "yes", "yes"}
	err := executeBinaries(binFile, "install_software.sh", scriptOption)
	if err != nil {
		Fatalf("Failed in installing the binaries, err: %v", err)
	}
}

// Extract GPPERFMON location
func (i *Installation) extractGPPERFMON() {
	Debugf("Extracting the gpperfmon home location")
	file, err := FilterDirsGlob("/usr/local", fmt.Sprintf("*%s*", cmdOptions.CCVersion))
	if err != nil || len(file) == 0 {
		Fatalf("Cannot find the gpcc installation or encountered err: %v", err)
	}
	i.GPCC.GpPerfmonHome = file[0]
}

// On Multi host we need to install the binaries on all the host using gpccinstall
func (i *Installation) InstallGPCCBinariesIfMultiHost() {
	if environment(i.EnvFile).SingleOrMulti == "multi" {
		Infof("Running seginstall to install the gpcc on all the host")
		gpccSegInstallFile := Config.CORE.TEMPDIR + "gpccseginstall.sh"
		generateBashFileAndExecuteTheBashFile(gpccSegInstallFile, "/bin/sh", []string{
			fmt.Sprintf("source %s", i.EnvFile),
			fmt.Sprintf("source %s/gpcc_path.sh", i.GPCC.GpPerfmonHome),
			fmt.Sprintf("gpccinstall -f %s/hostfile", os.Getenv("HOME")),
		})
	} else {
		Infof("Skipping gpccinstall since this is a single machine installation")
	}
}

// Install the web interface for GPCC below 4.x
func (i *Installation) InstallWebUIBelow4x() {
	Infof("Running the setup for installing GPCC WEB UI for command center version: %s", cmdOptions.CCVersion)
	i.GPCC.InstancePort = i.validatePort("GPCC_PORT", defaultGpccPort)
	i.GPCC.WebSocketPort = i.validatePort("WEBSOCKET_PORT", defaultWebSocket) // Safe Guard to prevent 4.x and below clash
	i.GPCC.InstanceName = commandCenterInstanceName()
	var scriptOption []string

	// CC Option for different version of cc installer
	// Not going to refactor this piece of code, since this is legacy and testing all the version option is a pain
	// so we will leave this as it is
	if strings.HasPrefix(cmdOptions.CCVersion, "1") { // CC Version 1.x
		scriptOption = []string{i.GPCC.InstanceName, "n", i.GPCC.InstanceName, i.GPInitSystem.MasterPort, i.GPCC.InstancePort, "n", "n", "n", "n", "EOF"}
	} else if strings.HasPrefix(cmdOptions.CCVersion, "2.5") || strings.HasPrefix(cmdOptions.CCVersion, "2.4") { // Option of CC 2.5 & after
		scriptOption = []string{i.GPCC.InstanceName, "n", i.GPCC.InstanceName, i.GPInitSystem.MasterPort, i.GPCC.InstancePort, strconv.Itoa(strToInt(i.GPCC.InstancePort) + 1), "n", "n", "n", "n", "EOF"}
	} else if strings.HasPrefix(cmdOptions.CCVersion, "2.1") || strings.HasPrefix(cmdOptions.CCVersion, "2.0") { // Option of CC 2.0 & after
		scriptOption = []string{i.GPCC.InstanceName, "n", i.GPCC.InstanceName, i.GPInitSystem.MasterPort, "n", i.GPCC.InstancePort, "n", "n", "n", "n", "EOF"}
	} else if strings.HasPrefix(cmdOptions.CCVersion, "2") { // Option for other version of cc 2.x
		scriptOption = []string{i.GPCC.InstanceName, "n", i.GPCC.InstanceName, i.GPInitSystem.MasterPort, "n", i.GPCC.InstancePort, strconv.Itoa(strToInt(i.GPCC.InstancePort) + 1), "n", "n", "n", "n", "EOF"}
	} else if strings.HasPrefix(i.GPCC.InstanceName, "3.0") { // Option for CC version 3.0
		scriptOption = []string{i.GPCC.InstanceName, i.GPCC.InstanceName, "n", i.GPInitSystem.MasterPort, i.GPCC.InstancePort, "n", "n", "EOF"}
	} else { // All the newer version option unless changed.
		scriptOption = []string{i.GPCC.InstanceName, i.GPCC.InstanceName, "n", i.GPInitSystem.MasterPort, "n", i.GPCC.InstancePort, "n", "n", "EOF"}
	}
	i.installGPCCUI(scriptOption)
}

// Install the Command Center Web UI
func (i *Installation) installGPCCUI(args []string) error {
	installGPCCWebFile := Config.CORE.TEMPDIR + "install_web_ui.sh"
	options := []string{
		"source " + i.EnvFile,
		"source " + i.GPCC.GpPerfmonHome + "/gpcc_path.sh",
		"echo",
		"gpcmdr --setup << EOF",
	}
	for _, arg := range args {
		options = append(options, arg)
	}
	options = append(options, "echo")
	generateBashFileAndExecuteTheBashFile(installGPCCWebFile, "/bin/sh", options)

	return nil
}
