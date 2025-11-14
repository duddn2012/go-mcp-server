package config

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	ServerPort     string
	AllowedOrigins []string
}

func Get() *Config {
	once.Do(func() {
		instance = load()
	})
	return instance
}

func load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("[Config] .env file not found, using environment variables")
	}

	allowedOriginsStr := getEnvOrDefault("ALLOWED_ORIGINS", "http://localhost")
	allowedOrigins := parseAllowedOrigins(allowedOriginsStr)

	config := &Config{
		DBHost:         getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:         getEnvOrDefault("DB_PORT", "5432"),
		DBUser:         getEnvOrDefault("DB_USER", "postgres"),
		DBPassword:     getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:         getEnvOrDefault("DB_NAME", "mcp_server"),
		ServerPort:     getEnvOrDefault("SERVER_PORT", "8080"),
		AllowedOrigins: allowedOrigins,
	}

	log.Printf("[Config] Loaded: DB=%s:%s, Server=:%s, AllowedOrigins=%v",
		config.DBHost, config.DBPort, config.ServerPort, config.AllowedOrigins)

	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseAllowedOrigins(originsStr string) []string {
	origins := strings.Split(originsStr, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	return origins
}
