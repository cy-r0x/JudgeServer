package users

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var reqUser User
	if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// basic validation
	if reqUser.Username == "" || reqUser.Password == "" {
		utils.SendResponse(w, http.StatusBadRequest, "username and password are required")
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		utils.SendResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	var allowedContest *string
	if reqUser.AllowedContest != nil {
		ac := *reqUser.AllowedContest
		allowedContest = &ac
	}

	role := reqUser.Role
	if role == "" {
		role = "user"
	}

	newUser := models.User{
		FullName:       reqUser.FullName,
		Username:       reqUser.Username,
		Password:       string(hashedPassword),
		Role:           role,
		Clan:           reqUser.Clan,
		RoomNo:         reqUser.RoomNo,
		PcNo:           reqUser.PcNo,
		AllowedContest: allowedContest,
		CreatedAt:      time.Now(),
	}

	err = h.db.Create(&newUser).Error
	if err != nil {
		log.Println("DB Create Error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	reqUser.Id = newUser.ID
	reqUser.Password = ""
	reqUser.Role = newUser.Role
	reqUser.CreatedAt = newUser.CreatedAt

	utils.SendResponse(w, http.StatusCreated, reqUser)
}
