package main

import (
	"fmt"
	"os"
)

type Installation struct {
	HostFileLocation           string
	WorkingHostFileLocation    string
	SegmentHostLocation        string
	SegInstallHostLocation     string
	BinaryInstallationLocation string
	BinaryOrRpm                string
	RpmFileName                string
	GpInitSystemConfigLocation string
	PortFileName               string
	EnvFile                    string
	Timestamp                  string
	SingleORMulti              string
	StandbyHostAvailable       bool
	GPInitSystem               GPInitSystemConfig
	GPCC                       GPCCConfig
}

type GPInitSystemConfig struct {
	MasterHostname        string
	ArrayName             string
	SegPrefix             string
	DBName                string
	MasterDir             string
	SegmentDir            string
	MirrorDir             string
	MasterPort            string
	SegmentPort           string
	MirrorPort            string
	ReplicationPort       string
	MirrorReplicationPort string
}

type GPCCConfig struct {
	InstanceName  string
	InstancePort  string
	GpPerfmonHome string
	WebSocketPort string
	GPCCBinaryLoc string
	UninstallFile string
}

const (
	defaultMasterPort          = 3000
	defaultPrimaryPort         = 30000
	defaultMirrorPort          = 35000
	defaultReplicatePort       = 40000
	defaultMirrorReplicatePort = 45000
	defaultGpccPort            = 28000
	defaultWebSocket           = 8899
)

func install() {
	Infof("Running the installation for the product: %s", cmdOptions.Product)

	// Initialize the struct
	i := new(Installation)

	// Checking if this is a single install VM or Mutli node VM
	var singleORmulti string
	noSegments := strToInt(os.Getenv("GPDB_SEGMENTS"))
	if noSegments > 0 {
		i.SingleORMulti = "multi"
	} else {
		i.SingleORMulti = "single"
	}
	Debugf("Is this single or multi node installation: %s", i.SingleORMulti)

	// Currently we will only support single VM installation of ubuntu, if there is a need we can
	// expand the tool for multi segment on different hosts
	if i.SingleORMulti == "multi" && getSystemInfoAndCheckIfItsUbuntu() {
		Fatalf("The installation of GPDB via %s tool, only support single VM at the moment", programName)
	}

	// Get or Generate the hostname file
	i.generateHostFile()

	// Run the installation
	if cmdOptions.Product == "gpdb" { // Install GPDB
		i.PortFileName = "gpdb_ports.save"
		i.installGPDB(singleORmulti)
	} else { // its a GPCC installation
		i.PortFileName = "gpcc_ports.save"
		i.installGPCC()
	}
}

// Install GPDB
func (i *Installation) installGPDB(singleOrMutli string) {
	Infof("Starting the program to install GPDB version: %s", cmdOptions.Version)

	// Pre check
	i.preGPDBChecks()

	// Install the product
	i.installGPDBProduct()

	// run post steps
	i.postGPDBInstall()

	Infof("Installation of GPDB with version %s is complete", cmdOptions.Version)
	Infof("exiting ....")
}

// Install GPCC
func (i *Installation) installGPCC() {
	Infof("Starting the program to install GPCC version \"%s\" on the GPDB Version \"%s\"", cmdOptions.CCVersion, cmdOptions.Version)

	// Run prechecks
	i.preGPCCChecks()

	// Install GPCC product
	i.installGPCCProduct()

	// run post steps
	i.postGPCCInstall()

	Infof("Installation of GPCC with version %s on GPDB with version %s is complete", cmdOptions.CCVersion, cmdOptions.Version)
	Infof("exiting ....")
}

// Get and generate host file if doesn't exists the hostname
func (i *Installation) generateHostFile() {
	i.HostFileLocation = fmt.Sprintf("%s/hostfile", os.Getenv("HOME"))
	Debugf("Generating the hostfile at %s", i.HostFileLocation)

	exists, _ := doesFileOrDirExists(i.HostFileLocation)
	if !exists {
		Infof("Host file doesn't exists, creating one: %s", i.HostFileLocation)
		// Read the contents from the /etc/hosts and generate a hostfile.
		hosts := contentExtractor(readFile("/etc/hosts"), "{if (NR!=1) {print $2}}", []string{})
		s := removeBlanks(hosts.String())
		// write to the file
		writeFile(i.HostFileLocation, []string{s})
	} else {
		Infof("Found host file at location: %s", i.HostFileLocation)
	}
}

// If no version is provided, prompt the list of the downloaded product to choose from
func chooseDownloadedProducts() string {
	totalProducts := displayDownloadedProducts("listandChoose")
	choice := PromptChoice(len(totalProducts))
	choosenProduct := totalProducts[choice-1]
	return extractVersionNumber(choosenProduct)
}
