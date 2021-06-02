package main

import (
	"bytes"
	"fmt"
)

func stopAllDb() {

	Infof("Stopping all the database running on this host, to free up semaphore for this installation")

	// Can't seems to find a simple way to stop all database, so we will built the below
	// simple shell script to execute the stop database command
	cleanupScript := "ps -ef | grep postgres | grep master | grep -v logger | while read line; " +
		"do " +
		"GPHOME=`echo $line|awk '{print $8}'| rev | cut -d'/' -f3- | rev`;" +
		"export MASTER_DATA_DIRECTORY=`echo $line|awk '{print $10}'`;" +
		"export PGPORT=`echo $line|awk '{print $12}'`;" +
		"export PGDATABASE=template1;" +
		"source $GPHOME/greenplum_path.sh;" +
		"gpstop -af;" +
		"done"
	StopScriptLoc := Config.CORE.TEMPDIR + "stop_all_db.sh"
	generateBashFileAndExecuteTheBashFile(StopScriptLoc, "/bin/sh", []string{
		cleanupScript,
	})
	cleanupGpccProcess()
	areAllProcessDown()
}

// Cleanup all the gpcc process
func cleanupGpccProcess() {
	Infof("Cleaning up gpcc process is found any")
	cleanupGPCCLoc := Config.CORE.TEMPDIR + "clean_all_gpcc.sh"
	generateBashFileAndExecuteTheBashFile(cleanupGPCCLoc, "/bin/sh", []string{
		"ps -ef | egrep \"gpsmon|gpmon|gpmonws|lighttpd\" | grep -v grep | awk '{print $2}' | xargs -n1 /bin/kill -11 >/dev/null 2>&1",
	})
}

// Check if the process are all down
func areAllProcessDown() {
	Debugf("Checking if the all the database process are down")
	// Send a warning message if the process is not completely stopped.
	cmdOut, _ := executeOsCommandOutput("pgrep", "postgres")
	var EmptyBytes []byte
	if !bytes.Equal(cmdOut, EmptyBytes) {
		Warn("Can't stop all postgres process, seems like some are left behind")
	} else {
		Info("All Postgres process are stopped on this server")
	}
}

// Check if the database is healthy
func isDbHealthy(sourcePath, port string) bool {
	Debugf("Checking if the database is healthy")
	// Query string
	queryString := "select 1"
	if sourcePath == "" {
		_, err := executeOsCommandOutput("psql", "-p", port, "-d", "template1", "-Atc", queryString)
		if err != nil {
			Fatalf("Error in checking database health, err: %v", err)
		}
	} else {
		p := Config.CORE.TEMPDIR + "db_health.sh"
		createFile(p)
		writeFile(p, []string{
			"source " + sourcePath,
			fmt.Sprintf("psql -d template1 -p %s -Atc \"%s\"", environment(sourcePath).PgPort, queryString),
		})
		_, err := executeOsCommandOutput("/bin/sh", p) // if command execute its healthy
		if err != nil {
			deleteFile(p)
			return false
		} else {
			deleteFile(p)
			return true
		}
	}

	return true
}

// Start the database if not started
func startDBifNotStarted(envFile string) {
	// is the database running , then return
	if isDbHealthy(envFile, "") { // Database is started and running
		Debugf("Database seems to be running, contining...")
	} else { // database is not running, lets start it up

		Warnf("Database is not started, attempting to start the database...")
		// Stop all database is not stopped unless asked not to stop it
		if !cmdOptions.Stop {
			stopAllDb()
		}

		// Start the database of concern
		startDB(envFile)

		// Check again if the database is healthy
		if !isDbHealthy(envFile, "") {
			Fatalf("Can't seems to start the database in the environment file \"%s\"exiting...", envFile)
		}
	}
}

// Start database
func startDB(envFile string) {
	Infof("Attempting to start the database as per the environment file: %s", envFile)

	// BashScript
	startFile := Config.CORE.TEMPDIR + "start.sh"
	generateBashFileAndExecuteTheBashFile(startFile, "/bin/bash", []string{
		"source " + envFile,
		"gpstart -a",
	})

	// Check if database is healthy
	if !isDbHealthy(envFile, "") { // Check if the database is healthy after start
		Fatalf("Can't seems to start the database in the environment file \"%s\" exiting...", envFile)
	}

	// Start Command Center WEB UI if available on this environment
	if isCommandCenterInstalled(envFile) {
		startGPCC(envFile)
	}
}

// Stop Database
func stopDB(envFile string) {
	Infof("Attempting to stop the database as per the environment file: %s", envFile)
	stopFile := Config.CORE.TEMPDIR + "stop.sh"
	generateBashFileAndExecuteTheBashFile(stopFile, "/bin/sh", []string{
		"source " + envFile,
		"gpstop -af",
	})

	// Stop Command Center WEB UI if available on this environment
	if isCommandCenterInstalled(envFile) {
		stopGPCC(envFile)
	}

	// Check if all the process are down
	areAllProcessDown()
}

// Start GPCC
func startGPCC(envFile string) {
	// start file
	Infof("Trying to start the GPCC Instance")
	startGPCCFile := Config.CORE.TEMPDIR + "gpcc_start.sh"
	if isThis4x() {
		generateBashFileAndExecuteTheBashFile(startGPCCFile, "/bin/sh", []string{
			fmt.Sprintf("source %s", envFile),
			"source $GPPERFMONHOME/gpcc_path.sh",
			"gpcc start ${GPCC_INSTANCE_NAME}",
		})
	} else {
		generateBashFileAndExecuteTheBashFile(startGPCCFile, "/bin/sh", []string{
			fmt.Sprintf("source %s", envFile),
			"source $GPPERFMONHOME/gpcc_path.sh",
			"gpcmdr --start ${GPCC_INSTANCE_NAME} << EOF",
			"y",
			"EOF",
		})
	}
}

// Stop GPPC
func stopGPCC(envFile string) {
	// stop file
	Infof("Trying to stop the GPCC Instance")
	stopGPCCFile := Config.CORE.TEMPDIR + "gpcc_stop.sh"
	if isThis4x() {
		generateBashFileAndExecuteTheBashFile(stopGPCCFile, "/bin/sh", []string{
			fmt.Sprintf("source %s", envFile),
			"source $GPPERFMONHOME/gpcc_path.sh",
			"gpcc stop ${GPCC_INSTANCE_NAME}",
		})
	} else {
		generateBashFileAndExecuteTheBashFile(stopGPCCFile, "/bin/sh", []string{
			fmt.Sprintf("source %s", envFile),
			"source $GPPERFMONHOME/gpcc_path.sh",
			"gpcmdr --stop ${GPCC_INSTANCE_NAME}",
		})
	}
	cleanupGpccProcess()
}

// Is command center installed on this version
func isCommandCenterInstalled(envFile string) bool {
	if IsValueEmpty(environment(envFile).GpccInstanceName) {
		return false
	} else {
		return true
	}
}
