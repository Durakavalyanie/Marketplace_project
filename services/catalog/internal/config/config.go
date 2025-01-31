package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	CatalogPort string
	SSLMode     string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return Config{
		DBHost:      getEnv("DB_HOST"),
		DBPort:      getEnv("DB_PORT"),
		DBUser:      getEnv("DB_USER"),
		DBPassword:  getEnv("DB_PASSWORD"),
		DBName:      getEnv("DB_NAME"),
		CatalogPort: getEnv("CATALOG_PORT"),
		SSLMode:     getEnv("SSL_MODE"),
	}
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return "NONE"
}
