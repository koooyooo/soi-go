package server

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/koooyooo/soi-go/pkg/soi"

	"github.com/koooyooo/soi-go/pkg/srv"

	"github.com/gin-gonic/gin"

	"github.com/koooyooo/soi-go/pkg/srv/repo"
)

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

func root(c *gin.Context) {
	c.JSON(200, gin.H{
		"Hello": "World",
	})
}

func listHandler(c *gin.Context) {
	ctx := createContext(context.Background(), c)
	repo := repo.NewRepository()
	sb, err := repo.LoadAll(ctx)
	if err != nil {
		_ = c.AbortWithError(500, err)
		return
	}
	c.JSON(200, sb)
}

func postHandler(c *gin.Context) {
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
		ctx = context.WithValue(ctx, srv.CtxKeyUserID, userID)
	}
	if soiBucketID := gc.Param("soi_bucket_id"); soiBucketID != "" {
		ctx = context.WithValue(ctx, srv.CtxKeySoiBucketID, soiBucketID)
	}
	// TODO 認証情報の埋め込み
	return ctx
}
