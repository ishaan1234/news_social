package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port         string
	DBUrl        string
	JWTSecret    string
	RateLimitRPS int
	Environment  string
	OpenAIAPIKey string
	NewsAPIKey   string
}

func Load() *Config {
	rateLimit, err := strconv.Atoi(getEnv("RATE_LIMIT_RPS", "10"))
	if err != nil {
		log.Println("Invalid RATE_LIMIT_RPS, defaulting to 10")
		rateLimit = 10
	}

	return &Config{
		Port:         getEnv("PORT", "8080"),
		DBUrl:        getEnv("DATABASE_URL", "postgres://localhost:5432/newsapp?sslmode=disable"),
		JWTSecret:    getEnv("JWT_SECRET", "supersecret"),
		RateLimitRPS: rateLimit,
		Environment:  getEnv("ENV", "development"),
		OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),
		NewsAPIKey:   getEnv("NEWS_API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return fallback
}
