package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/koooyooo/soi-go/pkg/cli/constant"
	"github.com/koooyooo/soi-go/pkg/cli/repository"
	"github.com/koooyooo/soi-go/pkg/common/file"
	"github.com/koooyooo/soi-go/pkg/model"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
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
	repo, err := repository.NewSQLiteRepository(bucket + ".db")
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
