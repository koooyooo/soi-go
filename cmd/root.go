/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/c-bata/go-prompt"
	"golang.org/x/net/context"
	"log"
	"os"
	"soi-go/pkg/cli/config"
	"soi-go/pkg/cli/repository"
	"soi-go/pkg/cli/service"
	"soi-go/pkg/cli/soiprompt"
	"soi-go/pkg/common/file"
	"soi-go/pkg/model"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "soi-go",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		control(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.soi-go.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
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
