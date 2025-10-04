/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/koooyooo/soi-go/pkg/common/file"
	"github.com/koooyooo/soi-go/pkg/config"
	"github.com/koooyooo/soi-go/pkg/model"
	"github.com/koooyooo/soi-go/pkg/repository"
	"github.com/koooyooo/soi-go/pkg/service"
	"github.com/koooyooo/soi-go/pkg/soiprompt"

	"github.com/c-bata/go-prompt"
	"golang.org/x/net/context"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "soi-go",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		control(cmd, args)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	exec.Command("reset").Run()
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func control(_ *cobra.Command, _ []string) {
	ctx := context.Background()
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
	repo, ok, err := repository.NewRepository(ctx, cfg)
	if err != nil {
		log.Fatalf("failed in creating repository: %v\n", err)
	}
	if !ok {
		log.Fatal("no repository found")
	}
	svc := service.NewService(ctx, cfg.DefaultBucket, repo)
	sp := soiprompt.NewPrompter(cfg, svc, br)

	var basicOpts = []prompt.Option{
		prompt.OptionTitle("soi input"),
		prompt.OptionPrefix("soi> "),
		prompt.OptionMaxSuggestion(15),
	}
	var themedOpts []prompt.Option
	if cfg.Theme == "" || cfg.Theme == "black" {
		themedOpts = blackBgTheme(basicOpts...)
	} else {
		themedOpts = whiteBgTheme(basicOpts...)
	}

	p := prompt.New(
		sp.Execute,
		sp.Complete,
		themedOpts...,
	)
	p.Run()
}

func blackBgTheme(baseOpts ...prompt.Option) []prompt.Option {
	theme := []prompt.Option{
		prompt.OptionPrefixTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.Blue),
		prompt.OptionSelectedSuggestionTextColor(prompt.White),
		prompt.OptionSelectedDescriptionBGColor(prompt.DarkGray),
		prompt.OptionSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionTextColor(prompt.Black),
		prompt.OptionDescriptionBGColor(prompt.LightGray),
		prompt.OptionScrollbarThumbColor(prompt.Black),
		prompt.OptionScrollbarBGColor(prompt.DarkGray),
		prompt.OptionInputTextColor(prompt.LightGray),
		prompt.OptionPreviewSuggestionTextColor(prompt.DarkBlue),
	}
	return append(baseOpts, theme...)
}

func whiteBgTheme(baseOpts ...prompt.Option) []prompt.Option {
	theme := []prompt.Option{
		prompt.OptionPrefixTextColor(prompt.DarkBlue),
		prompt.OptionSelectedSuggestionBGColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionTextColor(prompt.White),
		prompt.OptionSelectedDescriptionBGColor(prompt.DarkGray),
		prompt.OptionSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionTextColor(prompt.Black),
		prompt.OptionDescriptionBGColor(prompt.LightGray),
		prompt.OptionScrollbarThumbColor(prompt.Black),
		prompt.OptionScrollbarBGColor(prompt.DarkGray),
		prompt.OptionInputTextColor(prompt.Black),
		prompt.OptionPreviewSuggestionTextColor(prompt.DarkBlue),
	}
	return append(baseOpts, theme...)
}
