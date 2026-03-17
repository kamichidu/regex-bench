use regex::Regex;
use std::env;
use std::fs;
use std::time::Instant;

// LangArena LogParser patterns from https://kostya.github.io/LangArena/
// 13 real-world patterns for Apache log parsing.
struct Pattern {
    name: &'static str,
    pattern: &'static str,
}

const PATTERNS: &[Pattern] = &[
    Pattern { name: "errors", pattern: r#" [5][0-9]{2} | [4][0-9]{2} "#  },
    Pattern { name: "bots", pattern: r"(?i)bot|crawler|scanner|spider|indexing|crawl|robot|spider" },
    Pattern { name: "suspicious", pattern: r"(?i)etc/passwd|wp-admin|\.\./"},
    Pattern { name: "ips", pattern: r"\d+\.\d+\.\d+\.35" },
    Pattern { name: "api_calls", pattern: r#"/api/[^ "]+"# },
    Pattern { name: "post_requests", pattern: r"POST [^ ]* HTTP" },
    Pattern { name: "auth_attempts", pattern: r"(?i)/login|/signin" },
    Pattern { name: "methods", pattern: r"(?i)get|post|put" },
    Pattern { name: "emails", pattern: r"[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}" },
    Pattern { name: "passwords", pattern: r#"password=[^&\s"]+"# },
    Pattern { name: "tokens", pattern: r#"token=[^&\s"]+|api[_-]?key=[^&\s"]+"# },
    Pattern { name: "sessions", pattern: r#"session[_-]?id=[^&\s"]+"# },
    Pattern { name: "peak_hours", pattern: r"\[\d+/\w+/\d+:1[3-7]:\d+:\d+ [+\-]\d+\]" },
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
        eprintln!("Usage: langarena-benchmark <input-file>");
        std::process::exit(1);
    }

    let data = fs::read_to_string(&args[1]).expect("Failed to read file");
    let size_mb = data.len() as f64 / 1024.0 / 1024.0;

    println!("Rust regex LangArena LogParser (input: {:.2} MB)", size_mb);
    println!("{:<15} {:>10} {:>10} {:>6}", "pattern", "compile", "search", "matches");
    println!("─────────────────────────────────────────────────");

    for p in PATTERNS {
        measure(&data, p);
    }
}
