package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/koooyooo/soi-go/pkg/srv/constant"

	"github.com/gin-gonic/gin"
	"github.com/koooyooo/soi-go/pkg/soi"
	"github.com/koooyooo/soi-go/pkg/srv/auth"
	"github.com/koooyooo/soi-go/pkg/srv/repo"
)

// createContext は認証情報をContextに埋め込みます
func createContext(ctx context.Context, gc *gin.Context) context.Context {
	if userID := gc.Param("user_id"); userID != "" {
		ctx = context.WithValue(ctx, constant.CtxKeyUserID, userID)
	}
	if soiBucketID := gc.Param("soi_bucket_id"); soiBucketID != "" {
		ctx = context.WithValue(ctx, constant.CtxKeySoiBucketID, soiBucketID)
	}
	// TODO 認証情報の埋め込み
	return ctx
}

// rootHandler はルートパスのハンドラです
func rootHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"Result": "OK",
	})
}

// listHandler はsoiListを取得するハンドラです
func listHandler(c *gin.Context) {
	// 認証
	authResult, err := auth.Authorize(c)
	if !authResult {
		fmt.Printf("auth failed: %s", err.Error())
		_ = c.AbortWithError(404, err)
		return
	}
	if err != nil {
		fmt.Printf("auth error: %s", err.Error())
		_ = c.AbortWithError(500, err)
		return
	}

	ctx := createContext(context.Background(), c)
	repo := repo.NewRepository()
	sb, err := repo.LoadAll(ctx)
	if err != nil {
		fmt.Printf("load error: %s", err.Error())
		_ = c.AbortWithError(500, err)
		return
	}
	c.JSON(200, sb)
}

// postHandler はsoiを登録するハンドラです
func postHandler(c *gin.Context) {
	// 認証
	authResult, err := auth.Authorize(c)
	if !authResult {
		_ = c.AbortWithError(404, err)
		return
	}
	if err != nil {
		_ = c.AbortWithError(500, err)
		return
	}

	ctx := createContext(context.Background(), c)
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		_ = c.AbortWithError(500, err)
		return
	}
	var s soi.SoiVirtual
	if err = json.Unmarshal(b, &s); err != nil {
		_ = c.AbortWithError(500, err)
		return
	}
	repo := repo.NewRepository()
	if err = repo.Store(ctx, &s); err != nil {
		_ = c.AbortWithError(500, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Store": "OK",
	})
}

// replaceHandler は全てのsoisを洗い替えします
func replaceHandler(c *gin.Context) {
	// 認証
	authResult, err := auth.Authorize(c)
	if !authResult {
		_ = c.AbortWithError(404, err)
		return
	}
	if err != nil {
		_ = c.AbortWithError(500, err)
		return
	}

	ctx := createContext(context.Background(), c)
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		_ = c.AbortWithError(500, err)
		return
	}
	var svb soi.SoiVirtualBucket
	if err = json.Unmarshal(b, &svb); err != nil {
		_ = c.AbortWithError(500, err)
		return
	}
	repo := repo.NewRepository()
	if err = repo.StoreAll(ctx, &svb); err != nil {
		_ = c.AbortWithError(500, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"StoreAll": "OK",
	})
}
