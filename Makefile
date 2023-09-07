.PHONY: build test coverage lint clean help
SHELL := '/bin/bash'
.DEFAULT_GOAL := help


build: clean ## build WASM module
	@cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js ./src/assets/wasm_exec.js
	@GOARCH=wasm GOOS=js go build -ldflags="-s -w" -o ./src/assets/gastly.wasm ./main.go

test: clean ## run lib tests
	@go test -race -v -covermode=atomic -coverprofile=coverage.out ./lib/...

test-wasm: clean ## run wasm tests
	@GOOS=js GOARCH=wasm go test -v -covermode=atomic -coverprofile=coverage.out \
		-v -exec="$(shell go env GOROOT)/misc/wasm/go_js_wasm_exec" ./lib/wasm/...

coverage: ## check test coverage
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o coverage.html

lint: ## lint go files in current directory
	@go run honnef.co/go/tools/cmd/staticcheck@2023.1 github.com/chumaumenze/gastly/lib/...

clean: ## remove build artefacts
	@go clean
	@rm -f coverage.out coverage.html
	@rm -rf ./src/assets/gastly.wasm* ./src/assets/wasm_exec.js


# got from :https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
# but disallow underscore in command names as we want some private to have format "_command-name"
help:  ## print command reference
	@printf "  Welcome to \033[36mGastly\033[0m command reference.\n"
	@printf "  If you wish to contribute, please follow guide at top section of \033[36mMakefile\033[0m.\n\n"
	@printf "  Usage:\n    \033[36mmake <target> [..arguments]\033[0m\n\n  Targets:\n"
	@grep -E '^[a-zA-Z-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-20s\033[0m %s\n", $$1, $$2}'
