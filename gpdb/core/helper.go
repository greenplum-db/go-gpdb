package core

import (
	"strings"
	"os"
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
func DoesFileOrDirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}