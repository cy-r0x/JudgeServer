package usercsv

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

func safeIntParse(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func (h *Handler) AddUserCsv(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	contestID := r.FormValue("contest_id")

	file, _, err := r.FormFile("file")
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	defer file.Close()

	contestID = strings.TrimSpace(contestID)

	if contestID == "" {
		utils.SendResponse(w, http.StatusBadRequest, "invalid form data", nil)
		return
	}

	// Fetch Contest to get the prefix
	var contest models.Contest
	if err := h.db.Where("id = ?", contestID).First(&contest).Error; err != nil {
		utils.SendResponse(w, http.StatusNotFound, "Contest not found", nil)
		return
	}
	prefix := contest.UserPrefix

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "error reading csv file", nil)
		return
	}

	if len(records) < 2 {
		utils.SendResponse(w, http.StatusBadRequest, "CSV file is empty or missing data", nil)
		return
	}

	// We no longer need to write a static CSV, since GetCSV reads from DB.
	// h.NewWriteHandler(prefix, clanLengthInt, contestID)
	// h.WriteHeader()

	// Respond immediately indicating accepted processing
	utils.SendResponse(w, http.StatusAccepted, "Users creation started in the background", nil)

	// Process asynchronously
	go func() {
		type userResult struct {
			user User
			err  error
			idx  int
		}

		results := make(chan userResult, len(records))
		var wg sync.WaitGroup

		for idx, record := range records[1:] {
			wg.Add(1)
			go func(idx int, record []string) {
				defer wg.Done()

				user, err := h.GenerateUser(record, idx, prefix, contestID)
				if err != nil {
					results <- userResult{err: fmt.Errorf("error generating user: %w", err), idx: idx}
					return
				}

				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.UnHashedPassword), bcrypt.DefaultCost)
				if err != nil {
					results <- userResult{err: fmt.Errorf("error hashing password: %w", err), idx: idx}
					return
				}
				user.Password = string(hashedPassword)

				results <- userResult{user: user, idx: idx}
			}(idx, record)
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		users := make([]User, len(records)-1)
		hasErrors := false
		for result := range results {
			if result.err != nil {
				log.Println("Error processing user CSV async:", result.err)
				hasErrors = true
				continue
			}
			users[result.idx] = result.user
		}

		if hasErrors {
			log.Println("Failed to process some users from CSV. Check logs for details.")
			return
		}

		trx := h.db.Begin()
		if trx.Error != nil {
			log.Println("error starting transaction:", trx.Error)
			return
		}

		defer func() {
			if r := recover(); r != nil {
				trx.Rollback()
			}
		}()

		for _, user := range users {
			var allowedContest *string
			if user.AllowedContest != nil && *user.AllowedContest != "" {
				ac := *user.AllowedContest
				allowedContest = &ac
			}

			// Map correctly to updated models.User fields
			role := models.RoleUser
			if user.Role != "" {
				role = models.Role(user.Role)
			}

			newUser := models.User{
				Name:           user.Name,
				Username:       user.Username,
				Password:       user.Password,
				Role:           role,
				AllowedContest: allowedContest,
				AdditionalInfo: user.AdditionalInfo,
				RoomNo:         user.RoomNo,
				PcNo:           user.PcNo,
				CreatedAt:      time.Now(),
			}

			if err := trx.Create(&newUser).Error; err != nil {
				log.Println("error registering user:", err)
				trx.Rollback()
				return
			}

			if allowedContest != nil {
				creds := models.UserCreds{
					ContestId:     *allowedContest,
					UserId:        newUser.Id,
					PlainPassword: user.UnHashedPassword, // Save the UNHASHED plain password here
				}
				if err := trx.Create(&creds).Error; err != nil {
					log.Println("Failed to save user credentials:", err)
					trx.Rollback()
					return
				}
			}
		}

		if err := trx.Commit().Error; err != nil {
			log.Println("error committing transaction:", err)
			return
		}

		log.Println("Async CSV User processing successful")
	}()
}
