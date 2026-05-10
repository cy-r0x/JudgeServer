package submissions

import (
	"encoding/json"
	"log/slog"
)

func (h *Handler) submitToQueue(payload *QueueSubmission) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		slog.Error("Error marshaling queueData", "error", err)
		return err
	}

	err = h.queueClient.QueueMessage(jsonData)
	if err != nil {
		slog.Error("Error passing submission to queue", "error", err)
		return err
	}

	return nil
}
