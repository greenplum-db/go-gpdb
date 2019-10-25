package main

import (
	"fmt"
	"github.com/jinzhu/configor"
	"os"
	"strconv"
	"strings"
)

// Struct that store the configuration file for the program to run
var Config = struct {
	CORE struct {
		APPLICATIONNAME string `yaml:"APPLICATION_NAME"`
		OS              string `yaml:"OS"`
		ARCH            string `yaml:"ARCH"`
		GOBUILD         string `yaml:"GO_BUILD"`
		BASEDIR         string `yaml:"BASE_DIR"`
		TEMPDIR         string `yaml:"TEMP_DIR"`
	} `yaml:"CORE"`
	DOWNLOAD struct {
		APITOKEN    string `yaml:"API_TOKEN"`
		DOWNLOADDIR string `yaml:"DOWNLOAD_DIR"`
	} `yaml:"DOWNLOAD"`
	INSTALL struct {
		ENVDIR               string `yaml:"ENV_DIR"`
		UNINTSALLDIR         string `yaml:"UNINTSALL_DIR"`
		FUTUREREFDIR         string `yaml:"FUTUREREF_DIR"`
		MASTERUSER           string `yaml:"MASTER_USER"`
		MASTERPASS           string `yaml:"MASTER_PASS"`
		GPMONPASS            string `yaml:"GPMON_PASS"`
		MASTERDATADIRECTORY  string `yaml:"MASTER_DATA_DIRECTORY"`
		SEGMENTDATADIRECTORY string `yaml:"SEGMENT_DATA_DIRECTORY"`
		MIRRORDATADIRECTORY  string `yaml:"MIRROR_DATA_DIRECTORY"`
		TOTALSEGMENT         int    `yaml:"TOTAL_SEGMENT"`
	} `yaml:"INSTALL"`
}{}

// If the parameter not set then set the defaults
func setDefaults(para, defaultValue, whichPara string) string {
	if IsValueEmpty(para) {
		Warnf("%s parameter missing in the config file, setting to default", whichPara)
		return defaultValue
	}
	return para
}

// If the parameter not set then error out
func isMissing(para, whichPara string) string {
	if IsValueEmpty(para) {
		Fatalf("Mandatory '%s' parameter is missing in the config file, please set it", whichPara)
	}
	return para
}

// All the directory should end up with "/"
func endWithSlash(path string) string {
	if !strings.HasSuffix(path, "/") {
		return fmt.Sprintf("%s/", path)
	}
	return path
}

// Validate the token
func validateToken() {
	if !cmdOptions.Github && (IsValueEmpty(Config.DOWNLOAD.APITOKEN) || Config.DOWNLOAD.APITOKEN == "<API TOKEN>") {
		// Check if its set as environment variables
		token := os.Getenv("UAA_API_TOKEN")
		if token == "" { // No token set
			Fatal("The API Token is either missing or not provided, please update the config file and try again" +
				", if you don't have a pivnet account and wish to use the open source gpdb release, " +
				"please use the \"--github\" flag.")
		} else {
			Config.DOWNLOAD.APITOKEN = token
		}
	}
}

