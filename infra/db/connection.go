package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/judgenot0/judge-backend/config"
	_ "github.com/lib/pq"
)

func GetConnectionString(cfg *config.DBConfig) string {
	user := cfg.DB_USER
	password := cfg.DB_PASSWORD
	host := cfg.DB_HOST
	port := cfg.DB_PORT
	dbname := cfg.DB_NAME
	sslmode := "disable"
	if cfg.ENABLE_SSL_MODE == "true" {
		sslmode = "require"
	}
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s", user, password, host, port, dbname, sslmode)
}

func NewConnection(cfg *config.DBConfig) (*sqlx.DB, error) {
	dbSource := GetConnectionString(cfg)
	dbCon, err := sqlx.Connect("postgres", dbSource)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return dbCon, nil
}
