package bench

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"time"
)

type Scenario int

const (
	Standard Scenario = iota
	DNA
	Extreme
	LangArena
)

type Pattern struct {
	Name    string
	Pattern string
}

type Engine interface {
	Name() string
	Compile(expr string) (interface{}, error)
	Search(re interface{}, data []byte) int
}

var StandardPatterns = []Pattern{
	{"literal_alt", `error|warning|fatal|critical`},
	{"multi_literal", `apple|banana|cherry|date|elderberry|fig|grape|honeydew|kiwi|lemon|mango|orange`},
	{"anchored", `^HTTP/[12]\.[01]`},
	{"inner_literal", `.*@example\.com`},
	{"suffix", `.*\.(txt|log|md)`},
	{"char_class", `[\w]+`},
	{"email", `[\w.+-]+@[\w.-]+\.[\w.-]+`},
	{"uri", `[\w]+://[^/\s?#]+[^\s?#]+(?:\?[^\s#]*)?(?:#[^\s]*)?`},
	{"version", `\d+\.\d+\.\d+`},
	{"ip", `(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])`},
	{"alpha_digit", `[a-zA-Z]+\d+`},
	{"word_digit", `\w+[0-9]+`},
	{"http_methods", `(?m)^(GET|POST|PUT|DELETE|PATCH)`},
	{"anchored_php", `^/.*[\w-]+\.php`},
	{"multiline_php", `(?m)^/.*\.php`},
	{"word_repeat", `(\w{2,8})+`},
}

var DNAPatterns = []Pattern{
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

var ExtremePatterns = []Pattern{
	{"ip_nomatch", `(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])`},
	{"inner_nomatch", `.*error.*`},
	{"suffix_find", `[^\s]+\.txt`},
	{"phone_nomatch", `\d{3}-\d{3}-\d{4}`},
	{"word_repeat", `(\w{2,8})+`},
}

var LangArenaPatterns = []Pattern{
	{"errors", ` [5][0-9]{2} | [4][0-9]{2} `},
	{"bots", `(?i)bot|crawler|scanner|spider|indexing|crawl|robot|spider`},
	{"suspicious", `(?i)etc/passwd|wp-admin|\.\./`},
	{"ips", `\d+\.\d+\.\d+\.35`},
	{"api_calls", `/api/[^ "]+`},
	{"post_requests", `POST [^ ]* HTTP`},
	{"auth_attempts", `(?i)/login|/signin`},
	{"methods", `(?i)get|post|put`},
	{"emails", `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`},
	{"passwords", `password=[^&\s"]+`},
	{"tokens", `token=[^&\s"]+|api[_-]?key=[^&\s"]+`},
	{"sessions", `session[_-]?id=[^&\s"]+`},
	{"peak_hours", `\[\d+/\w+/\d+:1[3-7]:\d+:\d+ [+\-]\d+\]`},
}

func Main(e Engine, s Scenario, inputFile string) {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	var patterns []Pattern
	scenarioName := ""
	switch s {
	case Standard:
		patterns = StandardPatterns
		scenarioName = "Standard"
	case DNA:
		patterns = DNAPatterns
		scenarioName = "DNA"
	case Extreme:
		patterns = ExtremePatterns
		scenarioName = "Extreme"
	case LangArena:
		patterns = LangArenaPatterns
		scenarioName = "LangArena"
	}

	fmt.Printf("%s %s (input: %.2f MB)\n", e.Name(), scenarioName, float64(len(data))/1024/1024)
	fmt.Printf("%-15s %15s %15s %10s\n", "pattern", "compile(ns)", "search(ns)", "matches")
	fmt.Println("---------------------------------------------------------------------")
	for _, p := range patterns {
		measure(e, data, p)
	}
}

func measure(e Engine, data []byte, p Pattern) {
	runtime.GC()
	compileStart := time.Now()
	re, err := e.Compile(p.Pattern)
	compileElapsed := time.Since(compileStart)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error compiling %s: %v\n", p.Name, err)
		return
	}

	// Search Measurement (Median of 11 iterations)
	const iterations = 11
	searchSamples := make([]time.Duration, iterations)
	var matches int
	for i := 0; i < iterations; i++ {
		runtime.GC()
		start := time.Now()
		matches = e.Search(re, data)
		searchSamples[i] = time.Since(start)
	}
	sort.Slice(searchSamples, func(i, j int) bool { return searchSamples[i] < searchSamples[j] })
	medianSearch := searchSamples[iterations/2]

	fmt.Printf("%-15s %15d %15d %10d\n", p.Name, compileElapsed.Nanoseconds(), medianSearch.Nanoseconds(), matches)
}
