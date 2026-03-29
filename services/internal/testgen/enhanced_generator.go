package testgen

import (
	"context"
	"fmt"
	"strings"

	"github.com/ai-dev-assistant/services/internal/llm"
	"github.com/ai-dev-assistant/services/internal/rag"
)

type EnhancedGenerator struct {
	*Generator
}

func NewEnhancedGenerator(llmGateway *llm.Gateway, ragService *rag.Service) *EnhancedGenerator {
	return &EnhancedGenerator{
		Generator: NewGenerator(llmGateway, ragService),
	}
}

func (eg *EnhancedGenerator) GenerateWithContext(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
	// Parse source code to understand structure
	analysis := ParseTypeScript(req.SourceCode)

	// Query RAG for similar test patterns
	var testPatterns string
	if len(analysis.Functions) > 0 {
		ragResp, err := eg.ragService.Query(ctx, rag.QueryRequest{
			Query:      fmt.Sprintf("test patterns for %s", analysis.Functions[0].Name),
			Collection: "test_examples",
			TopK:       3,
		})
		if err == nil {
			testPatterns = ragResp.Answer
		}
	}

	// Build enhanced prompt with code analysis
	enhancedPrompt := eg.buildEnhancedPrompt(req, analysis, testPatterns)

	completion, err := eg.llmGateway.Complete(ctx, llm.CompletionRequest{
		Provider:     llm.ProviderClaude,
		Model:        "claude-3-sonnet-20240229",
		SystemPrompt: "You are an expert TypeScript test engineer specializing in legacy code.",
		Messages: []llm.Message{
			{Role: "user", Content: enhancedPrompt},
		},
		Temperature: 0.2,
		MaxTokens:   3500,
	})

	if err != nil {
		return nil, err
	}

	testCode := eg.extractCodeBlock(completion.Content)
	coverage := eg.generateDetailedCoverage(analysis, testCode)

	return &GenerateResponse{
		TestCode:   testCode,
		Coverage:   coverage,
		TokensUsed: completion.TokensUsed,
	}, nil
}

func (eg *EnhancedGenerator) buildEnhancedPrompt(req GenerateRequest, analysis *CodeAnalysis, patterns string) string {
	var builder strings.Builder
	
	builder.WriteString(fmt.Sprintf("Generate comprehensive %s tests for this TypeScript module.\n\n", req.TestFramework))
	builder.WriteString(fmt.Sprintf("File: %s\n\n", req.FilePath))
	builder.WriteString("Code Analysis:\n")
	builder.WriteString(fmt.Sprintf("- Functions: %d\n", len(analysis.Functions)))
	builder.WriteString(fmt.Sprintf("- Has async code: %v\n", analysis.HasAsync))
	builder.WriteString(fmt.Sprintf("- Has promises: %v\n\n", analysis.HasPromises))

	if patterns != "" {
		builder.WriteString("Relevant Test Patterns:\n")
		builder.WriteString(patterns)
		builder.WriteString("\n\n")
	}

	builder.WriteString("Source Code:\n")
	builder.WriteString(req.SourceCode)

	return builder.String()
}

func (eg *EnhancedGenerator) generateDetailedCoverage(analysis *CodeAnalysis, testCode string) []string {
	coverage := []string{}
	
	for _, fn := range analysis.Functions {
		if strings.Contains(testCode, fn.Name) {
			coverage = append(coverage, fmt.Sprintf("✓ %s() covered", fn.Name))
		} else {
			coverage = append(coverage, fmt.Sprintf("✗ %s() not covered", fn.Name))
		}
	}

	return coverage
}
