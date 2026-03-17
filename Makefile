.PHONY: all build run clean generate extreme extreme-3000x langarena

INPUT = input/data.txt
INPUT_NOMATCH = input/no-match-data.txt
INPUT_NODIGITS = input/no-digits-data.txt
INPUT_LANGARENA = input/langarena-data.txt

all: generate build run

generate:
	@echo "Generating input data..."
	@go run scripts/generate-input.go

generate-nomatch:
	@echo "Generating no-match input data..."
	@go run scripts/generate-no-match-input.go

generate-nodigits:
	@echo "Generating no-digits input data (for 3000x speedup)..."
	@go run scripts/generate-3000x-input.go

build: build-go-stdlib build-go-coregex

build-go-stdlib:
	@echo "Building go-stdlib..."
	@cd go-stdlib && go build -ldflags "-s -w" -o ../bin/go-stdlib.exe .

build-go-coregex:
	@echo "Building go-coregex..."
	@cd go-coregex && go mod tidy && go build -ldflags "-s -w" -o ../bin/go-coregex.exe .

build-extreme: build-go-stdlib-extreme build-go-coregex-extreme

build-go-stdlib-extreme:
	@echo "Building go-stdlib-extreme..."
	@cd go-stdlib-extreme && go build -ldflags "-s -w" -o ../bin/go-stdlib-extreme.exe .

build-go-coregex-extreme:
	@echo "Building go-coregex-extreme..."
	@cd go-coregex-extreme && go mod tidy && go build -ldflags "-s -w" -o ../bin/go-coregex-extreme.exe .

run: $(INPUT)
	@echo ""
	@echo "==============================================="
	@./bin/go-stdlib.exe $(INPUT)
	@echo ""
	@echo "==============================================="
	@./bin/go-coregex.exe $(INPUT)
	@echo ""

# EXTREME benchmarks: demonstrate 1000-3000x speedup on no-match data
# This is the worst case for stdlib (must scan entire file)
# and best case for coregex prefilters (skip quickly)
extreme: generate-nomatch build-extreme run-extreme

run-extreme: $(INPUT_NOMATCH)
	@echo ""
	@echo "==============================================================================="
	@echo "EXTREME BENCHMARKS: No-match data (worst case for stdlib, best for prefilters)"
	@echo "==============================================================================="
	@echo ""
	@./bin/go-stdlib-extreme.exe $(INPUT_NOMATCH)
	@echo ""
	@echo "==============================================================================="
	@./bin/go-coregex-extreme.exe $(INPUT_NOMATCH)
	@echo ""

# EXTREME-3000x: demonstrate 800-1000x speedup on no-digits data
# Uses data with NO DIGITS for maximum DigitPrefilter advantage
# (3000x achieved in go test -bench with 1MB input)
extreme-3000x: generate-nodigits build-extreme run-extreme-3000x

run-extreme-3000x: $(INPUT_NODIGITS)
	@echo ""
	@echo "==============================================================================="
	@echo "EXTREME-3000x: No-digits data (800-1000x on 6MB, 3000x on 1MB in go test)"
	@echo "==============================================================================="
	@echo ""
	@./bin/go-stdlib-extreme.exe $(INPUT_NODIGITS)
	@echo ""
	@echo "==============================================================================="
	@./bin/go-coregex-extreme.exe $(INPUT_NODIGITS)
	@echo ""

# LANGARENA benchmarks: 13 real-world LogParser patterns from kostya/LangArena
langarena: generate-langarena build-langarena run-langarena

generate-langarena:
	@echo "Generating LangArena input data..."
	@go run scripts/generate-langarena-input.go

build-langarena: build-go-stdlib-langarena build-go-coregex-langarena

build-go-stdlib-langarena:
	@echo "Building go-stdlib-langarena..."
	@cd go-stdlib-langarena && go build -ldflags "-s -w" -o ../bin/go-stdlib-langarena.exe .

build-go-coregex-langarena:
	@echo "Building go-coregex-langarena..."
	@cd go-coregex-langarena && go mod tidy && go build -ldflags "-s -w" -o ../bin/go-coregex-langarena.exe .

run-langarena: $(INPUT_LANGARENA)
	@echo ""
	@echo "==============================================================================="
	@echo "LANGARENA: 13 LogParser patterns (https://kostya.github.io/LangArena/)"
	@echo "==============================================================================="
	@echo ""
	@./bin/go-stdlib-langarena.exe $(INPUT_LANGARENA)
	@echo ""
	@echo "==============================================================================="
	@./bin/go-coregex-langarena.exe $(INPUT_LANGARENA)
	@echo ""

clean:
	@rm -rf bin/*.exe input/*.txt

$(INPUT):
	@$(MAKE) generate

$(INPUT_NOMATCH):
	@$(MAKE) generate-nomatch

$(INPUT_NODIGITS):
	@$(MAKE) generate-nodigits

$(INPUT_LANGARENA):
	@$(MAKE) generate-langarena
