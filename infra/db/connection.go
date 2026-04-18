package db

import (
	"log"

	"github.com/judgenot0/judge-backend/config"
	"github.com/judgenot0/judge-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	dns := cfg.DBURL

	dbCon, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
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
