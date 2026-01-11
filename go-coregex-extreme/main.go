package main

import (
	"fmt"
	"os"
	"time"

	"github.com/coregx/coregex"
)

// Extreme patterns - designed to show 1000-3000x speedup on no-match data
type Pattern struct {
	Name    string
	Pattern string
	Note    string
}

var patterns = []Pattern{
	{
		Name:    "ip_nomatch",
		Pattern: `(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])`,
		Note:    "DigitPrefilter: 3000x+ on no-match",
	},
	{
		Name:    "inner_nomatch",
		Pattern: `.*error.*`,
		Note:    "ReverseInner: 900x on no-match",
	},
	{
		Name:    "suffix_find",
		Pattern: `[^\s]+\.txt`,
		Note:    "ReverseSuffix: find .txt files",
	},
	{
		Name:    "phone_nomatch",
		Pattern: `\d{3}-\d{3}-\d{4}`,
		Note:    "DigitPrefilter: phone numbers",
	},
}

const iterations = 1000

func measure(data []byte, p Pattern) {
	re := coregex.MustCompile(p.Pattern)

	// Warmup
	var matched bool
	for i := 0; i < 10; i++ {
		matched = re.Match(data)
	}

	// Time all iterations together, then divide
	start := time.Now()
	for i := 0; i < iterations; i++ {
		re.Match(data)
	}
	totalNs := time.Since(start).Nanoseconds()
	avgNs := totalNs / iterations

	matchStr := "no"
	if matched {
		matchStr = "yes"
	}

	if avgNs >= 1000000 {
		ms := float64(avgNs) / 1000000.0
		fmt.Printf("%-15s %10.2f ms  match: %s\n", p.Name, ms, matchStr)
	} else if avgNs >= 1000 {
		us := float64(avgNs) / 1000.0
		fmt.Printf("%-15s %10.2f µs  match: %s\n", p.Name, us, matchStr)
	} else {
		fmt.Printf("%-15s %10d ns  match: %s\n", p.Name, avgNs, matchStr)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go-coregex-extreme <input-file>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Go coregex EXTREME (no-match worst case) - input: %.2f MB\n", float64(len(data))/1024/1024)
	fmt.Println("─────────────────────────────────────────────────────────────────────")

	for _, p := range patterns {
		measure(data, p)
	}
}
