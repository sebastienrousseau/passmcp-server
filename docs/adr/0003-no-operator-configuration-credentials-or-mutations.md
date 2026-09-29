<!-- SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com> -->
<!-- SPDX-License-Identifier: GPL-3.0-only -->

# 0003 — No operator configuration, no credentials, no mutations

**Status:** Accepted · **Decided:** 2026-09-28, with the first commit ·
**Recorded:** 2026-09-29

## Context

The agent chooses the endpoint (see [ADR 0002](0002-loopback-only-allowlist.md)).
Anything the operator configured for their own use of passmcp — a
profile with a token, a defaults block that allows mutating tools, a
report destination — would otherwise apply to a run an agent started,
against a target the operator did not pick for it.

## Decision

Every run is fixed three ways, whatever the tool arguments say:

- **No configuration.** `PASSMCP_CONFIG` points at an empty file.
- **No credentials.** `--auth none`. No tool argument or server flag
  forwards a token.
- **Read-only.** No flag that allows a mutating or destructive tool is
  passed, so passmcp invokes only the evaluated server's tools that
  declare `readOnlyHint`. passmcp-server's own tools are all annotated
  `readOnlyHint: true`.

Tool arguments are decoded strictly: an unknown argument is refused, not
ignored, so a misspelt or invented option cannot be mistaken for one
that took effect.

## Consequences

- An evaluation shows what an unauthenticated client sees. A server that
  needs credentials is evaluated with passmcp in a terminal; the README's
  "When not to use" section says so.
- Adding a credential of any kind, or a flag that allows mutations, is a
  breaking change under the Stability guarantees, and is expected to be
  declined.

## Where it is enforced

The package documentation and tests of `internal/runner`, which assert
the arguments and environment of every run, and `strictDecode` in
`internal/server/tools.go`.
