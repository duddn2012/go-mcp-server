package config

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var (
	configInstance *Config
	once           sync.Once
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
		configInstance = load()
	})

	return configInstance
}

func load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// ALLOWED_ORIGINS를 쉼표로 구분된 문자열에서 슬라이스로 변환
	allowedOriginsStr := getEnv("ALLOWED_ORIGINS", "http://localhost")
	allowedOrigins := strings.Split(allowedOriginsStr, ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}

	return &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "mcp_server"),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		AllowedOrigins: allowedOrigins,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
