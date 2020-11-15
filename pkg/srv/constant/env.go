package constant

import "os"

var (
	EnvKeySoiBucketName = ServerEnvKey("SOI_BUCKET_NAME")
)

type ServerEnvKey string

func (ek ServerEnvKey) Get() string {
	return os.Getenv(string(ek))
}
