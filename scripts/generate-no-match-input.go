//go:build ignore

// generate-no-match-input.go generates input data WITHOUT any pattern matches.
// This is the worst case for stdlib (scans entire file) and best case for
// coregex prefilters (skip quickly). Used to demonstrate 1000-3000x speedups.

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

const (
	targetSize = 6 * 1024 * 1024 // 6 MB
)

var (
	// Words that will NOT match any extreme patterns
	// Avoiding: digits, @, dots in filenames, "error/warning/fatal"
	words = []string{
		"the", "be", "to", "of", "and", "in", "that", "have",
		"it", "for", "not", "on", "with", "he", "as", "you", "do", "at",
		"this", "but", "his", "by", "from", "they", "we", "say", "her", "she",
		"function", "return", "if", "else", "while", "var", "const", "let",
		"import", "export", "class", "interface", "type", "struct", "package",
		"server", "client", "request", "response", "data", "file", "config",
		"hello", "world", "test", "main", "process", "handle", "create",
		"update", "delete", "find", "search", "query", "result", "output",
		"input", "stream", "buffer", "cache", "memory", "storage", "disk",
		"network", "socket", "connection", "session", "token", "auth",
		"user", "admin", "guest", "role", "permission", "access", "denied",
		"success", "failed", "pending", "complete", "active", "inactive",
	}
)

func randomWords(n int) string {
	parts := make([]string, n)
	for i := 0; i < n; i++ {
		parts[i] = words[rand.Intn(len(words))]
	}
	return strings.Join(parts, " ")
}

func generateContent() string {
	var builder strings.Builder
	builder.Grow(targetSize + 1024)

	for builder.Len() < targetSize {
		// Random line length 8-20 words
		wordCount := 8 + rand.Intn(12)
		line := randomWords(wordCount)
		builder.WriteString(line)
		builder.WriteByte('\n')
	}

	return builder.String()
}

func main() {
	rand.Seed(42) // Fixed seed for reproducibility

	data := generateContent()

	if err := os.MkdirAll("input", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating dir: %v\n", err)
		os.Exit(1)
	}

	err := os.WriteFile("input/no-match-data.txt", []byte(data), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated input/no-match-data.txt (%.2f MB)\n", float64(len(data))/1024/1024)
	fmt.Println("This file contains NO matches for IP, suffix, or inner patterns.")
	fmt.Println("Used to demonstrate 1000-3000x speedups (worst case for stdlib).")
}
