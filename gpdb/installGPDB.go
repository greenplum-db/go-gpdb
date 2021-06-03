package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Precheck before installing GPDB
func (i *Installation) preGPDBChecks() {
	Infof("Running precheck before installing the gpdb version: %s", cmdOptions.Version)
	// Validate the master & segment exists and is readable
	dirValidator()

	// Error out if there is multiple installation of GPDB
	allEnv := ListEnvironmentsInstalled("*")
	if len(allEnv) > 0 && getSystemInfoAndCheckIfItsUbuntu() {
		Warnf("Currently having multiple installation of GPDB isn't supported on Ubuntu")
		Warnf("Please use \"%s remove -v <version>\", to remove the old version " +
			"or to make room for this installation", programName)
		Fatalf("Already GPDB is installed on this box...")
	}

	// Check if there is already a version of GPDB installed
	installedEnvFiles(fmt.Sprintf("*%s*", cmdOptions.Version), "confirm", true)
}

func (i *Installation) installGPDBProduct() {
	Infof("Running Installation of gpdb version: %s", cmdOptions.Version)
	// Start the installation procedure
	i.installProduct()

	// Check ssh to host is working and enable password less login
	i.setUpHost()

	// Build & Execute the gpinitsystem configuration
	i.Timestamp = time.Now().Format("20060102150405")
	i.buildGpInitSystem()
}

func (i *Installation) postGPDBInstall() {
	Infof("Running post installation steps of the gpdb version: %s", cmdOptions.Version)

	// Store the last used port for future use
	i.savePort()

	// Create EnvFile
	i.createEnvFile()

	// Initialize standby master
	if i.StandbyHostAvailable && cmdOptions.Standby {
		i.activateStandby()
	} else if cmdOptions.Standby {
		Errorf("Cannot activate standby, please activate the standby manually")
	}

	// Create uninstall script
	i.createUninstallScript()

	// Installation complete, print on the screen the env file to source and cleanup temp files
	displayEnvFileToSource(i.EnvFile)
	defer deleteFile(i.WorkingHostFileLocation)
	defer deleteFile(i.SegmentHostLocation)
}

// Check if the directory provided is readable and writeable
func dirValidator() error {
	Debugf("Checking for the existences of master and segment directory")

	// Check if the master & segment location is readable and writable
	masterDirExists, err := doesFileOrDirExists(Config.INSTALL.MASTERDATADIRECTORY)
	if err != nil {
		Fatalf("Cannot get the information of directory %s, err: %v", Config.INSTALL.MASTERDATADIRECTORY, err)
	}

	segmentDirExists, err := doesFileOrDirExists(Config.INSTALL.SEGMENTDATADIRECTORY)
	if err != nil {
		Fatalf("Cannot get the information of directory %s, err: %v", Config.INSTALL.SEGMENTDATADIRECTORY, err)
	}

	// if the file doesn't exists then let try creating it ...
	if !masterDirExists || !segmentDirExists {
		CreateDir(Config.INSTALL.MASTERDATADIRECTORY)
		CreateDir(Config.INSTALL.SEGMENTDATADIRECTORY)
	}
	return nil
}

// Installing the greenplum binaries
func (i *Installation) installProduct() {
	// Installing the binaries.
	binFile, isItBinaryFile := getBinaryFile(cmdOptions.Version)

	// Is this a binary or rpm file
	if getSystemInfoAndCheckIfItsUbuntu() {
		// Located a deb file
		i.installDebFile(binFile)
	} else if isItBinaryFile {
		// Located a binary file
		i.installBinaryFile(binFile)
		i.BinaryOrRpm = "binary"
	} else {
		// Located a rpm file
		i.installRpmFile(binFile)
		i.BinaryOrRpm = "rpm"
		i.RpmFileName = binFile
	}

	// Now source the environment if installation went fine
	err := sourceGPDBPath(i.BinaryInstallationLocation)
	if err != nil {
		Fatalf("Failed to set the environment variable while sourcing the file")
	}

}

// We have located a binary file
func (i *Installation) installBinaryFile(binFile string) {
	// Location and name of the binaries
	Infof("Using the bin file to install the GPDB Product: %s", binFile)
	i.BinaryInstallationLocation = fmt.Sprintf("/usr/local/%[1]s/greenplum-db-%[1]s", cmdOptions.Version)

	// Execute the command to install the binaries
	var scriptOption = []string{"yes", i.BinaryInstallationLocation, "yes", "yes"}
	err := executeBinaries(binFile, "install_software.sh", scriptOption)
	if err != nil {
		Fatalf("Failed in installing the binaries, err: %v", err)
	}
}

