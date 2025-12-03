package middlewares

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/judgenot0/judge-backend/utils"
)

type EngineData struct {
	SubmissionId    int64    `json:"submission_id"`
	Verdict         string   `json:"verdict"`
	ExecutionTime   *float32 `json:"execution_time"`
	ExecutionMemory *float32 `json:"execution_memory"`
	Timestamp       int64    `json:"timestamp"`
}

type EnginePayload struct {
	Data        *EngineData `json:"payload"`
	AccessToken string      `json:"access_token"`
}

func VerifyToken(enginePayload EnginePayload, secret string) bool {
	if secret == "" || enginePayload.Data == nil {
		return false
	}

	if time.Since(time.Unix(enginePayload.Data.Timestamp, 0)) > 5*time.Minute {
		return false
	}

	message, err := json.Marshal(enginePayload.Data)
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	expectedHex := hex.EncodeToString(expectedMAC)

	return hmac.Equal([]byte(expectedHex), []byte(enginePayload.AccessToken))
}

func (m *Middlewares) AuthEngine(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var enginePayload EnginePayload
		err := decoder.Decode(&enginePayload)

		if err != nil {
			log.Println(err)
			utils.SendResponse(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
		ok := VerifyToken(enginePayload, m.config.EngineKey)
		if !ok {
			utils.SendResponse(w, http.StatusBadRequest, "Invalid Token")
			return
		}

		ctx := context.WithValue(r.Context(), "engineData", enginePayload.Data)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
