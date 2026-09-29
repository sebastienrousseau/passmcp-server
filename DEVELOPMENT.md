<!-- SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com> -->
<!-- SPDX-License-Identifier: GPL-3.0-only -->

# Development

The single entry point for working on passmcp-server: toolchain, how to
reproduce every CI gate locally, and how a release is cut. If a gate fails
in CI and you cannot reproduce it from this file, that is a bug in this
file.

## Requirements

| Tool | Version | Why |
|---|---|---|
| Go | 1.26.8 or later, the `go` directive in `go.mod` | `GOTOOLCHAIN=auto` downloads it; CI tests on that version and on latest stable |
| make | any | Task runner for everything below |
| passmcp | the version in `CHANGELOG.md` | `passmcp_check` runs it; `go install satellion.com/passmcp/cmd/passmcp@vX.Y.Z` |

Optional, only for the gate that uses it: `golangci-lint` (`make lint`),
`markdownlint-cli2`, `codespell` and `lychee` (the Docs Lint workflow and
`pre-commit`), `zsh` and `fish` (`make completions` syntax-checks their
scripts when present), `curl` and `python3` (`make family`, `make lockstep`,
`make server-json`), `uvx` or `pipx` (`make server-json`), `jq` (the
dogfood job), `docker` (`make image`), `goreleaser` (`goreleaser check`).

The [dev container](.devcontainer/devcontainer.json) has Go at the
`go.mod` floor and golangci-lint at the version CI pins, and runs
`make build test` when it is created.

## Reproducing every CI gate

| CI job | Local command |
|---|---|
| Test (three OSes × two Go versions) | `make test` |
| Race & Shuffled Tests | `make test-race` |
| Coverage Gate (85% per package but `cmd/passmcp-server`) | `make coverage` |
| Lint, with the complexity ceilings | `gofmt -l .` and `make lint` |
| Vulnerability Scan | `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` |
| Repository Checks | `make smoke completions install-smoke versions server-json family digest lockstep` |
| passmcp evaluates passmcp-server | `passmcp check --stdio --rps 0 --no-color --output json -- ./build/passmcp-server --passmcp "$(command -v passmcp)"` after `make build` |
| Licence Headers | `make spdx-check` and `make name-guard` |
| Markdown & Spelling | `make readme-check`, `markdownlint-cli2 '**/*.md'` and `codespell` |
| Link Check | `lychee --offline --include-fragments '**/*.md'` |
| DCO check | `git log --format=%B origin/main.. \| grep Signed-off-by` |
| Manual (build, and on `main` deploy with `coverage.json`) | `mkdocs build --strict` after `pip install --require-hashes -r docs/requirements.txt`, and `make coverage-json` |
| OpenSSF Scorecard | not reproducible locally; runs on push to `main` and weekly, and publishes to [scorecard.dev](https://scorecard.dev/viewer/?uri=github.com/sebastienrousseau/passmcp-server) |

`make` with no target runs the gates that need no network, in the order
they fail fastest.

The dogfood job fails below a score of 90 or on any failing check other
than `supply.provenance`, which judges how the binary was built rather
than how the server behaves. A local build from a dirty tree fails it;
that is expected.

## Complexity

`make lint` enforces the portfolio's per-function ceilings through
golangci-lint: cyclomatic complexity 10 (`gocyclo`), cognitive complexity
15 (`gocognit`) and 60 lines (`funlen`). Every function is under them, so
there is no baseline of existing offenders; a new one fails the Lint job.
Halstead difficulty has no golangci-lint analyser and is not gated.

## Coverage

The gate is 85% statement coverage in every package with statements but
`cmd/passmcp-server`, whose `main` only wires signals and the process exit
around the tested `run`. `make coverage` writes `coverage.out`; ci.yml's
Coverage Gate checks each package against the threshold.

The README's coverage badge is a separate, whole-module number:
`make coverage-json` runs the tests with a cover profile and
`scripts/coveragebadge` turns it into a
[shields.io endpoint document](https://shields.io/badges/endpoint-badge),
`build/coverage.json`. It is brightgreen from 90%, green from 85%, yellow
from 70% and red below, truncated rather than rounded so the badge never
shows the gate as met early. The Manual workflow publishes it on every
push to `main` at
<https://sebastienrousseau.com/passmcp-server/coverage.json>, which the
badge reads.

## Install contract

`GNUmakefile` includes `Makefile` and adds `install`, `uninstall` and
`install-smoke`. `make install` builds the binary and its completions and
installs them under `PREFIX` (default `/usr/local`), staged under
`DESTDIR`: the binary in `bin`, completions in
`share/bash-completion/completions`, `share/zsh/site-functions` and
`share/fish/vendor_completions.d`, and the README, changelog, licence and
security policy in `share/doc/passmcp-server`. The binary generates no
manual page, so none is installed. `make install-smoke` stages an install
in a temporary directory, runs the installed binary, uninstalls, and fails
if anything is missing or left behind.

## Test layout

`internal/server/server_test.go` drives the server through its stdio loop
with a fake runner: the handshake, every tool, every refusal. The
allowlist cases are there. `internal/runner/runner_test.go` stands a shell
script in for passmcp and asserts the arguments and environment every run
gets; it is `//go:build !windows`, so on Windows that package reports no
test files. `cmd/passmcp-server/main_test.go` covers `run`; `main` itself is
exempt from the coverage gate. `scripts/coveragebadge/main_test.go` covers
the badge document: its colour bands, and every malformed profile.

## Generated artefacts

None are committed. Release archives, checksums and the image are built by
goreleaser into `dist/`; `make build`, `make smoke`, `make completions`,
`make coverage-json` and `make image` write to `build/`. Both are ignored.

## Release model

The version is passmcp's latest release, exactly, and a release is cut after
passmcp's, on a `feat/vX.Y.Z` branch:

1. Date the `## [X.Y.Z]` section the release branch opened in
   `CHANGELOG.md` as `## [X.Y.Z] — date` (or move the `## [Unreleased]`
   entries under one), and update the version in the README's install
   lines and ecosystem sentence, in `server.json` (both `version` and the
   image tag) and in `CITATION.cff` (`version` and `date-released`), and
   move `go.mod` to passmcp-reporting's `vX.Y.Z`, which is tagged first.
2. Point the `Dockerfile`'s `FROM` at passmcp X.Y.Z's image digest
   (`docker buildx imagetools inspect ghcr.io/sebastienrousseau/passmcp:X.Y.Z`)
   and set `PASSMCP_VERSION=X.Y.Z`.
3. `make lockstep` and `scripts/verify-release-versions.sh vX.Y.Z`
   (`make versions` runs it for the newest CHANGELOG release).
4. `goreleaser check`, and optionally the Release workflow's dry run.
5. Push a signed annotated tag `vX.Y.Z` with the message
   `passmcp-server vX.Y.Z`. The Release workflow builds the binaries and the
   image, signs them, and attests the checksums.
6. Read the tag, the release page and the image's labels back before
   calling it done, then publish the registry listing:
   [docs/publishing.md](docs/publishing.md).

Steps 1 and 2 are what `.github/workflows/sync.yml` does on passmcp's
release dispatch: it opens a pull request that makes both edits. With a
`SYNC_TOKEN` secret, a fine-grained token with `pull-requests: write` on
this repository, that pull request's checks start on their own; without
it, close and reopen the pull request to start them, because one opened
with `GITHUB_TOKEN` triggers no workflows.

## Conventions

- Stdout is the MCP transport; nothing else is written there.
- Every exported identifier is documented.
- Anything a tool argument or a passmcp report carries is untrusted input.
