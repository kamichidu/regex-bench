package main

import (
	"fmt"
	"github.com/Jemmic/go-pcre2"
	"os"
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
	// PCRE2 doesn't use (?m) for multiline by default in some cases,
	// but the library should handle it if passed in the pattern.
	re, err := pcre2.Compile(p.Pattern, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error compiling pattern %s: %v\n", p.Name, err)
		return
	}
	compileElapsed := time.Since(compileStart)

	searchStart := time.Now()
	count := 0
	matcher := re.Matcher(data, 0)
	for matcher.Match(data, 0) {
		count++
		// Move to the next position.
		// Note: This is a simple implementation of FindAll for benchmarking.
		// pcre2 doesn't have a direct equivalent of FindAll.
		indices := matcher.Index()
		if len(indices) < 2 {
			break
		}
		end := indices[1]
		if end >= len(data) {
			break
		}
		data = data[end:]
		matcher = re.Matcher(data, 0)
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

	fmt.Printf("Go PCRE2 (input: %.2f MB)\n", float64(len(data))/1024/1024)
	fmt.Printf("%-15s %10s %10s %6s\n", "pattern", "compile", "search", "matches")
	fmt.Println("─────────────────────────────────────────────────")

	for _, p := range patterns {
		measure(data, p)
	}
}
