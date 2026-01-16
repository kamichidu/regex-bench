# regex-bench

Cross-language regex benchmark for **real-world patterns**.

Created to provide data for [golang/go#26623](https://github.com/golang/go/issues/26623) discussion on Go regex performance.

## Test Environment

All benchmarks run on **identical conditions**:
- **OS**: Linux (Ubuntu via WSL2 or GitHub Actions)
- **Input**: 6.0 MB generated text file
- **Method**: Each engine compiled natively, same input file, same patterns

> **Note**: Cross-compiled Go binaries run in WSL2 for fair comparison with Rust.

## Results (v0.11.4)

**GitHub Actions Ubuntu, 6.0 MB input** (using `FindAll` for fair comparison)

| Pattern | Go stdlib | Go coregex | Rust regex | vs stdlib | vs Rust | Winner |
|---------|-----------|------------|------------|-----------|---------|--------|
| inner_literal | 206 ms | **0.45 ms** | 0.59 ms | **457x** | **1.3x faster** | coregex |
| email | 249 ms | **1.30 ms** | 1.55 ms | **192x** | **1.2x faster** | coregex |
| uri | 242 ms | 1.65 ms | **1.05 ms** | **147x** | 1.6x slower | Rust |
| multiline_php | 103 ms | **~1 ms** | 0.78 ms | **100x+** | ~1.3x slower | **near parity** |
| literal_alt | 434 ms | 4.22 ms | **0.91 ms** | **103x** | 4.6x slower | Rust |
| multi_literal | 1300 ms | 13.06 ms | **4.56 ms** | **99x** | 2.9x slower | Rust |
| suffix | 203 ms | 2.12 ms | **1.60 ms** | **96x** | 1.3x slower | Rust |
| http_methods | 95 ms | 1.04 ms | **0.72 ms** | **92x** | 1.4x slower | Rust |
| ip | 468 ms | **6.37 ms** | 11.55 ms | **73x** | **1.8x faster** | coregex |
| version | 154 ms | 2.78 ms | **0.91 ms** | **55x** | 3.1x slower | Rust |
| char_class | 494 ms | **49.31 ms** | 52.44 ms | **10x** | **1.06x faster** | coregex |
| alpha_digit | 242 ms | 41.49 ms | **11.34 ms** | **5.8x** | 3.7x slower | Rust |
| word_digit | 251 ms | 41.60 ms | **11.77 ms** | **6.0x** | 3.5x slower | Rust |
| anchored | 0.02 ms | 0.03 ms | 0.07 ms | ~1x | ~1x | — |
| anchored_php | 0.02 ms | 0.23 ms | 0.41 ms | — | — | — |

> **coregex v0.11.4** — FindAll multiline fix: 78x faster, near Rust parity (Issue #102). Run `make extreme` for 1800x demo.

### Key Findings

**Go coregex v0.11.4 vs Go stdlib:**
- All patterns: **5.8-457x faster**
- Best: `inner_literal` **457x**, `email` **192x**, `uri` **147x**, `literal_alt` **103x**
- `multiline_php` **100x+** (MultilineReverseSuffix, v0.11.4 fix)
- `multi_literal` **99x** (Aho-Corasick)
- `suffix` **96x** (ReverseSuffix)
- `http_methods` **92x** (multiline log parsing with `(?m)^`)
- `ip` **73x** (DigitPrefilter)
- `version` **55x** (DigitPrefilter)
- `char_class` **10x** (CharClassSearcher)

**Go coregex faster than Rust (4 patterns):**
- `ip`: **coregex 1.8x faster** (6.4ms vs 11.6ms)
- `inner_literal`: **coregex 1.3x faster** (0.45ms vs 0.59ms)
- `email`: **coregex 1.2x faster** (1.30ms vs 1.55ms)
- `char_class`: **coregex 1.06x faster** (49ms vs 52ms)

**Near Rust parity:**
- `multiline_php`: **~1.3x slower** (~1ms vs 0.78ms) — was 84x slower in v0.11.1!

**Rust faster than coregex:**
- `literal_alt`: Rust 4.6x faster (Teddy Fat with more buckets)
- `alpha_digit`: Rust 3.7x faster
- `word_digit`: Rust 3.5x faster
- `version`: Rust 3.1x faster
- `multi_literal`: Rust 2.9x faster
- `uri`: Rust 1.6x faster
- `http_methods`: Rust 1.4x faster
- `suffix`: Rust 1.3x faster

> **Note**: Rust regex has 10+ years of development. coregex optimizations are targeted, not universal.

### Analysis

| Engine | Strengths | Weaknesses |
|--------|-----------|------------|
| **Go stdlib** | Simple, no dependencies | No optimizations, 5.8-457x slower |
| **Go coregex** | Reverse search, SIMD prefilters, Aho-Corasick, **4 patterns faster than Rust**, multiline near-parity | Teddy gap |
| **Rust regex** | DFA state acceleration, Teddy Fat, mature DFA | inner_literal, ip, email, char_class slower than coregex |

**v0.11.4 (Current):**
- FindAll multiline fix: **78x faster**, near Rust parity (Issue #102)
- `multiline_php`: was 84x slower → now **~1.3x slower** than Rust
- **4 patterns faster than Rust**: ip (1.8x), inner_literal (1.3x), email (1.2x), char_class (1.06x)

**Historical Improvements:**
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

**GitHub Actions Ubuntu results** (6 MB no-digits data, v0.10.10):

| Pattern | Go stdlib | Go coregex | Speedup |
|---------|-----------|------------|---------|
| ip_nomatch | 392 ms | 215 µs | **1820x** |
| suffix_find | 226 ms | 218 µs | **1037x** |
| inner_nomatch | 210 ms | 254 µs | **826x** |
| phone_nomatch | 132 ms | 216 µs | **613x** |

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
