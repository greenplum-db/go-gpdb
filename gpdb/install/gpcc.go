package install


import (
	"os/exec"
	"errors"
	"strconv"
	"strings"
	"github.com/ielizaga/piv-go-gpdb/core"
)


// Installed GPPERFMON software to the database
func InstallGpperfmon() error {

	log.Info("Installing GPPERFMON Database..")

	// Check if the password for gpmon is provided on config.yml, else set it to default "changeme"
	if core.IsValueEmpty(core.EnvYAML.Install.GpMonPass) {
		log.Warning("The value of GPMON_PASS is missing, setting the password of gpperfmon & gpmon defaulted to \"changeme\"")
		core.EnvYAML.Install.GpMonPass = "changeme"
	}

	// Installation script
	var install_script []string
	install_script = append(install_script, "source " + EnvFileName)
	install_script = append(install_script, "gpperfmon_install --enable --password "+ core.EnvYAML.Install.GpMonPass +" --port $PGPORT")
	install_script = append(install_script, "echo \"host      all     gpmon    0.0.0.0/0    md5\" >> $MASTER_DATA_DIRECTORY/pg_hba.conf")
	install_script = append(install_script, "echo \"host     all         gpmon         ::1/128       md5\" >> $MASTER_DATA_DIRECTORY/pg_hba.conf")

	// Execute the script
	temp_file := core.TempDir + "install_gpperfmon.sh"
	err := ExecuteBash(temp_file, install_script)
	if err != nil { return err } else {}

	return nil
}

// Check if the GPCC has installed correctly
func WasGPCCInstallationSucess() error {

	// Is gpmmon running
	log.Info("Checking if GPMMON process is running")
	_ , err := exec.Command("pgrep", "gpmmon").Output()
	if err != nil {
		return errors.New("GPMMON process is not running")
	} else {
		log.Info("GPMMON process is running, check if the gpperfmon database has all the required tables.")
	}

	// Can we access gpperfmon database
	log.Info("Checking if GPPERFMON database can be accessed")
	query_string := "select * from system_now"
	_, err = execute_db_query(query_string, "gpperfmon",false, "")
	if err != nil {
		return errors.New("Cannot read the gpperfmon database, installation of gpperfmon database failed..")
	} else {
		log.Info("GPPERFMON Database can be accessed, continuing the script...")
	}

	//

	return nil
}

// Install the Command Center Web UI
func InstallGPCCUI(args []string, cc_home string) error {

	log.Info("Installing command center WEB UI")
	// Check if the version of command center is of 1.x, since that version didn't have WLM.
	var gpcmdrArgs []string
	gpcmdrArgs = append(gpcmdrArgs, "source " + EnvFileName)
	gpcmdrArgs = append(gpcmdrArgs, "source " + cc_home + "/gpcc_path.sh")
	gpcmdrArgs = append(gpcmdrArgs, "echo")
	gpcmdrArgs = append(gpcmdrArgs, "gpcmdr --setup << EOF")
	for _, arg := range args {
		gpcmdrArgs = append(gpcmdrArgs, arg)
	}
	gpcmdrArgs = append(gpcmdrArgs, "echo")

	// Write it to the file and execute
	file := core.TempDir + "gpcmdr_setup.sh"
	err := ExecuteBash(file, gpcmdrArgs)
	if err != nil { return err }

	return nil
}

// Install GPCC WebUI
func InstallGPCCWEBUI(cc_name string, ccp int) error {

	log.Info("Running the setup for installing GPCC WEB UI")

	// CC OPtion for different version of cc installer
	if strings.HasPrefix(core.RequestedCCInstallVersion, "1") { // CC Version 1.x doesn't have WLM and the option used are
		var script_option = []string{cc_name, "n", cc_name, strconv.Itoa(ThisDBMasterPort), GPCC_PORT, "n", "n", "n", "n", "EOF"}
		InstallGPCCUI(script_option, BinaryInstallLocation)
	} else if strings.HasPrefix(core.RequestedCCInstallVersion, "2.5") || strings.HasPrefix(core.RequestedCCInstallVersion, "2.4") { // Option of CC 2.5 & after
		InstallWLM = true
		var script_option = []string{cc_name, "n", cc_name, strconv.Itoa(ThisDBMasterPort), GPCC_PORT, strconv.Itoa(ccp + 1), "n", "n", "n", "n", "EOF"}
		InstallGPCCUI(script_option, BinaryInstallLocation)
	} else if strings.HasPrefix(core.RequestedCCInstallVersion, "2.1") || strings.HasPrefix(core.RequestedCCInstallVersion, "2.0") { // Option of CC 2.0 & after
		InstallWLM = true
		var script_option = []string{cc_name, "n", cc_name, strconv.Itoa(ThisDBMasterPort), "n", GPCC_PORT, "n", "n", "n", "n", "EOF"}
		InstallGPCCUI(script_option, BinaryInstallLocation)
	} else if strings.HasPrefix(core.RequestedCCInstallVersion, "2") { // Option for other version of cc 2.x
		InstallWLM = true
		var script_option = []string{cc_name, "n", cc_name, strconv.Itoa(ThisDBMasterPort), "n", GPCC_PORT, strconv.Itoa(ccp + 1), "n", "n", "n", "n", "EOF"}
		InstallGPCCUI(script_option, BinaryInstallLocation)
	} else if strings.HasPrefix(core.RequestedCCInstallVersion, "3.0") { // Option for CC version 3.0
		InstallWLM = true
		var script_option = []string{cc_name, cc_name, "n", strconv.Itoa(ThisDBMasterPort), GPCC_PORT, "n", "n", "EOF"}
		InstallGPCCUI(script_option, BinaryInstallLocation)
	} else { // All the newer version option unless changed.
		InstallWLM = true
		var script_option = []string{cc_name, cc_name, "n", strconv.Itoa(ThisDBMasterPort), "n", GPCC_PORT, "n", "n", "EOF"}
		InstallGPCCUI(script_option, BinaryInstallLocation)
	}

	return nil
}

