package llm

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type Provider string

const (
	ProviderOpenAI Provider = "openai"
	ProviderClaude Provider = "claude"
)

type Message struct {
	Role    string
	Content string
}

type CompletionRequest struct {
	Provider     Provider
	Model        string
	Messages     []Message
	Temperature  float32
	MaxTokens    int
	SystemPrompt string
}

type CompletionResponse struct {
	Content      string
	TokensUsed   int
	FinishReason string
}

type Gateway struct {
	openaiClient *openai.Client
	claudeKey    string
}

func NewGateway(openaiKey, claudeKey string) *Gateway {
	var openaiClient *openai.Client
	if openaiKey != "" {
		openaiClient = openai.NewClient(openaiKey)
	}

	return &Gateway{
		openaiClient: openaiClient,
		claudeKey:    claudeKey,
	}
}

func (g *Gateway) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	switch req.Provider {
	case ProviderOpenAI:
		return g.completeOpenAI(ctx, req)
	case ProviderClaude:
		return g.completeClaude(ctx, req)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", req.Provider)
	}
}

func (g *Gateway) completeOpenAI(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	if g.openaiClient == nil {
		return nil, fmt.Errorf("OpenAI client not initialized")
	}

	messages := make([]openai.ChatCompletionMessage, 0, len(req.Messages)+1)
	
	if req.SystemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: req.SystemPrompt,
		})
	}

	for _, msg := range req.Messages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	resp, err := g.openaiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       req.Model,
		Messages:    messages,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
	})

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	return &CompletionResponse{
		Content:      resp.Choices[0].Message.Content,
		TokensUsed:   resp.Usage.TotalTokens,
		FinishReason: string(resp.Choices[0].FinishReason),
	}, nil
}

func (g *Gateway) completeClaude(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	// Claude API implementation using anthropic-sdk-go or HTTP client
	// Simplified for brevity - full implementation would use official SDK
	return nil, fmt.Errorf("Claude implementation requires anthropic-sdk-go")
}
