package core

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// Read the configuration and create directory if not exists
// Set default values if values are missing
func CreateDir() error {

	log.Info("Checking Directory Environment Varaibles")

	env_dir, _env_dir := os.LookupEnv("ENV_DIR")
	download_dir, _download_dir := os.LookupEnv("DOWNLOAD_DIR")
	uninstall_dir, _uninstall_dir := os.LookupEnv("UNINSTALL_DIR")

	// Check if the directory exists, else create one.
	base_dir, _base_dir := Clean(os.LookupEnv("BASE_DIR")) + "/" + os.LookupEnv("APPLICATION_NAME")

	if _, err := os.Stat(base_dir); os.IsNotExist(err) {
		log.Warning("Directory \"" + base_dir + "\" does not exists, creating one")
		os.Mkdir(base_dir, mode)
	}

	// Temp the files to
	TempDir = base_dir + EnvYAML.Core.TempDir
	tmp_bool, err := DoesFileOrDirExists(TempDir)
	if err != nil {
		return err
	}
	if !tmp_bool {
		log.Warning("Directory \"" + TempDir + "\" does not exists, creating one")
		err := os.MkdirAll(TempDir, 0755)
		if err != nil {
			return err
		}
	}

	// Download the files to
	DownloadDir = base_dir + EnvYAML.Download.DownloadDir
	dl_bool, err := DoesFileOrDirExists(DownloadDir)
	if err != nil {
		return err
	}
	if !dl_bool {
		log.Warning("Directory \"" + DownloadDir + "\" does not exists, creating one")
		err := os.MkdirAll(DownloadDir, 0755)
		if err != nil {
			return err
		}
	}

	// Environment location
	EnvFileDir = base_dir + EnvYAML.Install.EnvDir
	env_bool, err := DoesFileOrDirExists(EnvFileDir)
	if err != nil {
		return err
	}
	if !env_bool {
		log.Warning("Directory \"" + EnvFileDir + "\" does not exists, creating one")
		err := os.MkdirAll(EnvFileDir, 0755)
		if err != nil {
			return err
		}
	}

	// Uninstall location
	UninstallDir = base_dir + EnvYAML.Install.UnistallDir
	uninstall_bool, err := DoesFileOrDirExists(UninstallDir)
	if err != nil {
		return err
	}
	if !uninstall_bool {
		log.Warning("Directory \"" + UninstallDir + "\" does not exists, creating one")
		err := os.MkdirAll(UninstallDir, 0755)
		if err != nil {
			return err
		}
	}

	// Future Reference location
	FutureRefDir = base_dir + EnvYAML.Install.FutureRef
	futureref_bool, err := DoesFileOrDirExists(FutureRefDir)
	if err != nil {
		return err
	}
	if !futureref_bool {
		log.Warning("Directory \"" + FutureRefDir + "\" does not exists, creating one")
		err := os.MkdirAll(FutureRefDir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

// Configuration file reader.
func Config() error {

	home := os.Getenv("HOME")

	// Creating Directory needed for the program
	err = CreateDir()
	if err != nil {
		return err
	}

	return nil
}
