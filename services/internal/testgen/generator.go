package testgen

import (
	"context"
	"fmt"
	"strings"

	"github.com/ai-dev-assistant/services/internal/llm"
	"github.com/ai-dev-assistant/services/internal/rag"
)

type GenerateRequest struct {
	SourceCode   string `json:"source_code"`
	FilePath     string `json:"file_path"`
	TestFramework string `json:"test_framework"` // "jest", "vitest", "mocha"
}

type GenerateResponse struct {
	TestCode   string   `json:"test_code"`
	Coverage   []string `json:"coverage"`
	TokensUsed int      `json:"tokens_used"`
}

type Generator struct {
	llmGateway *llm.Gateway
	ragService *rag.Service
}

func NewGenerator(llmGateway *llm.Gateway, ragService *rag.Service) *Generator {
	return &Generator{
		llmGateway: llmGateway,
		ragService: ragService,
	}
}

func (g *Generator) Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
	// Step 1: Query RAG for testing best practices and patterns
	ragQuery := fmt.Sprintf("TypeScript unit testing patterns for %s framework", req.TestFramework)
	ragResp, err := g.ragService.Query(ctx, rag.QueryRequest{
		Query:      ragQuery,
		Collection: "testing_docs",
		TopK:       3,
	})

	var testingContext string
	if err == nil && ragResp != nil {
		testingContext = ragResp.Answer
	}

	// Step 2: Build system prompt with testing guidelines
	systemPrompt := g.buildSystemPrompt(req.TestFramework, testingContext)

	// Step 3: Generate tests using LLM
	userPrompt := fmt.Sprintf(`Generate comprehensive unit tests for the following TypeScript code.

File: %s

Source Code:
%s

Requirements:
- Use %s testing framework
- Cover all functions and edge cases
- Include mocks for external dependencies
- Follow testing best practices
- Add descriptive test names`, req.FilePath, req.SourceCode, req.TestFramework)

	completion, err := g.llmGateway.Complete(ctx, llm.CompletionRequest{
		Provider:     llm.ProviderOpenAI,
		Model:        "gpt-4-turbo-preview",
		SystemPrompt: systemPrompt,
		Messages: []llm.Message{
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.2,
		MaxTokens:   3000,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to generate tests: %w", err)
	}

	// Step 4: Extract test code and coverage info
	testCode := g.extractCodeBlock(completion.Content)
	coverage := g.analyzeCoverage(req.SourceCode, testCode)

	return &GenerateResponse{
		TestCode:   testCode,
		Coverage:   coverage,
		TokensUsed: completion.TokensUsed,
	}, nil
}

func (g *Generator) buildSystemPrompt(framework, context string) string {
	base := `You are an expert TypeScript test engineer. Generate high-quality unit tests that:
- Cover all code paths and edge cases
- Use proper mocking and stubbing
- Follow testing best practices
- Are maintainable and readable`

	if context != "" {
		base += fmt.Sprintf("\n\nAdditional Context:\n%s", context)
	}

	return base
}

func (g *Generator) extractCodeBlock(content string) string {
	// Extract code from markdown code blocks
	if strings.Contains(content, "```") {
		parts := strings.Split(content, "```")
		if len(parts) >= 3 {
			code := parts[1]
			// Remove language identifier
			if idx := strings.Index(code, "\n"); idx != -1 {
				code = code[idx+1:]
			}
			return strings.TrimSpace(code)
		}
	}
	return content
}

func (g *Generator) analyzeCoverage(sourceCode, testCode string) []string {
	coverage := []string{}
	
	// Simple heuristic analysis
	if strings.Contains(testCode, "describe(") {
		coverage = append(coverage, "Test suites defined")
	}
	if strings.Contains(testCode, "it(") || strings.Contains(testCode, "test(") {
		coverage = append(coverage, "Test cases implemented")
	}
	if strings.Contains(testCode, "mock") || strings.Contains(testCode, "jest.fn()") {
		coverage = append(coverage, "Mocks configured")
	}
	if strings.Contains(testCode, "expect(") {
		coverage = append(coverage, "Assertions present")
	}

	return coverage
}
