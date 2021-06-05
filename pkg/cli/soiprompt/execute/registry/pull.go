package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	config2 "github.com/koooyooo/soi-go/pkg/cli/config"

	"github.com/koooyooo/soi-go/pkg/cli/constant"
	"github.com/koooyooo/soi-go/pkg/fileio"
	"github.com/koooyooo/soi-go/pkg/soi"
)

// pull はレジストリよりデータをロードします
func Pull(_ string) error {
	soisDir, err := constant.LocalBucket.Path()
	if err != nil {
		return err
	}

	// リクエスト作成
	user, _, headerVal, err := generateAuthValues()
	if err != nil {
		return err
	}
	cfg, err := config2.LoadConfig()
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api/v1/%s/%s/sois", cfg.Server, user, "default"),
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
	empty, err := fileio.IsEmpty(soisDir)
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
	var sb soi.SoiVirtualBucket
	if err = json.Unmarshal(b, &sb); err != nil {
		return err
	}
	if !fileio.Exists(soisDir) {
		if err := os.Mkdir(soisDir, 0700); err != nil {
			return err
		}
	}
	for _, sv := range sb.Sois {
		fmt.Println("load:", sv.Path)
		b, err := json.Marshal(sv.SoiData)
		if err != nil {
			return err
		}
		relDir, file := path.Split(sv.Path)
		fullDir := filepath.Join(soisDir, relDir)
		if err := os.MkdirAll(fullDir, 0700); err != nil {
			return err
		}
		if err := ioutil.WriteFile(filepath.Join(fullDir, file), b, 0644); err != nil {
			return err
		}
	}
	fmt.Println("pulled")
	return nil
}
