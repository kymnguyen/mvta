package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ServiceName   string
	Port          string
	MongoURI      string
	MongoDB       string
	WorkflowDir   string
	KafkaBrokers  []string
	KafkaTopic    string
	KafkaGroupID  string
	KafkaDLQTopic string
	JWTSecret     string
	JaegerURL     string
	TimeoutWorker TimeoutWorkerConfig
}

type TimeoutWorkerConfig struct {
	Enabled   bool
	Interval  time.Duration
	BatchSize int
}

func Load() *Config {
	return &Config{
		ServiceName:   getEnv("SERVICE_NAME", "workflow-svc"),
		Port:          getEnv("PORT", "50003"),
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:       getEnv("MONGO_DB", "workflow_db"),
		WorkflowDir:   getEnv("WORKFLOW_DIR", "./config/workflows"),
		KafkaBrokers:  strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaTopic:    getEnv("KAFKA_TOPIC", "vehicle.events"),
		KafkaGroupID:  getEnv("KAFKA_GROUP_ID", "workflow-svc"),
		KafkaDLQTopic: getEnv("KAFKA_DLQ_TOPIC", "workflow.dlq"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		JaegerURL:     getEnv("JAEGER_URL", "http://localhost:14268/api/traces"),
		TimeoutWorker: TimeoutWorkerConfig{
			Enabled:   getEnvBool("TIMEOUT_WORKER_ENABLED", true),
			Interval:  getEnvDuration("TIMEOUT_WORKER_INTERVAL", 30*time.Second),
			BatchSize: getEnvInt("TIMEOUT_WORKER_BATCH_SIZE", 100),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}
