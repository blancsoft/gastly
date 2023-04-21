.PHONY: serve build-server build test lint clean setup help
SHELL := '/bin/bash'
.DEFAULT_GOAL := help

OUT_DIR="dist/"
PWA_MAIN="./cmd/pwa/main.go"

serve: build-server
	@go run $(PWA_MAIN) -serve

build-server: clean
	GOARCH=wasm GOOS=js go build -o web/app.wasm $(PWA_MAIN)
	@go run $(PWA_MAIN) -build
	@cp -r web/ $(OUT_DIR)

test:
	@go test ./...

lint: ## lint go files in current directory
	@printf "=========Linting=========\n\n"
	@golangci-lint run

build: ## build in snapshot mode
	@printf "=========Building binaries in snapshot mode=========\n\n"
	@goreleaser release --snapshot --rm-dist

clean: ## remove build artefacts
	@printf "=========Cleaning up=========\n\n"
	@rm -rf output/ target/ web/app.wasm $(OUT_DIR)
	@go clean

setup: ## install dev tools
	@printf "=========Getting dev tools=========\n\n"
	@go install github.com/goreleaser/goreleaser@latest
	@printf "To install \033[36mgolangci-lint\033[0m, see https://golangci-lint.run/usage/install/#local-installation\n\n"

# got from :https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
# but disallow underscore in command names as we want some private to have format "_command-name"
help:  ## print command reference
	@printf "  Welcome to \033[36mGastly\033[0m command reference.\n"
	@printf "  If you wish to contribute, please follow guide at top section of \033[36mMakefile\033[0m.\n\n"
	@printf "  Usage:\n    \033[36mmake <target> [..arguments]\033[0m\n\n  Targets:\n"
	@grep -E '^[a-zA-Z-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-20s\033[0m %s\n", $$1, $$2}'
