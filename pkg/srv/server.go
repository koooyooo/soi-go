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
	r.GET("/", root)
	r.GET("/show", showHandlerG)
	r.POST("/store", storeHandlerG)

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
	repo := repo.NewRepository()
	sb, err := repo.LoadAll(context.Background())
	if err != nil {
		_ = c.AbortWithError(500, err)
		return
	}
	c.JSON(200, sb)
}

func storeHandlerG(c *gin.Context) {
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
	if err = repo.Store(context.Background(), &s); err != nil {
		_ = c.AbortWithError(500, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Store": "OK",
	})
}
