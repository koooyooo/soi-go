package execute

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/koooyooo/soi-go/pkg/cli/constant"
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/common"
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/utils"
	"github.com/koooyooo/soi-go/pkg/soi"
)

// add はsoiの追加を行います
func add(in string) error {
	flags := flag.NewFlagSet("add", flag.PanicOnError)
	n := flags.String("n", "", "name of the uri")
	d := flags.String("d", "new", "soiRoot to which soi store")
	err := flags.Parse(strings.Split(in, " ")[1:])
	if err != nil {
		return err
	}

	uri := flags.Arg(0)

	name := *n
	if name == "" {
		title, ok, err := utils.ParseTitleByURL(uri)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("no name(title) found in the url: %s\nplease specify the name with -n option\n", uri)
		}
		name = title
	}

	s := soi.SoiData{
		Name:      name,
		URI:       uri,
		Hash:      common.Hash(uri),
		Tags:      []string{},
		CreatedAt: time.Now(), // .Format("2006-01-02T15:04:05Z07:00"),
	}
	b, err := json.Marshal(&s)
	if err != nil {
		return err
	}
	soiRoot, err := constant.LocalBucket.Path()
	if err != nil {
		return err
	}
	baseDir := filepath.Join(soiRoot, *d)
	if err = os.MkdirAll(baseDir, 0700); err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(baseDir, common.ToStorableName(name)), b, 0600)
}
