package config

import (
	"errors"
	"log"
	"os"

	env "github.com/joho/godotenv"
)

type DBConfig struct {
	DB_HOST         string
	DB_PORT         string
	DB_NAME         string
	DB_USER         string
	DB_PASSWORD     string
	ENABLE_SSL_MODE string
}

type Config struct {
	HttpPort  string
	SecretKey string
	DB        *DBConfig
	EngineKey string
}

var configuration *Config
var err error

func loadConfig() (*Config, error) {
	err := env.Load()
	if err != nil {
		log.Fatalln("Env Not Found, please make sure there is .env file in the root directory")
		return nil, errors.New("Env Not Found, please make sure there is .env file in the root directory")
	}

	var config Config

	http_port := os.Getenv("HTTP_PORT")
	if http_port == "" {
		log.Fatalln("HTTP_PORT not defined")
		return nil, errors.New("HTTP_PORT not defined")
	}

	config.HttpPort = http_port

	secret_key := os.Getenv("JWT_SECRET")
	if secret_key == "" {
		log.Fatalln("JWT_SECRET not defined")
		return nil, errors.New("JWT_SECRET not defined")
	}
	config.SecretKey = secret_key

	engine_key := os.Getenv("ENGINE_KEY")
	if engine_key == "" {
		log.Fatalln("ENGINE_KEY not defined")
		return nil, errors.New("ENGINE_KEY not defined")
	}
	config.EngineKey = engine_key

	// Configure database
	dbConfig := &DBConfig{
		DB_HOST:         os.Getenv("DB_HOST"),
		DB_PORT:         os.Getenv("DB_PORT"),
		DB_NAME:         os.Getenv("DB_NAME"),
		DB_USER:         os.Getenv("DB_USER"),
		DB_PASSWORD:     os.Getenv("DB_PASSWORD"),
		ENABLE_SSL_MODE: os.Getenv("DB_SSL_MODE"),
	}

	// Validate required DB settings
	if dbConfig.DB_HOST == "" {
		log.Fatalln("DB_HOST not defined")
		return nil, errors.New("DB_HOST not defined")
	}
	if dbConfig.DB_PORT == "" {
		log.Fatalln("DB_PORT not defined")
		return nil, errors.New("DB_PORT not defined")
	}
	if dbConfig.DB_NAME == "" {
		log.Fatalln("DB_NAME not defined")
		return nil, errors.New("DB_NAME not defined")
	}
	if dbConfig.DB_USER == "" {
		log.Fatalln("DB_USER not defined")
		return nil, errors.New("DB_USER not defined")
	}
	if dbConfig.DB_PASSWORD == "" {
		log.Fatalln("DB_PASSWORD not defined")
		return nil, errors.New("DB_PASSWORD not defined")
	}

	// Set default SSL mode if not provided
	if dbConfig.ENABLE_SSL_MODE == "" {
		dbConfig.ENABLE_SSL_MODE = "disable"
	}

	config.DB = dbConfig

	return &config, nil
}

func GetConfig() (*Config, error) {
	if configuration == nil {
		configuration, err = loadConfig()
		if err != nil {
			return nil, err
		}
	}
	return configuration, nil
}
