package config

import (
	"os"
)

type Config struct {
	Port      string
	MongoURI  string
	JWTSecret string
}

func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "50000"),
		MongoURI:  getEnv("MONGO_URI", "mongodb://mongo:password@localhost:27017/auth_db?authSource=admin"),
		JWTSecret: getEnv("JWT_SECRET", "supersecret"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
