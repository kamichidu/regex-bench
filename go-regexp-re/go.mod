module github.com/kolkov/regex-bench/go-regexp-re

go 1.25.4

require (
	github.com/kamichidu/go-regexp-re v0.9.0
	github.com/kolkov/regex-bench/internal/bench v0.0.0
)

replace github.com/kolkov/regex-bench/internal/bench => ../internal/bench
