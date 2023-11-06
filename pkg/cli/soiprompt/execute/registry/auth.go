package registry

import (
	"encoding/base64"
	"fmt"

	"soi-go/pkg/cli/constant"
)

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
