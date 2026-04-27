package standings

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetStandings(w http.ResponseWriter, r *http.Request) {
	const limit = 100
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	crrPage, err := strconv.Atoi(page)
	if err != nil {
		log.Println(err)
		utils.SendResponse(w, http.StatusInternalServerError, "Internal Server Error", nil)
		return
	}

	offset := (crrPage - 1) * limit

	contestId := r.PathValue("contestId")
	if contestId == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Contest ID is required", nil)
		return
	}

	currentTime := time.Now()

	h.mu.RLock()
	entry, exists := h.Last_standings[contestId]
	if exists && entry.timestamp != nil && currentTime.Sub(*entry.timestamp) < 15*time.Second {
		response := *entry.standings
		h.mu.RUnlock()

		totalStandings := len(response.Standings)

		start := offset
		end := offset + limit
		if start >= totalStandings {
			response.Standings = []UserStanding{}
		} else {
			if end > totalStandings {
				end = totalStandings
			}
			response.Standings = response.Standings[start:end]
		}

		response.TotalItem = totalStandings
		response.TotalPages = (totalStandings + limit - 1) / limit
		response.Limit = limit
		response.Page = crrPage

		utils.SendResponse(w, http.StatusOK, nil, response)
		return
	}
	h.mu.RUnlock()

	var contest models.Contest
	var contestProblems []models.ContestProblem
	var problemResults []models.ContestProblemResult
	var fetchErr error

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		fetchErr = h.db.Where("id = ?", contestId).First(&contest).Error
	}()

	go func() {
		defer wg.Done()
		fetchErr = h.db.Where("contest_id = ?", contestId).Order("\"index\" ASC").Find(&contestProblems).Error
	}()

	go func() {
		defer wg.Done()
		fetchErr = h.db.Where("contest_id = ?", contestId).
			Preload("User").
			Preload("Problem").
			Find(&problemResults).Error
	}()

	wg.Wait()

	if fetchErr != nil {
		log.Println(fetchErr)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch standings data", nil)
		return
	}

	// Build user standings from ContestProblemResult
	userData := make(map[string]*userStandingData)
	for _, pr := range problemResults {
		if _, exists := userData[pr.UserId]; !exists {
			userData[pr.UserId] = &userStandingData{
				User:     &pr.User,
				Problems: make(map[string]*models.ContestProblemResult),
			}
		}
		userData[pr.UserId].Problems[pr.ProblemId] = &pr
	}

	standings := make([]UserStanding, 0, len(userData))
	for userId, ud := range userData {
		solvedCount := 0
		totalPenalty := 0
		var lastSolvedAt *time.Time

		problemStatuses := make([]ProblemStatus, 0, len(contestProblems))
		for _, cp := range contestProblems {
			pr, hasResult := ud.Problems[cp.ProblemID]
			status := ProblemStatus{
				Solved:     false,
				Attempts:   0,
				Penalty:    0,
				FirstBlood: false,
			}

			if hasResult {
				if pr.IsSolved {
					solvedCount++
					totalPenalty += pr.Penalty
					if lastSolvedAt == nil || (pr.SolvedAt != nil && pr.SolvedAt.After(*lastSolvedAt)) {
						lastSolvedAt = pr.SolvedAt
					}
				}
				status.Solved = pr.IsSolved
				status.Attempts = pr.WrongAttempts
				status.Penalty = pr.Penalty
				status.FirstBlood = pr.IsFirstBlood
				if pr.SolvedAt != nil {
					status.FirstSolvedAt = pr.SolvedAt
				}
			}

			problemStatuses = append(problemStatuses, status)
		}

		standings = append(standings, UserStanding{
			UserId:       userId,
			Username:     ud.User.Username,
			Name:         ud.User.Name,
			SolvedCount:  solvedCount,
			TotalPenalty: totalPenalty,
			Problems:     problemStatuses,
			LastSolvedAt: lastSolvedAt,
		})
	}

	// Sort standings: by solved count desc, then penalty asc, then last solved asc, then user id asc
	sortUserStandings(standings)

	// Build problem mapping
	problemMapping := make(map[int]string)
	for i, cp := range contestProblems {
		problemMapping[i+1] = cp.ProblemID
	}

	// Compute problem solve status
	problemSolveStatus := make(map[int]ProblemSolveStatus)
	for i := range contestProblems {
		problemSolveStatus[i+1] = ProblemSolveStatus{Solved: 0, Attempted: 0}
	}
	for _, pr := range problemResults {
		cpIndex := 0
		for i, cp := range contestProblems {
			if cp.ProblemID == pr.ProblemId {
				cpIndex = i + 1
				break
			}
		}
		if ps, exists := problemSolveStatus[cpIndex]; exists {
			if pr.IsSolved {
				ps.Solved++
			}
			ps.Attempted++
			problemSolveStatus[cpIndex] = ps
		}
	}

	response := StandingsResponse{
		ContestId:          contestId,
		ContestTitle:       contest.Title,
		ProblemMapping:     problemMapping,
		Standings:          standings,
		StartTime:          contest.StartTime,
		DurationSeconds:    contest.DurationSeconds,
		ProblemSolveStatus: problemSolveStatus,
		TotalItem:          len(standings),
		TotalPages:         (len(standings) + limit - 1) / limit,
		Limit:              limit,
		Page:               crrPage,
	}

	h.mu.Lock()
	h.Last_standings[contestId] = struct {
		timestamp *time.Time
		standings *StandingsResponse
	}{
		timestamp: &currentTime,
		standings: &response,
	}
	h.mu.Unlock()

	totalStandings := len(response.Standings)
	start := offset
	end := offset + limit
	if start >= totalStandings {
		response.Standings = []UserStanding{}
	} else {
		if end > totalStandings {
			end = totalStandings
		}
		response.Standings = response.Standings[start:end]
	}

	utils.SendResponse(w, http.StatusOK, nil, response)
}

type userStandingData struct {
	User     *models.User
	Problems map[string]*models.ContestProblemResult
}

func sortUserStandings(standings []UserStanding) {
	for i := 0; i < len(standings); i++ {
		for j := i + 1; j < len(standings); j++ {
			if lessUserStanding(standings[j], standings[i]) {
				standings[i], standings[j] = standings[j], standings[i]
			}
		}
	}
}

func lessUserStanding(a, b UserStanding) bool {
	if a.SolvedCount != b.SolvedCount {
		return a.SolvedCount > b.SolvedCount
	}
	if a.TotalPenalty != b.TotalPenalty {
		return a.TotalPenalty < b.TotalPenalty
	}
	if a.LastSolvedAt != nil && b.LastSolvedAt != nil {
		if !a.LastSolvedAt.Equal(*b.LastSolvedAt) {
			return a.LastSolvedAt.Before(*b.LastSolvedAt)
		}
	}
	if a.LastSolvedAt == nil && b.LastSolvedAt != nil {
		return true
	}
	if a.LastSolvedAt != nil && b.LastSolvedAt == nil {
		return false
	}
	return a.UserId < b.UserId
}
