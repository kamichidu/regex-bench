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
	Matches string // Matches or Match status (yes/no)
	Unit    string
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
	reInput := regexp.MustCompile(`input: ([0-9.]+) MB`)
	reStandard := regexp.MustCompile(`^([a-z_0-9-]+)\s+([0-9.]+)\s+([0-9.]+)\s+([0-9]+)`)
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
			matches := m[4]
			res.Results[pattern] = Result{
				Pattern: pattern,
				Compile: compile,
				Search:  search,
				Matches: matches,
				Unit:    "ms",
			}
		} else if m := reExtreme.FindStringSubmatch(line); m != nil {
			pattern := m[1]
			search, _ := strconv.ParseFloat(m[2], 64)
			unit := m[3]
			matches := m[4]

			val := search
			if unit == "µs" {
				val /= 1000.0
			} else if unit == "ns" {
				val /= 1000000.0
			}

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

		if sc == "extreme" {
			fmt.Println("\n### Search Time & Match Status (Speedup vs stdlib)")
		} else {
			fmt.Println("\n### Search Time in ms (Speedup vs stdlib)")
		}
		fmt.Println(header)
		fmt.Println(sep)

		for _, p := range patterns {
			row := "| " + p + " | "

			// Get stdlib baseline
			var stdlibTime float64
			for _, er := range allEngineResults {
				if er.Name == "stdlib" {
					if r, ok := er.Results[p]; ok {
						stdlibTime = r.Search
					}
					break
				}
			}

			for _, er := range allEngineResults {
				if r, ok := er.Results[p]; ok {
					ratioStr := ""
					if stdlibTime > 0 && er.Name != "stdlib" {
						if r.Search > 0 {
							if r.Search < stdlibTime {
								// Faster
								ratio := stdlibTime / r.Search
								if ratio >= 1.05 {
									ratioStr = fmt.Sprintf(" (%.1fx)", ratio)
								}
							} else {
								// Slower
								ratio := r.Search / stdlibTime
								if ratio >= 1.05 {
									ratioStr = fmt.Sprintf(" (-%.1fx)", ratio)
								}
							}
						}
					}

					if sc == "extreme" {
						origVal := r.Search
						if r.Unit == "µs" {
							origVal *= 1000.0
						}
						if r.Unit == "ns" {
							origVal *= 1000000.0
						}
						row += fmt.Sprintf("%.2f %s%s | ", origVal, r.Unit, ratioStr)
					} else {
						row += fmt.Sprintf("%.2f%s | ", r.Search, ratioStr)
					}
				} else {
					row += "- | "
				}
			}
			fmt.Println(row)
		}

		// Mermaid Charts
		fmt.Printf("\n### Visualizations: %s Performance by Pattern\n", strings.Title(sc))
		for _, p := range patterns {
			fmt.Printf("\n#### Pattern: %s\n", p)
			fmt.Println("```mermaid")
			fmt.Println("xychart-beta")
			fmt.Printf("    title \"Search Time (ms) - %s\"\n", p)

			engineNames := []string{}
			for _, er := range allEngineResults {
				engineNames = append(engineNames, `"`+er.Name+`"`)
			}
			fmt.Printf("    x-axis [%s]\n", strings.Join(engineNames, ", "))
			fmt.Println("    y-axis \"Time (ms)\"")

			values := []string{}
			for _, er := range allEngineResults {
				if r, ok := er.Results[p]; ok {
					values = append(values, fmt.Sprintf("%.4f", r.Search))
				} else {
					values = append(values, "0")
				}
			}
			fmt.Printf("    bar [%s]\n", strings.Join(values, ", "))
			fmt.Println("```")
		}
	}
}
