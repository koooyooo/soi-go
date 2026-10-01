package execute

import (
	"encoding/json"
	"flag"
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

	from, err := pathutil.ResolveUnderRoot(baseDir, flags.Arg(0))
	if err != nil {
		return err
	}
	to, err := pathutil.ResolveUnderRoot(baseDir, flags.Arg(1))
	if err != nil {
		return err
	}

	toDir := filepath.Dir(to)
	toIsDir := false
	if file.Exists(to) {
		toIsDir, err = file.IsDir(to)
		if err != nil {
			return err
		}
	} else if strings.HasSuffix(flags.Arg(1), "/") {
		toIsDir = true
		toDir = to
	}

	if toIsDir {
		if err := os.MkdirAll(to, 0700); err != nil {
			return err
		}
		to = filepath.Join(to, filepath.Base(from))
	} else {
		if err := os.MkdirAll(toDir, 0700); err != nil {
			return err
		}
		if !strings.HasSuffix(to, ".json") && strings.HasSuffix(from, ".json") {
			to = to + ".json"
		}
	}

	if err := os.Rename(from, to); err != nil {
		return err
	}

	// 単一 JSON なら内部 Path / Name を実ファイル位置に合わせて書き戻す
	if strings.HasSuffix(to, ".json") {
		if err := rewriteSoiLocation(baseDir, to); err != nil {
			return err
		}
	}
	e.Cache.Clear()
	return nil
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
