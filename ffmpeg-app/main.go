package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"os/exec"
)

type VideoMessage struct {
	Filepath string `json:"filepath"`
	VideoID  string `json:"video_id"`
}

var baseUrl = os.Getenv("BASE_URL")

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

		cmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", "scale="+resolution.Size, "-c:v", "libx264", "-hls_time", "9", "-hls_list_size", "0", "-hls_base_url", hlsBaseUrl, "-hls_segment_filename", outputTSPath, outputM3U8Path)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error: %v, Output: %s\n", err, output)
			return err
		}
	}

	err := os.Remove(videoPath)
	if err != nil {
		return fmt.Errorf("failed to delete source video file %s: %v", videoPath, err)
	}

	return nil
}

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@" + os.Getenv("RABBITMQ_HOST") + ":5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"video_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	for msg := range msgs {
		var videoMessage VideoMessage
		err := json.Unmarshal(msg.Body, &videoMessage)
		if err != nil {
			log.Println(err)
			continue
		}

		err = ConvertVideo(videoMessage.Filepath, videoMessage.VideoID)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
