package compilerun

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) CompileRun(w http.ResponseWriter, r *http.Request) {
	const maxBodySize = 50 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	// Forward request to engine
	url := h.config.EngineUrl + "/run"
	resp, err := http.Post(url, "application/json", r.Body)
	if err != nil {
		// Check if the error is due to payload size limit
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			log.Printf("Payload too large: exceeded %d bytes", maxBodySize)
			utils.SendResponse(w, http.StatusRequestEntityTooLarge, map[string]string{
				"error": "Request payload too large",
			})
			return
		}

		log.Printf("Error forwarding request to engine: %v", err)
		utils.SendResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to connect to execution engine",
		})
		return
	}
	defer resp.Body.Close()

	// Decode response from engine
	var result struct {
		Result string `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decoding engine response: %v", err)
		utils.SendResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to parse execution result",
		})
		return
	}

	// Send successful response
	utils.SendResponse(w, http.StatusOK, result)
}
