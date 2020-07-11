/*
Package fileio offers file related functions
*/
package fileio

import "os"

// FileExists returns true if the file exists.
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func IsDir(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return f.IsDir(), nil
}
