package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/vicodevv/relay/internal/http"
	"github.com/vicodevv/relay/internal/storage"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	dbConfig := storage.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5433"),
		User:     getEnv("DB_USER", "relay"),
		Password: getEnv("DB_PASSWORD", "relay123"),
		DBName:   getEnv("DB_NAME", "relay_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	db, err := storage.NewPostgresDB(dbConfig)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := storage.NewWorkflowRepository(db)

	port := getEnv("APP_PORT", "8080")
	server := http.NewServer(repo, port)

	if err := server.Start(); err != nil {
		logrus.Fatalf("Server error: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
