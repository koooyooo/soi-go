package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/koooyooo/soi-go/pkg/cli/config"
	"github.com/koooyooo/soi-go/pkg/cli/loader"
	"github.com/koooyooo/soi-go/pkg/cli/repository"
	"github.com/koooyooo/soi-go/pkg/model"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
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
	dbPath := bucket + ".db"
	os.RemoveAll(dbPath)
	repo, err := repository.NewSQLiteRepository(dbPath)
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