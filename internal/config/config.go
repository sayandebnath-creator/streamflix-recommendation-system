package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBSSLMode           string
	EmbeddingServiceURL string
	EmbeddingBatchSize  int
	EmbeddingWorkers int
	PosterValidationBatchSize int
	PosterValidationWorkers   int
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	batchSize, err := strconv.Atoi(os.Getenv("EMBEDDING_BATCH_SIZE"))
	if err != nil {
		batchSize = 100
	}

	workers, err := strconv.Atoi(os.Getenv("EMBEDDING_WORKERS"))
	if err != nil {
		workers = 5
	}

	posterBatchSize, err := strconv.Atoi(os.Getenv("POSTER_VALIDATION_BATCH_SIZE"))
	if err != nil {
		posterBatchSize = 100
	}

	posterWorkers, err := strconv.Atoi(os.Getenv("POSTER_VALIDATION_WORKERS"))
	if err != nil {
		posterWorkers = 10
	}

	return &Config{
		Port:                os.Getenv("PORT"),
		DBHost:              os.Getenv("DB_HOST"),
		DBPort:              os.Getenv("DB_PORT"),
		DBUser:              os.Getenv("DB_USER"),
		DBPassword:          os.Getenv("DB_PASSWORD"),
		DBName:              os.Getenv("DB_NAME"),
		DBSSLMode:           os.Getenv("DB_SSLMODE"),
		EmbeddingServiceURL: os.Getenv("EMBEDDING_SERVICE_URL"),
		EmbeddingBatchSize:  batchSize,
		EmbeddingWorkers: workers,
		PosterValidationBatchSize: posterBatchSize,
		PosterValidationWorkers:   posterWorkers,
	}
}
