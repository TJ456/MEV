package utils

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file
func LoadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Println(".env file not found. Falling back to system environment variables.")
    } else {
        log.Println("Environment variables loaded from .env file.")
    }
}

// GetEnv returns the value of the given key
func GetEnv(key string) string {
    value := os.Getenv(key)
    if value == "" {
        log.Fatalf("Environment variable %s not set", key)
    }
    return value
}
