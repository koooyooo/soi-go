package execute

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func (e *Executor) edit(in string) error {
	flags := flag.NewFlagSet("edit", flag.PanicOnError)
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

	var cName, cArgs string
	switch runtime.GOOS {
	case "darwin", "linux", "freebsd":
		cName = "vim"
		cArgs = path
	case "windows":
		cName = "cmd"
		cArgs = "/c start notepad.exe " + path
	}
	c := exec.Command(cName, cArgs)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		log.Fatal(err)
	}

	return nil
}
