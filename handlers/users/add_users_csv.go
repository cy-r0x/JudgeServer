package users

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/judgenot0/judge-backend/handlers/structs"
	usercsv "github.com/judgenot0/judge-backend/handlers/user_csv"
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

	csvHandler, err := usercsv.NewHandler(prefix, clanLengthInt, contestIdInt)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "error creating csv file")
		return
	}
	csvHandler.WriteHeader()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "error reading csv file")
		return
	}

	// Process users concurrently
	type userResult struct {
		user structs.User
		err  error
		idx  int
	}

	results := make(chan userResult, len(records))
	var wg sync.WaitGroup

	// Process each record concurrently
	for idx, record := range records {
		wg.Add(1)
		go func(idx int, record []string) {
			defer wg.Done()

			user, err := csvHandler.WriteUser(record, idx)
			if err != nil {
				results <- userResult{err: fmt.Errorf("error writing to csv: %w", err), idx: idx}
				return
			}

			// Hash password (CPU-intensive operation)
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
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
	users := make([]structs.User, len(records))
	for result := range results {
		if result.err != nil {
			utils.SendResponse(w, http.StatusInternalServerError, result.err.Error())
			return
		}
		users[result.idx] = result.user
	}

	// Insert all users in a single transaction
	trx, err := h.db.Beginx()
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "error starting transaction")
		return
	}
	defer trx.Rollback()

	query := `
		INSERT INTO users (full_name, username, password, clan, room_no, pc_no, role, allowed_contest)
		VALUES (:full_name,:username, :password, :clan, :room_no, :pc_no, :role, :allowed_contest);
	`

	for _, user := range users {
		_, err = trx.NamedExec(query, user)
		if err != nil {
			utils.SendResponse(w, http.StatusInternalServerError, "error registering user")
			return
		}
	}

	query = `INSERT INTO filepath (contest_id, file_path) VALUES ($1, $2)`
	_, err = trx.Exec(query, contestIdInt, csvHandler.FilePath)
	if err != nil {
		fmt.Println(err)
		utils.SendResponse(w, http.StatusInternalServerError, "error saving file path")
		return
	}

	if err := trx.Commit(); err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "error committing transaction")
		return
	}
	utils.SendResponse(w, http.StatusOK, "Users added successfully")
}
