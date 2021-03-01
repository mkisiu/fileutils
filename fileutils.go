package fileutils

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// FileList returns slice of strings, containing filenames base on prefix and suffix
func FileList(path, prefix, suffix string) ([]string, error) {
	var listOfFiles []string

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return listOfFiles, err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), prefix) && strings.HasSuffix(file.Name(), suffix) {
			listOfFiles = append(listOfFiles, file.Name())
		}
	}

	return listOfFiles, nil
}

// CopyFile copies file from source location to destination.
func CopyFile(src, dst string) error {
	// Open the source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copying
	_, err = io.Copy(dstFile, srcFile)
	return err
}

// MoveFile moves a file from source location to destination.
func MoveFile(src, dst string) error {
	err := os.Rename(src, dst)
	if err != nil {
		return err
	}
	return err
}
