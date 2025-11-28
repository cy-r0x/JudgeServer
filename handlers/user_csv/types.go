package usercsv

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	Id               int64     `json:"id" db:"id"`
	FullName         string    `json:"full_name" db:"full_name"`
	Username         string    `json:"username" db:"username"`
	Password         string    `json:"password,omitempty" db:"password"`
	UnHashedPassword string    `json:"-" db:"-"`
	Role             string    `json:"role,omitempty" db:"role"`
	Clan             *string   `json:"clan,omitempty" db:"clan"`
	RoomNo           *string   `json:"room_no,omitempty" db:"room_no"`
	PcNo             *string   `json:"pc_no,omitempty" db:"pc_no"`
	AllowedContest   *int64    `json:"allowed_contest,omitempty" db:"allowed_contest"`
	CreatedAt        time.Time `json:"created_at,omitempty" db:"created_at"`
}

// type Handler struct {
// 	contestId  *int64
// 	clanLength int
// 	prefix     string
// 	mu         sync.Mutex
// 	FilePath   string
// }

type WriterHandler struct {
	contestId  *int64
	clanLength int
	prefix     string
	writer     *csv.Writer
	FilePath   string
}

type Handler struct {
	db     *sqlx.DB
	mu     sync.Mutex
	Writer *WriterHandler
}

func (h *Handler) NewWriteHandler(prefix string, clanLength int, contestId int64) error {
	const dir = "./generated_csv/"
	fileName := prefix + "_users.csv"
	filePath := filepath.Join(dir, fileName)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}

	h.Writer = &WriterHandler{
		writer:     csv.NewWriter(file),
		prefix:     prefix,
		clanLength: clanLength,
		contestId:  &contestId,
		FilePath:   filePath,
	}
	return nil
}
