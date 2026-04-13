package db

import (
	"fmt"
	"log"

	"github.com/judgenot0/judge-backend/config"
	"github.com/judgenot0/judge-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func NewConnection(cfg *config.DBConfig) (*gorm.DB, error) {
	dbSource := GetConnectionString(cfg)
	dbCon, err := gorm.Open(postgres.Open(dbSource), &gorm.Config{})
	if err != nil {
		log.Println("Database connection error:", err)
		return nil, err
	}
	return dbCon, nil
}

func Migrate(dbConn *gorm.DB) error {
	if err := dbConn.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto").Error; err != nil {
		log.Println("Failed to enable pgcrypto extension:", err)
		return err
	}

	err := dbConn.AutoMigrate(
		&models.User{},
		&models.Contest{},
		&models.Problem{},
		&models.ContestProblem{},
		&models.Submission{},
		&models.Testcase{},
		&models.Filepath{},
		&models.ContestStanding{},
		&models.ContestSolve{},
		&models.ContestUserProblem{},
		&models.ContestProblemStat{},
	)
	if err != nil {
		log.Println("Failed to AutoMigrate:", err)
		return err
	}

	// Create default admin user if it doesn't exist
	var adminCount int64
	if err := dbConn.Model(&models.User{}).Where("username = ?", "admin").Count(&adminCount).Error; err != nil {
		log.Println("Failed to check if admin user exists:", err)
		return err
	}

	if adminCount == 0 {
		adminUser := models.User{
			FullName: "admin",
			Username: "admin",
			Password: "$2a$12$Ncde3vjx7AbBXwyDlzgN5ue8PKgD1XexbvWdityKLbQHsHJAi1jKG",
			Role:     "admin",
		}
		if err := dbConn.Create(&adminUser).Error; err != nil {
			log.Println("Failed to create default admin user:", err)
			return err
		}
		log.Println("Default admin user created successfully")
	}

	log.Println("GORM AutoMigration Done")
	return nil
}
