package main

import (
	"github.com/koooyooo/soi-go/pkg/cli/config"
	"log"
	"os"

	"github.com/koooyooo/soi-go/pkg/cli/soiprompt"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/common/file"
	"github.com/koooyooo/soi-go/pkg/model"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed in loading config: %v", err)
	}
	b, err := model.NewBucket(cfg.DefaultBucket)
	if err != nil {
		log.Fatalf("failed in creating default bucket: %v", err)
	}
	soisDir, err := b.Path()
	if err != nil {
		log.Fatalf("failed in getting soisdir: %v", err)
	}
	if !file.Exists(soisDir) {
		if err := os.MkdirAll(soisDir, 0700); err != nil {
			log.Fatalf("failed in creating sois dir: %v", err)
		}
	}

	br := &model.BucketRef{
		Bucket: b,
	}
	sp := soiprompt.NewPrompter(cfg, br)

	p := prompt.New(
		sp.Execute,
		sp.Complete,
		prompt.OptionTitle("soi input"),
		prompt.OptionPrefix("soi> "),
		prompt.OptionPrefixTextColor(prompt.Black),
		//prompt.OptionSelectedSuggestionBGColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionBGColor(prompt.Blue),
		prompt.OptionSelectedSuggestionTextColor(prompt.White),
		prompt.OptionSelectedDescriptionBGColor(prompt.DarkGray),
		prompt.OptionSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionTextColor(prompt.Black),
		prompt.OptionDescriptionBGColor(prompt.LightGray),
		prompt.OptionScrollbarThumbColor(prompt.Black),
		prompt.OptionScrollbarBGColor(prompt.DarkGray),
		prompt.OptionInputTextColor(prompt.Black),
		prompt.OptionPreviewSuggestionTextColor(prompt.DarkBlue),
		prompt.OptionMaxSuggestion(15),
	)
	p.Run()
}
