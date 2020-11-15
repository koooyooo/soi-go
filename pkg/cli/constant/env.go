package constant

import "os"

var (
	EnvKeySoiUserName = ClientEnvKey("SOI_USER_NAME")
	EnvKeySoiUserPass = ClientEnvKey("SOI_USER_PASS")
)

type ClientEnvKey string

func (ek ClientEnvKey) Get() string {
	return os.Getenv(string(ek))
}
