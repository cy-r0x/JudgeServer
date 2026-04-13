package submissions

import (
	"encoding/json"
	"log"
)

func (h *Handler) submitToQueue(problem *Problem) error {
	jsonData, err := json.Marshal(problem)
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
