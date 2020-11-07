package srv

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/koooyooo/soi-go/pkg/srv/repo"

	"github.com/koooyooo/soi-go/pkg/cli"
)

func Run() {
	r := gin.Default()
	r.GET("/api/v1/", root)
	r.GET("/api/v1/:user_id/sois", showHandlerG)
	r.POST("/api/v1/:user_id/sois", storeHandlerG)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("run server fails: %v", err)
	}
}

func root(c *gin.Context) {
	c.JSON(200, gin.H{
		"Hello": "World",
	})
}

func showHandlerG(c *gin.Context) {
	ctx := createContext(context.Background(), c)
	repo := repo.NewRepository()
	sb, err := repo.LoadAll(ctx)
	if err != nil {
		_ = c.AbortWithError(500, err)
		return
	}
	c.JSON(200, sb)
}

func storeHandlerG(c *gin.Context) {
	ctx := createContext(context.Background(), c)
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		_ = c.AbortWithError(500, err)
		return
	}
	var s cli.SoiVirtual
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

func createContext(ctx context.Context, gc *gin.Context) context.Context {
	if userID := gc.Param("user_id"); userID != "" {
		ctx = context.WithValue(ctx, "user_id", userID)
	}
	// TODO 認証情報の埋め込み
	return ctx
}
