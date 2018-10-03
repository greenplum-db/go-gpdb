package install

import (
	"os/exec"
	"bytes"
	"errors"
	"github.com/ielizaga/piv-go-gpdb/core"
)

func StopAllDB() error {

	SourceGPDBPath()

	// Get all database running
	log.Info("Stopping all the database running on this host, to free up semaphore for this installation")

	// Can't seems to find a simple way to stop all database, so we will built the below
	// simple shell script to execute the stop database command
	cleanupScript := "ps -ef | grep silent | grep master | while read line; " +
		"do " +
		"GPHOME=`echo $line|awk '{print $8}'| rev | cut -d'/' -f3- | rev`;"+
		"export MASTER_DATA_DIRECTORY=`echo $line|awk '{print $10}'`;"+
		"export PGPORT=`echo $line|awk '{print $12}'`;"+
		"export PGDATABASE=template1;"+
		"source $GPHOME/greenplum_path.sh;"+
		"gpstop -af;" +
		"done"

	// Execute the command
	StopScriptLoc := core.TempDir + "stop_all_db.sh"
	var StopScript []string
	StopScript = append(StopScript, cleanupScript)
	StopScript = append(StopScript, "ps -ef | egrep \"gpmon|gpmonws|lighttpd\" | grep -v grep | awk '{print $2}' | xargs -n1 /bin/kill -11 &>/dev/null; echo > /dev/null")
	err := ExecuteBash(StopScriptLoc, StopScript)
	if err != nil { return err }

	// Send a warning message if the process is not completely stopped.
	cmdOut, _ := exec.Command("pgrep", "postgres").Output()
	var EmptyBytes []byte
	if !bytes.Equal(cmdOut, EmptyBytes) {
		log.Warning("Can't stop all postgres process, seems like some are left behind")
	} else {
		log.Info("All Postgres process are stopped on this server")
	}

	// Cleanup temp files.
	core.DeleteFile(StopScriptLoc)

	// Stopping all the WLM instance
	//StopAllWLM()

	return nil
}


// Start database
func StartDB() error {

	log.Info("Attempting to start the database as per the environment file: " + EnvFileName)

	// BashScript
	var BashSricpt []string
	BashSricpt = append(BashSricpt, "source " + EnvFileName)
	BashSricpt = append(BashSricpt, "gpstart -a")

	// Create the file
	temp_file := core.TempDir + "start.sh"

	// Execute the script
	err := ExecuteBash(temp_file, BashSricpt)
	if err != nil { return err }

	// Check if the database is healthy after start
	err = IsDBHealthy()
	if err != nil { return errors.New("Can't seems to start the database in the environment file \"" + EnvFileName + "\"exiting...") }

	// If WLM is installed on this environment then start it
	//if !core.IsValueEmpty(WLMInstallDir) {
	//	StartWLMService()
	//}

	// Start Command Center WEB UI if available on this environment
	if !core.IsValueEmpty(GPPERFMONHOME) {
		StartGPCC(GPCC_INSTANCE_NAME, GPPERFMONHOME)
	}

	return nil
}

// Stop Database
func StopDB() error {

	log.Info("Attempting to stop the database as per the environment file: " + EnvFileName)

	// BashScript
	var BashSricpt []string
	BashSricpt = append(BashSricpt, "source " + EnvFileName)
	BashSricpt = append(BashSricpt, "gpstop -af")

	// Create the file
	temp_file := core.TempDir + "stop.sh"

	// Execute the script
	err := ExecuteBash(temp_file, BashSricpt)
	if err != nil { return err }

	// If WLM is installed on this environment then stop it
	//if !core.IsValueEmpty(WLMInstallDir) {
	//	StopWLMService()
	//}

	// Start Command Center WEB UI if available on this environment
	if !core.IsValueEmpty(GPPERFMONHOME) {
		StopGPCC(GPCC_INSTANCE_NAME, GPPERFMONHOME)
	}

	return nil
}

// Start the database if not started
func StartDBifNotStarted() error {

	// is the database running , then return
	err := IsDBHealthy()
	if err == nil { // Database is started and running
		log.Info("Database seems to be running, contining...")
		return nil
	} else { // database is not running, lets start it up

		log.Warning("Database is not started, attempting to start the database...")

		// Stop all database is not stopped
		err := StopAllDB()
		if err != nil { return err }

		// Start the database of concern
		err = StartDB()
		if err != nil { return err }

		// Check again if the database is healthy
		err = IsDBHealthy()
		if err != nil { return errors.New("Can't seems to start the database in the environment file \"" + EnvFileName + "\"exiting...") }
	}

	return nil
}

