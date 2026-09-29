<!-- SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com> -->
<!-- SPDX-License-Identifier: GPL-3.0-only -->

# 0004 — The version is passmcp's latest release, and the image is passmcp's

**Status:** Accepted · **Decided:** 2026-09-28, with the first commit ·
**Recorded:** 2026-09-29

## Context

passmcp-server wraps a passmcp release: it runs the program (see
[ADR 0001](0001-run-passmcp-as-a-program.md)) and its image carries it.
A version chosen here independently would tell a reader nothing about
which passmcp they get. The passmcp family releases every component at
one version, together
([passmcp's docs/ecosystem.md](https://github.com/sebastienrousseau/passmcp/blob/main/docs/ecosystem.md)).

## Decision

- The version is passmcp's latest release, exactly. It lives in the
  newest `## [x.y.z]` heading of `CHANGELOG.md`; `main.Version` is
  stamped at build time and never hard-coded.
- The `Dockerfile` builds `FROM` passmcp's release image for that
  version, pinned by digest, and names the version in `PASSMCP_VERSION`.
- `go.mod` requires passmcp-reporting at the same release.
- passmcp's release dispatch runs `.github/workflows/sync.yml`, which opens
  a pull request moving every one of those together.

## Consequences

- A release here can contain nothing new; that is the rule working.
- Four checks hold it: `make lockstep` (the version is passmcp's
  latest), `make digest` (the pinned digest is passmcp's image for
  `PASSMCP_VERSION`), `make versions` (every file that names the version,
  and `go.mod`, agree), and `make family` (this repository's row in the
  family manifest is true).
- passmcp-reporting must be tagged before passmcp, so the sync can move
  `go.mod` to it.
