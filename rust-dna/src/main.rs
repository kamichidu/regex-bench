use regex::Regex;
use std::env;
use std::fs;
use std::time::Instant;

// regexdna patterns from The Computer Language Benchmarks Game
// DNA reverse complement matching
struct Pattern {
    name: &'static str,
    pattern: &'static str,
}

const PATTERNS: &[Pattern] = &[
    Pattern { name: "dna_1", pattern: r"agggtaaa|tttaccct" },
    Pattern { name: "dna_2", pattern: r"[cgt]gggtaaa|tttaccc[acg]" },
    Pattern { name: "dna_3", pattern: r"a[act]ggtaaa|tttacc[agt]t" },
    Pattern { name: "dna_4", pattern: r"ag[act]gtaaa|tttac[agt]ct" },
    Pattern { name: "dna_5", pattern: r"agg[act]taaa|ttta[agt]cct" },
    Pattern { name: "dna_6", pattern: r"aggg[acg]aaa|ttt[cgt]ccct" },
    Pattern { name: "dna_7", pattern: r"agggt[cgt]aa|tt[acg]accct" },
    Pattern { name: "dna_8", pattern: r"agggta[cgt]a|t[acg]taccct" },
    Pattern { name: "dna_9", pattern: r"agggtaa[cgt]|[acg]ttaccct" },
];

fn measure(data: &str, p: &Pattern) {
    let compile_start = Instant::now();
    let re = Regex::new(p.pattern).expect("Invalid regex");
    let compile_elapsed = compile_start.elapsed();

    let search_start = Instant::now();
    let count = re.find_iter(data).count();
    let search_elapsed = search_start.elapsed();

    let compile_ms = compile_elapsed.as_secs_f64() * 1000.0;
    let search_ms = search_elapsed.as_secs_f64() * 1000.0;

    println!("{:<15} {:>10.2} {:>10.2} {:>6}", p.name, compile_ms, search_ms, count);
}

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() != 2 {
        eprintln!("Usage: dna-benchmark <input-file>");
        std::process::exit(1);
    }

    let data = fs::read_to_string(&args[1]).expect("Failed to read file");
    let size_mb = data.len() as f64 / 1024.0 / 1024.0;

    println!("Rust regex regexdna (input: {:.2} MB)", size_mb);
    println!("{:<15} {:>10} {:>10} {:>6}", "pattern", "compile", "search", "matches");
    println!("─────────────────────────────────────────────────");

    for p in PATTERNS {
        measure(&data, p);
    }
}
