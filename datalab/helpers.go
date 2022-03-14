package main

import (
	"fmt"
	"github.com/ryanuber/columnize"
	"os"
	"os/exec"
	"strings"
)

// IsValueEmpty Check is the value is empty
func IsValueEmpty(v string) bool {
	if len(strings.TrimSpace(v)) == 0 {
		return true
	}
	return false
}

// Check if the Os executable exists
func isCommandAvailable(name string) {
	cmd := exec.Command(name, "-v")
	if err := cmd.Run(); err != nil {
		Fatalf("the command executable \"%s\" is not installed on your machine, please install it", name)
	}
}

// Check if we have the name already on our config
func nameInConfig(a string) (int, bool) {
	for index, b := range Config.Vagrants {
		if b.Name == a {
			return index, true
		}
	}
	return -1, false
}

// Execute the OS command
func executeOsCommand(command string, env []string, args ...string) {
	// Check if the command is installed on the local host
	isCommandAvailable(command)
	// Initialize the command
	cmd := exec.Command(command, args...)
	// Set all the env variable
	cmd.Env = os.Environ()
	for _, e := range env {
		cmd.Env = append(cmd.Env, e)
	}
	// Go to the directory
	cmd.Dir = Config.GoGPDBPath
	// Attach the stdout and std err
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	// Start and wait for the execution
	err := cmd.Start()
	if err != nil {
		Fatalf("Failed to start the command %s, with arguments %s, err: %v", command, args, err)
	}
	err = cmd.Wait()
	if err != nil {
		Fatalf("Failed in waiting for the command %s, with arguments %s, err: %v", command, args, err)
	}
}

// Print the data in tabular format
func printOnScreen(message string, content []string) {

	// Message before the table
	fmt.Printf("\n%s\n\n", message)

	// Print the table format
	result := columnize.SimpleFormat(content)

	// Print the results
	fmt.Println(result + "\n")

}

// Get the element from the config
func getVagrantKeyFromConfig(command string) VagrantKey {
	index, exists := nameInConfig(cmdOptions.Hostname)
	if !exists {
		Fatalf(missingVMInOurConfig, command, cmdOptions.Hostname, programName)
	}
	return Config.Vagrants[index]
}

// Generate a env based on what key we got
func generateEnvByVagrantKey(command string) []string {
	vk := getVagrantKeyFromConfig(command)
	return generateEnvArray(vk.CPU, vk.Memory, vk.Segment, vk.Os, vk.Subnet, vk.Name, vk.Standby, vk.Developer)
}

// Remove slash at suffix
func removeSlash(str string) string {
	if strings.HasSuffix(str, "/") {
		return str[0 : len(str)-1]
	} else {
		return str
	}
}
