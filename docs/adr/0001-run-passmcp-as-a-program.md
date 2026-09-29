<!-- SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com> -->
<!-- SPDX-License-Identifier: GPL-3.0-only -->

# 0001 — Run the passmcp program; do not link its engine

**Status:** Accepted · **Decided:** 2026-09-28, with the first commit ·
**Recorded:** 2026-09-29

## Context

passmcp-server exists to put passmcp's evaluation in an agent's hands.
There were two ways to get at that evaluation: import passmcp's engine as
a Go library, or run the `passmcp` program and read the report it prints.

passmcp's engine lives in its `internal` packages, which another module
cannot import. Exposing it would make those packages a public API that
passmcp would then have to keep stable for this one consumer. Linking a
copy would also fix this server to whichever build of the engine it was
compiled against, so the version it reported and the version an operator
runs by hand could differ.

## Decision

`internal/runner` runs the passmcp program, the same artefact an
operator runs, with `--output json`, and parses the report. The path is
`passmcp` on `PATH`, or `--passmcp` / `PASSMCP_SERVER_PASSMCP`. The
container image is built on passmcp's own release image, so it carries
the program at `/usr/local/bin/passmcp`.

## Consequences

- `passmcp_version` reports whatever `passmcp version` says; there is no
  second version to drift.
- Every safety property passmcp has holds here, because the server can only
  ask for what passmcp's own flags allow (see
  [ADR 0003](0003-no-operator-configuration-credentials-or-mutations.md)).
- `passmcp_check` does not work until passmcp is installed. The README's
  Install section says so first, and `passmcp_version` names the fix when
  the program is missing.
- A run costs a process start. The benchmarks measure it
  ([BENCHMARKS.md](../BENCHMARKS.md)).
- Attestation checks are the exception: they use the Apache-2.0 verifier
  from passmcp-reporting in process, because that is a published library
  with a stable API and needs no network.

## Where it is enforced

The package documentation of `internal/runner`, and
`internal/runner/runner_test.go`, which stands a shell script in for
passmcp and asserts the arguments and environment of every run.
