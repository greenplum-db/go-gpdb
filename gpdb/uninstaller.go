package main

import (
	"fmt"
	"strings"
)

// Build uninstall script
func (i *Installation) createUninstallScript() error {
	// Uninstall script location
	uninstallFile := Config.INSTALL.UNINTSALLDIR + "uninstall_" + cmdOptions.Version + "_" + i.Timestamp
	Infof("Creating Uninstall file for this installation at: " + uninstallFile)

	// Query
	queryString := `
select $$ssh $$ || hostname || $$ "ps -ef|grep postgres|grep -v grep|grep $$ ||  port || $$ | awk '{print \$2}'| xargs -n1 /bin/kill -11 &>/dev/null" $$ from gp_segment_configuration 
union
select $$ssh $$ || hostname || $$ "rm -rf /tmp/.s.PGSQL.$$ || port || $$*"$$ from gp_segment_configuration
union
`

	// From GPDB6 onwards there is no filespace, so we cannot know what is the default data directory from database
	// we will get that information from the gpinitsystem struct that we created earlier
	if isThisGPDB6xAndAbove() {
		queryString = queryString + fmt.Sprintf(`select $$ssh $$ || c.hostname || $$ "rm -rf %s*" $$ from gp_segment_configuration c`, i.GPInitSystem.MasterDir+"/"+i.GPInitSystem.ArrayName)
	} else {
		queryString = queryString + `select $$ssh $$ || c.hostname || $$ "rm -rf $$ || f.fselocation || $$"$$ from pg_filespace_entry f, gp_segment_configuration c where c.dbid = f.fsedbid`
	}

	// Execute the query
	cmdOut, err := executeOsCommandOutput("psql", "-p", i.GPInitSystem.MasterPort, "-d", "template1", "-Atc", queryString)
	if err != nil {
		Fatalf("Error in running uninstall command on database, err: %v", err)
	}

	// Create the file
	createFile(uninstallFile)
	writeFile(uninstallFile, []string{
		string(cmdOut),
		"rm -rf " + Config.INSTALL.ENVDIR + "env_" + cmdOptions.Version + "_" + i.Timestamp,
		"rm -rf " + uninstallFile,
	})
	return nil
}

// Uninstall gpcc
func (i *Installation) uninstallGPCCScript() error {
	i.GPCC.UninstallFile = Config.INSTALL.UNINTSALLDIR + fmt.Sprintf("uninstall_gpcc_%s_%s_%s", cmdOptions.Version, cmdOptions.CCVersion, i.Timestamp)
	Infof("Created uninstall script for this version of GPCC Installation: %s", i.GPCC.UninstallFile)
	writeFile(i.GPCC.UninstallFile, []string{
		"source " + i.EnvFile,
		"source " + i.GPCC.GpPerfmonHome + "/gpcc_path.sh",
		"gpcmdr --stop " + i.GPCC.InstanceName + " &>/dev/null",
		"gpcc stop " + i.GPCC.InstanceName + " &>/dev/null",
		"rm -rf " + i.GPCC.GpPerfmonHome + "/instances/" + i.GPCC.InstanceName,
		"gpconfig -c gp_enable_gpperfmon -v off &>/dev/null",
		"echo \"Stopping the database to cleanup any gpperfmon process\"",
		"gpstop -af &>/dev/null",
		"echo \"Starting the database\"",
		"gpstart -a &>/dev/null",
		"cp $MASTER_DATA_DIRECTORY/pg_hba.conf $MASTER_DATA_DIRECTORY/pg_hba.conf." + i.Timestamp + " &>/dev/null",
		"grep -v gpmon $MASTER_DATA_DIRECTORY/pg_hba.conf." + i.Timestamp + " > $MASTER_DATA_DIRECTORY/pg_hba.conf",
		"rm -rf $MASTER_DATA_DIRECTORY/pg_hba.conf." + i.Timestamp + " &>/dev/null",
		"psql -d template1 -p $PGPORT -Atc \"drop database gpperfmon\" &>/dev/null",
		"psql -d postgres -p $PGPORT -Atc \"drop schema gpcc_schema cascade\" &>/dev/null",
		"psql -d template1 -p $PGPORT -Atc \"drop role gpmon\" &>/dev/null",
		"rm -rf $MASTER_DATA_DIRECTORY/gpperfmon/*",
		"cp " + i.EnvFile + " " + i.EnvFile + "." + i.Timestamp,
		"egrep -v \"GPCC_UNINSTALL_LOC|GPCCVersion|GPPERFMONHOME|GPCC_INSTANCE_NAME|GPCCPORT\" " + i.EnvFile + "." + i.Timestamp + " > " + i.EnvFile,
		"rm -rf " + i.EnvFile + "." + i.Timestamp,
		"rm -rf " + i.GPCC.UninstallFile,
	})
	return nil
}

