package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/benhoyt/goawk/interp"
	"github.com/benhoyt/goawk/parser"
	"github.com/mholt/archiver"
	"github.com/ryanuber/columnize"
)

// Function that checks if the string is available on a array.
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Check is the value is empty
func IsValueEmpty(v string) bool {
	if len(strings.TrimSpace(v)) == 0 {
		return true
	}
	return false
}

// exists returns whether the given file or directory exists or not
func doesFileOrDirExists(path string) (bool, error) {
	Debugf("Checking if the directory %s exists", path)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Create directory
func CreateDir(path string) {
	// Check if the path or directory exists
	exists, err := doesFileOrDirExists(path)
	if err != nil {
		Fatalf("Failed to check the directory status, the error: %v", err)
	}
	// If not exists then create one
	if !exists {
		Debugf("Directory \"%s\" does not exists, creating one", path)
		err := os.MkdirAll(path, 0755)
		if err != nil {
			Fatalf("Failed to create the directory, the error: %v", err)
		}
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

// Progress of download
func PrintDownloadPercent(done chan int64, path string, total int64) {
	var stop bool = false

	for {
		select {
		case <-done:
			stop = true
		default:

			// Open the file
			file, err := os.Open(path)
			if err != nil {
				Fatalf("Error in opening the file, err: %v", err)
			}

			// Get stats of the file
			fi, err := file.Stat()
			if err != nil {
				Fatalf("Error in obtaining the stats of the file, err: %v", err)
			}

			// Size now
			size := fi.Size()

			// Display Progress of download
			if size == 0 {
				size = 1
			}

			var percent float64 = float64(size) / float64(total) * 100
			var bytesToMB float64 = 1024 * 1024

			fmt.Printf("Downloading file %.2f MB of %.2f MB: %.0f", float64(size)/bytesToMB, float64(total)/bytesToMB, percent)
			fmt.Println("% completed")
		}

		// Download is completed, time to terminate
		if stop {
			Info("Downloading completed ....")
			break
		}

		// Ask to sleep, before repainting the screen.
		time.Sleep(time.Second)
	}
}

// Remove files if exists so that unarchiver can succeed.
func removeOldBinFiles(search string) {
	if cmdOptions.Product == "gpdb" {
		removeFiles(Config.DOWNLOAD.DOWNLOADDIR, fmt.Sprintf("*%s*.bin*", cmdOptions.Version))
		removeFiles(Config.DOWNLOAD.DOWNLOADDIR, "*README_INSTALL*")
	} else if cmdOptions.Product == "gpcc" && !isThis4x() {
		removeFiles(Config.DOWNLOAD.DOWNLOADDIR, fmt.Sprintf("*%s*.bin*", cmdOptions.CCVersion))
		removeFiles(Config.DOWNLOAD.DOWNLOADDIR, "*README_INSTALL*")
	} else { // GPCC 4.x has folder into folders
		allfiles, _ := FilterDirsGlob(Config.DOWNLOAD.DOWNLOADDIR, fmt.Sprintf("%s", removeZip(search)))
		for _, v := range allfiles {
			if !strings.HasSuffix(v, ".zip") {
				deleteFile(v)
			}
		}
	}
}

// Unzip the binaries.
func locateAndExtractPackage(search string) (string, bool) {
	// Check if we can find the binaries in the directory
	allfiles, _ := FilterDirsGlob(Config.DOWNLOAD.DOWNLOADDIR, fmt.Sprintf("%s", search))

	// Did we find any
	if len(allfiles) > 0 {
		binary := detectFileName(allfiles)
		if strings.HasSuffix(binary, ".rpm") {
			return locatedRpmFile(search, binary)
		} else {
			return locatedBinaryFile(search, binary)
		}
	} else {
		if cmdOptions.Product == "gpdb" {
			Fatalf("No binary zip found for the product %s with version %s under directory %s", cmdOptions.Product, cmdOptions.Version, Config.DOWNLOAD.DOWNLOADDIR)
		} else if cmdOptions.Product == "gpcc" {
			Fatalf("No binary zip found for the product %s with version %s under directory %s", cmdOptions.Product, cmdOptions.CCVersion, Config.DOWNLOAD.DOWNLOADDIR)
		} else { // Should never reach here since we only accept gpdb and gpcc only, if it does then print the error below
			Fatalf("Don't know the installation tag for product provided: %s", cmdOptions.Product)
		}
	}
	return "", true
}

// Detect the file
// Starting from GPDB 6, the filename or GPDB and GPCC remains the
// same and breaks the GPCC installation, so we detect it here and
// select the appropriate file
func detectFileName(files []string) string {
	if len(files) == 1 {
		return files[0]
	} else {
		for _, file := range files {
			if cmdOptions.Product == "gpdb" {
				if !strings.Contains(file, "cc-web") {
					return file
				}
			} else {
				if strings.Contains(file, "cc-web") {
					return file
				}
			}
		}
	}
	return ""
}

// located a binary file
func locatedBinaryFile(search, binary string) (string, bool) {
	Infof("Found & unzip the binary for the version %s: %s", cmdOptions.Version, binary)
	removeOldBinFiles(search)
	err := archiver.Unarchive(binary, Config.DOWNLOAD.DOWNLOADDIR)
	if err != nil {
		Fatalf("Couldn't unzip the file, err: %v", err)
	}
	Debugf("Unzipped the file %s completed successfully", binary)
	return obtainExecutableFilename(search), true
}

// located a rpm file
func locatedRpmFile(search, binary string) (string, bool) {
	Infof("Found a rpm binary for the version %s: %s", cmdOptions.Version, binary)
	return binary, false
}

// Locate the full directory name where the rpm was installed
func locateGreenplumInstallationDirectory() string {
	// rpm usually installs the software in /usr/local
	// we need to check what is the directory name it has taken
	baseDir := "/usr/local/"
	folders, _ := FilterDirsGlob(baseDir, fmt.Sprintf("*%s*", cmdOptions.Version))
	if len(folders) > 0 {
		// We found one
		return detectFileName(folders)
	} else {
		Fatalf(fmt.Sprintf("Cannot locate the directory name at %s where the version %s is installed", baseDir, cmdOptions.Version))
	}

	return ""
}

// Find the binary file
func findBinaryFile(search, version string) string {
	binFile, _ := FilterDirsGlob(Config.DOWNLOAD.DOWNLOADDIR, fmt.Sprintf("%s.bin", search))
	if len(binFile) > 0 {
		return binFile[0]
	} else {
		Fatalf("No binaries found for the product %s with version %s under directory %s", cmdOptions.Product, version, Config.DOWNLOAD.DOWNLOADDIR)
	}
	return ""
}

func removeZip(search string) string {
	return strings.Replace(search, ".zip", "", -1)
}

// Get the execute file
func obtainExecutableFilename(search string) string {
	if cmdOptions.Product == "gpdb" { // Get the binary file name
		return findBinaryFile(search, cmdOptions.Version)
	} else if cmdOptions.Product == "gpcc" { // GPCC binaries
		if isThis4x() { // newer directory
			// Get the binary file name
			binFile, _ := FilterDirsGlob(Config.DOWNLOAD.DOWNLOADDIR, fmt.Sprintf("%[1]s/%[1]s", removeZip(search)))
			if len(binFile) > 0 {
				return binFile[0]
			} else {
				Fatalf("No binaries found for the product %s with version %s under directory %s", cmdOptions.Product, cmdOptions.CCVersion, Config.DOWNLOAD.DOWNLOADDIR)
			}
		} else { // older directory
			return findBinaryFile(search, cmdOptions.CCVersion)
		}
	} else { // Should never reach here since we only accept gpdb and gpcc only, if it does then print the error below
		Fatalf("Don't know the installation tag for product provided: %s", cmdOptions.Product)
	}
	return ""
}

// Extract the contents that we are interested
func contentExtractor(contents []byte, src string, vars []string) bytes.Buffer {
	// Create a parser
	prog, err := parser.ParseProgram([]byte(src), nil)
	if err != nil {
		Fatalf("Failed to parse the program: %s", src)
	}

	// The configuration
	var buf bytes.Buffer
	config := &interp.Config{
		Stdin:  bytes.NewReader([]byte(contents)),
		Vars:   vars,
		Output: &buf,
	}

	// Execute the program
	_, err = interp.ExecProgram(prog, config)
	if err != nil {
		Fatalf("Failure in executing the goawk script: %v", err)
	}
	return buf
}

// Check if the binaries exits and unzip the binaries.
func getBinaryFile(version string) (string, bool) {
	Debugf("Finding and unziping the binaries for the version %s", version)
	return locateAndExtractPackage(fmt.Sprintf("*%s*", version))
}

// Remove blank lines from the contentExtractor
func removeBlanks(s string) string {
	regex, err := regexp.Compile("\n$")
	if err != nil {
		Fatalf("Failure in removing blank lines, err: %v", err)
	}
	s = strings.TrimSpace(regex.ReplaceAllString(s, ""))
	return s
}

// is the port out of range
func outOfRangePort(port string) bool {
	if strToInt(port) > 63000 {
		return true
	}
	return false
}

// string to init
func strToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// Extract the version
func extractVersion(version string) float64 {
	var f float64
	re := regexp.MustCompile("\\d.\\d+")
	labNo := re.FindAllString(version, -1)
	if len(labNo) > 0 { // we found the lab number
		f, err := strconv.ParseFloat(labNo[0], 64)
		if err != nil {
			Errorf("Failed to convert the string(%s) to float, received the error: %v", labNo[0], err)
		}
		return f
	}
	return f
}

// Ensure the version is of the format
// eg. 4.3.30.10 etc
func isValidVersionFormat(version string) bool {
	match, _ := regexp.MatchString("[0-9]+\\.[0-9]+\\.[0-9]+", version)
	return match
}

// Extract the version from the name
func extractVersionNumber(filename string) string {
	r, _ := regexp.Compile(`(-[0-9]+.[0-9]+.[0-9]+-|-[0-9]+.[0-9]+.[0-9]+.[0-9]+-)`)
	version := r.FindString(filename)
	return strings.Replace(version, "-", "", -1)
}

// Size in MB
func sizeInMB(size int64) int64 {
	if size > 0 {
		size = size / 1024 / 1024
	}
	return size
}

// Is this GPDB 6 version
func isThisGPDB6xAndAbove() bool {
	v := extractVersion(cmdOptions.Version)
	if v >= 6 {
		return true
	} else {
		return false
	}
}

// Check if the Os executable exists
func isCommandAvailable(name string) bool {
	cmd := exec.Command("command", name, "-V")
	if err := cmd.Run(); err != nil {
		Fatalf("%s executable is not installed on this box, please run 'yum install -y %[1]s to install it'", name, name)
	}
	return true
}

// Check if the file exists
func DidWeDownloadThisVersionBefore(pattern, mesg string) bool {
	filePath, _ := FilterDirsGlob(Config.DOWNLOAD.DOWNLOADDIR, fmt.Sprintf(pattern, cmdOptions.Version))
	if len(filePath) > 0 && !cmdOptions.Always {
		if !strings.Contains(filePath[0], "cc-web") { // should not collide with GPCC 6
			Warnf("%s %s found, skipping download...", mesg, filePath[0])
			Warn("To force re-download of the file, use -a flag")
			return true
		}
	}
	return false
}
