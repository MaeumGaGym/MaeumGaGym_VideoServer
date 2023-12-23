package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"pokabook/go-file-server/model"
	"strings"
)

var baseUrl = os.Getenv("BASE_URL")

func UploadVideo(ctx *gin.Context) {
	file, err := ctx.FormFile("video")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No such video"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	if ext != ".mov" && ext != ".mp4" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file extension"})
		return
	}

	
	randomStr, _ := uuid.NewUUID()
	videoId := randomStr.String()[:8]

	tempDir := "/app/videos/" + videoId
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tempFilePath := filepath.Join(tempDir, file.Filename)
	log.Println("Saving uploaded file to: ", tempFilePath)

	if err := ctx.SaveUploadedFile(file, tempFilePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileData, err := os.Open(tempFilePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer fileData.Close()

	buffer := make([]byte, 512)
	_, err = fileData.Read(buffer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	mimeType := http.DetectContentType(buffer)
	if mimeType != "video/quicktime" && mimeType != "video/mp4" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}

	if err := model.ConvertVideo(tempFilePath, videoId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	url := baseUrl + videoId + "/index.m3u8"
	ctx.JSON(http.StatusOK, gin.H{"videoURL": url})
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

func RemoveVideo(ctx *gin.Context) {
	videoId := ctx.Param("id")

	path := "/app/videos/" + videoId
	log.Println("removing file to: ", path)

	err := os.RemoveAll(path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}
