package main

import (
	regexp "github.com/kamichidu/go-regexp-re"
	"github.com/kolkov/regex-bench/internal/bench"
)

type Engine struct{}

func (e Engine) Name() string                             { return "Go regexp-re" }
func (e Engine) Compile(expr string) (interface{}, error) { return regexp.Compile(expr) }
func (e Engine) Search(re interface{}, data []byte) int {
	return len(re.(*regexp.Regexp).FindAll(data, -1))
}
func (e Engine) Match(re interface{}, data []byte) bool {
	return re.(*regexp.Regexp).Match(data)
}

func main() {
	bench.Main(Engine{}, bench.Extreme)
}
