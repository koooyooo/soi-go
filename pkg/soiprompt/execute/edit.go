package execute

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func (e *Executor) edit(in string) error {
	flags := flag.NewFlagSet("edit", flag.ContinueOnError)
	if err := flags.Parse(strings.Split(in, " ")[1:]); err != nil {
		return err
	}
	s, err := findSoi(e.Cache.ListSoiCache, flags.Args())
	if err != nil {
		return err
	}
	bucketPath, err := e.Bucket.Path()
	if err != nil {
		return err
	}
	path := s.FilePath(bucketPath)

	var cName string
	var cArgs []string
	switch runtime.GOOS {
	case "darwin", "linux", "freebsd":
		cName = "vim"
		cArgs = []string{path}
	case "windows":
		cName = "cmd"
		cArgs = []string{"/c", "start", "notepad.exe", path}
	default:
		return fmt.Errorf("unsupported os: %s", runtime.GOOS)
	}
	c := exec.Command(cName, cArgs...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
