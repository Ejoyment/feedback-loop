package testgen

import (
	"testing"
)

func TestParseTypeScript(t *testing.T) {
	code := `
import { User } from './types';

export async function fetchUser(id: string): Promise<User> {
  const response = await fetch('/api/users/' + id);
  return response.json();
}

export function calculateTotal(items: number[]): number {
  return items.reduce((sum, item) => sum + item, 0);
}
`

	analysis := ParseTypeScript(code)

	if len(analysis.Functions) != 2 {
		t.Errorf("Expected 2 functions, got %d", len(analysis.Functions))
	}

	if !analysis.HasAsync {
		t.Error("Expected HasAsync to be true")
	}

	if !analysis.HasPromises {
		t.Error("Expected HasPromises to be true")
	}

	if len(analysis.Imports) == 0 {
		t.Error("Expected imports to be parsed")
	}
}

func TestParseParams(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"", []string{}},
		{"a: string", []string{"a"}},
		{"a: string, b: number", []string{"a", "b"}},
		{"x, y, z", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		result := parseParams(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("parseParams(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}
