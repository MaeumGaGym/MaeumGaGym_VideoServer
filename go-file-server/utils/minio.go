package utils

import (
	"github.com/minio/minio-go"
	"io"
	"log"
	"os"
)

var minioClient *minio.Client

var endpoint = os.Getenv("STORAGE_ENDPOINT")
var accessKey = os.Getenv("STORAGE_ACCESS_KEY")
var secretKey = os.Getenv("STORAGE_SECRET_KEY")
var bucketName = os.Getenv("STORAGE_BUCKET")

func InitMinio() {
	var err error

	minioClient, err = minio.New(endpoint, accessKey, secretKey, true)
	if err != nil {
		log.Fatalln(err)
	}
}

func RemoveFile(videoId string) error {
	err := minioClient.RemoveObject(bucketName, videoId)
	if err != nil {
		return err
	}
	return nil
}

func GetVideoFile(videoId, scale string) ([]byte, error) {
	objectName := videoId + "/" + scale + "/index.m3u8"
	object, err := minioClient.GetObject(bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, err
	}

	return data, nil
}
