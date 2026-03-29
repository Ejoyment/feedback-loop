package main

import (
	"log"
	"os"

	"github.com/ai-dev-assistant/services/internal/api"
	"github.com/ai-dev-assistant/services/internal/config"
	"github.com/ai-dev-assistant/services/internal/llm"
	"github.com/ai-dev-assistant/services/internal/rag"
	"github.com/ai-dev-assistant/services/internal/testgen"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	llmGateway := llm.NewGateway(cfg.OpenAIKey, cfg.ClaudeKey)
	vectorClient := rag.NewVectorClient(cfg.VectorEngineAddr)
	ragService := rag.NewService(vectorClient, llmGateway)
	testGenerator := testgen.NewGenerator(llmGateway, ragService)

	router := gin.Default()
	
	api.RegisterRAGRoutes(router, ragService)
	api.RegisterTestGenRoutes(router, testGenerator)
	api.RegisterHealthRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting AI Developer Assistant on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
