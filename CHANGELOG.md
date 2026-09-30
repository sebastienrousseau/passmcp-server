<!-- SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com> -->
<!-- SPDX-License-Identifier: GPL-3.0-only -->

# Changelog

All notable changes to passmcp-server are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and versions are
[Semantic Versioning](https://semver.org/) shaped.

**This repository carries passmcp's version.** It is in lockstep with
[passmcp](https://github.com/sebastienrousseau/passmcp): it wraps passmcp's
release — the container image is built on passmcp's image for the same
version — so its version is passmcp's latest release, exactly, and a release
here with nothing in it is the version rule working.

## [Unreleased]

## [0.0.3] — 2026-09-30

### Changed

- **In lockstep with passmcp 0.0.3.** The image builds on `ghcr.io/sebastienrousseau/passmcp@sha256:88668e403cb4fe2e63032e9d69550acb59a449701d1f2ab1b52d436e96806dad`, the multi-arch image passmcp's release published for 0.0.3, and every install line and the registry listing name the release. Opened by the sync workflow on the release's dispatch; passmcp's own changelog says what changed in the diagnostic.

## [0.0.2] — 2026-09-29

The family's second release. The three tools, their arguments and their
results are the 0.0.1 ones; what changed is how the repository is built,
checked and documented.

### Added

- **`make install` and `make uninstall`.** A `GNUmakefile` installs the
  binary and its bash, zsh and fish completions under `PREFIX` (default
  `/usr/local`), staged under `DESTDIR` for packagers. CI stages an
  install, runs the binary and checks that uninstall leaves nothing.
- **A coverage badge and an OpenSSF Scorecard.** The Manual workflow
  publishes `coverage.json`, measured from the commit it deploys, beside
  the manual on GitHub Pages; a Scorecard workflow publishes this
  repository's score.
- **Decision records** 0001 to 0004 for the decisions made with the first
  commit, and `docs/COMPARISON.md` and `docs/BENCHMARKS.md` with the
  evidence and method behind the README's tables.
- **A dev container** that boots to a working `make`.

### Changed

- **In lockstep with passmcp 0.0.2.** The image builds on `ghcr.io/sebastienrousseau/passmcp@sha256:fb15e3a3ff2bc2270ce308778ead54da3f10bd0f8160aaab4a06a60a2da6bbdb`, the multi-arch image passmcp's release published for 0.0.2, and every install line and the registry listing name the release. Written by the sync workflow's rewrite step on the release branch; passmcp's own changelog says what changed in the diagnostic.
- **passmcp-reporting at its 0.0.2 release.** `go.mod` required a
  pre-release pseudo-version; it now requires the tagged release, and
  `make versions` fails when the two disagree.
- **The README follows the family standard**: the seven-badge row, the
  family table linking all nine components, and a release for every
  capability rather than a vague status.
- **Complexity ceilings are a lint gate**: cyclomatic 10, cognitive 15 and
  60 lines per function. The five functions over them were split, with no
  change in behaviour, and so were the two in `scripts/spdx_sweep.go`,
  which the linter does not see because it builds only with `go run`.
- `ARCHITECTURE.md` moved to `docs/ARCHITECTURE.md`.

### Fixed

- **The install lines were never version-checked.**
  `scripts/verify-release-versions.sh` and the sync workflow looked for the
  module path the project had before its rename, so they matched nothing:
  a stale `go install` line passed the release check, and a sync would not
  have moved it. Both now match `satellion.com/...`, and the check fails
  if it finds no install line at all.
- **The sync workflow no longer stops on this file's shape.** It read the
  release's entries only from an `## [Unreleased]` section, which this
  file did not have, so passmcp's next release dispatch would have ended
  in an error. It now dates an undated `## [X.Y.Z]` section a release
  branch opened ahead of the release, or else moves the `## [Unreleased]`
  entries, and leaves an empty `## [Unreleased]` above either way.
- **`make family` accepts the family manifest's new status names.**
  Schema 2 of passmcp's `ecosystem.json` calls a released component
  `released` where schema 1 said `shipping`; the check accepts both, so it
  holds before and after passmcp's manifest moves.

## [0.0.1] — 2026-09-29

The first release.

### Added

- **passmcp as MCP tools.** An agent can evaluate an MCP server, or check an
  attestation about one, from inside the editor. The tools are read-only,
  the targets allowlisted, and no credential is ever sent.

[Unreleased]: https://github.com/sebastienrousseau/passmcp-server/compare/v0.0.3...HEAD
[0.0.3]: https://github.com/sebastienrousseau/passmcp-server/releases/tag/v0.0.3
[0.0.2]: https://github.com/sebastienrousseau/passmcp-server/releases/tag/v0.0.2
[0.0.1]: https://github.com/sebastienrousseau/passmcp-server/releases/tag/v0.0.1
