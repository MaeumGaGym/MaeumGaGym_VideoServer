package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"pokabook/go-file-server/model"
	"time"
)

func UploadVideo(ctx *gin.Context) {
	file, err := ctx.FormFile("video")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tempDir := "/app/videos"
	tempFilePath := filepath.Join(tempDir, file.Filename)
	log.Printf("Saving uploaded file to: %s\n", tempFilePath)

	if err := ctx.SaveUploadedFile(file, tempFilePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rand.NewSource(time.Now().UnixNano())
	randomStr := fmt.Sprintf("%06d", rand.Intn(1000000))

	if err := model.ConvertVideo(tempFilePath, randomStr); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"videoId": randomStr})
}

func GetM3U8(ctx *gin.Context) {
	videoId := ctx.Param("id")

	scale := ctx.Query("scale")
	if scale == "" {
		scale = "720p"
	}

	ctx.SetCookie("scale", scale, 3600, "", "", false, true)

	data, err := model.GetVideo(videoId + "/" + scale + "/index.m3u8")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File does not exist"})
		return
	}

	ctx.Data(http.StatusOK, "application/x-mpegURL", data)
}

func GetTS(ctx *gin.Context) {
	videoId := ctx.Param("id")
	videoFile := ctx.Param("ts")
	ext := filepath.Ext(videoFile)

	if ext != ".ts" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file extension"})
		return
	}

	scale, err := ctx.Cookie("scale")
	if err != nil {
		scale = "720p"
	}

	data, err := model.GetVideo(videoId + "/" + scale + "/" + videoFile)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File does not exist"})
		return
	}

	ctx.Data(http.StatusOK, "video/MP2T", data)
}