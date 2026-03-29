package testgen

import (
	"regexp"
	"strings"
)

type CodeAnalysis struct {
	Functions   []FunctionInfo
	Classes     []ClassInfo
	Imports     []string
	HasAsync    bool
	HasPromises bool
}

type FunctionInfo struct {
	Name       string
	Params     []string
	IsAsync    bool
	IsExported bool
}

type ClassInfo struct {
	Name    string
	Methods []string
}

func ParseTypeScript(code string) *CodeAnalysis {
	analysis := &CodeAnalysis{
		Functions: []FunctionInfo{},
		Classes:   []ClassInfo{},
		Imports:   []string{},
	}

	// Extract imports
	importRegex := regexp.MustCompile(`import\s+.*\s+from\s+['"](.+)['"]`)
	for _, match := range importRegex.FindAllStringSubmatch(code, -1) {
		if len(match) > 1 {
			analysis.Imports = append(analysis.Imports, match[1])
		}
	}

	// Extract functions
	funcRegex := regexp.MustCompile(`(?:export\s+)?(?:async\s+)?function\s+(\w+)\s*\(([^)]*)\)`)
	for _, match := range funcRegex.FindAllStringSubmatch(code, -1) {
		if len(match) > 2 {
			isAsync := strings.Contains(match[0], "async")
			isExported := strings.Contains(match[0], "export")
			
			analysis.Functions = append(analysis.Functions, FunctionInfo{
				Name:       match[1],
				Params:     parseParams(match[2]),
				IsAsync:    isAsync,
				IsExported: isExported,
			})

			if isAsync {
				analysis.HasAsync = true
			}
		}
	}

	// Detect promises
	if strings.Contains(code, "Promise") || strings.Contains(code, ".then(") {
		analysis.HasPromises = true
	}

	return analysis
}

func parseParams(paramStr string) []string {
	if paramStr == "" {
		return []string{}
	}
	
	params := strings.Split(paramStr, ",")
	result := make([]string, 0, len(params))
	
	for _, p := range params {
		p = strings.TrimSpace(p)
		if p != "" {
			// Extract parameter name (before colon if typed)
			if idx := strings.Index(p, ":"); idx != -1 {
				p = p[:idx]
			}
			result = append(result, strings.TrimSpace(p))
		}
	}
	
	return result
}
