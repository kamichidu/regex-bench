package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Result struct {
	Pattern string
	Compile float64
	Search  float64
	Matches int
}

type EngineResult struct {
	Name      string
	InputSize float64 // MB
	Results   map[string]Result
}

func parseFile(path string) (*EngineResult, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	engineName := strings.TrimSuffix(filepath.Base(path), ".txt")
	res := &EngineResult{
		Name:    engineName,
		Results: make(map[string]Result),
	}

	scanner := bufio.NewScanner(file)
	// Example: Go stdlib (input: 6.09 MB)
	reInput := regexp.MustCompile(`input: ([0-9.]+) MB`)
	// Example: literal_alt         0.01       0.45   1441
	reLine := regexp.MustCompile(`^([a-z_0-9-]+)\s+([0-9.]+)\s+([0-9.]+)\s+([0-9]+)`)

	for scanner.Scan() {
		line := scanner.Text()
		if m := reInput.FindStringSubmatch(line); m != nil {
			res.InputSize, _ = strconv.ParseFloat(m[1], 64)
		}
		if m := reLine.FindStringSubmatch(line); m != nil {
			pattern := m[1]
			compile, _ := strconv.ParseFloat(m[2], 64)
			search, _ := strconv.ParseFloat(m[3], 64)
			matches, _ := strconv.Atoi(m[4])
			res.Results[pattern] = Result{pattern, compile, search, matches}
		}
	}
	return res, nil
}

func main() {
	scenarios := []string{"standard", "dna", "extreme", "langarena"}
	engines := []string{"stdlib", "regexp-re", "coregex", "re2-wasm", "re2-cgo", "pcre2", "hyperscan"}

	fmt.Println("# Benchmark Report")
	fmt.Println("\nGenerated automatically from the latest run.")

	for _, sc := range scenarios {
		fmt.Printf("\n## Scenario: %s\n", strings.Title(sc))

		var allEngineResults []*EngineResult
		patternsMap := make(map[string]bool)

		for _, en := range engines {
			path := filepath.Join("results", sc, en+".txt")
			if res, err := parseFile(path); err == nil {
				allEngineResults = append(allEngineResults, res)
				for p := range res.Results {
					patternsMap[p] = true
				}
			}
		}

		if len(allEngineResults) == 0 {
			continue
		}

		patterns := make([]string, 0, len(patternsMap))
		for p := range patternsMap {
			patterns = append(patterns, p)
		}
		sort.Strings(patterns)

		// Table Header
		header := "| Pattern | "
		sep := "|---|"
		for _, er := range allEngineResults {
			header += er.Name + " | "
			sep += "---|"
		}
		fmt.Println("\n### Search Time (ms)")
		fmt.Println(header)
		fmt.Println(sep)

		for _, p := range patterns {
			row := "| " + p + " | "
			for _, er := range allEngineResults {
				if r, ok := er.Results[p]; ok {
					row += fmt.Sprintf("%.2f | ", r.Search)
				} else {
					row += "- | "
				}
			}
			fmt.Println(row)
		}

		// Mermaid Chart
		fmt.Println("\n### Visualization (Search Time)")
		fmt.Println("```mermaid")
		fmt.Println("xychart-beta")
		fmt.Printf("    title \"%s Scenario Performance\"\n", strings.Title(sc))

		// X-axis: Engines
		engineNames := []string{}
		for _, er := range allEngineResults {
			engineNames = append(engineNames, `"`+er.Name+`"`)
		}
		fmt.Printf("    x-axis [%s]\n", strings.Join(engineNames, ", "))

		// Y-axis: Average Search Time across all patterns
		fmt.Println("    y-axis \"Avg Search Time (ms)\"")

		avgTimes := []string{}
		for _, er := range allEngineResults {
			sum := 0.0
			count := 0
			for _, r := range er.Results {
				sum += r.Search
				count++
			}
			avg := 0.0
			if count > 0 {
				avg = sum / float64(count)
			}
			avgTimes = append(avgTimes, fmt.Sprintf("%.2f", avg))
		}
		fmt.Printf("    bar [%s]\n", strings.Join(avgTimes, ", "))
		fmt.Println("```")
	}
}
