package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Jemmic/go-pcre2"
	"github.com/kolkov/regex-bench/internal/bench"
)

type Engine struct{}

func (e Engine) Name() string                             { return "Go PCRE2" }
func (e Engine) Compile(expr string) (interface{}, error) { return pcre2.Compile(expr, 0) }
func (e Engine) Search(re interface{}, data []byte) int {
	count := 0
	r := re.(*pcre2.Regexp)
	matcher := r.Matcher(data, 0)
	for matcher.Match(data, 0) {
		count++
		indices := matcher.Index()
		if len(indices) < 2 {
			break
		}
		end := indices[1]
		if end >= len(data) {
			break
		}
		data = data[end:]
		matcher = r.Matcher(data, 0)
	}
	return count
}

func main() {
	var scenarioStr string
	flag.StringVar(&scenarioStr, "scenario", "standard", "benchmark scenario (standard, dna, extreme, langarena)")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Printf("Usage: %s [-scenario <scenario>] <input-file>\n", os.Args[0])
		os.Exit(1)
	}

	var s bench.Scenario
	switch scenarioStr {
	case "standard":
		s = bench.Standard
	case "dna":
		s = bench.DNA
	case "extreme":
		s = bench.Extreme
	case "langarena":
		s = bench.LangArena
	default:
		fmt.Fprintf(os.Stderr, "Unknown scenario: %s\n", scenarioStr)
		os.Exit(1)
	}

	bench.Main(Engine{}, s, args[0])
}
