package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"pokabook/go-file-server/dto"
	"pokabook/go-file-server/model"
	"pokabook/go-file-server/utils"
	"strings"
	"time"
)

var baseUrl = os.Getenv("BASE_URL")
var bucketName = os.Getenv("STORAGE_BUCKET")
var endpoint = os.Getenv("STORAGE_ENDPOINT")

func Generate(ctx *gin.Context) {
	var req dto.GenerateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !utils.VerifyToken(ctx.GetHeader("MaeumgaGym-Token")) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	randomStr, _ := uuid.NewUUID()
	videoId := randomStr.String()[:8]

	randomPassword := utils.GenerateRandomPassword()

	encryptedParams := utils.EncryptQueryParams(map[string]string{
		"fileType":   req.FileType,
		"TimeToLive": time.Now().Format(time.RFC3339),
		"videoId":    videoId,
	}, randomPassword)

	uploadURL := baseUrl + "upload?" + "params=" + url.QueryEscape(encryptedParams) + "&key=" + url.QueryEscape(string(randomPassword))

	ctx.JSON(http.StatusOK, gin.H{"uploadURL": uploadURL, "videoId": videoId})
}

func UploadVideo(ctx *gin.Context) {

	encryptedParams := ctx.Query("params")

	key := ctx.Query("key")

	decryptedParams, err := utils.CustomDecrypt(encryptedParams, []byte(key))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decrypt query params"})
		return
	}

	var params map[string]string
	if err := json.Unmarshal([]byte(decryptedParams), &params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to convert decrypted params to map"})
		return
	}

	requestTimeStr := params["TimeToLive"]
	TimeToLive, err := time.Parse(time.RFC3339, requestTimeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request time"})
		return
	}

	if time.Since(TimeToLive).Minutes() > 3 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Upload time exceeded 3 minutes"})
		return
	}

	file, err := ctx.FormFile("video")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No video file was received"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".mov" && ext != ".mp4" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file extension"})
		return
	}

	videoId := params["videoId"]

	tempDir := "/app/videos/" + videoId
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	tempFilePath := filepath.Join(tempDir, videoId+ext)
	if err := ctx.SaveUploadedFile(file, tempFilePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	fileData, err := os.Open(tempFilePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}
	defer fileData.Close()
	buffer := make([]byte, 261)
	_, err = fileData.Read(buffer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	kind, err := filetype.Match(buffer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	requestedFileType := params["fileType"]

	if !utils.IsFileTypeMatched(requestedFileType, kind.MIME.Value) {
		os.Remove(tempFilePath)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or mismatched file type"})
		return
	}

	if err := model.ConvertVideo(tempFilePath, videoId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	videoURL := baseUrl + videoId + "/index.m3u8"
	thumbnailURL := "https://" + endpoint + "/" + bucketName + "/" + videoId + "/thumbnail.png"
	ctx.JSON(http.StatusOK, gin.H{"videoURL": videoURL, "thumbnailURL": thumbnailURL})
}

func GetM3U8(ctx *gin.Context) {
	videoId := ctx.Param("id")

	scale := ctx.Query("scale")
	if scale == "" {
		scale = "720p"
	}

	ctx.SetCookie("scale", scale, 3600, "", "", false, true)

	data, err := utils.GetVideoFile(videoId, scale)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File does not exist"})
		return
	}

	ctx.Data(http.StatusOK, "application/x-mpegURL", data)
}

func RemoveVideo(ctx *gin.Context) {
	videoId := ctx.Param("id")

	if !utils.VerifyToken(ctx.GetHeader("MaeumgaGym-Token")) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	err := utils.RemoveFile(videoId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}

func HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
	return
}
