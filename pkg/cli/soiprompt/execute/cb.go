package execute

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/koooyooo/soi-go/pkg/cli/constant"
)

// cb はbucketの変更を行います
func cb(in string) error {
	flags := flag.NewFlagSet("cb", flag.PanicOnError)
	if err := flags.Parse(strings.Split(in, " ")[1:]); err != nil {
		return err
	}

	bucket := flags.Arg(0)
	// 引数なしの場合は現在のバケット表示
	if bucket == "" {
		fmt.Printf("current bucket: [%s]\n", constant.LocalBucket.GetName())
		return nil
	}
	buckets, err := constant.ListBuckets()
	if err != nil {
		return err
	}

	// 存在するBucketなら変更
	for _, b := range buckets {
		if b == bucket {
			constant.LocalBucket.SetName(bucket)
			fmt.Printf("change current bucket: %s \n", bucket)
			return nil
		}
	}
	soisDir, err := constant.SoisDir()
	if err != nil {
		return err
	}
	if err := os.Mkdir(filepath.Join(soisDir, bucket), 0700); err != nil {
		return err
	}
	constant.LocalBucket.SetName(bucket)
	fmt.Printf("create & change current bucket: %s \n", bucket)
	return nil
}
