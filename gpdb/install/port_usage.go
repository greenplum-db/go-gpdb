package install

import (
	"strconv"
	"net"
	"os/exec"
	"strings"
	"github.com/ielizaga/piv-go-gpdb/core"
)

// Check if the port is used, if yes then what is the next sequence program can use
func IsPortUsed(port int, iteration int) (int, error) {

	// Storing the base port
	BASE := port

	// Iterate through iteration to find how many port do we need
	for i:=1; i<=iteration; i++ {

		// Error out if the ports are not in the format needed
		_, err := strconv.ParseUint(strconv.Itoa(port), 10, 16)
		if err != nil { return 0, err }

		// Check if the port is available, if not find the next sequence
		ln, err := net.Listen("tcp", ":" + strconv.Itoa(port))
		if err != nil {
			log.Warning("PORT \"" + strconv.Itoa(port) + "\" is unavailable, finding the next sequence" )
			BASE =  BASE + iteration
			return IsPortUsed(BASE, iteration)
		}

		// Close listening to the port
		err = ln.Close()
		if err != nil {
			if err != nil { return 0, err }
		}

		// Iterate
		port =  port + i
	}

	// Return the collected information
	return BASE, nil

}

// Check if we have the last used ports
func DoWeHavePortBase(file string, name string, which_port string) (string, error) {

	portBaseFile := file + name
	return_code, err := core.DoesFileOrDirExists(portBaseFile)
	if err != nil { return "", err }
	if return_code {
		log.Info("Found port file: \"" + portBaseFile + "\"")
		cmdOut, _ := exec.Command("grep", which_port, portBaseFile).Output()
		s := string(cmdOut)
		if strings.Contains(s, which_port) {
			return strings.TrimSpace(s), nil
		}
	} else {
		core.CreateFile(portBaseFile)
		return "", nil
	}

	return "", nil
}


// Store the last used port
func StoreLastUsedPort() error {

	// Fully qualified filename
	FurtureRefFile := core.FutureRefDir + PortBaseFileName
	log.Info("Storing the last used ports for this installation at: " + FurtureRefFile)

	// Delete the file if already exists
	_ = core.DeleteFile(FurtureRefFile)
	err := core.CreateFile(FurtureRefFile)
	if err != nil {return err}

	// Contents to write to the file
	var FutureContents []string
	FutureContents = append(FutureContents, "PORT_BASE: " + strconv.Itoa(GpInitSystemConfig.PortBase + core.EnvYAML.Install.TotalSegments))
	FutureContents = append(FutureContents, "MASTER_PORT: " + strconv.Itoa(GpInitSystemConfig.MasterPort + 1))

	// Write to the file
	err = core.WriteFile(FurtureRefFile, FutureContents)
	if err != nil {return err}

	return nil
}
