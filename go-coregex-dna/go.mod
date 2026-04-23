module github.com/kolkov/regex-bench/go-coregex-dna

go 1.25.4

require github.com/coregx/coregex v0.12.5-0.20260308171116-cc5d92220dc2

require (
	github.com/coregx/ahocorasick v0.1.0 // indirect
	github.com/kolkov/regex-bench/internal/bench v0.0.0
	golang.org/x/sys v0.40.0 // indirect
)

replace github.com/kolkov/regex-bench/internal/bench => ../internal/bench