// Uninstall using gpdeletesystem
func removeEnvGpDeleteSystem(envFile string) error {
	Info("Starting the database if stopped to run the gpdeletesystem on the environment")

	// Start the database if not started
	startDBifNotStarted(envFile)

	Infof("Calling gpdeletesystem to remove the environment: %s", envFile)

	// Write it to the file.
	file := Config.CORE.TEMPDIR + "run_deletesystem.sh"
	createFile(file)
	writeFile(file, []string{
		"source " + envFile,
		"gpdeletesystem -d $MASTER_DATA_DIRECTORY -f << EOF",
		"y",
		"y",
		"EOF",
	})
	_, err := executeOsCommandOutput("/bin/sh", file)
	if err != nil {
		deleteFile(file)
		return err
	}
	deleteFile(file)
	return nil
}

// Uninstall GPCC
func removeGPCC(envFile string) {
	Infof("Uninstalling the version of command center that is currently installed on this environment.")
	gpccEnvFile := environment(envFile).GpccUninstallLoc
	if !IsValueEmpty(gpccEnvFile) {
		executeOsCommand("/bin/sh", gpccEnvFile)
	}
}

// Uninstall using manual method
func removeEnvManually(version, timestamp string) {
	uninstallScript := Config.INSTALL.UNINTSALLDIR + fmt.Sprintf("uninstall_%s_%s", version, timestamp)
	Infof("Cleaning up the extra files using the uninstall script: %s", uninstallScript)
	exists, err := doesFileOrDirExists(uninstallScript)
	if err != nil {
		Fatalf("error when trying to find the uninstaller file \"%s\", err: %v", uninstallScript, err)
	}
	if exists {
		executeOsCommandOutput("/bin/sh", uninstallScript)
	} else {
		Fatalf("Unable to find the uninstaller file \"%s\"", uninstallScript)
	}
}

// Main Remove method
func remove() {
	Infof("Starting program to uninstall the version: %s", cmdOptions.Version)

	// Check if the envfile for that version exists
	chosenEnvFile := installedEnvFiles(fmt.Sprintf("*%s*", cmdOptions.Version), "choose", true)

	// If we receive none, then display the error to user
	var timestamp, version string
	if IsValueEmpty(chosenEnvFile) {
		Fatalf("Cannot find any environment with the version: %s", cmdOptions.Version)
	} else { // Else store the value
		timestamp = strings.Split(chosenEnvFile, "_")[2]
		version = strings.Split(chosenEnvFile, "_")[1]
	}
	Infof("The choosen enviornment file to remove is: %s ", chosenEnvFile)
	Info("Uninstalling the environment")

	// If there is failure in gpstart, user can use force to force manual uninstallation
	if !cmdOptions.Force {
		err := removeEnvGpDeleteSystem(chosenEnvFile)
		if err != nil {
			Warnf("Failed to uninstall using gpdeletesystem, trying manual method..")
		}
	} else {
		Infof("Forcing uninstall of the environment: %s", chosenEnvFile)
	}

	// Uninstall GPPC
	removeGPCC(chosenEnvFile)

	// Run this to cleanup the file created by go-gpdb
	removeEnvManually(version, timestamp)

	Infof("Uninstallation of environment \"%s\" was a success", chosenEnvFile)
	Info("exiting ....")
}
