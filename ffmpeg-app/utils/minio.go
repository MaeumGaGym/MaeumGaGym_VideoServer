package utils

import (
	"github.com/minio/minio-go"
	"log"
	"os"
)

var minioClient *minio.Client

var endpoint = os.Getenv("STORAGE_ENDPOINT")
var accessKey = os.Getenv("STORAGE_ACCESS_KEY")
var secretKey = os.Getenv("STORAGE_SECRET_KEY")

func InitMinio() {
	var err error

	minioClient, err = minio.New(endpoint, accessKey, secretKey, true)
	if err != nil {
		log.Fatalln(err)
	}
}

func UploadFile(bucketName, objectName, filePath, fileType string) error {
	_, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: fileType})
	if err != nil {
		return err
	}
	return nil
}
