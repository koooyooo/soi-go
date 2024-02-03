package execute

import (
	"flag"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	fileio2 "soi-go/pkg/common/file"
)

// rm はsoiの削除を行います
func (e *Executor) rm(in string) error {
	baseDir, err := e.Bucket.Path()
	if err != nil {
		return err
	}
	flags := flag.NewFlagSet("rm", flag.PanicOnError)
	if err = flags.Parse(strings.Split(in, " ")[1:]); err != nil {
		return err
	}
	relDir := flags.Arg(0)
	if relDir == "" {
		fmt.Println("cannot delete bucket dir.")
		return nil
	}
	target := filepath.Join(baseDir, relDir)
	if !fileio2.Exists(target) {
		fmt.Println("No file or dir found.")
		return nil
	}
	return exec.Command("rm", "-rf", target).Start()
}
