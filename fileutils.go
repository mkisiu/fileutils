package fileutils

import (
	"encoding/csv"
	"fmt"
	"io"
	//"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileList returns slice of strings, containing filenames base on prefix and suffix
func FileList(path, prefix, suffix string) ([]string, error) {
	var listOfFiles []string

	files, err := os.ReadDir(path)
	if err != nil {
		return listOfFiles, err
	}

	//files, err := ioutil.ReadDir(path)
	//if err != nil {
	//	return listOfFiles, err
	//}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), prefix) && strings.HasSuffix(file.Name(), suffix) {
			listOfFiles = append(listOfFiles, file.Name())
		}
	}

	return listOfFiles, nil
}

func SubFileList(path, prefix, suffix string) ([]string, error) {
	var listOfFiles []string

	err := filepath.Walk(path,
		func(folderAndPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasPrefix(info.Name(), prefix) && strings.HasSuffix(info.Name(), suffix) {
				listOfFiles = append(listOfFiles, folderAndPath)
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return listOfFiles, nil
}

// CopyFile copies file from source location to destination.
func CopyFileDepr(src, dst string) error {
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

// CopyFile waits until src looks stable (size stops changing) and then copies it to dst.
// Simple guard against processing a still-growing file. No atomic rename involved.
func CopyFile(src, dst string) error {
	const attempts = 5
	const settle = 500 * time.Millisecond

	if !waitStableSize(src, attempts, settle) {
		return fmt.Errorf("source not stable: %s", src)
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst) // overwrite if exists
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	// Optional: _ = out.Sync()
	_ = out.Sync()

	return nil
}

// CsvToMap converts the csv data (with header) into map, with header parts as keys
func CsvToMap(fileName string, separator rune) ([]map[string]string, error) {
	data, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	r := csv.NewReader(data)
	r.Comma = separator
	rows := []map[string]string{}
	var header []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if header == nil {
			header = record
		} else {
			dict := map[string]string{}
			for i := range header {
				dict[header[i]] = record[i]
			}
			rows = append(rows, dict)
		}
	}
	return rows, nil
}

// waitStableSize returns true when the file size stops changing.
// It samples size up to `attempts` times with `settle` delay.
// If two consecutive samples are equal, the file is considered stable.
func waitStableSize(path string, attempts int, settle time.Duration) bool {
	if attempts < 3 {
		attempts = 3
	}

	var prev int64 = -1
	for i := 0; i < attempts; i++ {
		fi, err := os.Stat(path)
		if err != nil || !fi.Mode().IsRegular() {
			return false
		}
		size := fi.Size()

		// Not first sample and size didn’t change → stable
		if i > 0 && size == prev {
			return true
		}
		prev = size
		time.Sleep(settle)
	}
	return false
}
