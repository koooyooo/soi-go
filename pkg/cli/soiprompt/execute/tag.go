package execute

import (
	"context"
	"flag"
	"soi-go/pkg/model"
	"strings"
)

// tag はsoiのタグ付けを行います
func (e *Executor) tag(in string) error {
	ctx := context.Background()
	flags := flag.NewFlagSet("tag", flag.PanicOnError)
	if err := flags.Parse(strings.Split(in, " ")[1:]); err != nil {
		return err
	}

	args := flags.Args()
	hash := args[0]
	tags := args[1:]

	sois := e.Cache.ListSoiCache
	var tgt *model.SoiData
	for _, s := range sois {
		if s.Hash == hash {
			tgt = s
		}
	}
	tgt.Tags = removeTagHead(tags)
	// TODO 既存タグを全置き換え、KVタグは無視してしまっている

	if err := e.Service.Store(ctx, tgt); err != nil {
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
