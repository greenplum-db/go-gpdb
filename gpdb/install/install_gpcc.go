package install


import (
	"errors"
	"strings"
	"strconv"
	"time"
	"os"
	"github.com/ielizaga/piv-go-gpdb/core"
)

// GPCC Installer
func InstalSingleNodeGPCC()  error {

	// Check if the provided GPDB version environment file exists
	env, err := PrevEnvFile("choose")
	if err != nil { return err }
	if env == "" {
		return errors.New("Couldn't find any environment file for the database version: " + core.RequestedInstallVersion + ", exiting...")
	} else {
		EnvFileName = core.EnvFileDir + env
		log.Info("Using the environment file \"" + EnvFileName + "\" for installing the GPCC")
	}

	// store the this database port and the GPHOME location
	err = ExtractEnvVariables(EnvFileName)
	if err != nil { return err }

	// extract current date and time
	t := time.Now().Format("20060102150405")

	// If exists then is there a GPCC already installed, if yes then ask for confirmation
	if GPCC_INSTANCE_NAME != "" {
		log.Warning("Found a instance of GPCC already installed on this environment, please confirm if the existing GPCC can be uninstalled")
		// Now ask for the confirmation
		confirm := core.YesOrNoConfirmation()

		// What was the confirmation
		if confirm == "y" {  // yes, then uninstall the old GPCC installation
			err := UninstallGPCC(t, EnvFileName)
			if err != nil { return err }
		} else { // no then exit
			log.Info("Cancelling the installation...")
			os.Exit(0)
		}
	}

	// Check if the database is running, if not then start the database
	err = StartDBifNotStarted()
	if err != nil { return err }

	// Check if the binaries exists on the directory
	// if yes, Unzip the binaries if its file is zipped
	binary_file, err := UnzipBinary(core.RequestedCCInstallVersion)
	if err != nil { return err }

	// Extract the binaries.
	BinaryInstallLocation = "/usr/local/greenplum-cc-web-dbv-" + core.RequestedInstallVersion + "-ccv-" + core.RequestedCCInstallVersion
	var script_option = []string{"yes", BinaryInstallLocation, "yes", "yes"}
	err = ExecuteBinaries(binary_file, InstallGPDBBashFileName, script_option)
	if err != nil { return err }

	// If this the first time then GPPERFHOME would not be there
	// on the environment file, so we update the global variable here
	if core.IsValueEmpty(GPPERFMONHOME) {
		GPPERFMONHOME = BinaryInstallLocation
	}

	// Install the command center database (GPPERFMON)
	err = InstallGpperfmon()
	if err != nil { return err }

	// Restart the database to make the changes to take into effect
	log.Info("Restarting the database so that the command center agents can be started")
	err = StopDB()
	if err != nil { return err }
	err = StartDB()
	if err != nil { return err }

	// Verify the command center is properly installed
	err = WasGPCCInstallationSucess()
	if err != nil { return err }

	// Whats is the next available port
	gpcc_port, err := DoWeHavePortBase(core.FutureRefDir, GPCCPortBaseFileName, "GPCC_PORT")
	if err != nil { return err }
	if gpcc_port != "" {
		GPCC_PORT = string(gpcc_port)[11:]
	} else {
		log.Warning("Didn't find GPCC_PORT in the file, setting it to default value: " + strconv.Itoa(GPCC_PORT_BASE))
		GPCC_PORT = strconv.Itoa(GPCC_PORT_BASE)
	}

	// Is the port used
	ccp, err := strconv.Atoi(GPCC_PORT)
	if err != nil { return err }
	ccp, err = IsPortUsed(ccp, 1)
	if err != nil { return err }
	log.Info("Setting the GPCC_PORT has: "+ strconv.Itoa(ccp))
	GPCC_PORT = strconv.Itoa(ccp)

	// Install the GPCC Web UI without WLM
	cc_name := "gpcc_" + core.RequestedInstallVersion + "_" + core.RequestedCCInstallVersion + "_" + t
	cc_name = strings.Replace(cc_name, ".", "", -1)
	_ = InstallGPCCWEBUI(cc_name, ccp)

	// Start the GPCC Web UI
	err = StartGPCC(cc_name, BinaryInstallLocation)
	if err != nil { return err }

	// If success update the environment file
	err = UpdateEnvFile(cc_name)

	// Check if the port is not greater than 63000, since unix limit is 64000
	port, _ := strconv.Atoi(GPCC_PORT)
	if  port > 63000 {
		log.Warning("PORT has execeeded the unix port limit, setting it to default.")
		GPCC_PORT = strconv.Itoa(GPCC_PORT_BASE)
	}

	// Store the last used port
	err = StoreLastUsedGPCCPort()
	if err != nil { return err }

	// Install the WLM (it works but complication due to non uniformity of WLM structure, will work on it on future )
	//err = Installwlm(t)
	//if err != nil { return err }

	// Start the new browser for use ( Work's , but cant figure out way to run the firefox web in the background )
	//err = StartGPCCBrowser()
	//if err != nil { return err }

	// Open terminal after source the environment file
	err = SetVersionEnv(EnvFileName)
	if err != nil { return err }

	ip, _ := GetLocalIP()
	log.Info("GPCC Web URL: http://"+ ip +":" + GPCC_PORT + ", username / password: gpmon / " + core.EnvYAML.Install.GpMonPass)
	log.Info("Installation of GPCC/WLM software has been completed successfully")

	return nil
}
