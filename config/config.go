package config

import (
	"log/slog"
	"os"

	env "github.com/joho/godotenv"
)

type Config struct {
	HttpPort     string
	SecretKey    string
	DBURL        string
	EngineKey    string
	EngineUrl    string
	WorkerCount  int
	QueueName    string
	RabbitMQURL  string
	CookieSecure bool
}

var configuration *Config
var err error

func loadConfig() (*Config, error) {
	var config Config

	env.Load()

	http_port := os.Getenv("HTTP_PORT")
	if http_port == "" {
		slog.Error("HTTP_PORT not defined")
		os.Exit(1)
	}

	config.HttpPort = http_port

	secret_key := os.Getenv("JWT_SECRET")
	if secret_key == "" {
		slog.Error("JWT_SECRET not defined")
		os.Exit(1)
	}
	config.SecretKey = secret_key

	engine_key := os.Getenv("ENGINE_KEY")
	if engine_key == "" {
		slog.Error("ENGINE_KEY not defined")
		os.Exit(1)
	}
	config.EngineKey = engine_key

	engine_url := os.Getenv("ENGINE_URL")
	if engine_url == "" {
		slog.Error("ENGINE_URL not defined")
		os.Exit(1)
	}
	config.EngineUrl = engine_url

	// Configure database
	config.DBURL = os.Getenv("DB_URL")
	if config.DBURL == "" {
		slog.Error("DB_URL not defined")
		os.Exit(1)
	}

	config.QueueName = os.Getenv("QUEUE_NAME")
	if config.QueueName == "" {
		config.QueueName = "judge_queue"
		slog.Info("QUEUE_NAME not set, using default", "queue_name", "judge_queue")
	}

	config.RabbitMQURL = os.Getenv("RABBITMQ_URL")
	if config.RabbitMQURL == "" {
		config.RabbitMQURL = "amqp://guest:guest@localhost:5672/"
		slog.Info("RABBITMQ_URL not set, using default", "url", "amqp://guest:guest@localhost:5672/")
	}

	config.CookieSecure = os.Getenv("COOKIE_SECURE") == "true"

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
