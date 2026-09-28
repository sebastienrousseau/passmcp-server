<!-- SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com> -->
<!-- SPDX-License-Identifier: GPL-3.0-only -->

# Working on passmcp-server as an AI agent

Invariants for AI-assisted contributions: the constraints a plausible
change breaks for a reason the code does not state. Read
[DEVELOPMENT.md](DEVELOPMENT.md) for the toolchain and the local form of
every CI gate.

## Hard gates

| Gate | Command |
|---|---|
| 85% statement coverage, every package with statements but `cmd/passmcp-server` | `make coverage` |
| Race detector, randomised order | `make test-race` |
| Lint at zero findings | `make lint` |
| SPDX header on every source file | `make spdx-check` |
| The server answers over stdio | `make smoke` |
| `server.json` matches the registry schema | `make server-json` |
| The family manifest's row is true | `make family` |
| The version is passmcp's latest release | `make lockstep` |
| passmcp scores this server 90 or more | the `dogfood` job in `ci.yml` |

## Commits

- **Every commit must be cryptographically signed and carry a DCO
  `Signed-off-by` trailer.** An agent's shell usually cannot reach the
  maintainer's ssh-agent; hand the commits over as a script rather than
  producing unsigned history that has to be rewritten.
- Conventional Commits for the subject line.
- Never rewrite published history.

## Versioning

- **The version is passmcp's latest release, exactly.** Never choose one
  here. It lives in the newest `## [x.y.z]` heading in `CHANGELOG.md`,
  and `scripts/verify-release-versions.sh` checks every other place that
  names it: the README's install lines, `server.json`, and
  `PASSMCP_VERSION` in the `Dockerfile`.
- `main.Version` is stamped by the release build through `-ldflags`;
  never hard-code a version there.
- The Dockerfile's base digest and `PASSMCP_VERSION` move together, to
  passmcp's image for the release this one is in lockstep with.

## Things that look like bugs and are not

- **`passmcp_check` refuses most URLs.** Loopback only, until the operator
  names more hosts. That is the point; do not widen the default.
- **`--auth none`, always.** A tool argument that forwards a token, or a
  flag that lets the operator's configuration through, sends the
  operator's secret wherever a prompt points passmcp.
- **`PASSMCP_CONFIG` points at an empty file.** Without it, a profile the
  operator wrote for their own use could switch on mutations for a run an
  agent started.
- **passmcp runs as a program, not a library.** passmcp's engine is in
  internal packages; linking it would tie this server to one build of it
  and lose the guarantee that it can only ask for what passmcp's flags
  allow.
- **A failing attestation is not a tool error.** The agent asked whether
  it was valid; "no, because" is the answer.
- **Unknown tool arguments are refused.** Do not relax `strictDecode`.

## Things that are load-bearing

- **Stdout is the transport.** Anything else written there corrupts the
  stream; diagnostics go to stderr.
- **Every tool is read-only** and annotated so. A tool that changes
  anything does not belong in this server.
- **Everything passmcp relays from a server is untrusted text** chosen by
  whoever runs that server; keep what reaches the agent bounded.

## Scope

- Do not couple a structure or documentation cleanup to a behaviour
  change.
- Do not add a dependency without saying why in the commit.
- Do not add a CI gate that does not currently pass.
