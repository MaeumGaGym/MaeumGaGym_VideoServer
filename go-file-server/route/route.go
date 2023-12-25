package route

import (
	"github.com/gin-gonic/gin"
	"pokabook/go-file-server/controller"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/upload", controller.UploadVideo)
	r.POST("/generate", controller.Generate)
	r.GET("/:id/index.m3u8", controller.GetM3U8)
	r.GET("/:id/:ts", controller.GetTS)
	r.DELETE("/:id", controller.RemoveVideo)
}
