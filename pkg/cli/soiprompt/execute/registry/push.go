package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/koooyooo/soi-go/pkg/cli/config"
	"github.com/koooyooo/soi-go/pkg/common/hash"
	"github.com/koooyooo/soi-go/pkg/model"
)

func Push(cfg *config.Config, bucket *model.Bucket, _ string) error {
	if bucket.IsLocalOnly() {
		return fmt.Errorf("Bucketname %s is Local Bucket", bucket.Name)
	}
	soisDir, err := bucket.Path()
	if err != nil {
		return err
	}
	var sb model.ServerBucket
	if err := filepath.Walk(soisDir, func(path string, fi os.FileInfo, _ error) error {
		if fi.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".json") {
			return nil
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		var s model.SoiData
		if err = json.Unmarshal(b, &s); err != nil {
			return err
		}
		s.Path = strings.TrimPrefix(path, soisDir)
		sb.Sois = append(sb.Sois, &s)
		return nil
	}); err != nil {
		return err
	}

	// リクエスト作成
	user, pass, headerVal, err := generateAuthValues()
	if err != nil {
		return err
	}
	userHash, err := hash.Sha1(user + ":" + pass)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"%s/api/v1/%s/%s/sois:replace", cfg.Server, userHash, bucket.Name),
		strings.NewReader(sb.String()))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", headerVal)
	req.Header.Add("ContentType", "application/json")
	//go func() {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("push failed: %s", err)
	}
	if resp.StatusCode == http.StatusOK {
		fmt.Fprintf(os.Stderr, "pushed\n")
	}
	//}()
	return nil
}
