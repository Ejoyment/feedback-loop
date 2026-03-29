package rag

import (
	"context"
	"fmt"
	"strings"

	"github.com/ai-dev-assistant/services/internal/llm"
)

type Document struct {
	ID       string                 `json:"id"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata"`
	Type     string                 `json:"type"` // "documentation", "code_schema", "api_spec"
}

type QueryRequest struct {
	Query      string `json:"query"`
	Collection string `json:"collection"`
	TopK       int    `json:"top_k"`
}

type QueryResponse struct {
	Answer   string     `json:"answer"`
	Sources  []Document `json:"sources"`
	Metadata struct {
		TokensUsed int `json:"tokens_used"`
	} `json:"metadata"`
}

type Service struct {
	vectorClient *VectorClient
	llmGateway   *llm.Gateway
}

func NewService(vectorClient *VectorClient, llmGateway *llm.Gateway) *Service {
	return &Service{
		vectorClient: vectorClient,
		llmGateway:   llmGateway,
	}
}

func (s *Service) Query(ctx context.Context, req QueryRequest) (*QueryResponse, error) {
	// Step 1: Embed the query
	queryVector, err := s.vectorClient.Embed(ctx, req.Query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	// Step 2: Search for relevant documents
	topK := req.TopK
	if topK == 0 {
		topK = 5
	}

	searchResults, err := s.vectorClient.Search(ctx, req.Collection, queryVector, topK)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	// Step 3: Build context from retrieved documents
	var contextBuilder strings.Builder
	sources := make([]Document, 0, len(searchResults))

	for _, result := range searchResults {
		content, _ := result.Payload["content"].(string)
		docType, _ := result.Payload["type"].(string)
		
		contextBuilder.WriteString(fmt.Sprintf("## Document (Score: %.3f)\n%s\n\n", result.Score, content))
		
		sources = append(sources, Document{
			ID:       result.ID,
			Content:  content,
			Type:     docType,
			Metadata: result.Payload,
		})
	}

	// Step 4: Generate answer using LLM with retrieved context
	systemPrompt := `You are an AI assistant helping engineers query internal documentation and codebase schemas.
Use the provided context to answer the question accurately. If the context doesn't contain enough information, say so.`

	userPrompt := fmt.Sprintf("Context:\n%s\n\nQuestion: %s", contextBuilder.String(), req.Query)

	completion, err := s.llmGateway.Complete(ctx, llm.CompletionRequest{
		Provider:     llm.ProviderOpenAI,
		Model:        "gpt-4-turbo-preview",
		SystemPrompt: systemPrompt,
		Messages: []llm.Message{
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.3,
		MaxTokens:   1500,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to generate answer: %w", err)
	}

	return &QueryResponse{
		Answer:  completion.Content,
		Sources: sources,
		Metadata: struct {
			TokensUsed int `json:"tokens_used"`
		}{TokensUsed: completion.TokensUsed},
	}, nil
}

func (s *Service) IndexDocument(ctx context.Context, doc Document) error {
	// Embed document content
	vector, err := s.vectorClient.Embed(ctx, doc.Content)
	if err != nil {
		return fmt.Errorf("failed to embed document: %w", err)
	}

	// Store in vector database with metadata
	// Implementation would use Qdrant client directly
	_ = vector
	return nil
}
