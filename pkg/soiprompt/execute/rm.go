package execute

import (
	"context"
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
	parts := strings.Fields(in)
	if len(parts) < 2 {
		fmt.Println("usage: rm <path>")
		return nil
	}
	// パス内の空白を許容するため、コマンド名以降を結合する
	if err = flags.Parse([]string{strings.Join(parts[1:], " ")}); err != nil {
		return err
	}
	relDir := flags.Arg(0)
	if relDir == "" {
		fmt.Println("cannot delete bucket dir.")
		return nil
	}

	sois, loadErr := e.Service.LoadAll(context.Background())
	var target string
	if loadErr != nil {
		target, err = pathutil.ResolveExistingUnderRoot(baseDir, relDir)
	} else {
		target, err = pathutil.ResolveExistingUnderRootWithSois(baseDir, relDir, sois)
	}
	if err != nil {
		fmt.Printf("%v (bucket: %s)\n", err, e.Bucket.Name)
		return nil
	}
	if err := os.RemoveAll(target); err != nil {
		return err
	}
	e.Cache.Clear()
	return nil
}
