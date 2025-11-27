package usercsv

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

func NewHandler(prefix string, clanLength int, contestId int64) (*Handler, error) {
	const dir = "./generated_csv/"
	fileName := prefix + "_users.csv"
	filePath := filepath.Join(dir, fileName)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return nil, err
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil, err
	}

	h := &Handler{
		writer:     csv.NewWriter(file),
		prefix:     prefix,
		clanLength: clanLength,
		contestId:  &contestId,
		FilePath:   filePath,
	}
	return h, nil
}