// Store the last used port
func StoreLastUsedGPCCPort() error {

	// Fully qualified filename
	FurtureRefFile := core.FutureRefDir + GPCCPortBaseFileName
	log.Info("Storing the last used ports for this installation at: " + FurtureRefFile)

	// Delete the file if already exists
	_ = core.DeleteFile(FurtureRefFile)
	err := core.CreateFile(FurtureRefFile)
	if err != nil {return err}

	// Contents to write to the file
	var FutureContents []string
	port, _ := strconv.Atoi(GPCC_PORT)
	port = port + 3
	FutureContents = append(FutureContents, "GPCC_PORT: " + strconv.Itoa(port))

	// Write to the file
	err = core.WriteFile(FurtureRefFile, FutureContents)
	if err != nil {return err}

	return nil
}

// Uninstall gpcc
func UninstallGPCC(t string, env_file string) error {

	log.Info("Uninstalling the version of command center that is currently installed on this environment.")

	// Generate the script
	var uninstallGPCCArgs []string
	uninstallGPCCArgs = append(uninstallGPCCArgs, "source " + env_file)
	uninstallGPCCArgs = append(uninstallGPCCArgs, "source " + GPPERFMONHOME + "/gpcc_path.sh")
	uninstallGPCCArgs = append(uninstallGPCCArgs, "gpcmdr --stop " + GPCC_INSTANCE_NAME + " &>/dev/null" )
	uninstallGPCCArgs = append(uninstallGPCCArgs, "rm -rf " + GPPERFMONHOME + "/instances/" + GPCC_INSTANCE_NAME)
	uninstallGPCCArgs = append(uninstallGPCCArgs, "gpconfig -c gp_enable_gpperfmon -v off &>/dev/null")
	uninstallGPCCArgs = append(uninstallGPCCArgs, "cp $MASTER_DATA_DIRECTORY/pg_hba.conf $MASTER_DATA_DIRECTORY/pg_hba.conf." + t)
	uninstallGPCCArgs = append(uninstallGPCCArgs, "grep -v gpmon $MASTER_DATA_DIRECTORY/pg_hba.conf."+ t +" > $MASTER_DATA_DIRECTORY/pg_hba.conf")
	uninstallGPCCArgs = append(uninstallGPCCArgs, "rm -rf $MASTER_DATA_DIRECTORY/pg_hba.conf."+ t )
	uninstallGPCCArgs = append(uninstallGPCCArgs, "psql -d template1 -p $PGPORT -Atc \"drop database gpperfmon\" &>/dev/null")
	uninstallGPCCArgs = append(uninstallGPCCArgs, "psql -d template1 -p $PGPORT -Atc \"drop role gpmon\" &>/dev/null")
	uninstallGPCCArgs = append(uninstallGPCCArgs, "rm -rf $MASTER_DATA_DIRECTORY/gpperfmon/*")
	uninstallGPCCArgs = append(uninstallGPCCArgs, "cp "+ env_file +" "+ env_file + "." + t )
	uninstallGPCCArgs = append(uninstallGPCCArgs, "egrep -v \"GPPERFMONHOME|GPCC_INSTANCE_NAME|GPCCPORT\" " + env_file + "." + t + " > " + env_file)
	uninstallGPCCArgs = append(uninstallGPCCArgs, "rm -rf "+ env_file + "." + t)

	// Write it to the file.
	file := core.TempDir + "uninstall_gpcc.sh"
	err := ExecuteBash(file, uninstallGPCCArgs)
	if err != nil { return err }

	return nil
}