# SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com>
# SPDX-License-Identifier: GPL-3.0-only

.PHONY: all build test test-race coverage coverage-json vet lint format spdx-check smoke digest \
        readme-check server-json image lockstep family completions help name-guard versions

# Every gate CI runs that needs no network, in the order the cheap ones fail
# first.
all: format vet lint spdx-check test smoke

build:
	CGO_ENABLED=0 go build -trimpath -o build/passmcp-server ./cmd/passmcp-server

test:
	go test ./... -cover

# Shell completions, generated from the flag set by the binary itself, into
# build/completions. bash is syntax-checked here; zsh and fish when present.
completions:
	mkdir -p build/completions
	go run ./cmd/passmcp-server --completion bash > build/completions/passmcp-server.bash
	go run ./cmd/passmcp-server --completion zsh > build/completions/_passmcp-server
	go run ./cmd/passmcp-server --completion fish > build/completions/passmcp-server.fish
	bash -n build/completions/passmcp-server.bash
	if command -v zsh >/dev/null; then zsh -n build/completions/_passmcp-server; fi
	if command -v fish >/dev/null; then fish -n build/completions/passmcp-server.fish; fi

test-race:
	go test -race -shuffle=on -count=1 ./...

# The gate is 85% statement coverage in every package with statements,
# except cmd/passmcp-server; ci.yml says why.
coverage:
	@mkdir -p build
	go test -count=1 -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1

# The shields.io endpoint document behind the README's coverage badge:
# statement coverage across the module, as CI measured it. The Manual
# workflow publishes it with GitHub Pages as coverage.json.
coverage-json: coverage
	go run ./scripts/coveragebadge -profile coverage.out > build/coverage.json
	@cat build/coverage.json

vet:
	go vet ./...

lint:
	golangci-lint run ./...

format:
	gofmt -l -w .

# The README follows the portfolio template: headings in order, no
# unresolved {{VARIABLES}} (AGENTS.md §7.3).
readme-check:
	scripts/readme-check.sh

spdx-check:
	go run ./scripts/spdx_sweep.go

# The server answers the handshake and lists its three tools over stdio.
smoke:
	@mkdir -p build
	@printf '%s\n' \
	  '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"make","version":"0"}}}' \
	  '{"jsonrpc":"2.0","method":"notifications/initialized"}' \
	  '{"jsonrpc":"2.0","id":2,"method":"tools/list"}' \
	  | go run ./cmd/passmcp-server > build/smoke.out
	@grep -q '"protocolVersion":"2025-11-25"' build/smoke.out
	@for t in passmcp_check passmcp_verify_attestation passmcp_version; do grep -q "\"name\":\"$$t\"" build/smoke.out || { echo "smoke: $$t missing from tools/list"; exit 1; }; done
	@echo "smoke: initialize and tools/list answered with all three tools"

# server.json against the registry schema it names. Needs the network and
# check-jsonschema (through uvx, or pipx on a CI runner).
CHECK_JSONSCHEMA ?= $(shell command -v uvx >/dev/null 2>&1 && echo 'uvx --from check-jsonschema check-jsonschema' || echo 'pipx run check-jsonschema')
server-json:
	@mkdir -p build
	curl -fsSL "$$(python3 -c 'import json; print(json.load(open("server.json"))["$$schema"])')" -o build/server.schema.json
	$(CHECK_JSONSCHEMA) --schemafile build/server.schema.json server.json

# The container image for this machine's architecture, built the way
# goreleaser lays out its context, without goreleaser.
IMAGE ?= passmcp-server:dev
ARCH ?= $(shell go env GOARCH)
image:
	CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH) go build -trimpath -ldflags "-s -w -X main.Version=dev" \
	  -o build/image/linux/$(ARCH)/passmcp-server ./cmd/passmcp-server
	cp Dockerfile build/image/Dockerfile
	docker build --platform linux/$(ARCH) -t $(IMAGE) build/image

# This repository carries passmcp's version. See docs/ecosystem.md in passmcp.
# The Dockerfile's FROM digest is passmcp's image for its PASSMCP_VERSION.
digest:
	scripts/verify-digest.sh

lockstep:
	scripts/lockstep.sh

# The family manifest in passmcp is the single source of what this repository is.
family:
	scripts/family.sh

# Every file that names the version names the newest CHANGELOG release,
# and the passmcp-reporting module is required at that same release.
versions:
	scripts/verify-release-versions.sh "v$$(grep -Eo '^## \[[0-9]+\.[0-9]+\.[0-9]+\]' CHANGELOG.md | head -1 | tr -d '#[] ')"

help:
	@printf '%s\n' "targets: all build test test-race coverage coverage-json vet lint format spdx-check smoke" \
	  "         completions readme-check server-json image digest lockstep family versions name-guard" \
	  "GNUmakefile: install uninstall install-smoke (PREFIX, DESTDIR)"

# A retired product name may not appear anywhere in the tree
# (scripts/name-guard.sh).
name-guard:
	./scripts/name-guard.sh
