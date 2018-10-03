package install

import (
	"strconv"
	"os/exec"
	"os"
	"strings"
	"github.com/ielizaga/piv-go-gpdb/core"
)


func BuildGpInitSystemConfig(t string) error {

	log.Info("Building gpinitsystem config file for this installation")

	// Set the values of the below parameters
	GpInitSystemConfig.ArrayName = "gp_" + core.RequestedInstallVersion + "_" + t
	GpInitSystemConfig.SegPrefix = "gp_" + core.RequestedInstallVersion + "_" + t
	GpInitSystemConfig.DBName = DBName
	GpInitSystemConfig.MasterDir = strings.TrimSuffix(core.EnvYAML.Install.MasterDataDirectory, "/")
	core.EnvYAML.Install.SegmentDataDirectory = strings.TrimSuffix(core.EnvYAML.Install.SegmentDataDirectory, "/")

	// Store the hostname of the hostfile
	MasterHostFileList = core.TempDir + MachineFileList
	_ = core.DeleteFile(MasterHostFileList)
	err := core.CreateFile(MasterHostFileList)
	if err != nil { return err }
	var host []string
	host = append(host, GpInitSystemConfig.MasterHostName)
	err = core.WriteFile(MasterHostFileList, host)
	if err != nil { return err }

	// Check if we have the last used port base file
	log.Info("Obtaining ports to be set for PORT_BASE")
	PORTBASE , _ := DoWeHavePortBase(core.FutureRefDir, PortBaseFileName, "PORT_BASE")
	if PORTBASE != "" {
		PORTBASE = string(PORTBASE)[11:]
	} else {
		log.Warning("Didn't find PORT_BASE in the file, setting it to default value: " + strconv.Itoa(PORT_BASE))
		PORTBASE = strconv.Itoa(PORT_BASE)
	}

	// Check if the port is available or not
	pb, err := strconv.Atoi(PORTBASE)
	if err != nil { return err }
	pb, err = IsPortUsed(pb, core.EnvYAML.Install.TotalSegments)
	if err != nil { return err }
	log.Info("Setting the PORT_BASE has: "+ strconv.Itoa(pb))
	GpInitSystemConfig.PortBase = pb

	// Check if we have the last used master port file
	log.Info("Obtaining ports to be set for MASTER_PORT")
	MASTERPORT , _ := DoWeHavePortBase(core.FutureRefDir, PortBaseFileName, "MASTER_PORT")
	if MASTERPORT != "" {
		MASTERPORT = string(MASTERPORT)[13:]
	} else {
		log.Warning("Didn't find MASTER_PORT in the file, setting it to default value: " + strconv.Itoa(MASTER_PORT))
		MASTERPORT = strconv.Itoa(MASTER_PORT)
	}

	// MASTER PORT
	mp, err := strconv.Atoi(MASTERPORT)
	if err != nil { return err }
	mp, err = IsPortUsed(mp, 1)
	if err != nil { return err }
	log.Info("Setting the MASTER_PORT has: "+ strconv.Itoa(mp))
	GpInitSystemConfig.MasterPort = mp

	// Build gpinitsystem config file
	GpInitSystemConfigDir = core.TempDir + GpInitSystemConfigFile + "_" + core.RequestedInstallVersion + "_" + t
	log.Info("Creating the gpinitsystem config file at: " + GpInitSystemConfigDir)
	_ = core.DeleteFile(GpInitSystemConfigDir)
	err = core.CreateFile(GpInitSystemConfigDir)
	if err != nil { return err }

	// Write the below content to config file
	var ToWrite []string
	var primaryDir string
	ToWrite = append(ToWrite,"ARRAY_NAME=" + GpInitSystemConfig.ArrayName )
	ToWrite = append(ToWrite,"MACHINE_LIST_FILE=" + MasterHostFileList)
	ToWrite = append(ToWrite,"SEG_PREFIX=" + GpInitSystemConfig.SegPrefix )
	ToWrite = append(ToWrite,"MASTER_HOSTNAME=" + GpInitSystemConfig.MasterHostName )
	ToWrite = append(ToWrite,"MASTER_DIRECTORY=" + GpInitSystemConfig.MasterDir )
	ToWrite = append(ToWrite,"PORT_BASE=" + strconv.Itoa(GpInitSystemConfig.PortBase) )
	ToWrite = append(ToWrite,"MASTER_PORT=" + strconv.Itoa(GpInitSystemConfig.MasterPort) )
	ToWrite = append(ToWrite,"DATABASE_NAME=" + GpInitSystemConfig.DBName)
	for i:=1; i<=core.EnvYAML.Install.TotalSegments; i++ {
		primaryDir = primaryDir + " "+ core.EnvYAML.Install.SegmentDataDirectory
	}
	ToWrite = append(ToWrite,"declare -a DATA_DIRECTORY=(" + primaryDir + " )")
	err = core.WriteFile(GpInitSystemConfigDir, ToWrite)
	if err != nil { return err }

	return nil
}


func ExecuteGpInitSystem() error {

	log.Info("Executing gpinitsystem to install GPDB software")

	// Initalize the command
	cmdOut := exec.Command("gpinitsystem", "-c", GpInitSystemConfigDir , "-h", MasterHostFileList, "-a")

	// Attach the os output from the screen
	cmdOut.Stdout = os.Stdout
	cmdOut.Stderr = os.Stderr

	// Start the program
	err := cmdOut.Start()
	if err != nil { return err }

	// wait till it ends
	err = cmdOut.Wait()
	if err != nil { return err }

	return nil
}
