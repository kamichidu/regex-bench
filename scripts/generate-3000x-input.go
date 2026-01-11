//go:build ignore

// generate-3000x-input.go generates input data with NO DIGITS.
// This is designed to achieve 3000x+ speedup on IP/phone patterns because:
// - DigitPrefilter scans for first digit character
// - No digits → entire file skipped in microseconds
// - stdlib must scan byte-by-byte → hundreds of milliseconds

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
	// Words with NO DIGITS, NO "error/warning/fatal", NO email-like patterns
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
		"apple", "banana", "cherry", "orange", "grape", "lemon", "mango",
		"happy", "sad", "fast", "slow", "big", "small", "good", "bad",
		"run", "walk", "jump", "swim", "fly", "climb", "dance", "sing",
		"red", "blue", "green", "yellow", "black", "white", "purple", "pink",
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
		wordCount := 8 + rand.Intn(12)
		line := randomWords(wordCount)
		builder.WriteString(line)
		builder.WriteByte('\n')
	}

	return builder.String()
}

func main() {
	rand.Seed(42)

	data := generateContent()

	// Verify no digits
	for _, c := range data {
		if c >= '0' && c <= '9' {
			fmt.Fprintln(os.Stderr, "ERROR: Found digit in generated data!")
			os.Exit(1)
		}
	}

	if err := os.MkdirAll("input", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating dir: %v\n", err)
		os.Exit(1)
	}

	err := os.WriteFile("input/no-digits-data.txt", []byte(data), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated input/no-digits-data.txt (%.2f MB)\n", float64(len(data))/1024/1024)
	fmt.Println("This file contains NO DIGITS - designed for 3000x+ IP pattern speedup.")
	fmt.Println("DigitPrefilter will skip entire file in microseconds.")
}
