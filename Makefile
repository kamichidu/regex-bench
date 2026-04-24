.PHONY: all build build-all run report clean generate extreme extreme-3000x langarena

INPUT = input/data.txt
INPUT_NOMATCH = input/no-match-data.txt
INPUT_NODIGITS = input/no-digits-data.txt
INPUT_LANGARENA = input/langarena-data.txt
INPUT_DNA = input/dna-data.txt

ENGINES = stdlib regexp-re coregex re2-wasm re2-cgo pcre2 hyperscan

all: generate build run

generate:
	@echo "Generating input data..."
	@go run scripts/generate-input.go
	@go run scripts/generate-no-match-input.go
	@go run scripts/generate-langarena-input.go
	@go run scripts/generate-dna-input.go

build: $(foreach eng,$(ENGINES),build-go-$(eng))

# Generic build rules
build-go-stdlib:
	@echo "Building go-stdlib..."
	@cd go-stdlib && go build -ldflags "-s -w" -o ../bin/go-stdlib.exe .

build-go-regexp-re:
	@echo "Building go-regexp-re..."
	@cd go-regexp-re && go build -ldflags "-s -w" -o ../bin/go-regexp-re.exe .

build-go-coregex:
	@echo "Building go-coregex..."
	@cd go-coregex && go mod tidy && go build -ldflags "-s -w" -o ../bin/go-coregex.exe .

build-go-re2-wasm:
	@echo "Building go-re2-wasm..."
	@cd go-re2-wasm && go build -ldflags "-s -w" -o ../bin/go-re2-wasm.exe .

build-go-re2-cgo:
	@echo "Building go-re2-cgo..."
	@cd go-re2-cgo && go build -tags re2_cgo -ldflags "-s -w" -o ../bin/go-re2-cgo.exe .

build-go-pcre2:
	@echo "Building go-pcre2..."
	@cd go-pcre2 && go build -ldflags "-s -w" -o ../bin/go-pcre2.exe .

build-go-hyperscan:
	@echo "Building go-hyperscan..."
	@cd go-hyperscan && go build -ldflags "-s -w" -o ../bin/go-hyperscan.exe .

report: generate build
	@chmod +x scripts/run-all.sh
	@./scripts/run-all.sh
	@go run scripts/generate-report.go > REPORT.md
	@echo "Report generated in REPORT.md"

run: $(INPUT)
	@echo ""
	@echo "==============================================="
	@for eng in $(ENGINES); do ./bin/go-$$eng.exe $(INPUT); done
	@echo ""

clean:
	@rm -rf bin/*.exe input/*.txt results/ REPORT.md

$(INPUT):
	@$(MAKE) generate
