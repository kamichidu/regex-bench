Add **regexdna** benchmark suite — 9 patterns from [The Computer Language Benchmarks Game](https://benchmarksgame-team.pages.debian.net/benchmarksgame/) (DNA reverse complement matching).

### Patterns

```
agggtaaa|tttaccct
[cgt]gggtaaa|tttaccc[acg]
a[act]ggtaaa|tttacc[agt]t
ag[act]gtaaa|tttac[agt]ct
agg[act]taaa|ttta[agt]cct
aggg[acg]aaa|ttt[cgt]ccct
agggt[cgt]aa|tt[acg]accct
agggta[cgt]a|t[acg]taccct
agggtaa[cgt]|[acg]ttaccct
```

### Motivation

- @kostya found coregex Issue #116 via regexdna — classic bioinformatics benchmark
- Alternation + char class patterns exercise different strategies (Teddy, UseBoth, UseDFA)
- Local results (100KB DNA, coregex v0.12.2): Teddy 28x, UseBoth 5-14x, UseDFA ~1x
- Interesting comparison with Rust on DFA-heavy patterns

### Implementation

Follow extreme benchmark pattern:
- `go-stdlib-dna/`, `go-coregex-dna/`, `rust-dna/`
- `scripts/generate-dna-input.go` (100K-1M nucleotides, seed=42)
- `.github/workflows/dna-benchmark.yml`

### Future

When 4+ suites exist, refactor to suite-based structure:
```
engines/{go-stdlib,go-coregex,rust}/
suites/{standard,extreme,regexdna}.json
```