// Got a rpm file
func (i *Installation) installRpmFile(binFile string) {
	// Pre rpm setup: Install dependencies
	i.preRpmInstallationSetup()

	// Execute the command to install the rpm
	Infof("Using the rpm file to install the GPDB Product: %s, this might take several minutes....", binFile)
	baseDir := "/usr/local/" + cmdOptions.Version
	CreateDir(baseDir)
	executeOsCommand("sudo", "rpm", "--install", binFile, "--prefix=" + baseDir ,"--force")

	// Post Rpm setup: Change ownership
	Infof("Changing the ownership of the folder %s", baseDir)
	executeOsCommand("sudo", "chown", "-R", "gpadmin:gpadmin", baseDir + "/greenplum*")

	// Find the directory where the rpm was installed
	i.BinaryInstallationLocation = locateGreenplumInstallationDirectory(baseDir)
}

// Setup the infrastructure for installing the rpm binary
func (i *Installation) preRpmInstallationSetup() {
	// Removing previous installed rpm's
	Infof("Removing previous greenplum rpm's from rpm database")
	cleanupPreviousRpms := Config.CORE.TEMPDIR + "cleanup_previous_rpms.sh"
	options := []string{
		"sudo rpm -e --justdb $(rpm -qa | grep green) ",
	}
	generateBashFileAndExecuteTheBashFile(cleanupPreviousRpms, "/bin/sh", options)

	// Install all dependencies
	Infof("Installing all the dependencies")
	var dependencies = []string{"apr-util", "bash", "bzip2", "curl", "krb5", "libcurl", "libevent", "libxml2",
		"libyaml", "zlib", "openldap", "openssh", "openssl", "openssl-libs", "perl", "readline", "rsync", "R",
		"sed", "tar", "zip", "krb5-devel"}

	for _, i := range dependencies {
		Debugf("Installing pre rpm installation package: %s", i)
		executeOsCommand("sudo", "yum",  "install", i, "-y", "-q")
	}
}

// Check if the provided hostnames are valid
func (i *Installation) setUpHost() {
	Infof("Setting up & Checking if the host is reachable")

	// Get the hostname of the master
	i.GPInitSystem.MasterHostname, _ = os.Hostname()
	if i.GPInitSystem.MasterHostname == "" {
		Fatalf("The environment variable 'HOSTNAME' for master host is not set")
	}

	// Check is host is reachable
	i.isHostReachable()

	// Enable passwordless login
	i.executeGpsshExkey()

	// If this is a multi installation then install the software on all the segment host
	if i.SingleORMulti == "multi" {
		if i.BinaryOrRpm == "rpm" {
			i.rpmInstallOnAllSegmentHost()
		} else {
			i.runSegInstall()
		}
	}
}

// Check if all the host are working from the hostfile
func (i *Installation) isHostReachable() {
	// Get all the hostname from the hostfile and check if the host if reachable
	var saveReachableHost []string
	var segmentHost []string
	for _, host := range strings.Split(string(readFile(i.HostFileLocation)), "\n") {
		if host != "" && checkHostReachability(fmt.Sprintf("%s:22", host), true) {
			Debugf("The host %s is reachable", host)
			// Save all the host that is reachable
			saveReachableHost = append(saveReachableHost, host)
			// Save all the segment host list
			if host != i.GPInitSystem.MasterHostname && !strings.HasSuffix(host, "-s") {
				segmentHost = append(segmentHost, host)
			}
			// If we detect a standby host then record that we have one
			if strings.HasSuffix(host, "-s") {
				i.StandbyHostAvailable = true
			}
		}
	}

	// Save the reachable host
	Debugf("Total Host reachable is: %d", len(saveReachableHost))
	if len(saveReachableHost) == 0 { // If none of the host is reachable
		Fatalf("No hosts are reachable from the hostfile %s, check the host", i.HostFileLocation)
	} else if len(saveReachableHost) == 1 && saveReachableHost[0] == i.GPInitSystem.MasterHostname { // if only one host is reachable and its master then its a single install
		i.SingleORMulti = "single"
	} else if len(saveReachableHost) == 1 && saveReachableHost[0] != i.GPInitSystem.MasterHostname { // if only one host and its not master, then we can't continue
		Fatalf("Master hosts are not reachable from the hostfile %s, check the host", i.HostFileLocation)
	}

	// Check if we have equal number of hosts
	if i.SingleORMulti == "multi" && len(segmentHost)%2 == 1 { // if its multi and we have odd number of host, then we can't create mirror
		Fatalf("There is odd number of segment host, installation cannot continue")
	} else if i.SingleORMulti == "multi" && len(segmentHost) == 0 { // if multi and no segment host then we can't continue
		Fatalf("No segment host found, cannot continue")
	}

	// Create hostfile based on the what we have collected so far
	i.hostfileCreator(saveReachableHost, segmentHost)
}

