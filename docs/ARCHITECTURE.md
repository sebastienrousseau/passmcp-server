<!-- SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com> -->
<!-- SPDX-License-Identifier: GPL-3.0-only -->

# Architecture

passmcp-server is an MCP server over stdio that gives an agent three
read-only tools: evaluate an MCP server with
[passmcp](https://github.com/sebastienrousseau/passmcp), check a passmcp
attestation offline, and report versions. It holds no diagnostic logic:
evaluation runs the passmcp program, and attestation checks use the
Apache-2.0 verifier from
[passmcp-reporting](https://github.com/sebastienrousseau/passmcp-reporting).

## The two decisions

Both are recorded, with the alternatives, in
[ADR 0001](adr/0001-run-passmcp-as-a-program.md) and
[ADR 0003](adr/0003-no-operator-configuration-credentials-or-mutations.md);
the allowlist's default is [ADR 0002](adr/0002-loopback-only-allowlist.md).

**Run the passmcp program; do not link its engine.** passmcp's engine lives
in its `internal` packages, which cannot be imported, and linking a copy
would tie this server to one build of it. passmcp-server runs the same binary
an operator runs by hand and reads the JSON it prints. The version the
server reports is whatever `passmcp version` says, and every safety
property passmcp has carries over, because the server can only ask for
what passmcp's own flags allow.

**The agent chooses the endpoint, so nothing of the operator's travels
with it.** Every run is fixed three ways, whatever the caller asks:

- **No configuration.** `PASSMCP_CONFIG` points at an empty file, so no
  profile the operator wrote for their own use can switch on mutations,
  add credentials or redirect the report.
- **No credentials.** `--auth none`, always.
- **Read-only.** No flag that allows a mutating or destructive tool is
  ever passed.

On top of that, an allowlist decides where passmcp may be pointed:
loopback only by default, widened with `--allow` or `PASSMCP_SERVER_ALLOW`.

## Flow of one call

```text
MCP client ──stdio, JSON-RPC──► internal/server
                                  │ initialize / server/discover / ping
                                  │ tools/list → three tools, all readOnlyHint
                                  │ tools/call:
                                  │   passmcp_check   ─► allowlist ─┬─► internal/runner
                                  │   passmcp_version ──────────────┘         │
                                  │                                       ▼
                                  │                         passmcp check <endpoint>
                                  │                           --auth none --output json
                                  │                           PASSMCP_CONFIG=<empty file>
                                  │   ◄──── score, failing checks, guidance ─┘
                                  │
                                  │   passmcp_verify_attestation ─► passmcp-reporting/attestation
                                  │                                 Parse, Validate, Covers
                                  ◄──── verdict, in process, no network
```

`passmcp_check` refuses an endpoint the allowlist does not cover before
passmcp starts. Requests are answered in turn; a passmcp run can take a
minute, and a client that wants two at once can start two servers.

## Packages

| Path | Role |
| :--- | :--- |
| `cmd/passmcp-server` | Flags (`--allow`, `--passmcp`, `--version`, `--completion`), shell completions generated from the flag set, and the stdio loop |
| `scripts/coveragebadge` | Turns a Go cover profile into the shields.io endpoint document behind the README's coverage badge |
| `internal/server` | JSON-RPC handling, the handshake revisions and `server/discover`, the tool definitions, the allowlist, and attestation checks through passmcp-reporting's verifier |
| `internal/runner` | Runs the passmcp program with the fixed flags and environment, and parses its report |

Stdout is the MCP transport, so nothing else is written there; logs go
to stderr.

## Protocol revisions

The server speaks the handshake revisions 2025-11-25, 2025-06-18 and
2025-03-26 through `initialize`, and the stateless 2026-07-28 revision
through `server/discover`, which carries the server's identity in
`_meta`.

## Distribution

Releases ship binaries and a container image built on passmcp's own image,
pinned by digest, so the image carries the passmcp it runs. The listing in
the MCP Registry is `com.sebastienrousseau/passmcp-server`. Releases
follow passmcp's in lockstep: passmcp's release dispatch opens a sync pull
request that moves the image digest, `PASSMCP_VERSION`, the install lines,
`server.json`, `CITATION.cff` and the passmcp-reporting requirement in
`go.mod` together ([ADR 0004](adr/0004-version-in-lockstep-with-passmcp.md)).
