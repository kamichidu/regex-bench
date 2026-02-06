use regex::Regex;
use std::env;
use std::fs;
use std::time::Instant;

struct Pattern {
    name: &'static str,
    pattern: &'static str,
}

const PATTERNS: &[Pattern] = &[
    Pattern { name: "literal_alt", pattern: r"error|warning|fatal|critical" },
    Pattern { name: "multi_literal", pattern: r"apple|banana|cherry|date|elderberry|fig|grape|honeydew|kiwi|lemon|mango|orange" },
    Pattern { name: "anchored", pattern: r"^HTTP/[12]\.[01]" },
    Pattern { name: "inner_literal", pattern: r".*@example\.com" },
    Pattern { name: "suffix", pattern: r".*\.(txt|log|md)" },
    Pattern { name: "char_class", pattern: r"[\w]+" },
    Pattern { name: "email", pattern: r"[\w.+-]+@[\w.-]+\.[\w.-]+" },
    Pattern { name: "uri", pattern: r"[\w]+://[^/\s?#]+[^\s?#]+(?:\?[^\s#]*)?(?:#[^\s]*)?" },
    Pattern { name: "version", pattern: r"\d+\.\d+\.\d+" },
    Pattern { name: "ip", pattern: r"(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])" },
    // Composite patterns (concatenated char classes)
    Pattern { name: "alpha_digit", pattern: r"[a-zA-Z]+\d+" },
    Pattern { name: "word_digit", pattern: r"\w+[0-9]+" },
    // Branch dispatch patterns (anchored alternations) - multiline mode for log parsing
    Pattern { name: "http_methods", pattern: r"(?m)^(GET|POST|PUT|DELETE|PATCH)" },
    // Issue #79: Anchored patterns with wildcards (single-string URL matching)
    Pattern { name: "anchored_php", pattern: r"^/.*[\w-]+\.php" },
    // Issue #97: Multiline suffix patterns - UseMultilineReverseSuffix strategy (v0.11.1)
    Pattern { name: "multiline_php", pattern: r"(?m)^/.*\.php" },
    // Issue #105: Word quantifiers in capture groups - was 7M x slower before v0.11.5
    Pattern { name: "word_repeat", pattern: r"(\w{2,8})+" },
];

fn measure(data: &str, p: &Pattern) {
    // Compile (measured separately)
    let compile_start = Instant::now();
    let re = Regex::new(p.pattern).expect("Invalid regex");
    let compile_elapsed = compile_start.elapsed();

    // Search only
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
        eprintln!("Usage: benchmark <input-file>");
        std::process::exit(1);
    }

    let data = fs::read_to_string(&args[1]).expect("Failed to read file");
    let size_mb = data.len() as f64 / 1024.0 / 1024.0;

    println!("Rust regex (input: {:.2} MB)", size_mb);
    println!("{:<15} {:>10} {:>10} {:>6}", "pattern", "compile", "search", "matches");
    println!("─────────────────────────────────────────────────");

    for p in PATTERNS {
        measure(&data, p);
    }
}
