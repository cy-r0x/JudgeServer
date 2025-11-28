package usercsv

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"

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
	contest_id := r.FormValue("contest_id")

	file, _, err := r.FormFile("file")
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	clanLengthInt := safeIntParse(clanLen)
	contestIdInt := int64(safeIntParse(contest_id))

	if contestIdInt == 0 || clanLengthInt == 0 || prefix == "" {
		utils.SendResponse(w, http.StatusBadRequest, "invalid form data")
		return
	}

	err = h.NewWriteHandler(prefix, clanLengthInt, contestIdInt)
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
	trx, err := h.db.Beginx()
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "error starting transaction")
		return
	}
	defer trx.Rollback()

	query := `
		INSERT INTO users (full_name, username, password, clan, room_no, pc_no, role, allowed_contest)
		VALUES ($1,$2, $3, $4, $5, $6, $7, $8) RETURNING id;
	`

	for _, user := range users {

		err = h.WriteUser(user)
		if err != nil {
			utils.SendResponse(w, http.StatusInternalServerError, "error writing user to CSV")
			return
		}

		// Password already hashed in parallel goroutines above
		_, err = trx.Exec(query, user.FullName, user.Username, user.Password, user.Clan, user.RoomNo, user.PcNo, user.Role, user.AllowedContest)
		if err != nil {
			utils.SendResponse(w, http.StatusInternalServerError, "error registering user")
			return
		}
	}

	query = `INSERT INTO filepath (contest_id, file_path) VALUES ($1, $2)`
	_, err = trx.Exec(query, contestIdInt, h.Writer.FilePath)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "error saving file path")
		return
	}

	if err := trx.Commit(); err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "error committing transaction")
		return
	}
	utils.SendResponse(w, http.StatusOK, "Users added successfully")
}
