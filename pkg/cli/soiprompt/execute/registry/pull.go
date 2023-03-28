package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/koooyooo/soi-go/pkg/cli/config"
	"github.com/koooyooo/soi-go/pkg/common/file"
	"github.com/koooyooo/soi-go/pkg/common/hash"
	"github.com/koooyooo/soi-go/pkg/model"
)

// Pull はレジストリよりデータをロードします
func Pull(cfg *config.Config, bucket *model.Bucket, _ string) error {
	soisDir, err := bucket.Path()
	if err != nil {
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
		"GET",
		fmt.Sprintf("%s/api/v1/%s/%s/sois", cfg.Server, userHash, bucket.Name),
		nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", headerVal)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("pull response not OK: %d", resp.StatusCode)
	}

	// バックアップディレクトリを作成
	empty, err := file.IsEmpty(soisDir)
	if err != nil {
		return err
	}
	if !empty {
		if err := os.RemoveAll(soisDir + ".bk"); err != nil {
			return err
		}
		if err := os.Rename(soisDir, soisDir+".bk"); err != nil {
			return err
		}
	}

	// レスポンス処理
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var sb model.ServerBucket
	if err = json.Unmarshal(b, &sb); err != nil {
		return err
	}
	if !file.Exists(soisDir) {
		if err := os.Mkdir(soisDir, 0700); err != nil {
			return err
		}
	}
	for _, sv := range sb.Sois {
		b, err := json.Marshal(sv)
		if err != nil {
			return err
		}
		dir, file := path.Split(sv.FilePath(bucket.Name))
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, file), b, 0644); err != nil {
			return err
		}
	}
	fmt.Println("pulled")
	return nil
}
