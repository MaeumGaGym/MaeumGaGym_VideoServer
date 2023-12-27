package utils

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/exec"
	"path/filepath"
	"pokabook/ffmepg-app/database"
	"time"
)

var baseUrl = os.Getenv("BASE_URL")
var webhookUrl = os.Getenv("WEBHOOK_URL")
var bucketName = os.Getenv("STORAGE_BUCKET")

func ConvertVideo(videoPath string, randomStr string) error {
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		return fmt.Errorf("source video file does not exist: %s", videoPath)
	}

	outputBasePath := fmt.Sprintf("./videos/%s/", randomStr)

	resolutions := []struct {
		Name string
		Size string
	}{
		{"144p", "256x144"},
		{"240p", "426x240"},
		{"360p", "640x360"},
		{"480p", "854x480"},
		{"720p", "1280x720"},
		{"1080p", "1920x1080"},
	}

	for _, resolution := range resolutions {
		outputDirPath := outputBasePath + resolution.Name + "/"

		if _, err := os.Stat(outputDirPath); !os.IsNotExist(err) {
			continue
		}

		if err := os.MkdirAll(outputDirPath, 0755); err != nil {
			return fmt.Errorf("failed to create output directory %s: %v", outputDirPath, err)
		}

		u, _ := uuid.NewUUID()
		outputM3U8Path := outputDirPath + "index.m3u8"
		outputTSPath := outputDirPath + u.String() + "-%04d.ts"
		hlsBaseUrl := fmt.Sprintf("https://%s/%s/%s/%s/", endpoint, bucketName, randomStr, resolution.Name)

		start := time.Now()

		cmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", "scale="+resolution.Size, "-c:v", "libx264", "-hls_time", "9", "-hls_list_size", "0", "-hls_base_url", hlsBaseUrl, "-hls_segment_filename", outputTSPath, outputM3U8Path)

		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error: %v, Output: %s\n", err, output)
			return err
		}
		err = filepath.Walk(outputDirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				contentType := ""
				if filepath.Ext(path) == ".ts" {
					contentType = "video/MP2T"
				} else if filepath.Ext(path) == ".m3u8" {
					contentType = "application/x-mpegURL"
				}

				err = UploadFile(bucketName, filepath.Join(randomStr, resolution.Name, info.Name()), path, contentType)
				if err != nil {
					return fmt.Errorf("Failed to upload file to Minio: %v", err)
				}

				err = os.Remove(path)
				if err != nil {
					return fmt.Errorf("Failed to delete file: %v", err)
				}
			}

			return nil
		})
		if err != nil {
			return err
		}

		err = sendDiscordNotification(webhookUrl, randomStr, resolution.Name, time.Since(start).String(), baseUrl+randomStr+"/index.m3u8?scale="+resolution.Name)
		if err != nil {
			return fmt.Errorf("Failed to send Discord notification: %v", err)
		}

		err = database.SaveVideo(randomStr, resolution.Name, time.Since(start).String(), baseUrl+randomStr+"/index.m3u8?scale="+resolution.Name)
		if err != nil {
			return fmt.Errorf("Failed to save video data to sqlite: %v", err)
		}
	}

	err := os.Remove(videoPath)
	if err != nil {
		return fmt.Errorf("failed to delete source video file %s: %v", videoPath, err)
	}

	err = os.RemoveAll(outputBasePath)
	if err != nil {
		return fmt.Errorf("Failed to delete directory: %v", err)
	}

	return nil
}
