#!/bin/bash
# Generate comparison table from benchmark results
# Usage: ./scripts/generate-summary.sh stdlib.txt coregex.txt rust.txt

STDLIB=$1
COREGEX=$2
RUST=$3

# Parse benchmark output into associative arrays
declare -A stdlib_ms stdlib_matches
declare -A coregex_ms coregex_matches
declare -A rust_ms rust_matches

parse_results() {
    local file=$1
    local -n ms_arr=$2
    local -n match_arr=$3

    while IFS= read -r line; do
        # Skip header lines
        [[ "$line" =~ ^Go|^Rust|^─ ]] && continue
        [[ -z "$line" ]] && continue

        # Parse: "pattern_name    123.45 ms    1234 matches"
        if [[ "$line" =~ ^([a-z_]+)[[:space:]]+([0-9.]+)[[:space:]]*ms[[:space:]]+([0-9]+)[[:space:]]*matches ]]; then
            name="${BASH_REMATCH[1]}"
            ms="${BASH_REMATCH[2]}"
            matches="${BASH_REMATCH[3]}"
            ms_arr[$name]=$ms
            match_arr[$name]=$matches
        fi
    done < "$file"
}

parse_results "$STDLIB" stdlib_ms stdlib_matches
parse_results "$COREGEX" coregex_ms coregex_matches
parse_results "$RUST" rust_ms rust_matches

# Calculate speedup (returns formatted string like "123x" or "1.5x")
calc_speedup() {
    local base=$1
    local fast=$2

    if (( $(echo "$fast == 0" | bc -l) )); then
        echo "—"
        return
    fi

    local ratio=$(echo "scale=2; $base / $fast" | bc -l)

    if (( $(echo "$ratio >= 10" | bc -l) )); then
        printf "**%.0fx**" "$ratio"
    elif (( $(echo "$ratio >= 2" | bc -l) )); then
        printf "**%.1fx**" "$ratio"
    elif (( $(echo "$ratio >= 1.1" | bc -l) )); then
        printf "%.1fx" "$ratio"
    else
        echo "~1x"
    fi
}

# Determine winner
get_winner() {
    local coregex=$1
    local rust=$2

    if (( $(echo "$coregex < $rust * 0.95" | bc -l) )); then
        echo "coregex"
    elif (( $(echo "$rust < $coregex * 0.95" | bc -l) )); then
        echo "Rust"
    else
        echo "—"
    fi
}

# Bold the winner's time
format_time() {
    local ms=$1
    local is_winner=$2

    if [[ "$is_winner" == "true" ]]; then
        printf "**%.2f ms**" "$ms"
    else
        printf "%.2f ms" "$ms"
    fi
}

# Generate table header
echo "## Comparison Table"
echo ""
echo "| Pattern | Go stdlib | Go coregex | Rust regex | vs stdlib | vs Rust | Winner |"
echo "|---------|-----------|------------|------------|-----------|---------|--------|"

# Get all patterns (sorted by stdlib time descending for impact)
patterns=$(echo "${!stdlib_ms[@]}" | tr ' ' '\n' | sort)

for pattern in $patterns; do
    std="${stdlib_ms[$pattern]:-0}"
    cor="${coregex_ms[$pattern]:-0}"
    rus="${rust_ms[$pattern]:-0}"

    # Skip if missing data
    [[ -z "$std" || "$std" == "0" ]] && continue

    # Calculate speedups
    vs_stdlib=$(calc_speedup "$std" "$cor")

    # vs Rust comparison
    if (( $(echo "$cor < $rus" | bc -l) )); then
        ratio=$(echo "scale=1; $rus / $cor" | bc -l)
        vs_rust="**${ratio}x faster**"
    elif (( $(echo "$rus < $cor" | bc -l) )); then
        ratio=$(echo "scale=1; $cor / $rus" | bc -l)
        vs_rust="${ratio}x slower"
    else
        vs_rust="~1x"
    fi

    # Determine winner
    winner=$(get_winner "$cor" "$rus")

    # Format times (bold winner)
    if [[ "$winner" == "coregex" ]]; then
        cor_fmt="**${cor} ms**"
        rus_fmt="${rus} ms"
    elif [[ "$winner" == "Rust" ]]; then
        cor_fmt="${cor} ms"
        rus_fmt="**${rus} ms**"
    else
        cor_fmt="${cor} ms"
        rus_fmt="${rus} ms"
    fi

    echo "| $pattern | ${std} ms | $cor_fmt | $rus_fmt | $vs_stdlib | $vs_rust | $winner |"
done

echo ""
echo "**Legend**: Bold = fastest, vs stdlib = coregex speedup over stdlib, vs Rust = coregex comparison to Rust"
