package install

import (
	"io/ioutil"
	"regexp"
	"strings"
	"os/exec"
	"github.com/ielizaga/piv-go-gpdb/core"
)

// FInd WLM Binaries
func findWLMBinaries(path string, keyword string, info_text string, warn_text string) (string, error) {

	// Listing all the files in the directory
	output, err := ioutil.ReadDir(path)
	if err != nil { return "", err }
	for _, f := range output {

		// Hunting for file that matches the condition
		re := regexp.MustCompile(``+ keyword +``)
		matches := re.FindStringSubmatch(f.Name())

		// Found the file, return the file name
		if len(matches) != 0 {
			log.Info(info_text + matches[0])
			return matches[0], nil
		}
	}

	// Did we find the file.
	if !core.IsValueEmpty(warn_text) {
		log.Warning(warn_text)
	} else {
		log.Info(warn_text)
	}

	return "", nil
}


// Install WLM Binaries
func Installwlm(t string) error {

	log.Info("Installing workload manager")

	// WLM Binary file name
	log.Info("Checking the directory \"" + GPPERFMONHOME + "\" for WLM binaries.")
	WLMBinaryFile, err := findWLMBinaries(GPPERFMONHOME,
		"gp-wlm.*.bin",
		"Found the WLM Binary file: ",
		"Cannot find the WLM Binaries on the directory \"" + GPPERFMONHOME + "\", skipping the WLM Installation")
	if err != nil { return err }

	// Version and Directory
	WLMVersion := strings.Replace(WLMBinaryFile, ".bin", "", 1)
	WLMDir := core.EnvYAML.Install.MasterDataDirectory + "wlm/" + WLMVersion

	// If no binaries are found then return or exit from this function.
	if WLMBinaryFile == "" {
		return nil
	}

	// Check if the binary we are installing is already installed on this server
	HasThisWLMAlreadyInstalled, err := findWLMBinaries( core.EnvYAML.Install.MasterDataDirectory + "wlm/",
		WLMVersion,
		"Found the WLM Version already installed: ",
		"Cannot find any previous installation of WLM version \""+ WLMVersion +"\" installed on this server")

	// Yes found a previous installation, so lets uninstall and reinstall it (in case if its corrupted).
	if !core.IsValueEmpty(HasThisWLMAlreadyInstalled) {
		err := UninstallWLM(t, WLMDir)
		if err != nil { return err }
	}

	// Installing the workload manager.
	log.Info("Creating script to install the workload manager")
	var InstallWLMArg []string
	InstallWLMArg = append(InstallWLMArg, "source " + EnvFileName)
	InstallWLMArg = append(InstallWLMArg, "chmod +x " + GPPERFMONHOME + "/" + WLMBinaryFile)
	InstallWLMArg = append(InstallWLMArg, "mkdir -p " + WLMDir)
	InstallWLMArg = append(InstallWLMArg, GPPERFMONHOME + "/" + WLMBinaryFile + " --install=" + WLMDir)

	// Write it to the file and execute.
	file := core.TempDir + "install_wlm.sh"
	err = ExecuteBash(file, InstallWLMArg)
	if err != nil { return err }

	// Update the environment file
	err = UpdateWlmEnvFile(WLMDir, WLMVersion)
	if err != nil { return err }

	log.Info("Installation of WLM manager is complete")

	return nil
}

// Uninstall WLM
func UninstallWLM(t string, WLMDir string) error {

	log.Info("Found a old version of WLM: " + WLMDir)

	// Stop all WLM services
	err := StopWLMService()
	if err != nil { return err }

	// Uninstall Arguments
	var UninstallWLMArgs []string
	UninstallWLMArgs = append(UninstallWLMArgs, "source " + EnvFileName)
	UninstallWLMArgs = append(UninstallWLMArgs, WLMDir + "/gp-wlm/bin/uninstall --symlink " + WLMInstallDir + "/gp-wlm")
	UninstallWLMArgs = append(UninstallWLMArgs, "cp " + EnvFileName + " " + EnvFileName + "." + t)
	UninstallWLMArgs = append(UninstallWLMArgs, "egrep -v \"wlm|WLM\" " + EnvFileName + "." + t + " > " + EnvFileName)
	UninstallWLMArgs = append(UninstallWLMArgs, "rm " + EnvFileName + "." + t)

	// Write it to the file.
	file := core.TempDir + "uninstall_wlm.sh"
	err = ExecuteBash(file, UninstallWLMArgs)
	if err != nil { return err }

	return nil
}


// Get WLM version
func GetWLMVersion(path string) (string, error) {

	log.Info("Obtaining the version of the WLM installed")

	// extract the version of the WLM
	wlmv, err := exec.Command(path+"/gp-wlm/bin/gp-wlm -v").Output()
	if err != nil { return "", err }

	return string(wlmv), nil
}
