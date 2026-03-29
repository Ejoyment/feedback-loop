package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	fmt.Println("AI Developer Assistant CLI")
	fmt.Println("Commands: query <text>, test <file>, exit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		parts := strings.SplitN(input, " ", 2)
		
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		
		switch command {
		case "query":
			if len(parts) < 2 {
				fmt.Println("Usage: query <your question>")
				continue
			}
			handleQuery(parts[1])
			
		case "test":
			if len(parts) < 2 {
				fmt.Println("Usage: test <file_path>")
				continue
			}
			handleTestGen(parts[1])
			
		case "exit":
			fmt.Println("Goodbye!")
			return
			
		default:
			fmt.Println("Unknown command. Use: query, test, exit")
		}
	}
}

func handleQuery(query string) {
	req := map[string]interface{}{
		"query":      query,
		"collection": "internal_docs",
		"top_k":      5,
	}

	resp, err := makeRequest("POST", "http://localhost:8080/api/v1/rag/query", req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(resp), &result)
	
	fmt.Println("\nAnswer:", result["answer"])
	fmt.Println("\nSources:")
	if sources, ok := result["sources"].([]interface{}); ok {
		for i, src := range sources {
			if srcMap, ok := src.(map[string]interface{}); ok {
				fmt.Printf("  %d. %s\n", i+1, srcMap["id"])
			}
		}
	}
	fmt.Println()
}

func handleTestGen(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	req := map[string]interface{}{
		"source_code":    string(content),
		"file_path":      filePath,
		"test_framework": "jest",
	}

	resp, err := makeRequest("POST", "http://localhost:8080/api/v1/testgen/generate", req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(resp), &result)
	
	fmt.Println("\nGenerated Test Code:")
	fmt.Println(result["test_code"])
	fmt.Println()
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
