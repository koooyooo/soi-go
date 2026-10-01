package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/koooyooo/soi-go/pkg/common/hash"
	"github.com/koooyooo/soi-go/pkg/config"
	"github.com/koooyooo/soi-go/pkg/model"
)

func Push(cfg *config.Config, bucket *model.Bucket, _ string) error {
	if bucket.IsLocalOnly() {
		return fmt.Errorf("bucket name %s is a local bucket", bucket.Name)
	}
	soisDir, err := bucket.Path()
	if err != nil {
		return err
	}
	var sb model.ServerBucket
	if err := filepath.Walk(soisDir, func(path string, fi os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if fi.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".json") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var s model.SoiData
		if err = json.Unmarshal(b, &s); err != nil {
			return err
		}
		sb.Sois = append(sb.Sois, &s)
		return nil
	}); err != nil {
		return err
	}

	user, pass, headerVal, err := generateAuthValues(cfg)
	if err != nil {
		return err
	}
	userHash, err := hash.Sha1(user + ":" + pass)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/v1/%s/%s/sois:replace", cfg.Server, userHash, bucket.Name),
		strings.NewReader(sb.String()))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", headerVal)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("push failed: %w", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("push failed: status %d", resp.StatusCode)
	}
	fmt.Fprintf(os.Stderr, "pushed\n")
	return nil
}
