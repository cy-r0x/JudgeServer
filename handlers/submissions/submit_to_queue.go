package submissions

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (h *Handler) submitToQueue(problem *Problem) error {
	jsonData, err := json.Marshal(problem)
	if err != nil {
		log.Println("Error marshaling queueData:", err)
		return err
	}

	url := h.config.EngineUrl + "/submit"

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating queue request:", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error posting to queue:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Queue responded with status: %s\n", resp.Status)
		return nil
	}

	return nil
}
