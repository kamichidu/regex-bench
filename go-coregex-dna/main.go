package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/coregx/coregex"
	"github.com/coregx/coregex/meta"
)

func getCoregexVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	for _, dep := range info.Deps {
		if dep.Path == "github.com/coregx/coregex" {
			return dep.Version
		}
	}
	return "unknown"
}

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
	re := coregex.MustCompile(p.Pattern)
	compileElapsed := time.Since(compileStart)

	// Get strategy from meta engine
	engine, _ := meta.Compile(p.Pattern)
	strategy := engine.Strategy().String()

	searchStart := time.Now()
	matches := re.FindAll(data, -1)
	searchElapsed := time.Since(searchStart)

	compileMs := float64(compileElapsed) / float64(time.Millisecond)
	searchMs := float64(searchElapsed) / float64(time.Millisecond)

	fmt.Printf("%-15s %10.2f %10.2f %6d %s\n", p.Name, compileMs, searchMs, len(matches), strategy)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go-coregex-dna <input-file>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Go coregex %s regexdna (input: %.2f MB)\n", getCoregexVersion(), float64(len(data))/1024/1024)
	fmt.Printf("%-15s %10s %10s %6s %s\n", "pattern", "compile", "search", "matches", "strategy")
	fmt.Println("─────────────────────────────────────────────────")

	for _, p := range patterns {
		measure(data, p)
	}
}