// Read the configuration file and create directory if not exists
// or set default values if values are missing
func validateConfiguration() {
	Debug("Checking if the directories needed for the program exists")

	// Default parameter
	Config.CORE.BASEDIR = endWithSlash(setDefaults(Config.CORE.BASEDIR, "/home/gpadmin/", "BASE_DIR"))                                               // Base Dir
	Config.CORE.APPLICATIONNAME = setDefaults(Config.CORE.APPLICATIONNAME, "gpdbinstall", "APPLICATION_NAME")                                        // App name
	Config.CORE.TEMPDIR = endWithSlash(setDefaults(Config.CORE.TEMPDIR, "/temp/", "TEMP_DIR"))                                                       // Temp Directory
	Config.DOWNLOAD.DOWNLOADDIR = endWithSlash(setDefaults(Config.DOWNLOAD.DOWNLOADDIR, "/download/", "DOWNLOAD_DIR"))                               // Download Directory
	Config.INSTALL.ENVDIR = endWithSlash(setDefaults(Config.INSTALL.ENVDIR, "/env/", "ENV_DIR"))                                                     // Env Directory
	Config.INSTALL.UNINTSALLDIR = endWithSlash(setDefaults(Config.INSTALL.UNINTSALLDIR, "/uninstall/", "UNINTSALL_DIR"))                             // Uninstall Directory
	Config.INSTALL.FUTUREREFDIR = endWithSlash(setDefaults(Config.INSTALL.FUTUREREFDIR, "/future_reference/", "FUTUREREF_DIR"))                      // Future reference Directory
	Config.INSTALL.MASTERDATADIRECTORY = endWithSlash(setDefaults(Config.INSTALL.MASTERDATADIRECTORY, "/data/master/", "MASTER_DATA_DIRECTORY"))     // Master Directory
	Config.INSTALL.SEGMENTDATADIRECTORY = endWithSlash(setDefaults(Config.INSTALL.SEGMENTDATADIRECTORY, "/data/primary/", "SEGMENT_DATA_DIRECTORY")) // Segment Directory
	Config.INSTALL.MIRRORDATADIRECTORY = endWithSlash(setDefaults(Config.INSTALL.MIRRORDATADIRECTORY, "/data/mirror/", "MIRROR_DATA_DIRECTORY"))     // Segment Directory
	Config.INSTALL.MASTERPASS = setDefaults(Config.INSTALL.MASTERPASS, "changeme", "MASTER_PASS")                                                    // Master password
	Config.INSTALL.GPMONPASS = setDefaults(Config.INSTALL.GPMONPASS, "changeme", "GPMON_PASS")                                                       // Gpmon password
	Config.INSTALL.MASTERUSER = setDefaults(Config.INSTALL.MASTERUSER, "gpadmin", "MASTER_USER")                                                     // Master userv
	Config.INSTALL.TOTALSEGMENT = strToInt(setDefaults(strconv.Itoa(Config.INSTALL.TOTALSEGMENT), "2", "TOTAL_SEGMENT"))                             // Total Segments

	// Fail if these parameter is missing
	Config.CORE.OS = isMissing(Config.CORE.OS, "OS")                 // Go build OS
	Config.CORE.ARCH = isMissing(Config.CORE.ARCH, "ARCH")           // Go build ARCH
	Config.CORE.GOBUILD = isMissing(Config.CORE.GOBUILD, "GO_BUILD") // Go build version

	// Check if API Token
	validateToken()

	// Setup Path
	setupPath()
}

// Setup the directory location
func setupPath() {
	// Check if the directory exists, else create one.
	base_dir := Config.CORE.BASEDIR + Config.CORE.APPLICATIONNAME

	// Temp the files to
	Config.CORE.TEMPDIR = base_dir + Config.CORE.TEMPDIR
	CreateDir(Config.CORE.TEMPDIR)

	// Download the files to
	Config.DOWNLOAD.DOWNLOADDIR = base_dir + Config.DOWNLOAD.DOWNLOADDIR
	CreateDir(Config.DOWNLOAD.DOWNLOADDIR)

	// Environment location
	Config.INSTALL.ENVDIR = base_dir + Config.INSTALL.ENVDIR
	CreateDir(Config.INSTALL.ENVDIR)

	// Uninstall location
	Config.INSTALL.UNINTSALLDIR = base_dir + Config.INSTALL.UNINTSALLDIR
	CreateDir(Config.INSTALL.UNINTSALLDIR)

	// Future Reference location
	Config.INSTALL.FUTUREREFDIR = base_dir + Config.INSTALL.FUTUREREFDIR
	CreateDir(Config.INSTALL.FUTUREREFDIR)
}

// Load the configuration file to the memory
func config() {
	// Load the configuration
	configFile := os.Getenv("HOME") + "/config.yml"
	Debugf("Reading the configuration file and loading to the memory: %s", configFile)
	configor.Load(&Config, configFile)

	// Validate the configuration
	validateConfiguration()
}
