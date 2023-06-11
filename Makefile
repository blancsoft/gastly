.PHONY: serve build-wasm build test lint clean help
SHELL := '/bin/bash'
.DEFAULT_GOAL := help



serve: build-server
	@go run $(PWA_MAIN) -serve

build-wasm: clean
	@GOARCH=wasm GOOS=js go build -o "./pwa/public/gastly.wasm" "./cmd/wasm/main.go"
	@cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js "./pwa/public/wasm_exec.js"

test:
	@go test ./...
	@GOOS=js GOARCH=wasm go test -v -exec="$(shell go env GOROOT)/misc/wasm/go_js_wasm_exec" ./...

lint: ## lint go files in current directory
	@go run honnef.co/go/tools/cmd/staticcheck@2023.1 ./...

build: ## build in snapshot mode
	@goreleaser release --snapshot --rm-dist

clean: ## remove build artefacts
	@go clean
	@rm -rf "./pwa/public/gastly.wasm" "./pwa/public/wasm_exec.js"


# got from :https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
# but disallow underscore in command names as we want some private to have format "_command-name"
help:  ## print command reference
	@printf "  Welcome to \033[36mGastly\033[0m command reference.\n"
	@printf "  If you wish to contribute, please follow guide at top section of \033[36mMakefile\033[0m.\n\n"
	@printf "  Usage:\n    \033[36mmake <target> [..arguments]\033[0m\n\n  Targets:\n"
	@grep -E '^[a-zA-Z-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-20s\033[0m %s\n", $$1, $$2}'
