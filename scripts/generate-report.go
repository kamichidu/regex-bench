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
	Compile int64
	Search  int64
	Matches string
}

type EngineResult struct {
	Name      string
	InputSize float64
	Results   map[string]Result
}

func parseFile(path string) (*EngineResult, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		panic(err)
	}
	defer file.Close()

	engineName := strings.TrimSuffix(filepath.Base(path), ".txt")
	res := &EngineResult{
		Name:    engineName,
		Results: make(map[string]Result),
	}

	scanner := bufio.NewScanner(file)
	reInput := regexp.MustCompile(`input: ([0-9.]+) MB`)
	// New format: pattern compile(ns) search(ns) matches
	reStandard := regexp.MustCompile(`^([a-z_0-9-]+)\s+([0-9]+)\s+([0-9]+)\s+([0-9]+)`)

	for scanner.Scan() {
		line := scanner.Text()
		if m := reInput.FindStringSubmatch(line); m != nil {
			val, err := strconv.ParseFloat(m[1], 64)
			if err != nil {
				panic(err)
			}
			res.InputSize = val
		}
		if m := reStandard.FindStringSubmatch(line); m != nil {
			pattern := m[1]
			compile, err := strconv.ParseInt(m[2], 10, 64)
			if err != nil {
				panic(err)
			}
			search, err := strconv.ParseInt(m[3], 10, 64)
			if err != nil {
				panic(err)
			}
			matches := m[4]
			res.Results[pattern] = Result{
				Pattern: pattern,
				Compile: compile,
				Search:  search,
				Matches: matches,
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return res, nil
}

func main() {
	scenarios := []string{"standard", "dna", "extreme", "langarena"}
	engines := []string{"stdlib", "regexp-re", "coregex", "re2-wasm", "re2-cgo", "pcre2", "hyperscan"}

	fmt.Println("# Benchmark Report")

	for _, sc := range scenarios {
		fmt.Printf("\n## Scenario: %s\n", strings.Title(sc))
		var allEngineResults []*EngineResult
		patternsMap := make(map[string]bool)
		for _, en := range engines {
			res, err := parseFile(filepath.Join("results", sc, en+".txt"))
			if err != nil {
				continue
			}
			allEngineResults = append(allEngineResults, res)
			for p := range res.Results {
				patternsMap[p] = true
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

		// Table for Search Performance
		fmt.Printf("\n### Search Performance (ns)\n")
		printTable(allEngineResults, patterns)
	}
}

func printTable(engines []*EngineResult, patterns []string) {
	header := "| Pattern | "
	sep := "|---|"
	for _, er := range engines {
		header += er.Name + " | "
		sep += "---|"
	}
	fmt.Println(header)
	fmt.Println(sep)

	for _, p := range patterns {
		row := "| " + p + " | "
		var baseline int64
		for _, er := range engines {
			if er.Name == "stdlib" {
				if r, ok := er.Results[p]; ok {
					baseline = r.Search
				}
				break
			}
		}

		for _, er := range engines {
			if r, ok := er.Results[p]; ok {
				val := r.Search
				ratioStr := ""
				if baseline > 0 && er.Name != "stdlib" && val > 0 {
					if val < baseline {
						ratioStr = fmt.Sprintf(" (%.1fx)", float64(baseline)/float64(val))
					} else {
						ratioStr = fmt.Sprintf(" (-%.1fx)", float64(val)/float64(baseline))
					}
				}
				row += fmt.Sprintf("%d%s | ", val, ratioStr)
			} else {
				row += "- | "
			}
		}
		fmt.Println(row)
	}
}
