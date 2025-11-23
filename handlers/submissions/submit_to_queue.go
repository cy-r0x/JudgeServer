package submissions

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

func (h *Handler) submitToQueue(problem *Problem) error {

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		jsonData, err := json.Marshal(problem)
		if err != nil {
			log.Println("Error marshaling queueData:", err)
			return
		}
		url := h.config.EngineUrl + "/submit"
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("Error posting to queue:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Queue responded with status: %s\n", resp.Status)
		}
	}()

	wg.Wait()
	return nil
}
