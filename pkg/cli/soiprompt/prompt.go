// Package soiprompt は `complete`,`execute` の責務を別パッケージで管理できる Prompterを提供
// Executor, Completer は共に `*model.Bucket` を参照する `*model.BucketRef` を保持し、同一の参照を共有する
// これは片系統での Bucket入れ替えを検知可能にするため
package soiprompt

import (
	"github.com/koooyooo/soi-go/pkg/cli/config"
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/complete"
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/execute"
	"github.com/koooyooo/soi-go/pkg/model"
)

type Prompter struct {
	execute.Executor
	complete.Completer
}

func NewPrompter(conf *config.Config, b *model.BucketRef) *Prompter {
	return &Prompter{
		Executor: execute.Executor{
			Conf:      conf,
			BucketRef: b,
		},
		Completer: complete.Completer{
			BucketRef: b,
		},
	}
}