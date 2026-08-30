package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL       string
	Port              string
	SessionCookieName string
	AppEnv            string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	DatabaseURL := os.Getenv("DATABASE_URL")
	if DatabaseURL == "" {
		return nil, errors.New("Missing database url")
	}
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	SessionCookieName := os.Getenv("SESSION_COOKIE_NAME")
	if SessionCookieName == "" {
		return nil, errors.New("Missing session cookie name")
	}
	AppEnv := os.Getenv("APP_ENV")
	if AppEnv == "" {
		AppEnv = "development"
	}
	config := Config{
		DatabaseURL:       DatabaseURL,
		Port:              PORT,
		SessionCookieName: SessionCookieName,
		AppEnv:            AppEnv,
	}
	return &config, nil
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}
