package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Authorize はヘッダに認証情報を含めて送信します
// Basic認証を用います
func Authorize(c *gin.Context) (bool, error) {
	authHeader := c.Request.Header.Get("Authorization")
	fmt.Println(authHeader)
	//if authHeader == "" {
	//	return false, fmt.Errorf("authorization header not found")
	//}
	//authValue64 := strings.TrimPrefix(authHeader, "Basic ")
	//authValue, err := base64.StdEncoding.DecodeString(authValue64)
	//if err != nil {
	//	return false, fmt.Errorf("fail in decoding base64 value")
	//}
	//sp := strings.Split(string(authValue), ":")
	//user := sp[0]
	//pass := sp[1]
	//fmt.Println(user, pass)
	return true, nil
}
