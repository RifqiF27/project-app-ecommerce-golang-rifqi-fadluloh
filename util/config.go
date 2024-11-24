package util

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Configuration struct {
	AppName string
	Port    string
	Debug   bool
	DB      DatabaseConfig
}

type DatabaseConfig struct {
	Name     string
	Username string
	Password string
	Host     string
}

func ReadConfiguration() Configuration {
	if err := godotenv.Load(); err != nil {
        log.Println("No .env file found. Using system environment variables.")
    }
	debug := false
	if os.Getenv("DEBUG") == "true" {
		debug = true
	}

	return Configuration{
		AppName: os.Getenv("APP_NAME"),
		Port:    os.Getenv("PORT"),
		Debug:   debug,
		DB: DatabaseConfig{
			Name:     os.Getenv("DATABASE_NAME"),
			Username: os.Getenv("DATABASE_USERNAME"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			Host:     os.Getenv("DATABASE_HOST"),
		},
	}
}
