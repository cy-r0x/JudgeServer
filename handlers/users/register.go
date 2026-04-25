package users

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var reqUser models.User
	if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	if reqUser.Username == "" || reqUser.Password == "" {
		utils.SendResponse(w, http.StatusBadRequest, "username and password are required", nil)
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		utils.SendResponse(w, http.StatusInternalServerError, "Internal Server Error", nil)
		return
	}

	var allowedContest *string
	if reqUser.AllowedContest != nil {
		ac := *reqUser.AllowedContest
		allowedContest = &ac
	}

	role := reqUser.Role
	if role == "" {
		role = models.RoleUser
	}

	newUser := models.User{
		Name:           reqUser.Name,
		Username:       reqUser.Username,
		Password:       string(hashedPassword),
		Role:           role,
		AdditionalInfo: reqUser.AdditionalInfo,
		RoomNo:         reqUser.RoomNo,
		PcNo:           reqUser.PcNo,
		AllowedContest: allowedContest,
		CreatedAt:      time.Now(),
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newUser).Error; err != nil {
			return err
		}

		if newUser.AllowedContest != nil {
			creds := models.UserCreds{
				ContestId:     *newUser.AllowedContest,
				UserId:        newUser.Id,
				PlainPassword: reqUser.Password,
			}
			if err := tx.Create(&creds).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Println("DB Transaction Error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Internal Server Error", nil)
		return
	}

	reqUser.Id = newUser.Id
	reqUser.Password = ""
	reqUser.Role = newUser.Role
	reqUser.CreatedAt = newUser.CreatedAt

	utils.SendResponse(w, http.StatusCreated, "User created successfully", reqUser)
}
