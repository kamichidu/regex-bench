# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

---

## [2026-01-16] - Multiline Pattern Benchmark (v0.11.1)

### Added
- **New pattern: `multiline_php`** (`(?m)^/.*\.php`)
  - Tests UseMultilineReverseSuffix strategy (Issue #97)
  - Multiline suffix matching for log files

### Changed
- **Updated coregex v0.11.0 → v0.11.1**
  - UseMultilineReverseSuffix: 1.4-5.7x faster than stdlib for `(?m)^.*suffix`
  - Known gap: Rust is 84x faster (DFA state acceleration, tracked in Issue #99)

### Results (GitHub Actions Ubuntu, 6MB input)

| Pattern | Go stdlib | Go coregex | Rust regex | vs stdlib | Winner |
|---------|-----------|------------|------------|-----------|--------|
| inner_literal | 206 ms | **0.45 ms** | 0.59 ms | **457x** | coregex |
| email | 249 ms | **1.30 ms** | 1.55 ms | **192x** | coregex |
| uri | 242 ms | 1.65 ms | **1.05 ms** | **147x** | Rust |
| literal_alt | 434 ms | 4.22 ms | **0.91 ms** | **103x** | Rust |
| multi_literal | 1300 ms | 13.06 ms | **4.56 ms** | **99x** | Rust |
| suffix | 203 ms | 2.12 ms | **1.60 ms** | **96x** | Rust |
| http_methods | 95 ms | 1.04 ms | **0.72 ms** | **92x** | Rust |
| ip | 468 ms | **6.37 ms** | 11.55 ms | **73x** | coregex |
| version | 154 ms | 2.78 ms | **0.91 ms** | **55x** | Rust |
| char_class | 494 ms | **49.31 ms** | 52.44 ms | **10x** | coregex |
| multiline_php | 93 ms | 66.48 ms | **0.79 ms** | **1.4x** | Rust |

**Summary:**
- coregex wins: inner_literal, email, ip, char_class (4 patterns)
- Rust wins: uri, literal_alt, multi_literal, suffix, http_methods, version, multiline_php (7 patterns)
- Multiline gap: Rust 84x faster due to DFA state acceleration

---

## [2026-01-15] - ReverseSuffix CharClass Plus Fix

### Changed
- **Updated coregex v0.10.9 → v0.10.10**
  - Fix: CharClass Plus patterns (`[^\s]+`, `[\w]+`) now use ReverseSuffix optimization
  - Bug: `[^\s]+\.txt` caused extreme benchmark to hang (266ms/MB instead of µs)
  - Result: All extreme patterns now complete in µs

### Results (GitHub Actions, 6MB no-digits data)
| Pattern | stdlib | coregex | Speedup |
|---------|--------|---------|---------|
| ip_nomatch | 392 ms | 215 µs | **1820x** |
| suffix_find | 226 ms | 218 µs | **1037x** |
| inner_nomatch | 210 ms | 254 µs | **826x** |
| phone_nomatch | 132 ms | 216 µs | **613x** |

---

## [2026-01-15] - UTF-8 Optimization + Fuzz Bug Fixes

### Changed
- **Updated coregex v0.10.8 → v0.10.9**
  - UTF-8 suffix sharing reduces dot NFA states 39→28
  - Anchored suffix prefilter for O(1) rejection
  - CharClassSearcher excludes `*` patterns (zero-width match fix)
  - Invalid UTF-8 handling for negated char classes (stdlib compat)
  - ReverseInner/ReverseSuffix whitelist (strategy safety)

### Results
- No regressions, all speedups maintained
- Local benchmarks: 113x-389x+ on various patterns

---

## [2026-01-15] - FindAll Anchored Pattern Fix

### Changed
- **Updated coregex v0.10.5 → v0.10.8**
  - v0.10.8: FindAll 600x faster for anchored patterns (#92)
  - v0.10.7: UTF-8 fixes + 100% stdlib API compatibility
  - v0.10.6: CompositeSequenceDFA for overlapping patterns

### Fixed
- **Anchored patterns (`^...`) allocation fix**
  - Before: 0.21 ms (huge allocation for 6MB input)
  - After: 0.03 ms (smart allocation: cap=1 for anchored)
  - Start-anchored patterns can only match at position 0

### Results
- **5 patterns now faster than Rust** (was 3):
  - `inner_literal`: 1.8x faster (0.31ms vs 0.56ms)
  - `ip`: 1.8x faster (6.7ms vs 12.1ms)
  - `http_methods`: 1.5x faster
  - `char_class`: 1.4x faster
  - `email`: 1.2x faster

---

## [2026-01-14] - Overlapping Char Classes Fix

### Changed
- **Updated coregex v0.10.4 → v0.10.5**
  - Critical fix: `\w+[0-9]+` patterns now work correctly (#81)
  - Bug: Greedy `\w+` consumed all characters (including digits)
  - Fix: Recursive backtracking in CompositeSearcher
  - `word_digit` pattern now returns correct matches

### Results
- `word_digit` (`\w+[0-9]+`): Now finds 3575 matches (was 0)
- No performance regression (+0.10% geomean)

---

## [2026-01-14] - Thread-safety Release

### Changed
- **Updated coregex v0.10.3 → v0.10.4**
  - Critical fix: Panic on concurrent usage of compiled Regexp (#78)
  - Implements `sync.Pool` pattern (same as Go stdlib `regexp`)
  - Thread-safe concurrent access to `*Regexp` instances
  - 32-bit platform compatibility (atomic operations alignment)
  - No performance regression (-3.84% improvement in geomean)

### Added
- **New benchmark patterns:**

  | Pattern | Regex | Category | Purpose |
  |---------|-------|----------|---------|
  | `alpha_digit` | `[a-zA-Z]+\d+` | Composite | Concatenated char classes |
  | `word_digit` | `\w+[0-9]+` | Composite | Word followed by digits |
  | `http_methods` | `^(GET\|POST\|PUT\|DELETE\|PATCH)` | Anchored | HTTP method dispatch |

- Added patterns to all 3 benchmark suites:
  - `go-coregex/main.go`
  - `go-stdlib/main.go`
  - `rust/src/main.rs`

### Technical
- These patterns test upcoming CompositeSearcher optimization (#72)
- Branch dispatch patterns test O(1) first-byte optimization
- Prepares for future coregex v0.11.0 release

---

## [2026-01-12] - Capture Group Fix

### Changed
- **Updated coregex v0.10.2 → v0.10.3**
  - Fixed: FindStringSubmatch returned incorrect captures for `.+` patterns
  - Bug: `^(.+)-(\d+)$` on "hello-123" returned wrong `matches[1]`
  - Root cause: StateSplit in PikeVM passed captures without cloning

---

## [2026-01-07] - Version Pattern Hotfix

### Changed
- **Updated coregex v0.10.1 → v0.10.2**
  - Restored DigitPrefilter for version patterns (`\d+\.\d+\.\d+`)
  - v0.10.1 incorrectly chose ReverseInner with "." as inner literal
  - Performance restored: 8.2ms → 2.15ms (3.8x speedup)

---

## [2026-01-07] - Fat Teddy Release

### Changed
- **Updated coregex v0.9.5 → v0.10.0**
  - Fat Teddy 16-bucket SIMD (33-64 patterns, 9+ GB/s)
  - AVX2 assembly implementation
  - Pure Go scalar fallback

### Results
- `multi_literal` (12 patterns): 11.62 ms (Aho-Corasick)
- 5 patterns now faster than Rust regex

---

## [2026-01-05] - Initial Public Release

### Added
- Cross-language regex benchmark suite
- **Go stdlib** benchmarks
- **Go coregex** benchmarks  
- **Rust regex** benchmarks
- 10 real-world patterns
- GitHub Actions CI/CD
- Extreme benchmark mode (3000x speedup demo)

### Patterns
| Pattern | Description |
|---------|-------------|
| literal_alt | 4-literal alternation |
| multi_literal | 12-literal alternation |
| anchored | Start anchor |
| inner_literal | Inner literal search |
| suffix | Suffix matching |
| char_class | Character class |
| email | Email validation |
| uri | URI parsing |
| version | Version numbers |
| ip | IPv4 validation |

---

## Links

- **coregex**: https://github.com/coregx/coregex
- **Benchmark repo**: https://github.com/kolkov/regex-bench
- **Go regex issue**: https://github.com/golang/go/issues/26623
