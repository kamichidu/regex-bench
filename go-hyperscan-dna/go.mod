module github.com/kolkov/regex-bench/go-hyperscan-dna

go 1.25.4

require (
	github.com/flier/gohs v1.2.3
	github.com/kolkov/regex-bench/internal/bench v0.0.0
)

replace github.com/kolkov/regex-bench/internal/bench => ../internal/bench
