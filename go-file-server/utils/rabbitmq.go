package utils

import (
	"context"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

type VideoMessage struct {
	Filepath string `json:"filepath"`
	VideoID  string `json:"video_id"`
}

var conn *amqp091.Connection
var ch *amqp091.Channel
var q amqp091.Queue

func InitRabbitMQ() {
	var err error
	conn, err = amqp091.Dial("amqp://guest:guest@" + os.Getenv("RABBITMQ_HOST") + ":5672/")
	if err != nil {
		log.Fatal(err)
	}
	ch, err = conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	q, err = ch.QueueDeclare(
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
}

func PublishMessage(tempFilePath, randomStr string) error {

	videoMessage := VideoMessage{
		Filepath: tempFilePath,
		VideoID:  randomStr,
	}
	body, err := json.Marshal(videoMessage)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		context.Background(),
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}
