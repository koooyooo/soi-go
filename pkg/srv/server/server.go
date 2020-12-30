package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/koooyooo/soi-go/pkg/srv/auth"

	"github.com/koooyooo/soi-go/pkg/srv/constant"

	"github.com/koooyooo/soi-go/pkg/soi"

	"github.com/gin-gonic/gin"

	"github.com/koooyooo/soi-go/pkg/srv/repo"
)

// Run はルーティングを行います
func Run() {
	r := gin.Default()
	r.GET("/api/v1/", root)
	r.GET("/api/v1/:user_id/:soi_bucket_id/sois", listHandler)
	r.POST("/api/v1/:user_id/:soi_bucket_id/sois", postHandler)
	r.POST("/api/v1/:user_id/:soi_bucket_id/sois:replace", replaceHandler)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("run server fails: %v", err)
	}
}

// root はルートパスのルーティングです
func root(c *gin.Context) {
	c.JSON(200, gin.H{
		"Result": "OK",
	})
}

// listHandler はsoiのリストを取得します
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

// postHandler は
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
