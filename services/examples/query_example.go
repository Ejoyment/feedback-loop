package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// Example 1: Query internal documentation
	queryReq := map[string]interface{}{
		"query":      "How do I implement authentication in our microservices?",
		"collection": "internal_docs",
		"top_k":      5,
	}

	resp, err := makeRequest("POST", "http://localhost:8080/api/v1/rag/query", queryReq)
	if err != nil {
		panic(err)
	}
	fmt.Println("RAG Query Response:", resp)

	// Example 2: Generate unit tests
	testGenReq := map[string]interface{}{
		"source_code": `
export function calculateDiscount(price: number, discountPercent: number): number {
  if (price < 0 || discountPercent < 0 || discountPercent > 100) {
    throw new Error("Invalid input");
  }
  return price * (1 - discountPercent / 100);
}`,
		"file_path":      "src/utils/pricing.ts",
		"test_framework": "jest",
	}

	testResp, err := makeRequest("POST", "http://localhost:8080/api/v1/testgen/generate", testGenReq)
	if err != nil {
		panic(err)
	}
	fmt.Println("Test Generation Response:", testResp)
}

func makeRequest(method, url string, body interface{}) (string, error) {
	jsonData, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	return string(respBody), nil
}
