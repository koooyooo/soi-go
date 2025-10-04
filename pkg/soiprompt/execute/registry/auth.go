package registry

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/koooyooo/soi-go/pkg/config"
)

func generateAuthValues(cfg *config.Config) (user string, pass string, authValue string, err error) {
	user = cfg.UserName
	if user == "" {
		fmt.Println("username?")
		fmt.Print("> ")
		fmt.Scan(&user)
		if user == "" {
			return "", "", "", errors.New("no user name found")
		}
	}
	pass = cfg.UserPass
	if pass == "" {
		fmt.Println("password?")
		fmt.Print("> ")
		fmt.Scan(&pass)
		if pass == "" {
			return "", "", "", errors.New("not user pass found")
		}
	}
	return user, pass, "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+pass)), nil
}
