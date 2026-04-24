package main

import (
	"github.com/coregx/coregex"
	"github.com/kolkov/regex-bench/internal/bench"
)

type Engine struct{}

func (e Engine) Name() string                             { return "Go coregex" }
func (e Engine) Compile(expr string) (interface{}, error) { return coregex.Compile(expr) }
func (e Engine) Search(re interface{}, data []byte) int {
	return len(re.(*coregex.Regexp).FindAll(data, -1))
}
func (e Engine) Match(re interface{}, data []byte) bool {
	return re.(*coregex.Regexp).Match(data)
}

func main() {
	bench.Main(Engine{}, bench.DNA)
}
