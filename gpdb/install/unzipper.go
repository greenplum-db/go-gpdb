package install

import (
	"os"
	"path/filepath"
	"io"
	"regexp"
	"errors"
	"archive/zip"
	"strings"
	"github.com/ielizaga/piv-go-gpdb/core"
)

// Provides all the files in the directory
func ListFilesInDir(dest string) ([]string, error) {
	log.Info("Listing all the files in the download directory: " + dest)
	files, err := filepath.Glob(dest + "*")
	if err != nil { return files, err }
	return files, nil
}


// Find the file that matches the search string
func GetBinaryOfMatchingVersion(files []string, version string) (string, error) {

	// Loop through all the files and see if we can find a binaries that matches with the version
	log.Info("Checking if there is binaries that matches the version to install: " + version)
	for _, f := range files {
		pattern := "greenplum.*" + version
		matched, _ := regexp.MatchString(pattern, f)

		// If we find a match then return the file name
		if matched {
			return f, nil
		}
	}
	return "", errors.New("Unable to find any file that matches the version: " + version )
}

// Unzip the provided file.
func Unzip(src, dest string) error {

	// unzip the file
	log.Info("Unzipping the zip file: " + src)
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		// there is no consistency in the name of binaries, so here we will ensure the
		// binary name is equal to the zip file name
		filename := strings.Replace(src, ".zip", ".bin", 1)
		path := filename

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func UnzipBinary(version string) (string, error) {

	// List all the files in the download directory
	all_files, err := ListFilesInDir(core.DownloadDir)
	if err != nil { return "", err }

	// Check if any of the file matches to requested version to install
	binary_file, err := GetBinaryOfMatchingVersion(all_files, version)
	if err != nil { return "", err }

	// If we cannot find a match then error out
	if core.IsValueEmpty(binary_file) {
		return "", errors.New("Cannot find any binaries that matches the version: " + version)
	}

	// Check if the file is a zip file found or Unzip and do work accordingly
	if strings.HasSuffix(binary_file, ".zip") {
		err := Unzip(binary_file, core.DownloadDir)
		binary_file = strings.Replace(binary_file, ".zip", ".bin", 1)
		return binary_file, err
	} else {
		log.Info("Using binaries found in the download directory: " + binary_file)
		return binary_file, nil
	}

	return binary_file, nil
}