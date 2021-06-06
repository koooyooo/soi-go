package http

import (
	"log"

	"github.com/gin-gonic/gin"
)

// RunServer はルーティングを行います
func RunServer() {
	r := gin.Default()
	r.GET("/api/v1/", rootHandler)
	r.GET("/api/v1/:user_id/:soi_bucket_id/sois", listHandler)
	r.POST("/api/v1/:user_id/:soi_bucket_id/sois", postHandler)
	r.POST("/api/v1/:user_id/:soi_bucket_id/sois:replace", replaceHandler)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("run server fails: %v", err)
	}
}
