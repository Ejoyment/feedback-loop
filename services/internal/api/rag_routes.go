package api

import (
	"net/http"

	"github.com/ai-dev-assistant/services/internal/rag"
	"github.com/gin-gonic/gin"
)

func RegisterRAGRoutes(router *gin.Engine, ragService *rag.Service) {
	ragGroup := router.Group("/api/v1/rag")
	{
		ragGroup.POST("/query", func(c *gin.Context) {
			var req rag.QueryRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			resp, err := ragService.Query(c.Request.Context(), req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, resp)
		})

		ragGroup.POST("/index", func(c *gin.Context) {
			var doc rag.Document
			if err := c.ShouldBindJSON(&doc); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if err := ragService.IndexDocument(c.Request.Context(), doc); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "indexed", "id": doc.ID})
		})
	}
}
