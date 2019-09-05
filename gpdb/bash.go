package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Build executable shell script
func buildExecutableBashScript(filename string, executableFilename string, args []string) error {
	Debugf("Starting program to build executable script to install the binaries")

	// Create the file
	createFile(filename)

	// Create arguments on what needs to be written to the file
	Infof("Generating installer arguments to be passed to the file: %s", filename)
	executableLine := fmt.Sprintf("/bin/sh %s << EOF", executableFilename)

	// Load the contents to the file
	var passArgs []string
	passArgs = append(passArgs, executableLine)
	for _, v := range args {
		passArgs = append(passArgs, v)
	}
	passArgs = append(passArgs, "EOF")

	// Write to the file
	writeFile(filename, passArgs)
	return nil
}

// Execute shell script when called
func executeBinaries(binaryFile string, bashfilename string, scriptOptions []string) error {
	Debugf("Executing the binaries with contents in the file %s", binaryFile)
	// Build a quick shell script to install binaries, Filename name
	filename := Config.CORE.TEMPDIR + bashfilename

	// Cleanup the file if already exists (ignore error if any)
	deleteFile(filename)

	// Create the shell script
	err := buildExecutableBashScript(filename, binaryFile, scriptOptions)
	if err != nil {
		return err
	}

	// Execute the installer script
	Infof("Executing the bash script: %s", filename)
	_, err = exec.Command("/bin/sh", filename).Output()
	if err != nil {
		return err
	}

	// Cleanup the tempFile
	Infof("Cleaning up the file \"%s\", if found", filename)
	deleteFile(filename)
	return nil
}

// Execute Os commands
func executeOsCommand(command string, arguments ...string) {
	Debugf("Executing the command \"%s\" using the arguments %s", command, arguments)
	// Execute the command
	cmd := exec.Command(command, arguments...)

	// Attach the os error/output from the screen
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Start the command
	err := cmd.Start()
	if err != nil {
		Fatalf("Failed to start the start command %s, err: %v", command, err)
	}

	// Wait for it to finish
	// For some reason the gpinitsystem produces exit code 1 even then command ran successfully , so we ignore the exit code 1 here for gpinitsystem
	err = cmd.Wait()
	if err != nil && (strings.HasSuffix(command, "gpinitsystem") && err.Error() != "exit status 1") {
		Fatalf("Failed while waiting for the command %s err: %v", command, err)
	}
}

// Same as executeOsCommand but this one return output
func executeOsCommandOutput(command string, args ...string) ([]byte, error) {
	Debugf("Executing the command \"%s\" with output using the arguments %s", command, args)
	cmdOut, err := exec.Command(command, args...).Output()
	if err != nil {
		return cmdOut, fmt.Errorf("failed when executing OS command output, err: %v", err)
	}
	return cmdOut, nil
}

// Generate the file and execute the command
func generateBashFileAndExecuteTheBashFile(filename, command string, contents []string) {
	Debugf("Generate the bash file \"%s\" and executing using the command \"%s\"", filename, command)
	deleteFile(filename)
	createFile(filename)
	writeFile(filename, contents)
	executeOsCommand(command, filename)
	defer deleteFile(filename)
}
