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
	Match   float64
	Matches string
	Unit    string
}

type EngineResult struct {
	Name      string
	InputSize float64
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
	reInput := regexp.MustCompile(`input: ([0-9.]+) MB`)
	// Updated Standard: pattern compile search match matches
	reStandard := regexp.MustCompile(`^([a-z_0-9-]+)\s+([0-9.]+)\s+([0-9.]+)\s+([0-9.]+)\s+([0-9]+)`)
	reExtreme := regexp.MustCompile(`^([a-z_0-9-]+)\s+([0-9.]+)\s+(ms|µs|ns)\s+match:\s+(yes|no)`)

	for scanner.Scan() {
		line := scanner.Text()
		if m := reInput.FindStringSubmatch(line); m != nil {
			res.InputSize, _ = strconv.ParseFloat(m[1], 64)
		}
		if m := reStandard.FindStringSubmatch(line); m != nil {
			pattern := m[1]
			compile, _ := strconv.ParseFloat(m[2], 64)
			search, _ := strconv.ParseFloat(m[3], 64)
			match, _ := strconv.ParseFloat(m[4], 64)
			matches := m[5]
			res.Results[pattern] = Result{
				Pattern: pattern,
				Compile: compile,
				Search:  search,
				Match:   match,
				Matches: matches,
				Unit:    "ms",
			}
		} else if m := reExtreme.FindStringSubmatch(line); m != nil {
			pattern := m[1]
			search, _ := strconv.ParseFloat(m[2], 64)
			unit := m[3]
			matches := m[4]
			val := search
			if unit == "µs" { val /= 1000.0 } else if unit == "ns" { val /= 1000000.0 }
			res.Results[pattern] = Result{
				Pattern: pattern,
				Search:  val,
				Matches: matches,
				Unit:    unit,
			}
		}
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
			if res, err := parseFile(filepath.Join("results", sc, en+".txt")); err == nil {
				allEngineResults = append(allEngineResults, res)
				for p := range res.Results { patternsMap[p] = true }
			}
		}
		if len(allEngineResults) == 0 { continue }
		patterns := make([]string, 0, len(patternsMap)); for p := range patternsMap { patterns = append(patterns, p) }; sort.Strings(patterns)

		// Table for Search (All Matches)
		fmt.Printf("\n### Search Performance (FindAll)\n")
		printTable(allEngineResults, patterns, sc, false)

		// Table for Match (Boolean) - only for non-extreme
		if sc != "extreme" {
			fmt.Printf("\n### Match Performance (Boolean)\n")
			printTable(allEngineResults, patterns, sc, true)
		}
	}
}

func printTable(engines []*EngineResult, patterns []string, sc string, isMatch bool) {
	header := "| Pattern | "
	sep := "|---|"
	for _, er := range engines { header += er.Name + " | "; sep += "---|" }
	fmt.Println(header); fmt.Println(sep)

	for _, p := range patterns {
		row := "| " + p + " | "
		var baseline float64
		for _, er := range engines {
			if er.Name == "stdlib" {
				if r, ok := er.Results[p]; ok {
					if isMatch { baseline = r.Match } else { baseline = r.Search }
				}
				break
			}
		}

		for _, er := range engines {
			if r, ok := er.Results[p]; ok {
				val := r.Search
				if isMatch { val = r.Match }
				
				ratioStr := ""
				if baseline > 0 && er.Name != "stdlib" && val > 0 {
					if val < baseline {
						ratioStr = fmt.Sprintf(" (%.1fx)", baseline/val)
					} else {
						ratioStr = fmt.Sprintf(" (-%.1fx)", val/baseline)
					}
				}
				if sc == "extreme" {
					orig := r.Search; if r.Unit == "µs" { orig *= 1000 }; if r.Unit == "ns" { orig *= 1000000 }
					row += fmt.Sprintf("%.2f %s%s | ", orig, r.Unit, ratioStr)
				} else {
					row += fmt.Sprintf("%.2f%s | ", val, ratioStr)
				}
			} else { row += "- | " }
		}
		fmt.Println(row)
	}
}
