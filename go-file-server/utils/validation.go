package utils

import (
	"os"
)

var secretToken = os.Getenv("SECRET_TOKEN")

func VerifyToken(token string) bool {
	if token == secretToken {
		return true
	}
	return false
}

func IsFileTypeMatched(requestedFileType string, mimeType string) bool {
	switch requestedFileType {
	case "video/quicktime":
		return mimeType == "video/quicktime"
	case "video/mp4":
		return mimeType == "video/mp4"
	default:
		return false
	}
}
