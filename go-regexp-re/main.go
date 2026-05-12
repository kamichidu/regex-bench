package main

import (
	"flag"
	"fmt"
	"os"

	regexp "github.com/kamichidu/go-regexp-re"
	"github.com/kolkov/regex-bench/internal/bench"
)

type Engine struct{}

func (e Engine) Name() string                             { return "Go regexp-re" }
func (e Engine) Compile(expr string) (interface{}, error) { return regexp.Compile(expr) }
func (e Engine) Search(re interface{}, data []byte) int {
	return len(re.(*regexp.Regexp).FindAll(data, -1))
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
