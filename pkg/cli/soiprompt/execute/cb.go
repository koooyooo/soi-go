package execute

import (
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
	"strings"

	"soi-go/pkg/model"

	"soi-go/pkg/cli/constant"
)

// cb はbucketの変更を行います
func (e *Executor) cb(in string) error {
	ctx := context.Background()
	flags := flag.NewFlagSet("cb", flag.PanicOnError)
	if err := flags.Parse(strings.Split(in, " ")[1:]); err != nil {
		return err
	}

	bucketName := flags.Arg(0)
	// 引数なしの場合は現在のバケット表示
	if bucketName == "" {
		fmt.Printf("current bucket: [%s]\n", e.BucketRef.Bucket.Name)
		return nil
	}
	buckets, err := model.ListBuckets()
	if err != nil {
		return err
	}

	// 存在するBucketなら変更
	for _, b := range buckets {
		if b.Name == bucketName {
			e.Cache.Clear()
			e.Bucket = b
			if err := e.Service.ChangeBucket(ctx, bucketName); err != nil {
				return err
			}
			fmt.Printf("change current bucket: %s \n", b.Name)
			return nil
		}
	}
	soisDir, err := constant.SoisDir()
	if err != nil {
		return err
	}
	if err := os.Mkdir(filepath.Join(soisDir, bucketName), 0700); err != nil {
		return err
	}

	b, err := model.NewBucket(bucketName)
	if err != nil {
		return err
	}
	e.Cache.Clear()
	e.Bucket = b

	if err := e.Service.ChangeBucket(ctx, bucketName); err != nil {
		return err
	}
	fmt.Printf("create & change current bucket: %s \n", b.Name)
	return nil
}
