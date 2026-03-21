# regex-bench

Cross-language regex benchmark for **real-world patterns**.

Created to provide data for [golang/go#26623](https://github.com/golang/go/issues/26623) discussion on Go regex performance.

## Test Environment

All benchmarks run on **identical conditions**:
- **OS**: Linux (Ubuntu via WSL2 or GitHub Actions)
- **Input**: 6.0 MB generated text file
- **Method**: Each engine compiled natively, same input file, same patterns

> **Note**: Cross-compiled Go binaries run in WSL2 for fair comparison with Rust.

## Results (v0.12.15)

**GitHub Actions Ubuntu (AMD EPYC), 6.0 MB input** (using `FindAll` for fair comparison)

| Pattern | Go stdlib | Go coregex | Rust regex | vs stdlib | vs Rust | Winner |
|---------|-----------|------------|------------|-----------|---------|--------|
| inner_literal | 234 ms | **0.35 ms** | 0.28 ms | **667x** | 1.2x slower | Rust |
| email | 262 ms | **0.50 ms** | 0.23 ms | **523x** | 2.1x slower | Rust |
| uri | 258 ms | 0.61 ms | **0.35 ms** | **422x** | 1.7x slower | Rust |
| ip | 497 ms | **2.17 ms** | 12.0 ms | **229x** | **5.5x faster** | coregex |
| multiline_php | 104 ms | **0.50 ms** | 0.68 ms | **207x** | **1.3x faster** | coregex |
| suffix | 236 ms | 1.83 ms | **1.08 ms** | **128x** | 1.6x slower | Rust |
| literal_alt | 475 ms | 4.37 ms | **0.61 ms** | **108x** | 7.1x slower | Rust |
| multi_literal | 1412 ms | 12.76 ms | **4.70 ms** | **110x** | 2.7x slower | Rust |
| version | 169 ms | 1.62 ms | **0.71 ms** | **104x** | 2.2x slower | Rust |
| http_methods | 107 ms | 1.56 ms | **0.70 ms** | **68x** | 2.2x slower | Rust |
| char_class | 579 ms | **41.0 ms** | 50.1 ms | **14x** | **1.2x faster** | coregex |
| alpha_digit | 262 ms | 25.74 ms | **11.96 ms** | **10x** | 2.1x slower | Rust |
| word_digit | 271 ms | 26.09 ms | **11.96 ms** | **10x** | 2.1x slower | Rust |
| word_repeat | 652 ms | 186 ms | **48 ms** | **3x** | 3.8x slower | Rust |
| anchored | 0.00 ms | 0.02 ms | 0.01 ms | ~1x | 2.0x slower | — |
| anchored_php | 0.00 ms | 0.00 ms | 0.01 ms | — | ~same | — |

> **coregex v0.12.15** — Per-goroutine DFA cache, 7 correctness fixes, 38/38 stdlib compat test. Run `make extreme` for 2500x demo.

### Key Findings

**Go coregex v0.12.15 vs Go stdlib:**
- All patterns: **3-667x faster**
- Best: `inner_literal` **667x**, `email` **523x**, `uri` **422x**
- `ip` **229x** (DigitPrefilter)
- `multiline_php` **207x** (MultilineReverseSuffix, **faster than Rust!**)
- `suffix` **128x**, `literal_alt` **108x**, `multi_literal` **110x**
- `http_methods` **68x** (UseNFA with prefilter, correct (?m)^ anchor handling)
- `char_class` **14x** (CharClassSearcher)
- `word_repeat` **3x** (bidirectional DFA fallback)

**Go coregex faster than Rust (4 patterns):**
- `ip`: **coregex 5.5x faster** (2.2ms vs 12.0ms)
- `multiline_php`: **coregex 1.3x faster** (0.50ms vs 0.68ms)
- `char_class`: **coregex 1.2x faster** (41ms vs 50ms)
- `anchored_php`: **~same** (0.00ms vs 0.01ms)

**Rust faster than coregex:**
- `literal_alt`: Rust 7.1x faster (Teddy with more buckets)
- `word_repeat`: Rust 3.8x faster (DFA state acceleration)
- `multi_literal`: Rust 2.7x faster
- `version`: Rust 2.2x faster
- `http_methods`: Rust 2.2x faster (correctness fix: NFA verifies (?m)^ anchor)
- `alpha_digit`, `word_digit`: Rust 2.1x faster
- `email`: Rust 2.1x faster
- `uri`: Rust 1.7x faster
- `suffix`: Rust 1.6x faster

> **Note**: Rust regex has 10+ years of development. coregex optimizations are targeted, not universal.

### Analysis

| Engine | Strengths | Weaknesses |
|--------|-----------|------------|
| **Go stdlib** | Simple, no dependencies | No optimizations, 3.6-926x slower |
| **Go coregex** | Reverse search, SIMD prefilters, Aho-Corasick, bidirectional DFA, **4 patterns faster than Rust**, multiline Rust parity | Teddy gap, word_repeat |
| **Rust regex** | DFA state acceleration, Teddy Fat, mature DFA | inner_literal, ip, suffix, char_class slower than coregex |

**v0.12.15 (Current):**
- Per-goroutine DFA cache (Rust approach) — concurrent thread safety
- 7 correctness fixes (anchor verification, newline boundaries, alternation overflow)
- Stdlib compatibility test: 38/38 patterns match Go stdlib regexp
- **4 patterns faster than Rust**: ip (5.5x), multiline_php (1.3x), char_class (1.2x), anchored_php
- Pre-computed word boundary flags: 30% → 0.3% CPU
- Integrated prefilter+DFA loop (Rust approach)

**Historical Improvements:**
- v0.12.15: Per-goroutine DFA cache, 7 correctness fixes, stdlib compat test (38/38)
- v0.12.14: Concurrent isMatchDFA safety fix (#137)
- v0.12.13: FatTeddy AVX2 fix, prefilter acceleration, AC v0.2.1
- v0.12.1: Bidirectional DFA fallback, bounded repetitions fix (#115), AVX2 Teddy fix (#74)
- v0.12.0: Anti-quadratic guard, DFA loop unrolling, DFA cache clear & continue
- v0.11.4: FindAll multiline fix, 78x faster (Issue #102)
- v0.11.3: UseMultilineReverseSuffix prefix fast path 319-552x (Issue #99)
- v0.11.1: UseMultilineReverseSuffix for multiline patterns (Issue #97)
- v0.11.0: UseAnchoredLiteral 32-133x speedup (Issue #79)
- v0.10.10: ReverseSuffix CharClass Plus fix
- v0.10.9: UTF-8 optimization + fuzz-found bug fixes
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

**GitHub Actions Ubuntu results** (6 MB no-digits data, v0.12.1):

| Pattern | Go stdlib | Go coregex | Speedup |
|---------|-----------|------------|---------|
| ip_nomatch | 422 ms | 166 µs | **2542x** |
| suffix_find | 245 ms | 126 µs | **1945x** |
| phone_nomatch | 143 ms | 166 µs | **863x** |
| inner_nomatch | 229 ms | 382 µs | **598x** |

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
| http_methods | `(?m)^(GET\|POST\|PUT\|DELETE\|PATCH)` | Multiline log parsing | BranchDispatch |
| anchored_php | `^/.*[\w-]+\.php` | URL path matching | UseAnchoredLiteral |
| multiline_php | `(?m)^/.*\.php` | Multiline PHP paths | UseMultilineReverseSuffix |
| word_repeat | `(\w{2,8})+` | Word quantifiers | BoundedBacktracker + DFA fallback |

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

**Auto-generated comparison table** in Job Summary:
- Side-by-side results for all 3 engines
- Speedup calculations (vs stdlib, vs Rust)
- Winner column with bold formatting
- Raw output in collapsible section

## Links

- **coregex**: https://github.com/coregx/coregex
- **Go issue**: https://github.com/golang/go/issues/26623
- **Rust regex**: https://github.com/rust-lang/regex

## License

MIT
