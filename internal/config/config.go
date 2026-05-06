package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppEnv                      string
	ServiceName                 string
	LogLevel                    string
	AWSRegion                   string
	SQSQueueURL                 string
	SQSWaitTimeSeconds          int
	SQSMaxMessages              int
	SQSVisibilityTimeoutSeconds int
	PostgresDSN                 string
	FCMEnabled                  bool
	FCMProjectID                string
	FCMCredentialsJSON          string
	APNSEnabled                 bool
}

func Load() Config {
	return Config{
		AppEnv:                      getEnv("APP_ENV", "local"),
		ServiceName:                 getEnv("SERVICE_NAME", "notification-service"),
		LogLevel:                    getEnv("LOG_LEVEL", "info"),
		AWSRegion:                   getEnv("AWS_REGION", "ca-central-1"),
		SQSQueueURL:                 getEnv("SQS_QUEUE_URL", ""),
		SQSWaitTimeSeconds:          getEnvInt("SQS_WAIT_TIME_SECONDS", 20),
		SQSMaxMessages:              getEnvInt("SQS_MAX_MESSAGES", 10),
		SQSVisibilityTimeoutSeconds: getEnvInt("SQS_VISIBILITY_TIMEOUT_SECONDS", 30),
		PostgresDSN:                 getEnv("POSTGRES_DSN", ""),
		FCMEnabled:                  getEnvBool("FCM_ENABLED", true),
		FCMProjectID:                getEnv("FCM_PROJECT_ID", ""),
		FCMCredentialsJSON:          getEnv("FCM_CREDENTIALS_JSON", ""),
		APNSEnabled:                 getEnvBool("APNS_ENABLED", false),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}
