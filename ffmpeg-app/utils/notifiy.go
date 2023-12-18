package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendDiscordNotification(webhookURL, message string) error {
	payload := map[string]string{
		"content": message,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("POST request failed with status: %d", resp.StatusCode)
	}

	return nil
}
