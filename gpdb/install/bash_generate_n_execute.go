package install

import (
	"os/exec"
	"github.com/ielizaga/piv-go-gpdb/core"
)

// Build executable shell script
func BuildExecutableBashScript (filename string, executable_filename string, args []string) error {

	log.Info("Starting program to build executable script to install the binaries")

	// Create the file
	err := core.CreateFile(filename)
	if err != nil { return err }

	// Create arguments on what needs to be written to the file
	log.Info("Generating installer arguments to be passed to the file: "+ filename)
	executable_line := "/bin/sh " +  executable_filename + " &>/dev/null << EOF"
	var passArgs []string
	passArgs = append(passArgs, executable_line)
	for _ , v := range args {
		passArgs = append(passArgs, v)
	}
	passArgs = append(passArgs, "EOF")

	// Write to the file
	err = core.WriteFile(filename, passArgs)
	if err != nil { return err }

	return nil
}


// Execute shell script when called
func ExecuteBinaries(binary_file string, bashfilename string, script_options []string) error {

	// Build a quick shell script to install binaries
	// Filename name
	filename := core.TempDir + bashfilename

	// Cleanup the file if already exists (ignore error if any)
	_ = CleanupTempFile(filename)

	// Create the shell script
	err := BuildExecutableBashScript(filename, binary_file, script_options)
	if err != nil { return err }

	// Execute the installer script
	log.Info("Executing the bash script: "+ filename )
	_, err = exec.Command("/bin/sh", filename).Output()
	if err != nil { return err}

	// Cleanup the tempFile
	log.Info("Cleaning up the file \""+ filename + "\", if found")
	err = CleanupTempFile(filename)
	if err != nil { return err}

	return nil
}


// Cleanup the temp files once done
func CleanupTempFile (filename string) error {

	// Remove any tempfile
	err := core.DeleteFile(filename)
	if err != nil { return err }

	return nil
}