package config

import "os"

type Config struct {
	OpenAIKey        string
	ClaudeKey        string
	VectorEngineAddr string
	QdrantAddr       string
	Port             string
}

func Load() *Config {
	return &Config{
		OpenAIKey:        os.Getenv("OPENAI_API_KEY"),
		ClaudeKey:        os.Getenv("ANTHROPIC_API_KEY"),
		VectorEngineAddr: getEnvOrDefault("VECTOR_ENGINE_ADDR", "localhost:50051"),
		QdrantAddr:       getEnvOrDefault("QDRANT_ADDR", "localhost:6334"),
		Port:             getEnvOrDefault("PORT", "8080"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
