package execute

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/koooyooo/soi-go/pkg/model"
)

// tag はsoiのタグ付けを行います
func (e *Executor) tag(in string) error {
	ctx := context.Background()
	flags := flag.NewFlagSet("tag", flag.ContinueOnError)
	if err := flags.Parse(strings.Split(in, " ")[1:]); err != nil {
		return err
	}

	args := flags.Args()
	if len(args) == 0 {
		return fmt.Errorf("usage: tag <hash> [#tag...]")
	}
	hash := args[0]
	tags := args[1:]

	sois := e.Cache.ListSoiCache
	var tgt *model.SoiData
	for _, s := range sois {
		if s.Hash == hash || strings.HasPrefix(s.Hash, hash) {
			tgt = s
			break
		}
	}
	if tgt == nil {
		return fmt.Errorf("soi not found: %s", hash)
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
