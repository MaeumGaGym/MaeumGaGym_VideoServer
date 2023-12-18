package utils

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/exec"
	"time"
)

var baseUrl = os.Getenv("BASE_URL")
var webhookUrl = os.Getenv("WEBHOOK_URL")

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
		if err := os.MkdirAll(outputDirPath, 0755); err != nil {
			return fmt.Errorf("failed to create output directory %s: %v", outputDirPath, err)
		}

		u, _ := uuid.NewUUID()
		outputM3U8Path := outputDirPath + "index.m3u8"
		outputTSPath := outputDirPath + u.String() + "-%04d.ts"
		hlsBaseUrl := baseUrl + randomStr + "/"

		start := time.Now()

		cmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", "scale="+resolution.Size, "-c:v", "libx264", "-hls_time", "9", "-hls_list_size", "0", "-hls_base_url", hlsBaseUrl, "-hls_segment_filename", outputTSPath, outputM3U8Path)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error: %v, Output: %s\n", err, output)
			return err
		}

		message := fmt.Sprintf("### Video conversion completed!\n videoId: %s\n scale: %s\n time: %s\n Url: %s", randomStr, resolution.Name, time.Since(start), baseUrl+randomStr+"index.m3u8?scale="+resolution.Name)
		err = SendDiscordNotification(webhookUrl, message)
		if err != nil {
			return fmt.Errorf("Failed to send Discord notification: %v", err)
		}
	}

	err := os.Remove(videoPath)
	if err != nil {
		return fmt.Errorf("failed to delete source video file %s: %v", videoPath, err)
	}

	return nil
}
