package main

import (
	"fmt"
	"os"
	regexp "github.com/wasilibs/go-re2"
	"time"
)

// LangArena LogParser patterns from https://kostya.github.io/LangArena/
// 13 real-world patterns for Apache log parsing.
type Pattern struct {
	Name    string
	Pattern string
}

var patterns = []Pattern{
	{"errors", ` [5][0-9]{2} | [4][0-9]{2} `},
	{"bots", `(?i)bot|crawler|scanner|spider|indexing|crawl|robot|spider`},
	{"suspicious", `(?i)etc/passwd|wp-admin|\.\./`},
	{"ips", `\d+\.\d+\.\d+\.35`},
	{"api_calls", `/api/[^ "]+`},
	{"post_requests", `POST [^ ]* HTTP`},
	{"auth_attempts", `(?i)/login|/signin`},
	{"methods", `(?i)get|post|put`},
	{"emails", `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`},
	{"passwords", `password=[^&\s"]+`},
	{"tokens", `token=[^&\s"]+|api[_-]?key=[^&\s"]+`},
	{"sessions", `session[_-]?id=[^&\s"]+`},
	{"peak_hours", `\[\d+/\w+/\d+:1[3-7]:\d+:\d+ [+\-]\d+\]`},
}

func measure(data []byte, p Pattern) {
	compileStart := time.Now()
	re := regexp.MustCompile(p.Pattern)
	compileElapsed := time.Since(compileStart)

	searchStart := time.Now()
	matches := re.FindAll(data, -1)
	searchElapsed := time.Since(searchStart)

	compileMs := float64(compileElapsed) / float64(time.Millisecond)
	searchMs := float64(searchElapsed) / float64(time.Millisecond)

	fmt.Printf("%-15s %10.2f %10.2f %6d\n", p.Name, compileMs, searchMs, len(matches))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go-stdlib-langarena <input-file>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Go RE2 CGO LangArena LogParser (input: %.2f MB)\n", float64(len(data))/1024/1024)
	fmt.Printf("%-15s %10s %10s %6s\n", "pattern", "compile", "search", "matches")
	fmt.Println("─────────────────────────────────────────────────")

	for _, p := range patterns {
		measure(data, p)
	}
}
