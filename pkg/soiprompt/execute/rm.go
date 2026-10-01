package execute

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/koooyooo/soi-go/pkg/common/pathutil"
)

// rm はsoiの削除を行います
func (e *Executor) rm(in string) error {
	baseDir, err := e.Bucket.Path()
	if err != nil {
		return err
	}
	flags := flag.NewFlagSet("rm", flag.ContinueOnError)
	if err = flags.Parse(strings.Split(in, " ")[1:]); err != nil {
		return err
	}
	relDir := flags.Arg(0)
	if relDir == "" {
		fmt.Println("cannot delete bucket dir.")
		return nil
	}
	target, err := pathutil.ResolveExistingUnderRoot(baseDir, relDir)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	if err := os.RemoveAll(target); err != nil {
		return err
	}
	e.Cache.Clear()
	return nil
}
