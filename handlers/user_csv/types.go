package usercsv

import (
	"encoding/csv"
	"sync"
)

type Handler struct {
	contestId  *int64
	clanLength int
	prefix     string
	writer     *csv.Writer
	mu         sync.Mutex
	FilePath   string
}
