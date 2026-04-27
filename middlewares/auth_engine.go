package middlewares

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/judgenot0/judge-backend/utils"
)

type EnginePayload struct {
	SubmissionId    int64    `json:"submission_id"`
	Status          string   `json:"verdict"`
	ExecutionTime   *float32 `json:"execution_time"`
	ExecutionMemory *float32 `json:"execution_memory"`
	Timestamp       int64    `json:"timestamp"`
}

func VerifyToken(enginePayload *EnginePayload, accessToken string, secret string) bool {
	if secret == "" || enginePayload == nil {
		return false
	}

	if time.Since(time.Unix(enginePayload.Timestamp, 0)) > 5*time.Minute {
		return false
	}

	message, err := json.Marshal(enginePayload)
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	expectedHex := hex.EncodeToString(expectedMAC)

	return hmac.Equal([]byte(expectedHex), []byte(accessToken))
}

func (m *Middlewares) AuthEngine(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		const maxBodySize = 10 * 1024
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

		header := r.Header.Get("Authorization")
		if header == "" {
			utils.SendResponse(w, http.StatusUnauthorized, "Authorization header required", nil)
			return
		}

		headerArr := strings.Split(header, " ")
		if len(headerArr) != 2 {
			utils.SendResponse(w, http.StatusUnauthorized, "Invalid token format", nil)
			return
		}

		accessToken := headerArr[1]

		decoder := json.NewDecoder(r.Body)
		var enginePayload EnginePayload
		err := decoder.Decode(&enginePayload)

		if err != nil {
			log.Println(err)
			utils.SendResponse(w, http.StatusBadRequest, "Invalid JSON", nil)
			return
		}

		ok := VerifyToken(&enginePayload, accessToken, m.config.EngineKey)
		if !ok {
			utils.SendResponse(w, http.StatusBadRequest, "Invalid Token", nil)
			return
		}

		ctx := context.WithValue(r.Context(), "enginePayload", enginePayload)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
