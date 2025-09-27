package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

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

func Migrate(dbConn *sqlx.DB) error {
	// Get all SQL schema files
	schemaDir := "infra/db/schema"
	files, err := os.ReadDir(schemaDir)
	if err != nil {
		log.Printf("Failed to read schema directory: %v", err)
		return err
	}

	// Filter and sort SQL files
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles) // Execute in order

	log.Printf("Found %d SQL schema files to check", len(sqlFiles))

	// Iterate over SQL files
	for _, filename := range sqlFiles {
		filePath := filepath.Join(schemaDir, filename)

		// Derive table name from file name (strip extension)
		tableName := strings.TrimSuffix(filename, filepath.Ext(filename))[4:]
		// Check if table exists in PostgreSQL
		var exists bool
		err = dbConn.Get(&exists,
			`SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_schema = 'public'
				AND table_name = $1
			)`, tableName,
		)
		if err != nil {
			log.Printf("Failed to check if table %s exists: %v", tableName, err)
			return err
		}

		if exists {
			log.Printf("Skipping %s (table '%s' already exists)", filename, tableName)
			continue
		}

		log.Printf("Executing schema file: %s", filename)

		// Read the SQL file
		sqlContent, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to read SQL file %s: %v", filename, err)
			return err
		}

		// Execute the SQL
		_, err = dbConn.Exec(string(sqlContent))
		if err != nil {
			log.Printf("Failed to execute SQL from file %s: %v", filename, err)
			return err
		}

		log.Printf("Successfully executed: %s", filename)
	}

	log.Println("Database migration completed successfully!")
	log.Printf("Checked %d schema files in total", len(sqlFiles))
	return nil
}
