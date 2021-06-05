package execute

import (
	"flag"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	fileio2 "github.com/koooyooo/soi-go/pkg/common/file"

	"github.com/koooyooo/soi-go/pkg/cli/constant"
)

// rm はsoiの削除を行います
func rm(in string) error {
	baseDir, err := constant.LocalBucket.Path()
	if err != nil {
		return err
	}
	flags := flag.NewFlagSet("rm", flag.PanicOnError)
	if err = flags.Parse(strings.Split(in, " ")[1:]); err != nil {
		return err
	}

	target := filepath.Join(baseDir, flags.Arg(0))
	if !fileio2.Exists(target) {
		fmt.Println("No file or dir found.")
		return nil
	}
	return exec.Command("rm", "-rf", target).Start()
}
