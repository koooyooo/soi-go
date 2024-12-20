package execute

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/koooyooo/soi-go/pkg/common/file"
)

// mv はsoiの移動を行います
func (e *Executor) mv(in string) error {
	baseDir, err := e.Bucket.Path()
	if err != nil {
		return err
	}
	flags := flag.NewFlagSet("mv", flag.PanicOnError)
	if err = flags.Parse(strings.Split(in, " ")[1:]); err != nil {
		return err
	}
	from := filepath.Join(baseDir, flags.Arg(0))
	to := filepath.Join(baseDir, flags.Arg(1))

	// 移動先のディレクトリが存在しない場合は作成
	toDir := to[0:strings.LastIndex(to, "/")]
	if !file.Exists(toDir) {
		err = os.MkdirAll(toDir, 0700)
		if err != nil {
			return err
		}
	}

	// 末尾JSONの付与
	toIsDir, err := file.IsDir(to)
	if err != nil {
		switch err.(type) {
		case *os.PathError:
		default:
			return err
		}
	}
	if !toIsDir && !strings.HasSuffix(to, ".json") {
		to = to + ".json"
	}
	return exec.Command("mv", from, to).Start()
}
