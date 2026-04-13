package usercsv

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"sort"
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
		utils.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	prefix := r.FormValue("prefix")
	clanLen := r.FormValue("clan_length")
	contestID := r.FormValue("contest_id")

	file, _, err := r.FormFile("file")
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	clanLengthInt := safeIntParse(clanLen)
	contestID = strings.TrimSpace(contestID)

	if contestID == "" || clanLengthInt == 0 || prefix == "" {
		utils.SendResponse(w, http.StatusBadRequest, "invalid form data")
		return
	}

	err = h.NewWriteHandler(prefix, clanLengthInt, contestID)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "error creating csv file")
		return
	}
	h.WriteHeader()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "error reading csv file")
		return
	}

	// Process users concurrently
	type userResult struct {
		user User
		err  error
		idx  int
	}

	results := make(chan userResult, len(records))
	var wg sync.WaitGroup

	// Process each record concurrently
	for idx, record := range records[1:] {
		wg.Add(1)
		go func(idx int, record []string) {
			defer wg.Done()

			user, err := h.GenerateUser(record, idx)
			if err != nil {
				results <- userResult{err: fmt.Errorf("error generating user: %w", err), idx: idx}
				return
			}

			// Hash password in parallel (CPU-intensive operation)
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.UnHashedPassword), bcrypt.DefaultCost)
			if err != nil {
				results <- userResult{err: fmt.Errorf("error hashing password: %w", err), idx: idx}
				return
			}
			user.Password = string(hashedPassword)

			results <- userResult{user: user, idx: idx}
		}(idx, record)
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results in order
	users := make([]User, len(records)-1) // Skip header row
	for result := range results {
		if result.err != nil {
			utils.SendResponse(w, http.StatusInternalServerError, result.err.Error())
			return
		}
		users[result.idx] = result.user
	}

	//sort the user by room no and pc no
	sort.Slice(users, func(i, j int) bool {
		if users[i].RoomNo == nil || users[j].RoomNo == nil {
			return false
		}
		if *users[i].RoomNo == *users[j].RoomNo {
			if users[i].PcNo == nil || users[j].PcNo == nil {
				return false
			}
			return safeIntParse(*users[i].PcNo) < safeIntParse(*users[j].PcNo)
		}
		return *users[i].RoomNo < *users[j].RoomNo
	})

	// Insert all users in a single transaction
	trx := h.db.Begin()
	if trx.Error != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "error starting transaction")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			trx.Rollback()
		}
	}()

	for _, user := range users {
		err = h.WriteUser(user)
		if err != nil {
			trx.Rollback()
			utils.SendResponse(w, http.StatusInternalServerError, "error writing user to CSV")
			return
		}

		var allowedContest *string
		if user.AllowedContest != nil {
			ac := *user.AllowedContest
			allowedContest = &ac
		}

		newUser := models.User{
			FullName:       user.FullName,
			Username:       user.Username,
			Password:       user.Password,
			Role:           user.Role,
			AllowedContest: allowedContest,
			Clan:           user.Clan,
			RoomNo:         user.RoomNo,
			PcNo:           user.PcNo,
			CreatedAt:      time.Now(),
		}

		if err := trx.Create(&newUser).Error; err != nil {
			trx.Rollback()
			utils.SendResponse(w, http.StatusInternalServerError, "error registering user")
			return
		}
	}

	filePathEntry := models.Filepath{
		ContestID: contestID,
		FilePath:  h.Writer.FilePath,
	}

	if err := trx.Create(&filePathEntry).Error; err != nil {
		trx.Rollback()
		utils.SendResponse(w, http.StatusInternalServerError, "error saving file path")
		return
	}

	if err := trx.Commit().Error; err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "error committing transaction")
		return
	}

	utils.SendResponse(w, http.StatusOK, "Users added successfully")
}