// Generate hostname file based on installation.
func (i *Installation) hostfileCreator(saveReachableHost, segmentHost []string) {
	Debug("Creating hostfile")
	// Save the working host on a separate file
	i.WorkingHostFileLocation = Config.CORE.TEMPDIR + "hostfile"
	deleteFile(i.WorkingHostFileLocation)
	writeFile(i.WorkingHostFileLocation, saveReachableHost)

	// Save segment host file
	i.SegmentHostLocation = Config.CORE.TEMPDIR + "hostfile_segment"
	deleteFile(i.SegmentHostLocation)
	writeFile(i.SegmentHostLocation, segmentHost)

	// Save gpseginstall host file
	i.SegInstallHostLocation = Config.CORE.TEMPDIR + "hostfile_seginstall"
	deleteFile(i.SegInstallHostLocation)
	if i.StandbyHostAvailable {
		segmentHost = append(segmentHost, buildStandbyHostName(i.GPInitSystem.MasterHostname))
	}
	writeFile(i.SegInstallHostLocation, segmentHost)
}

// Run keyless access to the server
func (i *Installation) executeGpsshExkey() {
	Infof("Running gpssh-exkeys to enable keyless access on this server")

	if isThisGPDB6xAndAbove() {
		// GPDB 6 has changed the way we generate ssh key, so we now use the newer way
		// to initialize the ssh key and then use gpssh-exkeys to exchange to all host
		i.generatePasswordLessKeys()
	}
	// Exchange the keys b/w host
	executeOsCommand(fmt.Sprintf("%s/bin/gpssh-exkeys", os.Getenv("GPHOME")), "-f", i.WorkingHostFileLocation)
}

// Create password ssh key
func (i *Installation) generatePasswordLessKeys() {
	Debug("Creating passwordless ssh keys")

	// Create the ssh-keygen if there doesn't exists one
	pubFile := fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME"))
	sshKeyGenFile := Config.CORE.TEMPDIR + "init_ssh_keygen.sh"
	if !fileExists(pubFile) {
		generateBashFileAndExecuteTheBashFile(sshKeyGenFile, "/bin/sh", []string{
			fmt.Sprintf("ssh-keygen -b 2048 -t rsa -q -N '' -f %s", pubFile),
		})
	}

	// Now run ssh-copy-id for all the host
	// Check if the sshpass is installed on the OS
	if isCommandAvailable("sshpass") {
		i.setupSshCopyID()
	}
}

// Now generate ssh copy for all the segment host
func (i *Installation) setupSshCopyID() {
	Infof("Setting up the ssh Copy ID for all the host(s)")
	file, err := os.Open(i.HostFileLocation)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		host := removeBlanks(scanner.Text())
		Debugf("Setting up the ssh Copy ID for all the host %s", host)
		sshCopyIDFileName := Config.CORE.TEMPDIR + "ssh_copy_id.sh"
		generateBashFileAndExecuteTheBashFile(sshCopyIDFileName, "/bin/sh", []string{
			fmt.Sprintf("SSHPASS=%s sshpass -e ssh-copy-id %s -o StrictHostKeyChecking=no &> /dev/null", Config.INSTALL.MASTERPASS, host),
		})
	}

	if err := scanner.Err(); err != nil {
		Fatalf("Error in reading the file %s, err: %v", i.HostFileLocation, err)
	}
}

// Run the gpseginstall on all host
func (i *Installation) runSegInstall() {
	Infof("Running seg install to install the software on all the host")
	executeOsCommand(fmt.Sprintf("%s/bin/gpseginstall", os.Getenv("GPHOME")), "-f", i.SegInstallHostLocation)
	i.createSoftLink()
}

