package compilerun

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) CompileRun(w http.ResponseWriter, r *http.Request) {

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		url := h.config.EngineUrl + "/run"
		resp, err := http.Post(url, "application/json", r.Body)
		if err != nil {
			log.Println(err)
			utils.SendResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
		decoder := json.NewDecoder(resp.Body)
		var result struct {
			Result string `json:"result"`
		}
		decoder.Decode(&result)
		utils.SendResponse(w, http.StatusOK, result)
	}()

	wg.Wait()

}
