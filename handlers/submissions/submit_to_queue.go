package submissions

import (
	"encoding/json"
	"log"
)

func (h *Handler) submitToQueue(payload *QueueSubmission) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshaling queueData:", err)
		return err
	}

	err = h.queueClient.QueueMessage(jsonData)
	if err != nil {
		log.Println("Error passing submission to queue:", err)
		return err
	}

	return nil
}
