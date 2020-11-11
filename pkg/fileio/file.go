/*
Package fileio offers file related functions
*/
package fileio

import (
	"io"
	"os"
)

// Exists returns true if the file or dir exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func IsDir(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return f.IsDir(), nil
}
