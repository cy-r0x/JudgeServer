package config

import (
	"errors"
	"log"
	"os"

	env "github.com/joho/godotenv"
)

type Config struct {
	HttpPort    string
	SecretKey   string
	DBURL       string
	EngineKey   string
	EngineUrl   string
	WorkerCount int
	QueueName   string
	RabbitMQURL string
}

var configuration *Config
var err error

func loadConfig() (*Config, error) {
	var config Config

	env.Load()

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

	engine_url := os.Getenv("ENGINE_URL")
	if engine_url == "" {
		log.Fatalln("ENGINE_URL not defined")
		return nil, errors.New("ENGINE_URL not defined")
	}
	config.EngineUrl = engine_url

	// Configure database
	config.DBURL = os.Getenv("DB_URL")
	if config.DBURL == "" {
		log.Fatalln("DB_URL not defined")
		return nil, errors.New("DB_URL not defined")
	}

	config.QueueName = os.Getenv("QUEUE_NAME")
	if config.QueueName == "" {
		config.QueueName = "judge_queue"
		log.Println("QUEUE_NAME not set, using default: judge_queue")
	}

	config.RabbitMQURL = os.Getenv("RABBITMQ_URL")
	if config.RabbitMQURL == "" {
		config.RabbitMQURL = "amqp://guest:guest@localhost:5672/"
		log.Println("RABBITMQ_URL not set, using default: amqp://guest:guest@localhost:5672/")
	}

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
