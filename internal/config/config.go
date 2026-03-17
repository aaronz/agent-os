package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Graph    GraphConfig
	VectorDB VectorDBConfig
	Kafka    KafkaConfig
	Tracing  TracingConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type GraphConfig struct {
	URI      string
	User     string
	Password string
}

type VectorDBConfig struct {
	Provider string // "pinecone", "weaviate", "local"
	Endpoint string
	APIKey   string
	Index    string
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

type TracingConfig struct {
	Enabled     bool
	Endpoint    string
	ServiceName string
}

type AuthConfig struct {
	APIKeyPrefix string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "agentos"),
			Password: getEnv("DB_PASSWORD", "agentos"),
			DBName:   getEnv("DB_NAME", "agentos"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Graph: GraphConfig{
			URI:      getEnv("GRAPH_URI", "neo4j://localhost:7687"),
			User:     getEnv("GRAPH_USER", "neo4j"),
			Password: getEnv("GRAPH_PASSWORD", "password"),
		},
		VectorDB: VectorDBConfig{
			Provider: getEnv("VECTOR_PROVIDER", "local"),
			Endpoint: getEnv("VECTOR_ENDPOINT", ""),
			APIKey:   getEnv("VECTOR_API_KEY", ""),
			Index:    getEnv("VECTOR_INDEX", "memories"),
		},
		Kafka: KafkaConfig{
			Brokers: []string{getEnv("KAFKA_BROKER", "localhost:9092")},
			Topic:   getEnv("KAFKA_TOPIC", "agentos-events"),
		},
		Tracing: TracingConfig{
			Enabled:     getEnvBool("TRACING_ENABLED", true),
			Endpoint:    getEnv("TRACING_ENDPOINT", "http://localhost:14268"),
			ServiceName: "agent-os",
		},
		Auth: AuthConfig{
			APIKeyPrefix: "sk-",
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}
