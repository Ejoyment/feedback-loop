package langgraph

import (
	"context"
	"fmt"

	"github.com/ai-dev-assistant/services/internal/llm"
	"github.com/ai-dev-assistant/services/internal/rag"
)

type AgentWorkflow struct {
	graph      *Graph
	llmGateway *llm.Gateway
	ragService *rag.Service
}

func NewAgentWorkflow(llmGateway *llm.Gateway, ragService *rag.Service) *AgentWorkflow {
	aw := &AgentWorkflow{
		graph:      NewGraph(),
		llmGateway: llmGateway,
		ragService: ragService,
	}

	aw.setupGraph()
	return aw
}

func (aw *AgentWorkflow) setupGraph() {
	aw.graph.AddNode("analyze", aw.analyzeNode)
	aw.graph.AddNode("retrieve", aw.retrieveNode)
	aw.graph.AddNode("generate", aw.generateNode)
	aw.graph.AddNode("validate", aw.validateNode)

	aw.graph.AddEdge("analyze", "retrieve")
	aw.graph.AddEdge("retrieve", "generate")
	aw.graph.AddEdge("generate", "validate")

	aw.graph.SetEntry("analyze")
}

func (aw *AgentWorkflow) analyzeNode(ctx context.Context, state *State) error {
	query, _ := state.Data["query"].(string)
	state.Data["analyzed_query"] = query
	state.Data["intent"] = "documentation_search"
	return nil
}

func (aw *AgentWorkflow) retrieveNode(ctx context.Context, state *State) error {
	query, _ := state.Data["analyzed_query"].(string)
	
	resp, err := aw.ragService.Query(ctx, rag.QueryRequest{
		Query:      query,
		Collection: "internal_docs",
		TopK:       5,
	})
	
	if err != nil {
		return err
	}

	state.Data["retrieved_docs"] = resp.Sources
	return nil
}

func (aw *AgentWorkflow) generateNode(ctx context.Context, state *State) error {
	docs, _ := state.Data["retrieved_docs"].([]rag.Document)
	query, _ := state.Data["query"].(string)

	context := ""
	for _, doc := range docs {
		context += doc.Content + "\n\n"
	}

	completion, err := aw.llmGateway.Complete(ctx, llm.CompletionRequest{
		Provider:     llm.ProviderOpenAI,
		Model:        "gpt-4-turbo-preview",
		SystemPrompt: "You are a helpful assistant for engineers.",
		Messages: []llm.Message{
			{Role: "user", Content: fmt.Sprintf("Context:\n%s\n\nQuestion: %s", context, query)},
		},
		Temperature: 0.3,
		MaxTokens:   1000,
	})

	if err != nil {
		return err
	}

	state.Data["answer"] = completion.Content
	return nil
}

func (aw *AgentWorkflow) validateNode(ctx context.Context, state *State) error {
	answer, _ := state.Data["answer"].(string)
	if answer == "" {
		return fmt.Errorf("no answer generated")
	}
	return nil
}

func (aw *AgentWorkflow) Run(ctx context.Context, query string) (string, error) {
	state := &State{
		Data: map[string]interface{}{
			"query": query,
		},
	}

	if err := aw.graph.Execute(ctx, state); err != nil {
		return "", err
	}

	answer, _ := state.Data["answer"].(string)
	return answer, nil
}
