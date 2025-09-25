package config

import (
	"errors"
	"log"
	"os"

	env "github.com/joho/godotenv"
)

type Config struct {
	HttpPort  string
	SecretKey string
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
