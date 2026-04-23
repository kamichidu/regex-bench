module github.com/kolkov/regex-bench/go-regexp-re-extreme

go 1.25.4

require (
	github.com/kamichidu/go-regexp-re v0.0.0-20260423043213-f2e5d53fa0c4
	github.com/kolkov/regex-bench/internal/bench v0.0.0
)

replace github.com/kolkov/regex-bench/internal/bench => ../internal/bench