// Started GPCC WebUI
func StartGPCC (cc_name string, cc_home string) error {

	log.Info("Starting command center WEB UI")

	// Start the command center web UI
	var gpcmdrStartarg []string
	gpcmdrStartarg = append(gpcmdrStartarg, "source " + EnvFileName)
	gpcmdrStartarg = append(gpcmdrStartarg, "source " + cc_home + "/gpcc_path.sh")
	gpcmdrStartarg = append(gpcmdrStartarg, "gpcmdr --start " + cc_name + " &>/dev/null << EOF")
	gpcmdrStartarg = append(gpcmdrStartarg, "y")
	gpcmdrStartarg = append(gpcmdrStartarg, "EOF")

	// Write it to the file.
	file := core.TempDir + "gpcmdr_start.sh"
	err := ExecuteBash(file, gpcmdrStartarg)
	if err != nil { return err }

	return nil
}


// Stop GPCC Instance
func StopGPCC(cc_name string, cc_home string) error {

	log.Info("Stop command center WEB UI")

	// Start the command center web UI
	var gpcmdrStartarg []string
	gpcmdrStartarg = append(gpcmdrStartarg, "source " + EnvFileName)
	gpcmdrStartarg = append(gpcmdrStartarg, "source " + cc_home + "/gpcc_path.sh")
	gpcmdrStartarg = append(gpcmdrStartarg, "gpcmdr --stop " + cc_name)

	// Write it to the file.
	file := core.TempDir + "gpcmdr_stop.sh"
	err := ExecuteBash(file, gpcmdrStartarg)
	if err != nil { return err }

	return nil

}

// Start GPCC browser
// Not in use, since calling firefox halt the screen until the user closes the browser.
// will enable once the issue is fixed.
func StartGPCCBrowser() error {

	log.Info("Starting the GPCC Web Console")

	// Starting the browser for GPCC environment
	var StartGPCCWeb []string
	StartGPCCWeb = append(StartGPCCWeb, "LD_LIBRARY_PATH=/usr/lib64 firefox http://127.0.0.1:" + GPCC_PORT)

	// Write it to the file.
	file := core.TempDir + "start_gpcc_web.sh"
	err := ExecuteBash(file, StartGPCCWeb)
	if err != nil { return err }

	return nil
}

func StartWLMService() error {

	log.Info("Starting the workload manager for this environment")

	// Stopping all the WLM core.
	var StartWLMArgs []string
	StartWLMArgs = append(StartWLMArgs, "if [ -f " + WLMInstallDir + "/gp-wlm/bin/svc-mgr.sh ]; then ")
	StartWLMArgs = append(StartWLMArgs, WLMInstallDir + "/gp-wlm/bin/svc-mgr.sh --service=all --action=cluster-start")
	StartWLMArgs = append(StartWLMArgs, "fi")

	// Write it to the file.
	file := core.TempDir + "start_wlm.sh"
	err := ExecuteBash(file, StartWLMArgs)
	if err != nil { return err }

	return nil

}

func StopWLMService() error {

	log.Info("Stopping the workload manager services running on the environment")

	// Stopping all the WLM core.
	var StopWLMArgs []string
	StopWLMArgs = append(StopWLMArgs, "if [ -f " + WLMInstallDir + "/gp-wlm/bin/svc-mgr.sh ]; then ")
	StopWLMArgs = append(StopWLMArgs, WLMInstallDir + "/gp-wlm/bin/svc-mgr.sh --service=all --action=cluster-stop")
	StopWLMArgs = append(StopWLMArgs, "fi")

	// Write it to the file.
	file := core.TempDir + "stop_wlm.sh"
	err := ExecuteBash(file, StopWLMArgs)
	if err != nil { return err }

	return nil
}

// Stop all WLM instance if running
func StopAllWLM() error {

	// Stop all WLM instance running on the host
	log.Info("Stopping all WLM instance running on this host")
	var cmdString []string
	cmdScript := "ls " + core.EnvYAML.Install.MasterDataDirectory + "/wlm | while read line; " +
		"do " +
		core.EnvYAML.Install.MasterDataDirectory + "wlm/$line/gp-wlm/bin/svc-mgr.sh --service=all --action=cluster-stop;" +
		"done &>/dev/null"

	// Execute the command
	StopScriptLoc := core.TempDir + "stop_all_wlm.sh"
	cmdString = append(cmdString, cmdScript)
	err := ExecuteBash(StopScriptLoc, cmdString)
	if err != nil { return err }

	return nil

}
