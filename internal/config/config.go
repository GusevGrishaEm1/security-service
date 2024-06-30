package config

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port        int
	StoragePath string
	SecretKey   string
	TokenTTL    time.Duration
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var (
		port        = flag.Int("port", getEnvAsInt("PORT", 8080), "Server port")
		storagePath = flag.String("storage-path", getEnv("STORAGE_PATH", "./data"), "Path to storage")
		secretKey   = flag.String("secret-key", getEnv("SECRET_KEY", "default_secret"), "Secret key")
		tokenTTL    = flag.Duration("token-ttl", getEnvAsDuration("TOKEN_TTL", time.Hour), "Token time-to-live")
	)

	flag.Parse()

	config := &Config{
		Port:        *port,
		StoragePath: *storagePath,
		SecretKey:   *secretKey,
		TokenTTL:    *tokenTTL,
	}

	return config, nil
}

// Helper functions to read environment variables with default values.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(name string, defaultValue int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(name string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(name, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
