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
	"github.com/koooyooo/soi-go/pkg/cli/constant"
	"github.com/koooyooo/soi-go/pkg/soi"
)

func Push(_ string) error {
	bucket := constant.LocalBucket
	if bucket.IsLocalOnly() {
		return fmt.Errorf("Bucketname %s is Local Bucket", bucket.GetName())
	}
	soisDir, err := bucket.Path()
	if err != nil {
		return err
	}
	var sb soi.SoiVirtualBucket
	if err := filepath.Walk(soisDir, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".json") {
			return nil
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		var s soi.SoiData
		if err = json.Unmarshal(b, &s); err != nil {
			return err
		}
		sb.Sois = append(sb.Sois, &soi.SoiVirtual{
			SoiData: &s,
			Path:    strings.TrimPrefix(path, soisDir+"/"),
		})
		return nil
	}); err != nil {
		return err
	}

	// リクエスト作成
	user, _, headerVal, err := generateAuthValues()
	if err != nil {
		return err
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"%s/api/v1/%s/%s/sois:replace", cfg.Server, user, "default"),
		strings.NewReader(sb.String()))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", headerVal)
	req.Header.Add("ContentType", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		fmt.Println("pushed")
	}
	return nil
}
