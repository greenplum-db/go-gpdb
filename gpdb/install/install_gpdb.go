package install

import (
	"../core"
	"github.com/op/go-logging"
	"time"
)

var (
	log = logging.MustGetLogger("gpdb")
)

// GPDB Single Node Installer
func InstallSingleNodeGPDB() error {

	// If the install is called from download command the set default values
	if !core.IsValueEmpty(core.RequestedDownloadVersion) {
		core.RequestedInstallVersion = core.RequestedDownloadVersion
	}

	// Validate the master & segment exists and is readable
	err := DirValidator(core.EnvYAML.Install.MasterDataDirectory, core.EnvYAML.Install.SegmentDataDirectory)
	if err != nil {
		return err
	}

	// Check if there is already a previous version of the same version
	_, err = PrevEnvFile("confirm")
	if err != nil {
		return err
	}

	// Unzip the binaries, if its file is zipped
	binary_file, err := UnzipBinary(core.RequestedInstallVersion)
	if err != nil {
		return err
	}

	// execute the binaries.
	BinaryInstallLocation = "/usr/local/greenplum-db-" + core.RequestedInstallVersion
	var script_option = []string{"yes", BinaryInstallLocation, "yes", "yes"}
	err = ExecuteBinaries(binary_file, InstallGPDBBashFileName, script_option)
	if err != nil {
		return err
	}

	// Check ssh to host is working and enable password less login
	err = CheckHostnameIsValid()
	if err != nil {
		return err
	}

	// Generate gpinitsystem config file
	t := time.Now().Format("20060102150405")
	err = BuildGpInitSystemConfig(t)
	if err != nil {
		return err
	}

	// Store the Master Port to the variable
	ThisDBMasterPort = GpInitSystemConfig.MasterPort

	// Stop any database if any
	err = StopAllDB()
	if err != nil {
		return err
	}

	// Execute gpinitsystem ( For some reason the gpinitsystem produces exit code
	// 1 even then command ran successfully , so we ignore the exit code 1 here )
	err = ExecuteGpInitSystem()
	if err != nil && err.Error() != "exit status 1" {
		return err
	}

	// Check if the database is healthy
	err = IsDBHealthy()
	if err != nil {
		return err
	}

	// Create Environment file for this installation
	err = CreateEnvFile(t)
	if err != nil {
		return err
	}

	// Uninstall script for this installation
	err = CreateUnistallScript(t)
	if err != nil {
		return err
	}

	// Check if the port is not greater than 63000, since unix limit is 64000
	if GpInitSystemConfig.PortBase > 63000 || GpInitSystemConfig.MasterPort > 63000 {
		log.Warning("PORT has execeeded the unix port limit, setting it to default.")
		GpInitSystemConfig.PortBase = PORT_BASE
		GpInitSystemConfig.MasterPort = MASTER_PORT
	}

	// Store that last used port
	err = StoreLastUsedPort()
	if err != nil {
		return err
	}

	// Open terminal after source the environment file
	err = SetVersionEnv(EnvFileName)
	if err != nil {
		return err
	}

	log.Info("Installation of GPDB software has been completed successfully")

	return nil

}
