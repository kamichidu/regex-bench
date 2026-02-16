//go:build ignore

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

// Generate DNA sequence matching regexdna benchmark from
// The Computer Language Benchmarks Game.
// Uses IUB frequencies for homo sapiens.

const targetSize = 6 * 1024 * 1024 // 6 MB (same as standard benchmark)

func main() {
	rand.Seed(42) // Fixed seed for reproducibility

	// IUB frequencies for homo sapiens (same as benchmarks game)
	type Gene struct {
		ch   byte
		prob float64
	}
	genes := []Gene{
		{'a', 0.302954942668},
		{'c', 0.5009432431601},
		{'g', 0.6984905497992},
		{'t', 1.0},
	}

	var sb strings.Builder
	sb.Grow(targetSize + 1024)

	col := 0
	for sb.Len() < targetSize {
		r := rand.Float64()
		for _, g := range genes {
			if r < g.prob {
				sb.WriteByte(g.ch)
				break
			}
		}
		col++
		if col == 60 {
			sb.WriteByte('\n')
			col = 0
		}
	}
	if col > 0 {
		sb.WriteByte('\n')
	}

	if err := os.MkdirAll("input", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating dir: %v\n", err)
		os.Exit(1)
	}

	err := os.WriteFile("input/dna-data.txt", []byte(sb.String()), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated input/dna-data.txt (%.2f MB)\n", float64(sb.Len())/1024/1024)
}
