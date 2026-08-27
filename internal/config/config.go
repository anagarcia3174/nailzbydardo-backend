package config

// imports you'll need:
// - "os" for reading env vars
// - "log" or "fmt" for error messages if something required is missing
// - "strconv" if you convert PORT from string to int (optional — decide if you want PORT as string or int)
// - "github.com/joho/godotenv" for loading the .env file
import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Define a struct that represents your application's configuration.
// Fields should mirror the env vars you decided on:
//   - DatabaseURL (string)
//   - Port (string or int — decide based on how you'll use it later, e.g. http.ListenAndServe(":"+port) wants a string like ":8080")
//   - SessionCookieName (string)

type Config struct {
	DatabaseURL       string
	Port              string
	SessionCookieName string
}

// Write a function, e.g. Load() (*Config, error)
// Responsibilities, in order:
//
// 1. Attempt to load the .env file using godotenv.
//    - Think about what should happen if the .env file doesn't exist.
//      In production there won't be one — you don't want that to be treated as fatal.
//      Only log/handle it, don't necessarily return an error for a missing file.
//
// 2. Read each required env var using os.Getenv.
//
// 3. For DatabaseURL and SessionCookieName (your "required, no default" fields):
//    - Check if the value is an empty string.
//    - If empty, this should be treated as a fatal configuration error.
//    - Decide how you want to surface that — returning an error from Load()
//      that main.go checks and exits on, vs. calling log.Fatal directly inside
//      Load(). Returning an error is generally the cleaner pattern, since it
//      keeps config.go free of "exit the whole program" decisions — that's
//      arguably main()'s job.
//
// 4. For Port (your "has a sensible default" field):
//    - Check if it's empty.
//    - If empty, fall back to a hardcoded default ("8080").
//    - If present, use the provided value.
//
// 5. Populate your Config struct with the resolved values.
//
// 6. Return the populated struct (and nil error, if using the error-return pattern).

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

	config := Config{
		DatabaseURL: DatabaseURL,
		Port: PORT,
		SessionCookieName: SessionCookieName,
	}
	return &config, nil
}

// Optional: a small private helper function for "get env var or fail" logic,
// since you'll repeat that pattern for both DatabaseURL and SessionCookieName.
// Something like: getRequiredEnv(key string) (string, error)
// that centralizes "look it up, check if empty, return a descriptive error
// naming which variable is missing" so you're not duplicating that logic twice.
