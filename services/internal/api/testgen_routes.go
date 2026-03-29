package api

import (
	"net/http"

	"github.com/ai-dev-assistant/services/internal/testgen"
	"github.com/gin-gonic/gin"
)

func RegisterTestGenRoutes(router *gin.Engine, generator *testgen.Generator) {
	testGroup := router.Group("/api/v1/testgen")
	{
		testGroup.POST("/generate", func(c *gin.Context) {
			var req testgen.GenerateRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if req.TestFramework == "" {
				req.TestFramework = "jest"
			}

			resp, err := generator.Generate(c.Request.Context(), req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, resp)
		})
	}
}
