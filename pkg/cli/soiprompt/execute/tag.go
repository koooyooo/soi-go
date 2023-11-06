package execute

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"soi-go/pkg/cli/loader"
)

// tag はsoiのタグ付けを行います
func (e *Executor) tag(in string) error {
	baseDir, err := e.Bucket.Path()
	if err != nil {
		return err
	}
	flags := flag.NewFlagSet("rm", flag.PanicOnError)
	if err = flags.Parse(strings.Split(in, " ")[1:]); err != nil {
		return err
	}

	args := flags.Args()
	tags := args[1:]

	target := filepath.Join(baseDir, args[0])
	if !loader.Exists(target) {
		fmt.Println("No file or dir found.")
		return nil
	}

	soi, err := loader.LoadSoiData(target)
	soi.Tags = removeTagHead(tags)
	// TODO 既存タグを全置き換え、KVタグは無視してしまっている

	if err := loader.StoreSoiData(target, soi); err != nil {
		return err
	}
	return nil
}

func removeTagHead(tags []string) []string {
	var nonHead []string
	for _, tag := range tags {
		nonHead = append(nonHead, strings.TrimLeft(tag, "#"))
	}
	return nonHead
}
