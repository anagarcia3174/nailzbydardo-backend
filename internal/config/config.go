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
	FrontendURL       string
	CalendarSecret    string
	BaseURL           string
}

func Load(path string) (*Config, error) {
	err := godotenv.Load(path)
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
	FrontendURL := os.Getenv("FRONTEND_URL")
	if FrontendURL == "" {
		return nil, errors.New("Missing frontend url")
	}
	CalendarSecret := os.Getenv("CALENDAR_SECRET")
	if CalendarSecret == "" {
		return nil, errors.New("Missing calendar secret")
	}
	BaseURL := os.Getenv("BASE_URL")
	if BaseURL == "" {
		return nil, errors.New("Missing base url")
	}
	config := Config{
		DatabaseURL:       DatabaseURL,
		Port:              PORT,
		SessionCookieName: SessionCookieName,
		AppEnv:            AppEnv,
		FrontendURL:       FrontendURL,
		CalendarSecret:    CalendarSecret,
		BaseURL:           BaseURL,
	}
	return &config, nil
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}