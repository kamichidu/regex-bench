# Go Regex Benchmark Project

This project is a benchmark suite designed to fairly and comprehensively compare the performance of various regular expression engines available in the Go ecosystem.

## Project Principles

- **Go-Centric**: Focused on the Go standard `regexp` library and alternative libraries accessible from Go.
- **Diverse Implementations**: Compares Pure Go implementations, WebAssembly (WASM) builds, and native libraries via CGO (RE2, PCRE2, Hyperscan).
- **Standardized Interfacing**: Measures performance using logic consistent with the standard `regexp` package wherever possible, while allowing for engine-specific optimizations.

## Project Structure

- `go-stdlib*`: Go Standard Library `regexp`.
- `go-regexp-re*`: Optimized Go-native implementation (`github.com/kamichidu/go-regexp-re`).
- `go-coregex*`: Speedup via pre-filters (`github.com/m-mizutani/coregex`).
- `go-re2-wasm*`: RE2 via WASM (`github.com/wasilibs/go-re2`).
- `go-re2-cgo*`: RE2 via CGO.
- `go-pcre2*`: PCRE2 via CGO (`github.com/Jemmic/go-pcre2`).
- `go-hyperscan*`: Intel Hyperscan via CGO (`github.com/flier/gohs/hyperscan`).

## Development Workflow

When adding a new regex engine, create a new directory by copying `go-stdlib` or an existing CGO-based engine as a template.

### Build Dependencies
The following system libraries are required to build engines using CGO:
- `libpcre2-dev`
- `libhyperscan-dev`
- `libre2-dev` (for RE2 CGO)

## Benchmark Scenarios

1. **Standard**: Common regular expression patterns.
2. **DNA**: DNA sequence pattern matching.
3. **Extreme**: Measurement of worst-case scenarios (e.g., no-match data) to test pre-filter efficiency.
4. **LangArena**: 13 real-world LogParser patterns from the LangArena benchmark.
