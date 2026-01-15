# regex-bench

Cross-language regex benchmark for **real-world patterns**.

Created to provide data for [golang/go#26623](https://github.com/golang/go/issues/26623) discussion on Go regex performance.

## Test Environment

All benchmarks run on **identical conditions**:
- **OS**: Linux (Ubuntu via WSL2 or GitHub Actions)
- **Input**: 6.0 MB generated text file
- **Method**: Each engine compiled natively, same input file, same patterns

> **Note**: Cross-compiled Go binaries run in WSL2 for fair comparison with Rust.

## Results

**GitHub Actions Ubuntu, 6.0 MB input** (using `FindAll` for fair comparison)

| Pattern | Go stdlib | Go coregex | Rust regex | vs stdlib | vs Rust |
|---------|-----------|------------|------------|-----------|---------|
| inner_literal | 232 ms | **0.31 ms** | 0.56 ms | **750x** | **1.8x faster** |
| email | 259 ms | **1.26 ms** | 1.45 ms | **206x** | **1.2x faster** |
| uri | 257 ms | 1.92 ms | **0.93 ms** | **134x** | 2.1x slower |
| suffix | 237 ms | 1.87 ms | **1.33 ms** | **127x** | 1.4x slower |
| multi_literal | 1404 ms | 13.10 ms | **5.08 ms** | **107x** | 2.6x slower |
| literal_alt | 495 ms | 4.83 ms | **0.84 ms** | **102x** | 5.8x slower |
| ip | 492 ms | **6.72 ms** | 12.13 ms | **73x** | **1.8x faster** |
| version | 168 ms | 2.45 ms | **0.68 ms** | **68x** | 3.6x slower |
| char_class | 523 ms | **38.43 ms** | 52.04 ms | **14x** | **1.4x faster** |
| alpha_digit | 260 ms | 40.71 ms | **11.96 ms** | **6x** | 3.4x slower |
| word_digit | 270 ms | 40.90 ms | **12.32 ms** | **7x** | 3.3x slower |
| anchored | 0.05 ms | **0.03 ms** | 0.03 ms | **1.7x** | ~1x |
| http_methods | 0.02 ms | 0.04 ms | 0.06 ms | 0.5x | **1.5x faster** |

> **coregex v0.10.8** — FindAll 600x faster for anchored patterns. Run `make extreme` for 3000x demo.

### Key Findings

**Go coregex v0.10.8 vs Go stdlib:**
- Most patterns: **6-750x faster**
- Best: `inner_literal` **750x**, `email` **206x**, `uri` **134x**, `suffix` **127x**
- `multi_literal` **107x** (Aho-Corasick)
- `literal_alt` **102x** (Teddy SIMD)
- `ip` **73x** (DigitPrefilter)
- `version` **68x** (DigitPrefilter)
- `char_class` **14x** (CharClassSearcher)

**Go coregex faster than Rust (5 patterns):**
- `inner_literal`: **coregex 1.8x faster** (0.31ms vs 0.56ms)
- `ip`: **coregex 1.8x faster** (6.7ms vs 12.1ms)
- `http_methods`: **coregex 1.5x faster** (0.04ms vs 0.06ms)
- `char_class`: **coregex 1.4x faster** (38ms vs 52ms)
- `email`: **coregex 1.2x faster** (1.26ms vs 1.45ms)

**Rust faster than coregex:**
- `literal_alt`: Rust 5.8x faster (Teddy Fat with more buckets)
- `version`: Rust 3.6x faster
- `alpha_digit`: Rust 3.4x faster
- `word_digit`: Rust 3.3x faster
- `multi_literal`: Rust 2.6x faster
- `uri`: Rust 2.1x faster

> **Note**: Rust regex has 10+ years of development. coregex optimizations are targeted, not universal.

### Analysis

| Engine | Strengths | Weaknesses |
|--------|-----------|------------|
| **Go stdlib** | Simple, no dependencies | No optimizations, 6-750x slower |
| **Go coregex** | Reverse search, SIMD prefilters, Aho-Corasick, **5 patterns faster than Rust** | Teddy gap vs Rust |
| **Rust regex** | Aho-Corasick (any count), mature DFA, overall fastest | inner_literal, ip, char_class, email, http_methods slower than coregex |

**v0.10.8 (Current):**
- **FindAll 600x faster** for anchored patterns (`^...`)
- **5 patterns faster than Rust**: inner_literal (1.8x), ip (1.8x), http_methods (1.5x), char_class (1.4x), email (1.2x)
- Gap vs Rust narrowing on most patterns

**Historical Improvements:**
- v0.10.8: FindAll allocation fix for anchored patterns
- v0.10.7: UTF-8 fixes + 100% stdlib API compatibility
- v0.10.5: CompositeSearcher backtracking fix
- v0.10.0: Fat Teddy AVX2 (33-64 patterns, 9+ GB/s)
- v0.9.5: Aho-Corasick integration, Teddy 32 patterns

## Extreme Speedups (1000-3000x)

The "3-3000x faster" claim refers to **specific edge cases** where coregex prefilters can skip entire input:

```bash
make extreme       # Run on no-match data (~300-560x)
make extreme-3000x # Run on no-digits data (1000-3000x)
```

**GitHub Actions Ubuntu results** (6 MB no-digits data):

| Pattern | Go stdlib | Go coregex | Speedup |
|---------|-----------|------------|---------|
| ip_nomatch | 389 ms | 214 µs | **1816x** |
| suffix_find | 225 ms | 217 µs | **1039x** |
| inner_nomatch | 208 ms | 253 µs | **823x** |
| phone_nomatch | 131 ms | 218 µs | **604x** |

[![Extreme Benchmark](https://github.com/kolkov/regex-bench/actions/workflows/extreme-benchmark.yml/badge.svg)](https://github.com/kolkov/regex-bench/actions/workflows/extreme-benchmark.yml)

> **Note**: Results vary between runs (±30%) due to CI VM load and OS scheduling.
> The key insight: coregex operates in **microseconds**, stdlib in **hundreds of milliseconds**.

**When do we see 3000x?**

The 3000x speedup occurs in coregex's own benchmark suite (`go test -bench`) under specific conditions:
- **Pattern**: IP regex on data with NO IP addresses
- **Size**: 1 MB of pure text
- **Measurement**: `go test -bench` with multiple iterations

```go
// In coregex repo:
BenchmarkIPRegex_Find/stdlib_1MB_no_ips    74.5ms
BenchmarkIPRegex_Find/coregex_1MB_no_ips   22.4µs  // 3324x
```

The extreme speedup happens because:
1. **DigitPrefilter** scans for first digit character
2. No digits in input → entire 1 MB skipped in ~20µs
3. stdlib must scan byte-by-byte → 74ms

**Verified speedups** (from coregex repo, `docs/dev/SPEEDUP_VERIFICATION.md`):

| Pattern | Strategy | Max Speedup |
|---------|----------|-------------|
| IP no-match (1MB) | DigitPrefilter | **3324x** |
| `.*\.txt$` (1MB) | ReverseSuffix | **1124x** |
| `.*error.*` (32KB) | ReverseInner | **909x** |

> The speedup depends on input characteristics. Real-world mixed data shows 15-560x.

## Patterns Tested

| Name | Pattern | Type | Optimization |
|------|---------|------|--------------|
| literal_alt | `error\|warning\|fatal\|critical` | 4-literal alternation | Teddy SIMD |
| multi_literal | `apple\|banana\|...\|orange` | 12-literal alternation | **Aho-Corasick** |
| anchored | `^HTTP/[12]\.[01]` | Start anchor | — |
| inner_literal | `.*@example\.com` | Inner literal | Reverse search |
| suffix | `.*\.(txt\|log\|md)` | Suffix match | Reverse search |
| char_class | `[\w]+` | Character class | CharClassSearcher |
| email | `[\w.+-]+@[\w.-]+\.[\w.-]+` | Complex real-world | Memmem SIMD |
| uri | `[\w]+://[^/\s?#]+[^\s?#]+...` | URL with query/fragment | Memmem SIMD |
| version | `\d+\.\d+\.\d+` | Version numbers | DigitPrefilter |
| ip | `(?:(?:25[0-5]\|2[0-4][0-9]\|...)\.){3}...` | IPv4 validation | DigitPrefilter + LazyDFA |

## Running Benchmarks

```bash
# Generate input data (6 MB)
go run scripts/generate-input.go

# Build for Linux
cd go-stdlib && GOOS=linux GOARCH=amd64 go build -o ../bin/go-stdlib-linux . && cd ..
cd go-coregex && GOOS=linux GOARCH=amd64 go build -o ../bin/go-coregex-linux . && cd ..

# Run all in WSL/Linux for fair comparison
wsl ./bin/go-stdlib-linux input/data.txt
wsl ./bin/go-coregex-linux input/data.txt
wsl ./bin/rust-benchmark input/data.txt
```

## CI Benchmarks

Benchmarks run automatically on GitHub Actions (Ubuntu) for reproducible results.

[![Benchmark](https://github.com/kolkov/regex-bench/actions/workflows/benchmark.yml/badge.svg)](https://github.com/kolkov/regex-bench/actions/workflows/benchmark.yml)

## Links

- **coregex**: https://github.com/coregx/coregex
- **Go issue**: https://github.com/golang/go/issues/26623
- **Rust regex**: https://github.com/rust-lang/regex

## License

MIT
