#!/bin/bash
set -e -u -o pipefail

# Default values
__engines=()
__scenarios=()

# Argument parsing with getopts
while getopts "e:s:h" __opt; do
    case "${__opt}" in
        e)
            __engines+=("${OPTARG}")
            ;;
        s)
            __scenarios+=("${OPTARG}")
            ;;
        h)
            echo "Usage: ${0} [-e engine] [-s scenario]"
            echo "Options:"
            echo "  -e    Select engine (stdlib, regexp-re, coregex, re2-wasm, re2-cgo, pcre2, hyperscan)"
            echo "  -s    Select scenario (standard, dna, extreme, langarena)"
            echo "  -h    Show this help"
            exit 0
            ;;
        *)
            exit 1
            ;;
    esac
done
shift $(("${OPTIND}" - 1))

# Fallback to all if not specified
if [[ "${#__engines[@]}" -eq 0 ]]; then
    __engines=("stdlib" "regexp-re" "coregex" "re2-wasm" "re2-cgo" "pcre2" "hyperscan")
fi
if [[ "${#__scenarios[@]}" -eq 0 ]]; then
    __scenarios=("standard" "dna" "extreme" "langarena")
fi

mkdir -p results/standard results/dna results/extreme results/langarena

for __scenario in "${__scenarios[@]}"; do
    echo "Processing scenario: ${__scenario}"

    case "${__scenario}" in
        "dna")
            __input="input/dna-data.txt"
            ;;
        "extreme")
            __input="input/no-match-data.txt"
            ;;
        "langarena")
            __input="input/langarena-data.txt"
            ;;
        *)
            __input="input/data.txt"
            ;;
    esac

    for __engine in "${__engines[@]}"; do
        __suffix=""
        [[ "${__scenario}" != "standard" ]] && __suffix="-${__scenario}"

        __bin="bin/go-${__engine}${__suffix}.exe"

        if [[ -f "${__bin}" ]]; then
            echo "  Running ${__engine}..."
            # Use tee to both display and save results
            "${__bin}" "${__input}" 2>/dev/null | tee "results/${__scenario}/${__engine}.txt"
        else
            echo "  Binary ${__bin} not found, skipping."
        fi
    done
done
