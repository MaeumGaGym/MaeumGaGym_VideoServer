package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func sendDiscordNotification(webhookURL, videoId, scale, duration, url string) error {
	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title": "Video Conversion Completed",
				"color": 8900331,
				"fields": []map[string]string{
					{
						"name":   "Video ID",
						"value":  videoId,
						"inline": "true",
					},
					{
						"name":   "Scale",
						"value":  scale,
						"inline": "true",
					},
					{
						"name":   "Duration",
						"value":  duration,
						"inline": "true",
					},
					{
						"name":   "URL",
						"value":  url,
						"inline": "false",
					},
				},
			},
		},
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
