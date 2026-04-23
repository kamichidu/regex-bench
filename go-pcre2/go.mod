module github.com/kolkov/regex-bench/go-pcre2

go 1.25.4

require github.com/Jemmic/go-pcre2 v0.0.0-20190111114109-bd52ad5f7098

require (
	github.com/kolkov/regex-bench/internal/bench v0.0.0
	github.com/stretchr/testify v1.11.1 // indirect
)

replace github.com/kolkov/regex-bench/internal/bench => ../internal/bench
