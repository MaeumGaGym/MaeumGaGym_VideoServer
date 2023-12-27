package main

import (
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"pokabook/ffmepg-app/database"
	"pokabook/ffmepg-app/utils"
)

type VideoMessage struct {
	Filepath string `json:"filepath"`
	VideoID  string `json:"video_id"`
}

func main() {
	log.Println("Ffmpeg Start!")
	utils.InitMinio()
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

	unfinishedVideos, err := database.GetUnfinishedVideos()
	if err != nil {
		log.Fatal(err)
	}
	for _, videoId := range unfinishedVideos {
		videoMessage := VideoMessage{
			Filepath: "./videos/" + videoId,
			VideoID:  videoId,
		}

		err = utils.ConvertVideo(videoMessage.Filepath, videoMessage.VideoID)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	for msg := range msgs {
		var videoMessage VideoMessage
		err := json.Unmarshal(msg.Body, &videoMessage)
		if err != nil {
			log.Println(err)
			continue
		}

		err = utils.ConvertVideo(videoMessage.Filepath, videoMessage.VideoID)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
