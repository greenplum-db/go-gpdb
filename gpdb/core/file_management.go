package core

import (
	"os"
	"io"
)

func CreateFile(path string) error {
	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil { return err}
		defer file.Close()
	}

	return nil
}

func  WriteFile(path string, args []string) error {
	// open file using READ & WRITE permission
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil { return err}
	defer file.Close()

	// write some text to file
	for _, k := range args {
		text := string(k)
		_, err = file.WriteString(text + "\n")
		if err != nil { return err}
	}

	// save changes
	err = file.Sync()
	if err != nil { return err}

	return nil
}

func AppendFile(path string, args []string) error {

	// open file using READ & APPEND permission
	var file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil { return err}
	defer file.Close()

	// append to the file
	for _, a := range args {

		if _, err = file.WriteString(a + "\n"); err != nil {
			return err
		}

	}

	// save changes
	err = file.Sync()
	if err != nil { return err}

	return nil
}

func ReadFile(path string) ([]byte, error) {
	// re-open file
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil { return []byte(""), err}
	defer file.Close()

	// read file
	var text = make([]byte, 1024)
	for {
		n, err := file.Read(text)
		if err != io.EOF {
			if err != nil { return []byte(""), err}
		}
		if n == 0 {
			break
		}
	}
	if err != nil { return []byte(""), err }

	return text, nil
}

func DeleteFile(path string) error {

	// delete file
	var err = os.Remove(path)
	if err != nil { return err}

	return nil
}