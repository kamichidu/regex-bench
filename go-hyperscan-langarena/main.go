package main

import (
	"fmt"
	"github.com/flier/gohs/hyperscan"
	"github.com/kolkov/regex-bench/internal/bench"
)

type Engine struct{}

func (e Engine) Name() string { return "Go Hyperscan" }
func (e Engine) Compile(expr string) (interface{}, error) { return hyperscan.Compile(expr) }
func (e Engine) Search(re interface{}, data []byte) int {
	db := re.(hyperscan.BlockDatabase)
	scratch, err := hyperscan.NewScratch(db)
	if err != nil {
		return 0
	}
	defer scratch.Free()
	count := 0
	db.Scan(data, scratch, func(id uint, from, to uint64, flags uint, context interface{}) error {
		count++
		return nil
	}, nil)
	return count
}
func (e Engine) Match(re interface{}, data []byte) bool {
	db := re.(hyperscan.BlockDatabase)
	scratch, err := hyperscan.NewScratch(db)
	if err != nil {
		return false
	}
	defer scratch.Free()
	matched := false
	db.Scan(data, scratch, func(id uint, from, to uint64, flags uint, context interface{}) error {
		matched = true
		return fmt.Errorf("found")
	}, nil)
	return matched
}

func main() {
	bench.Main(Engine{}, bench.LangArena)
}
