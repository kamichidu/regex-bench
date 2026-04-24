package main

import (
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
func (e Engine) Match(re interface{}, data []byte) bool {
	return re.(*pcre2.Regexp).Matcher(data, 0).Match(data, 0)
}

func main() {
	bench.Main(Engine{}, bench.Standard)
}
