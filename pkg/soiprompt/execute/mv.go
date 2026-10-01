package execute

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/koooyooo/soi-go/pkg/common/file"
	"github.com/koooyooo/soi-go/pkg/common/pathutil"
	"github.com/koooyooo/soi-go/pkg/model"
)

// mv はsoiの移動を行います
func (e *Executor) mv(in string) error {
	baseDir, err := e.Bucket.Path()
	if err != nil {
		return err
	}
	flags := flag.NewFlagSet("mv", flag.ContinueOnError)
	if err = flags.Parse(strings.Split(in, " ")[1:]); err != nil {
		return err
	}
	if flags.NArg() < 2 {
		return flag.ErrHelp
	}

	from, err := pathutil.ResolveExistingUnderRoot(baseDir, flags.Arg(0))
	if err != nil {
		return fmt.Errorf("mv from: %w", err)
	}

	to, err := resolveMvDestination(baseDir, flags.Arg(1), from)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(to), 0700); err != nil {
		return err
	}
	if err := os.Rename(from, to); err != nil {
		return err
	}
	if strings.HasSuffix(to, ".json") {
		if err := rewriteSoiLocation(baseDir, to); err != nil {
			return err
		}
	}
	e.Cache.Clear()
	return nil
}

func resolveMvDestination(baseDir, toArg, fromAbs string) (string, error) {
	if strings.HasSuffix(toArg, "/") {
		toDir, err := pathutil.ResolveUnderRoot(baseDir, toArg)
		if err != nil {
			return "", err
		}
		if err := os.MkdirAll(toDir, 0700); err != nil {
			return "", err
		}
		return filepath.Join(toDir, filepath.Base(fromAbs)), nil
	}

	if existing, err := pathutil.ResolveExistingUnderRoot(baseDir, toArg); err == nil {
		if ok, dirErr := file.IsDir(existing); dirErr == nil && ok {
			return filepath.Join(existing, filepath.Base(fromAbs)), nil
		}
		return existing, nil
	}

	return pathutil.ResolveDestUnderRoot(baseDir, toArg)
}

func rewriteSoiLocation(bucketRoot, absFile string) error {
	b, err := os.ReadFile(absFile)
	if err != nil {
		return err
	}
	var s model.SoiData
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	rel, err := filepath.Rel(bucketRoot, absFile)
	if err != nil {
		return err
	}
	rel = filepath.ToSlash(rel)
	dir := filepath.ToSlash(filepath.Dir(rel))
	if dir == "." {
		dir = ""
	}
	s.Path = dir
	s.Name = strings.TrimSuffix(filepath.Base(rel), ".json")
	ub, err := json.Marshal(&s)
	if err != nil {
		return err
	}
	return os.WriteFile(absFile, ub, 0600)
}
