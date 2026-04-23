package main

import (
	"fmt"
	"os"
	regexp "github.com/wasilibs/go-re2"
	"time"
)

// regexdna patterns from The Computer Language Benchmarks Game
// DNA reverse complement matching
type Pattern struct {
	Name    string
	Pattern string
}

var patterns = []Pattern{
	{"dna_1", `agggtaaa|tttaccct`},
	{"dna_2", `[cgt]gggtaaa|tttaccc[acg]`},
	{"dna_3", `a[act]ggtaaa|tttacc[agt]t`},
	{"dna_4", `ag[act]gtaaa|tttac[agt]ct`},
	{"dna_5", `agg[act]taaa|ttta[agt]cct`},
	{"dna_6", `aggg[acg]aaa|ttt[cgt]ccct`},
	{"dna_7", `agggt[cgt]aa|tt[acg]accct`},
	{"dna_8", `agggta[cgt]a|t[acg]taccct`},
	{"dna_9", `agggtaa[cgt]|[acg]ttaccct`},
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
		fmt.Println("Usage: go-stdlib-dna <input-file>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Go RE2 CGO regexdna (input: %.2f MB)\n", float64(len(data))/1024/1024)
	fmt.Printf("%-15s %10s %10s %6s\n", "pattern", "compile", "search", "matches")
	fmt.Println("─────────────────────────────────────────────────")

	for _, p := range patterns {
		measure(data, p)
	}
}
