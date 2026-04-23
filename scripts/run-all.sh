#!/bin/bash
# scripts/run_all.sh
mkdir -p results/standard results/dna results/extreme results/langarena

SCENARIOS=("standard" "dna" "extreme" "langarena")
ENGINES=("stdlib" "regexp-re" "coregex" "re2-wasm" "re2-cgo" "pcre2" "hyperscan")

for s in "${SCENARIOS[@]}"; do
    echo "Processing scenario: $s"
    for e in "${ENGINES[@]}"; do
        suffix=""
        [ "$s" != "standard" ] && suffix="-$s"

        # Binary name construction
        bin="bin/go-${e}${suffix}.exe"

        # Input selection
        input="input/data.txt"
        [ "$s" == "dna" ] && input="input/dna-data.txt"
        [ "$s" == "extreme" ] && input="input/no-match-data.txt"
        [ "$s" == "langarena" ] && input="input/langarena-data.txt"

        if [ -f "$bin" ]; then
            echo "  Running $e..."
            $bin "$input" > "results/$s/$e.txt" 2>/dev/null
        else
            echo "  Binary $bin not found, skipping."
        fi
    done
done
