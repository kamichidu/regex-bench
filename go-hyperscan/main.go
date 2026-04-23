package main

import (
	"fmt"
	"os"
	"github.com/flier/gohs/hyperscan"
	"time"
)

type Pattern struct {
	Name    string
	Pattern string
}

var patterns = []Pattern{
	{"literal_alt", `error|warning|fatal|critical`},
	{"multi_literal", `apple|banana|cherry|date|elderberry|fig|grape|honeydew|kiwi|lemon|mango|orange`},
	{"anchored", `^HTTP/[12]\.[01]`},
	{"inner_literal", `.*@example\.com`},
	{"suffix", `.*\.(txt|log|md)`},
	{"char_class", `[\w]+`},
	{"email", `[\w.+-]+@[\w.-]+\.[\w.-]+`},
	{"uri", `[\w]+://[^/\s?#]+[^\s?#]+(?:\?[^\s#]*)?(?:#[^\s]*)?`},
	{"version", `\d+\.\d+\.\d+`},
	{"ip", `(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])`},
	{"alpha_digit", `[a-zA-Z]+\d+`},
	{"word_digit", `\w+[0-9]+`},
	{"http_methods", `(?m)^(GET|POST|PUT|DELETE|PATCH)`},
	{"anchored_php", `^/.*[\w-]+\.php`},
	{"multiline_php", `(?m)^/.*\.php`},
	{"word_repeat", `(\w{2,8})+`},
}

func measure(data []byte, p Pattern) {
	compileStart := time.Now()
	// Hyperscan may not support all Go regexp syntax (e.g. (?m)) directly 
	// without specific flags, but we'll try basic compilation.
	db, err := hyperscan.Compile(p.Pattern)
	if err != nil {
		// Try with some flags if direct compile fails
		// For this benchmark, we'll report errors to stderr
		fmt.Fprintf(os.Stderr, "Error compiling pattern %s: %v\n", p.Name, err)
		return
	}
	defer db.Close()
	scratch, err := hyperscan.NewScratch(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating scratch for %s: %v\n", p.Name, err)
		return
	}
	defer scratch.Free()
	compileElapsed := time.Since(compileStart)

	searchStart := time.Now()
	count := 0
	err = db.(hyperscan.BlockDatabase).Scan(data, scratch, func(id uint, from, to uint64, flags uint, context interface{}) error {
		count++
		return nil
	}, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning with %s: %v\n", p.Name, err)
	}
	searchElapsed := time.Since(searchStart)

	compileMs := float64(compileElapsed) / float64(time.Millisecond)
	searchMs := float64(searchElapsed) / float64(time.Millisecond)

	fmt.Printf("%-15s %10.2f %10.2f %6d\n", p.Name, compileMs, searchMs, count)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: benchmark <input-file>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Go Hyperscan (input: %.2f MB)\n", float64(len(data))/1024/1024)
	fmt.Printf("%-15s %10s %10s %6s\n", "pattern", "compile", "search", "matches")
	fmt.Println("─────────────────────────────────────────────────")

	for _, p := range patterns {
		measure(data, p)
	}
}
