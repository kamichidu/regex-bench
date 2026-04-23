module github.com/kolkov/regex-bench/go-coregex-extreme

go 1.25.4

require github.com/coregx/coregex v0.12.19

require (
	github.com/coregx/ahocorasick v0.2.1 // indirect
	github.com/kolkov/regex-bench/internal/bench v0.0.0
	golang.org/x/sys v0.40.0 // indirect
)

replace github.com/kolkov/regex-bench/internal/bench => ../internal/bench