// Run yum install on all the segment host to install the GPDB software
func (i *Installation) rpmInstallOnAllSegmentHost() {
	Infof("Running rpm install to install the software on all the segment host, this might take several minutes....")
	gpsshExecutable := fmt.Sprintf("%s/bin/gpssh", os.Getenv("GPHOME"))
	gpscpExecutable := fmt.Sprintf("%s/bin/gpscp", os.Getenv("GPHOME"))
	destinationFileName := fmt.Sprintf("/tmp/gpdb-%s.rpm", cmdOptions.Version)

	// Copy the rpm to all the segment host
	gpscpFilename := Config.CORE.TEMPDIR + "gpscp_gpdb_rpm.sh"
	generateBashFileAndExecuteTheBashFile(gpscpFilename, "/bin/sh", []string{
		fmt.Sprintf("%s -f %s %s =:%s", gpscpExecutable, i.SegInstallHostLocation, i.RpmFileName, destinationFileName),
	})

	// Now run the yum install on all the segments
	gpsshFilename := Config.CORE.TEMPDIR + "gpssh_install_gpdb_rpm.sh"
	generateBashFileAndExecuteTheBashFile(gpsshFilename, "/bin/sh", []string{
		fmt.Sprintf("%s -f %s \"sudo yum install -y %s\" &> /dev/null", gpsshExecutable, i.SegInstallHostLocation, destinationFileName),
	})
}

// Install deb related GPDB software
func (i *Installation) installDebFile(debFile string) {
	Infof("Using deb file %s to install GPDB, this might take several minutes ....", debFile)
	baseDir := "/usr/local/"
	executeOsCommand("sudo", "apt", "install", debFile, "-y")

	// Find the directory where the rpm was installed
	i.BinaryInstallationLocation = locateGreenplumInstallationDirectory(baseDir)
}

// On the newer version of the GPDB they need the soft link for
// database to work
func (i *Installation) createSoftLink() {
	Debugf("Creating the softlink for the binaries on all the host")
	contents := readFile(i.SegInstallHostLocation)
	for _, v := range strings.Split(removeBlanks(string(contents)), "\n") {
		softLinkFile := Config.CORE.TEMPDIR + "soft_link.sh"
		generateBashFileAndExecuteTheBashFile(softLinkFile, "/bin/sh", []string{
			fmt.Sprintf("ssh %s \"rm -rf /usr/local/greenplum-db\"", v),
			fmt.Sprintf("ssh %s \"ln -s %s /usr/local/greenplum-db\"", v, i.BinaryInstallationLocation),
		})
	}
}

// Activate the master standby if requested.
func (i *Installation) activateStandby() {
	Infof("Activating the master standby for this installation")
	standbyHostLoc := Config.CORE.TEMPDIR + "activate_standby.sh"
	// We use remove last 3 line rather than replace function to avoid situation where user have
	// created a host with the name say host-standby and then we attach -s with it like host-standby-s
	// replace -s would replace -s at two places and thus causing error , so we just worry about the last character
	generateBashFileAndExecuteTheBashFile(standbyHostLoc, "/bin/sh", []string{
		fmt.Sprintf("source %s", i.EnvFile),
		fmt.Sprintf("gpinitstandby -s %s -a", buildStandbyHostName(i.GPInitSystem.MasterHostname)),
	})
}

// Build standby hostname
func buildStandbyHostName(masterHostname string) string {
	return masterHostname[0:len(masterHostname)-2] + "-s"
}

// Source greenplum path
func sourceGPDBPath(binLoc string) error {
	Debugf("Sourcing the greenplum binary location")

	// Setting up greenplum path
	err := os.Setenv("GPHOME", binLoc)
	if err != nil {
		return err
	}
	err = os.Setenv("PYTHONPATH", os.Getenv("GPHOME")+"/lib/python")
	if err != nil {
		return err
	}
	err = os.Setenv("PYTHONHOME", os.Getenv("GPHOME")+"/ext/python")
	if err != nil {
		return err
	}
	err = os.Setenv("PATH", os.Getenv("GPHOME")+"/bin:"+os.Getenv("PYTHONHOME")+"/bin:"+os.Getenv("PATH"))
	if err != nil {
		return err
	}
	err = os.Setenv("LD_LIBRARY_PATH", os.Getenv("GPHOME")+"/lib:"+os.Getenv("PYTHONHOME")+"/lib:"+os.Getenv("LD_LIBRARY_PATH"))
	if err != nil {
		return err
	}
	err = os.Setenv("OPENSSL_CONF", os.Getenv("GPHOME")+"/etc/openssl.cnf")
	if err != nil {
		return err
	}
	return nil
}
