package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// Create the file
func createFile(path string) {
	Debugf("Creating the file: %s", path)

	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			Fatalf("Error in creating a file: %v", err)
		}
		defer file.Close()
	}
}

// Open the file
func openFile(path string, openType int, permission os.FileMode) *os.File {
	Debugf("Opening the file: %s", path)

	// re-open file
	var file, err = os.OpenFile(path, openType, permission)
	if err != nil {
		Fatalf("Error in opening the file: %v", err)
	}
	return file
}

// write to the file
func writeFile(path string, contents []string) {
	Debugf("Writing to the file: %s", path)

	// Create file is not exists
	createFile(path)

	// open file using READ & WRITE permission
	var file = openFile(path, os.O_RDWR, 0644)
	defer file.Close()

	// write some text line-by-line to file
	for _, k := range contents {
		text := string(k)
		_, err := file.WriteString(text + "\n")
		if err != nil {
			Fatalf("Error in writing to the file: %v", err)
		}
	}

	// save changes
	err := file.Sync()
	if err != nil {
		Fatalf("Error in saving the write content to the file: %v", err)
	}
}

// Append the file
func appendFile(path string, args []string) {
	Debugf("Appending to the file: %s", path)

	// Create file is not exists
	createFile(path)

	// open file using READ & APPEND permission
	var file = openFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	defer file.Close()

	// append to the file
	for _, a := range args {
		if _, err := file.WriteString(a + "\n"); err != nil {
			Fatalf("Error in appending to the file: %v", err)
		}
	}

	// save changes
	err := file.Sync()
	if err != nil {
		Fatalf("Error in saving the append content to the file: %v", err)
	}
}

// Read the file
func readFile(path string) []byte {
	Debugf("Reading to the file: %s", path)

	content, err := ioutil.ReadFile(path)
	if err != nil {
		Fatalf("Error in reading the file \"%s\", err: %v", path, err)
	}

	return content
}

// Delete the file
func deleteFile(path string) {
	Debugf("Deleting the file: %s", path)

	// delete file
	var err = os.RemoveAll(path)
	if err != nil && !os.IsNotExist(err) {
		Fatalf("Error in deleting the file: %v", err)
	}
}

// Search the directory for the matching files
func FilterDirsGlob(dir, search string) ([]string, error) {
	return filepath.Glob(filepath.Join(dir, search))
}

// Remove all the file based on search
func removeFiles(path, search string) {
	Debugf("Removing the file with search %s from path %s", search, path)
	allFiles, _ := FilterDirsGlob(path, search)
	for _, f := range allFiles {
		if err := os.RemoveAll(f); err != nil {
			Fatalf("Failed to remove the file from path %s%s, err: %v", path, search, err)
		}
	}
}

// Get the size of the directory
func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

// Does the file exists
func fileExists(filename string) bool {
	if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
		return false
	}
	return true
}