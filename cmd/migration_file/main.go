package main

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
	"soi-go/pkg/cli/constant"
	"soi-go/pkg/cli/repository"
	"soi-go/pkg/common/file"
	"soi-go/pkg/model"
	"strings"
)

func main() {
	bucket := os.Args[1]
	fmt.Printf("bucket: [%s]\n", bucket)
	if bucket == "" {
		log.Fatal("bucket is required")
	}
	if err := migrateFile(bucket); err != nil {
		log.Fatal(err)
	}
}

func migrateFile(bucket string) error {
	ctx := context.Background()
	repo, err := repository.NewSQLiteRepository(ctx, bucket+".db", bucket)
	if err != nil {
		return err
	}
	sois, err := repo.LoadAll(context.Background(), bucket)
	if err != nil {
		return err
	}
	soisDir, err := constant.SoisDir()
	if err != nil {
		return err
	}
	for _, s := range sois {
		fmt.Println(s.Path) // TODO
		if err := storeSoi(soisDir, bucket, s); err != nil {
			return err
		}
	}
	return nil
}

func storeSoi(soiRoot, bucket string, s *model.SoiData) error {
	b, err := json.Marshal(&s)
	if err != nil {
		return err
	}
	baseDir := filepath.Join(soiRoot, bucket, strings.TrimPrefix(s.Path, soiRoot+"/"+bucket+"/"))
	if err = os.MkdirAll(baseDir, 0700); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(baseDir, file.ToStorableName(s.Name)), b, 0600)
}
