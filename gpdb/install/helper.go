package install

import (
	"../core"
	"errors"
	"os"
	"os/exec"
)

// Check if the directory provided is readable and writeable
func DirValidator(master string, segment string) error {

	// Check if the master & segment location is provided.
	if core.IsValueEmpty(master) {
		return errors.New("MASTER_DATA_DIRECTORY parameter missing in the config file, please set it and try again")
	} else {
		MasterDIR = master
	}

	if core.IsValueEmpty(segment) {
		return errors.New("SEGMENT_DATA_DIRECTORY parameter missing in the config file, please set it and try again")
	} else {
		SegmentDIR = segment
	}

	// Check if the master & segment location is readable and writable
	master_dir, err := core.DoesFileOrDirExists(MasterDIR)
	if err != nil {
		return err
	}

	segment_dir, err := core.DoesFileOrDirExists(SegmentDIR)
	if err != nil {
		return err
	}

	// if the file doesn't exists then let try creating it ...
	if !master_dir || !segment_dir {
		err := os.MkdirAll(MasterDIR, 0775)
		if err != nil {
			return err
		}
		err = os.MkdirAll(SegmentDIR, 0775)
		if err != nil {
			return err
		}
	}

	return nil
}

// Check if the provided hostnames are valid
func CheckHostnameIsValid() error {

	// Check Master host parameter is set
	if core.IsValueEmpty(core.EnvYAML.Install.MasterHost) {
		return errors.New("MASTER_HOST parameter missing in the config file, please set it and try again")
	} else {
		GpInitSystemConfig.MasterHostName = core.EnvYAML.Install.MasterHost
	}

	// Check if the provided hostname can be ssh'ed
	hostname := GpInitSystemConfig.MasterHostName
	log.Info("Checking connectivity to host \"" + hostname + "\" can be established")

	_, err := exec.Command("ssh", hostname, "-o", "ConnectTimeout=5", "echo 1").Output()
	if err != nil {
		return err
	}

	// Enable passwordless login
	err = ExecuteGpsshExkey()
	if err != nil {
		return err
	}

	return nil
}

// Run keyless access to the server
func ExecuteGpsshExkey() error {

	// Source GPDB PATH
	err := SourceGPDBPath()
	if err != nil {
		return err
	}

	// Checking if the username and password parameters are passed
	if core.IsValueEmpty(core.EnvYAML.Install.MasterUser) {
		return errors.New("MASTER_USER parameter missing in the config file, please set it and try again")
	}
	if core.IsValueEmpty(core.EnvYAML.Install.MasterPass) {
		return errors.New("MASTER_PASS parameter missing in the config file, please set it and try again")
	}

	// Execute gpssh script to enable keyless access
	log.Info("Running gpssh-exkeys to enable keyless access on this server")
	cmd := exec.Command(os.Getenv("GPHOME")+"/bin/gpssh-exkeys", "-h", GpInitSystemConfig.MasterHostName)
	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

// Source greenplum path
func SourceGPDBPath() error {

	// Setting up greenplum path
	err := os.Setenv("GPHOME", BinaryInstallLocation)
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

// Execute Bash Script
func ExecuteBash(filename string, BashScript []string) error {

	// Delete the file if exists, and then create the file
	_ = core.DeleteFile(filename)
	err := core.CreateFile(filename)
	if err != nil {
		return err
	}

	// Write the script to the file
	err = core.WriteFile(filename, BashScript)
	if err != nil {
		return err
	}

	// Execute the script
	log.Info("Executing the file: " + filename)
	cmd := exec.Command("/bin/sh", filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}

	// Cleanup temp files.
	err = core.DeleteFile(filename)
	if err != nil {
		return err
	}

	return nil

}
