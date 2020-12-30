package soiprompt

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/koooyooo/soi-go/pkg/config"

	"github.com/koooyooo/soi-go/pkg/cli/constant"

	"github.com/koooyooo/soi-go/pkg/soi"

	"github.com/koooyooo/soi-go/pkg/fileio"
)

// Executor は入力されたコマンドに応じた処理を行う
func Executor(in string) {
	in = strings.Trim(in, " ")
	cmd := strings.Split(in, " ")[0]
	switch cmd {
	case "add", "a":
		if err := add(in); err != nil {
			fmt.Println(err)
			return
		}
	case "mv":
		if err := mv(in); err != nil {
			fmt.Println(err)
			return
		}
	case "rm":
		if err := rm(in); err != nil {
			fmt.Println(err)
			return
		}
	case "open", "o", "list", "l", "dig", "d":
		if err := open(in); err != nil {
			fmt.Println(err)
			return
		}
	case "help", "h":
		if err := help(in); err != nil {
			fmt.Println(err)
			return
		}
	case "pull":
		if err := pull(in); err != nil {
			fmt.Println(err)
			return
		}
	case "push":
		if err := push(in); err != nil {
			fmt.Println(err)
			return
		}
	case "quit", "q", "exit":
		if err := quit(in); err != nil {
			fmt.Println(err)
			return
		}
	}
}

// add はsoiの追加を行う
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
		title, ok, err := parseTitleByURL(uri)
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
		Tags:      []string{},
		CreatedAt: time.Now(), // .Format("2006-01-02T15:04:05Z07:00"),
	}
	b, err := json.Marshal(&s)
	if err != nil {
		return err
	}
	soiRoot, err := fileio.SoisDirPath(constant.BucketName())
	if err != nil {
		return err
	}
	baseDir := filepath.Join(soiRoot, *d)
	if err = os.MkdirAll(baseDir, 0700); err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(baseDir, toStorableName(name)), b, 0600)
}

// mv はsoiの移動を行う
func mv(in string) error {
	baseDir, err := fileio.SoisDirPath(constant.BucketName())
	if err != nil {
		return err
	}
	flags := flag.NewFlagSet("mv", flag.PanicOnError)
	if err = flags.Parse(strings.Split(in, " ")[1:]); err != nil {
		return err
	}
	from := filepath.Join(baseDir, flags.Arg(0))
	to := filepath.Join(baseDir, flags.Arg(1))

	// 移動先のディレクトリが存在しない場合は作成
	toDir := to[0:strings.LastIndex(to, "/")]
	if !fileio.Exists(toDir) {
		err = os.MkdirAll(toDir, 0700)
		if err != nil {
			return err
		}
	}

	// 末尾JSONの付与
	toIsDir, err := fileio.IsDir(to)
	if err != nil {
		switch err.(type) {
		case *os.PathError:
		default:
			return err
		}
	}
	if !toIsDir && !strings.HasSuffix(to, ".json") {
		to = to + ".json"
	}
	return exec.Command("mv", from, to).Start()
}

// rm はsoiの削除を行う
func rm(in string) error {
	baseDir, err := fileio.SoisDirPath(constant.BucketName())
	if err != nil {
		return err
	}
	flags := flag.NewFlagSet("rm", flag.PanicOnError)
	if err = flags.Parse(strings.Split(in, " ")[1:]); err != nil {
		return err
	}

	target := filepath.Join(baseDir, flags.Arg(0))
	if !fileio.Exists(target) {
		fmt.Println("No file or dir found.")
		return nil
	}
	return exec.Command("rm", "-rf", target).Start()
}

// open は指定されたSoiを元にブラウザを開きます
func open(in string) error {
	flags := flag.NewFlagSet("open", flag.PanicOnError)
	firefox := flags.Bool("f", false, "use firefox")
	safari := flags.Bool("s", false, "use safari")
	if err := flags.Parse(strings.Split(in, " ")[1:]); err != nil {
		return err
	}

	// Soiファイルを特定
	soisDir, err := fileio.SoisDirPath(constant.BucketName())
	if err != nil {
		return err
	}

	relPath := addJSON(removeHeader(filepath.Join(flags.Args()...)))

	fullPath := filepath.Join(soisDir, relPath)

	// Soiファイルを読み込み
	b, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}
	var s soi.SoiData
	err = json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	// 閲覧履歴を追記
	s.NumViews++

	// 利用ログを記載
	s.UsageLogs = append(s.UsageLogs, soi.UsageLog{
		Type:   soi.UsageTypeOpen,
		UsedAt: time.Now(),
	})

	// Soiファイルを再登録
	ub, err := json.Marshal(s)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(fullPath, ub, 0600); err != nil {
		return err
	}

	// 環境設定に応じてブラウザオープン
	if *firefox {
		return exec.Command("open", "-a", "Firefox", s.URI).Start()
	}
	if *safari {
		return exec.Command("open", "-a", "Safari", s.URI).Start()
	}
	return exec.Command("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", s.URI).Start()
}

func addJSON(path string) string {
	if !strings.HasSuffix(path, ".json") {
		path = path + ".json"
	}
	return path
}

func pull(_ string) error {
	soisDir, err := fileio.SoisDirPath(constant.BucketName())
	if err != nil {
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

func generateAuthValues() (user string, pass string, authValue string, err error) {
	user = constant.EnvKeySoiUserName.Get()
	if user == "" {
		fmt.Println("username?")
		fmt.Print("> ")
		fmt.Scan(&user)
		if user == "" {
			return "", "", "", fmt.Errorf("not environment variable found: %s", constant.EnvKeySoiUserName)
		}
	}
	pass = constant.EnvKeySoiUserPass.Get()
	if pass == "" {
		fmt.Println("password?")
		fmt.Print("> ")
		fmt.Scan(&pass)
		if pass == "" {
			return "", "", "", fmt.Errorf("not environment variable found: %s", constant.EnvKeySoiUserPass)
		}
	}
	return user, pass, "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+pass)), nil
}

func push(_ string) error {
	soisDir, err := fileio.SoisDirPath(constant.BucketName())
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

func help(in string) error {
	fmt.Println(`
General:
      Soi is a url management CLI system. which could add url, find url and rename url.

Commands:

  [add]: 
      Desc:  add url to soi
      Usage: add (URL)
      Option:
        -n: logical name of the url    (default: <title> of the URL)
        -d: directory to store the url (default: "new")

  [dig]: 
      Desc:  dig url directory with [Tab] key completion and [→] key listing next suggestions
      Usage: dig (URL Completion with [Tab] and [→] key listing next suggestions)

  [list]:
      Desc:  list all urls with filtering
      Usage: list (free words)
		
  [tag]:
      Desc:  not implemented now

  [mv]:
      Desc:  move file to dir 
      Usage: mv (current path) to (dir)

  [quit]:
      Desc:  quit soi> and go back to console. Ctrl+D works too.
      Usage: quit`)
	return nil
}

func quit(in string) error {
	fmt.Println("bye!")
	os.Exit(0)
	return nil
}
