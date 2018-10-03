package core

import (
	"io/ioutil"
	"os"
	"gopkg.in/yaml.v2"
)

// Read the configuration file and create directory if not exists
// or set default values if values are missing
func CreateDir() error {

	log.Info("Checking if the directories needed for the program exists")
	if IsValueEmpty(EnvYAML.Core.BaseDir) {
		log.Warning("BASE_DIR parameter missing in the config file, setting to default")
		EnvYAML.Core.BaseDir = "/home/gpadmin/"
	}

	// App name
	if IsValueEmpty(EnvYAML.Core.AppName) {
		log.Warning("APPLICATION_NAME parameter missing in the config file, setting to default")
		EnvYAML.Core.AppName = "gpdbinstall"
	}

	// Temp Directory
	if IsValueEmpty(EnvYAML.Core.TempDir) {
		log.Warning("TEMP_DIR parameter missing in the config file, setting to default")
		EnvYAML.Core.TempDir = "/temp/"
	}

	// Download Directory
	if IsValueEmpty(EnvYAML.Download.DownloadDir) {
		log.Warning("DOWNLOAD_DIR parameter missing in the config file, setting to default")
		EnvYAML.Download.DownloadDir = "/download/"
	}

	// Env Directory
	if IsValueEmpty(EnvYAML.Install.EnvDir) {
		log.Warning("ENV_DIR parameter missing in the config file, setting to default")
		EnvYAML.Install.EnvDir = "/env/"
	}

	// Uninstall Directory
	if IsValueEmpty(EnvYAML.Install.UnistallDir) {
		log.Warning("UNINTSALL_DIR parameter missing in the config file, setting to default")
		EnvYAML.Install.UnistallDir = "/uninstall/"
	}

	// Future Reference Directory
	if IsValueEmpty(EnvYAML.Install.FutureRef) {
		log.Warning("UNINTSALL_DIR parameter missing in the config file, setting to default")
		EnvYAML.Install.FutureRef = "/future_reference/"
	}

	// Check if the directory exists, else create one.
	base_dir := EnvYAML.Core.BaseDir + EnvYAML.Core.AppName

	// Temp the files to
	TempDir =  base_dir + EnvYAML.Core.TempDir
	tmp_bool, err := DoesFileOrDirExists(TempDir)
	if err != nil {return err}
	if !tmp_bool {
		log.Warning("Directory \""+ TempDir + "\" does not exists, creating one")
		err:= os.MkdirAll(TempDir, 0755)
		if err != nil {return err}
	}

	// Download the files to
	DownloadDir =  base_dir + EnvYAML.Download.DownloadDir
	dl_bool, err := DoesFileOrDirExists(DownloadDir)
	if err != nil {return err}
	if !dl_bool {
		log.Warning("Directory \""+ DownloadDir + "\" does not exists, creating one")
		err:= os.MkdirAll(DownloadDir, 0755)
		if err != nil {return err}
	}

	// Environment location
	EnvFileDir = base_dir + EnvYAML.Install.EnvDir
	env_bool, err := DoesFileOrDirExists(EnvFileDir)
	if err != nil {return err}
	if !env_bool {
		log.Warning("Directory \""+ EnvFileDir + "\" does not exists, creating one")
		err := os.MkdirAll(EnvFileDir, 0755)
		if err != nil {return err}
	}

	// Uninstall location
	UninstallDir = base_dir + EnvYAML.Install.UnistallDir
	uninstall_bool, err := DoesFileOrDirExists(UninstallDir)
	if err != nil {return err}
	if !uninstall_bool {
		log.Warning("Directory \""+ UninstallDir + "\" does not exists, creating one")
		err := os.MkdirAll(UninstallDir, 0755)
		if err != nil {return err}
	}

	// Future Reference location
	FutureRefDir = base_dir + EnvYAML.Install.FutureRef
	futureref_bool, err := DoesFileOrDirExists(FutureRefDir)
	if err != nil {return err}
	if !futureref_bool {
		log.Warning("Directory \""+ FutureRefDir + "\" does not exists, creating one")
		err := os.MkdirAll(FutureRefDir, 0755)
		if err != nil {return err}
	}

	return nil
}


// Configuration file reader.
func Config() error {

	home := os.Getenv("HOME")
	configFile := home + "/.config.yml"

	log.Infof("Reading the configuration file: %s", configFile)

	// Read the config file and store the value on a struct
	source, err := ioutil.ReadFile(configFile)
	if err != nil {return err}
	yaml.Unmarshal(source, &EnvYAML)

	// Creating Directory needed for the program
	err = CreateDir()
	if err != nil {return err}

	return nil
}
