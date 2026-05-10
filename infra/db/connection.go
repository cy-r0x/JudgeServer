package db

import (
	"log/slog"

	"github.com/judgenot0/judge-backend/config"
	"github.com/judgenot0/judge-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	dns := cfg.DBURL

	dbCon, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		slog.Error("Database connection error", "error", err)
		return nil, err
	}
	return dbCon, nil
}

func Migrate(dbConn *gorm.DB) error {
	if err := dbConn.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto").Error; err != nil {
		slog.Error("Failed to enable pgcrypto extension", "error", err)
		return err
	}

	err := dbConn.AutoMigrate(
		&models.User{},
		&models.Contest{},
		&models.Problem{},
		&models.ContestProblem{},
		&models.Submission{},
		&models.Testcase{},
		&models.UserCreds{},
		&models.ContestProblemResult{},
	)
	if err != nil {
		slog.Error("Failed to AutoMigrate", "error", err)
		return err
	}

	// Create default admin user if it doesn't exist
	var adminCount int64
	if err := dbConn.Model(&models.User{}).Where("username = ?", "admin").Count(&adminCount).Error; err != nil {
		slog.Error("Failed to check if admin user exists", "error", err)
		return err
	}

	if adminCount == 0 {
		adminUser := models.User{
			Name:     "admin",
			Username: "admin",
			Password: "$2a$12$Ncde3vjx7AbBXwyDlzgN5ue8PKgD1XexbvWdityKLbQHsHJAi1jKG",
			Role:     models.RoleAdmin,
		}
		if err := dbConn.Create(&adminUser).Error; err != nil {
			slog.Error("Failed to create default admin user", "error", err)
			return err
		}
		slog.Info("Default admin user created successfully")
	}

	slog.Info("GORM AutoMigration Done")
	return nil
}
