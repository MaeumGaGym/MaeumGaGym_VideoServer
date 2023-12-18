package model

import (
	"os"
	"pokabook/go-file-server/utils"
)

func ConvertVideo(videoPath string, randomStr string) error {
	err := utils.PublishMessage(videoPath, randomStr)
	if err != nil {
		return err
	}
	return nil
}

func GetVideo(videoPath string) ([]byte, error) {
	return os.ReadFile("./videos/" + videoPath)
}
