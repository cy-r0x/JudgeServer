package standings

import (
	"log"
	"net/http"
	"sync"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

type ExportedUserStanding struct {
	Name       string `json:"name"`
	SolveCount int    `json:"solve_count"`
	Solved     []int  `json:"solved"`
}

func (h *Handler) ExportStandings(w http.ResponseWriter, r *http.Request) {
	contestIdStr := r.PathValue("contestId")
	if contestIdStr == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Contest ID is required", nil)
		return
	}
	contestId := contestIdStr

	var contestProblems []models.ContestProblem
	var problemResults []models.ContestProblemResult
	var fetchErr error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		fetchErr = h.db.Where("contest_id = ?", contestId).Order("\"index\" ASC").Find(&contestProblems).Error
	}()

	go func() {
		defer wg.Done()
		fetchErr = h.db.Where("contest_id = ?", contestId).
			Preload("User").
			Find(&problemResults).Error
	}()

	wg.Wait()

	if fetchErr != nil {
		log.Println(fetchErr)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch standings data", nil)
		return
	}

	// Build user -> solved problems mapping
	type userExportedData struct {
		Name       string
		SolvedCount int
		Solved      map[int]bool // index -> solved
	}

	userData := make(map[string]*userExportedData)
	for _, pr := range problemResults {
		if _, exists := userData[pr.UserId]; !exists {
			userData[pr.UserId] = &userExportedData{
				Name:  pr.User.Name,
				Solved: make(map[int]bool),
			}
		}
		if pr.IsSolved {
			userData[pr.UserId].SolvedCount++
		}
	}

	// Map problem ID to index for each user
	for _, pr := range problemResults {
		if pr.IsSolved {
			for i, cp := range contestProblems {
				if cp.ProblemID == pr.ProblemId {
					userData[pr.UserId].Solved[i+1] = true
					break
				}
			}
		}
	}

	exportedStandings := make([]ExportedUserStanding, 0, len(userData))
	for _, ud := range userData {
		solved := make([]int, 0, len(ud.Solved))
		for idx := range ud.Solved {
			solved = append(solved, idx)
		}

		exportedStandings = append(exportedStandings, ExportedUserStanding{
			Name:       ud.Name,
			SolveCount: ud.SolvedCount,
			Solved:     solved,
		})
	}

	utils.SendResponse(w, http.StatusOK, nil, exportedStandings)
}