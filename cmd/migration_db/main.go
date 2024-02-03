package main

import (
	"context"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"soi-go/pkg/cli/config"
	"soi-go/pkg/cli/loader"
	"soi-go/pkg/cli/repository"
	"soi-go/pkg/model"
)

func main() {
	flag.Parse()
	bucket := flag.Arg(0)
	os.RemoveAll(bucket + ".db")

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
	sois, err := loader.LoadSois(soisDir)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	os.RemoveAll("" + bucket + ".db")
	repo, err := repository.NewSQLiteRepository(ctx, "", bucket)
	if err != nil {
		log.Fatal(err)
	}
	if err := repo.Init(ctx); err != nil {
		log.Fatal(err)
	}
	for _, s := range sois {
		repo.Store(ctx, bucket, s)
	}
	fmt.Println("done")
}
