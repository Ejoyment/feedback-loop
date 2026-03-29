package langgraph

import (
	"context"
	"fmt"
)

type NodeFunc func(ctx context.Context, state *State) error

type State struct {
	Data map[string]interface{}
}

type Graph struct {
	nodes map[string]NodeFunc
	edges map[string][]string
	entry string
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[string]NodeFunc),
		edges: make(map[string][]string),
	}
}

func (g *Graph) AddNode(name string, fn NodeFunc) {
	g.nodes[name] = fn
}

func (g *Graph) AddEdge(from, to string) {
	g.edges[from] = append(g.edges[from], to)
}

func (g *Graph) SetEntry(name string) {
	g.entry = name
}

func (g *Graph) Execute(ctx context.Context, initialState *State) error {
	if g.entry == "" {
		return fmt.Errorf("no entry node set")
	}

	return g.executeNode(ctx, g.entry, initialState)
}

func (g *Graph) executeNode(ctx context.Context, nodeName string, state *State) error {
	fn, exists := g.nodes[nodeName]
	if !exists {
		return fmt.Errorf("node not found: %s", nodeName)
	}

	if err := fn(ctx, state); err != nil {
		return fmt.Errorf("node %s failed: %w", nodeName, err)
	}

	// Execute next nodes
	for _, nextNode := range g.edges[nodeName] {
		if err := g.executeNode(ctx, nextNode, state); err != nil {
			return err
		}
	}

	return nil
}
