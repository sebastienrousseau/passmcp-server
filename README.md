<!-- SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com> -->
<!-- SPDX-License-Identifier: GPL-3.0-only -->

<p align="center">
  <img src="https://raw.githubusercontent.com/sebastienrousseau/passmcp/main/.github/logo.svg" alt="passmcp-server logo" width="128" />
</p>

<h1 align="center">passmcp-server</h1>

<p align="center">
  passmcp, the Model Context Protocol server diagnostic, as MCP tools — so an agent can evaluate a server, or check an attestation about one, from inside the editor. Read-only, allowlisted, and it never sends a credential.
</p>

<p align="center">
  <a href="https://github.com/sebastienrousseau/passmcp-server/actions/workflows/ci.yml"><img src="https://img.shields.io/github/actions/workflow/status/sebastienrousseau/passmcp-server/ci.yml?branch=main&style=for-the-badge&logo=github&label=Build" alt="Build" /></a>
  <a href="https://github.com/sebastienrousseau/passmcp-server/blob/main/DEVELOPMENT.md#coverage"><img src="https://img.shields.io/endpoint?url=https%3A%2F%2Fsebastienrousseau.com%2Fpassmcp-server%2Fcoverage.json&style=for-the-badge&logo=codecov&logoColor=white" alt="Coverage" /></a>
  <a href="https://github.com/sebastienrousseau/passmcp-server/releases"><img src="https://img.shields.io/github/v/release/sebastienrousseau/passmcp-server?style=for-the-badge&color=fc8d62&logo=github&label=Release" alt="Release" /></a>
  <a href="https://pkg.go.dev/satellion.com/passmcp-server"><img src="https://img.shields.io/badge/go.dev-reference-007d9c?style=for-the-badge&labelColor=555555&logo=go&logoColor=white" alt="Docs" /></a>
  <a href="https://scorecard.dev/viewer/?uri=github.com/sebastienrousseau/passmcp-server"><img src="https://img.shields.io/ossf-scorecard/github.com/sebastienrousseau/passmcp-server?style=for-the-badge&label=OpenSSF%20Scorecard&logo=openssf" alt="OpenSSF Scorecard" /></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-GPL--3.0--only-blue.svg?style=for-the-badge" alt="License: GPL-3.0-only" /></a>
  <a href="https://github.com/sebastienrousseau/passmcp-server/blob/main/DEVELOPMENT.md#requirements"><img src="https://img.shields.io/badge/go-1.26.8%2B-93450a.svg?style=for-the-badge&logo=go" alt="Go 1.26.8+" /></a>
</p>

---

## Contents

**Getting started**

- [Install](#install) — release archives, `go install`, `make install`, the container image, and an MCP host configuration
- [Requirements](#requirements) — passmcp itself, and the Go floor to build from source
- [Quick Start](#quick-start) — ask the agent to evaluate a local server

**The passmcp-server ecosystem**

- [The passmcp-server ecosystem](#the-passmcp-server-ecosystem) — `passmcp`, `passmcp-reporting`, `passmcp-server`, `passmcp-action`, `passmcp-graph`, `passmcp-registry`, `passmcp-lsp`, `passmcp-census`, `satellion.com`

**Library reference**

- [Capabilities at a glance](#capabilities-at-a-glance) — the current surface by theme
- [Ecosystem comparison](#ecosystem-comparison) — short matrix; full table at [`docs/COMPARISON.md`](docs/COMPARISON.md)
- [Benchmarks](#benchmarks) — headline numbers; full table at [`docs/BENCHMARKS.md`](docs/BENCHMARKS.md)
- [Features](#features) — the allowlist, no credentials, read-only
- [Configuration](#configuration) — two flags and their environment variables
- [Examples](#examples) — tool calls and results

**Operational**

- [When not to use passmcp-server](#when-not-to-use-passmcp-server) — limitations
- [Development](#development) — make targets, the install contract, CI
- [Security](#security) — what an agent can and cannot make it do
- [Documentation](#documentation) — all reference docs
- [Stability guarantees](#stability-guarantees) — tool names, arguments and results
- [License](#license)

---

## Install

`passmcp_check` runs the `passmcp` program. **passmcp must be installed and on
`PATH`, or named with `--passmcp`**, for that tool to work; the container
image carries it. `passmcp_verify_attestation` needs nothing but passmcp-server.

### As a Go program

```sh
go install satellion.com/passmcp-server/cmd/passmcp-server@v0.0.1
go install satellion.com/passmcp/cmd/passmcp@v0.0.1
```

Release binaries for Linux, macOS and Windows on amd64 and arm64 are on
the [releases page](https://github.com/sebastienrousseau/passmcp-server/releases),
with signed checksums and SLSA provenance.

### From source, with `make install`

```sh
git clone https://github.com/sebastienrousseau/passmcp-server
cd passmcp-server
make install PREFIX="$HOME/.local"
```

`make install` builds the binary and its bash, zsh and fish completions and
installs them under `PREFIX` (default `/usr/local`), staged under `DESTDIR`
when a packager sets it. `make uninstall` removes them. The
[GNUmakefile](GNUmakefile) holds the contract, and CI checks the staged tree
on every push.

### As a container image

```sh
docker pull ghcr.io/sebastienrousseau/passmcp-server:0.0.1
```

The image is passmcp's own release image with passmcp-server added: distroless,
non-root, linux/amd64 and linux/arm64, with passmcp at `/usr/local/bin/passmcp`
and passmcp-server started with `--passmcp` pointing at it.

### In an MCP host

passmcp-server speaks MCP over stdio, so a host starts it as a child process.
For Claude Desktop (`claude_desktop_config.json`), Claude Code (`.mcp.json`)
and other hosts that read an `mcpServers` block:

```json
{
  "mcpServers": {
    "passmcp": {
      "command": "passmcp-server",
      "args": ["--allow", ".internal.example.com"]
    }
  }
}
```

With the container image instead:

```json
{
  "mcpServers": {
    "passmcp": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "ghcr.io/sebastienrousseau/passmcp-server:0.0.1"]
    }
  }
}
```

Inside a container, loopback is the container itself. To evaluate a server
on the host, run the container with `--network host` on Linux, or point at
`host.docker.internal` and add `-e PASSMCP_SERVER_ALLOW=host.docker.internal`
to the arguments.

The server is listed in the official MCP Registry as
`com.sebastienrousseau/passmcp-server`; [`server.json`](server.json) is
that listing.

---

## Requirements

| Requirement | Floor | Enforced by |
|---|---|---|
| passmcp | on `PATH`, or `--passmcp`; the version this one is in lockstep with | CI builds passmcp at that tag and evaluates this server with it |
| Go (building from source) | the `go` directive in [`go.mod`](go.mod) | CI tests on that version and on latest stable, on Linux, macOS and Windows |
| An MCP host | any that starts stdio servers and speaks revision 2025-03-26 or later | the handshake negotiates 2025-11-25, 2025-06-18 or 2025-03-26, and 2026-07-28 through `server/discover` |
| The server under test | a Streamable HTTP endpoint on the allowlist | `passmcp_check` refuses any other URL before passmcp runs |

The Go floor is raised only when a release needs a language feature, on a
patch release like everything else pre-1.0, and the changelog says so.

---

## Quick Start

```sh
go install satellion.com/passmcp-server/cmd/passmcp-server@v0.0.1
go install satellion.com/passmcp/cmd/passmcp@v0.0.1
claude mcp add passmcp -- passmcp-server
```

Then, with an MCP server of your own listening on
`http://127.0.0.1:3000/mcp`, ask the agent:

> Evaluate my MCP server at <http://127.0.0.1:3000/mcp> with passmcp and fix what fails.

The agent calls `passmcp_check`, which runs
`passmcp check http://127.0.0.1:3000/mcp --auth none --output json` and
returns the score, the grade and every failing check with its detail and a
link to the fix. Loopback needs no configuration; any other host must be
named with `--allow` first.

---

## The passmcp-server ecosystem

Every component is released at **0.0.1** and moves in lockstep: one version across the family, released together ([docs/ecosystem.md](https://github.com/sebastienrousseau/passmcp/blob/main/docs/ecosystem.md)).

| Component | Purpose | Use case |
| :--- | :--- | :--- |
| [passmcp](https://github.com/sebastienrousseau/passmcp) | The MCP server diagnostic: checks in nine phases, every finding tied to the request that showed it, signed attestations | Test a server before your agents trust it, and gate it in CI |
| [passmcp-reporting](https://github.com/sebastienrousseau/passmcp-reporting) | The attestation format, its JSON Schemas and offline verifier, the graph model, and the agentgateway processor | Verify an attestation in a gateway, registry or pipeline |
| [passmcp-server](https://github.com/sebastienrousseau/passmcp-server) | passmcp's diagnostics as read-only MCP tools | Evaluate a server, or check an attestation, from inside the agent |
| [passmcp-action](https://github.com/sebastienrousseau/passmcp-action) | passmcp in GitHub Actions and GitLab CI, the image pinned by digest | Fail a build on the findings you choose |
| [passmcp-graph](https://github.com/sebastienrousseau/passmcp-graph) | A local graph of agents, servers, tools and identities built from attestations | Find inherited risk and over-privilege, and gate on policy |
| [passmcp-registry](https://github.com/sebastienrousseau/passmcp-registry) | A signed public scorecard of the MCP Registry's remote servers | Check a public server's standing before connecting to it |
| [passmcp-lsp](https://github.com/sebastienrousseau/passmcp-lsp) | A language server for MCP artefacts, with check-id hover from the guidance catalogue | Catch mistakes in server.json, tool schemas and client configuration while editing |
| [passmcp-census](https://github.com/sebastienrousseau/passmcp-census) | The published reliability census: dataset, methodology, disclosure log and reproduction command | Cite ecosystem-wide reliability figures, and reproduce them |
| [satellion.com](https://github.com/sebastienrousseau/satellion.github.io) | The website, the Go module paths and the format URIs | Read the manual, and resolve `satellion.com/...` imports |

This repository is the distribution surface: its deliverable is a registry
listing, so passmcp is where agents look for tools. It wraps passmcp's
release, so its version is passmcp's latest, exactly; `make lockstep`
checks that, `make family` checks this repository's row in the family
manifest, and `make versions` checks that every file naming the version,
and the passmcp-reporting module in `go.mod`, agree on it.

---

## Capabilities at a glance

| Area | Capability | Status |
| :--- | :--- | :--- |
| Evaluate | `passmcp_check`: score, grade, counts and up to 25 failing checks for an allowlisted Streamable HTTP endpoint, optionally narrowed to some of passmcp's nine phases | [Released in 0.0.1](https://github.com/sebastienrousseau/passmcp-server/releases/tag/v0.0.1) |
| Verify | `passmcp_verify_attestation`: structure, subject digest and target of a passmcp attestation, offline | [Released in 0.0.1](https://github.com/sebastienrousseau/passmcp-server/releases/tag/v0.0.1) |
| Identify | `passmcp_version`: passmcp-server's version and the passmcp it runs | [Released in 0.0.1](https://github.com/sebastienrousseau/passmcp-server/releases/tag/v0.0.1) |
| Results | Text for the agent, plus `structuredContent` matching each tool's `outputSchema` | [Released in 0.0.1](https://github.com/sebastienrousseau/passmcp-server/releases/tag/v0.0.1) |
| Protocol | stdio; handshake 2025-11-25, 2025-06-18, 2025-03-26; `server/discover` for 2026-07-28; `ping` | [Released in 0.0.1](https://github.com/sebastienrousseau/passmcp-server/releases/tag/v0.0.1) |
| Servers that are programs (`--stdio`) | not through a tool; run `passmcp check --stdio` yourself | Out of scope |
| Credentials | never sent; every run is `--auth none` | Out of scope by design |

---

## Ecosystem comparison

The alternative is running passmcp in a terminal and pasting the report into
the conversation, or poking the server by hand in an inspector. passmcp-server
makes the decisions an agent should not: which hosts it may reach, that
no credential travels, and a result sized for a context window.

| Approach | An agent can call it | Targets limited by the operator | Scored, with remediation links |
| :--- | :---: | :---: | :---: |
| **passmcp-server** | yes | yes — loopback unless `--allow` | yes |
| `passmcp check` in a terminal | no | the operator types the URL | yes |
| [MCP Inspector](https://github.com/modelcontextprotocol/inspector) | no — a UI for a person | the operator types the URL | no |

See [`docs/COMPARISON.md`](docs/COMPARISON.md) for the evidence and complete matrix.

---

## Benchmarks

The server adds a process start and a JSON round trip to a run; the run
itself is passmcp's, and bounded at five minutes. Measured with
[hyperfine](https://github.com/sharkdp/hyperfine) on binaries built from
this tree and passmcp 0.0.1, on a machine that was running other builds at
the time, so the spread is wide and the minimum is the better guide.

| Scenario | Result | Environment |
| :--- | ---: | :--- |
| Start, `initialize`, `tools/list`, `passmcp_version`, exit | 27.2 ms ± 22.2 ms mean, 6.2 ms min (50 runs) | Apple A18 Pro, Go 1.27.1, 2026-09-29, load average 19 |
| passmcp's full stdio evaluation of passmcp-server | 242.3 ms ± 80.5 ms mean, 148.5 ms min (30 runs) | same |
| `passmcp_check` against a server | passmcp's own timings, in its report | the target server |

See [`docs/BENCHMARKS.md`](docs/BENCHMARKS.md) for methodology and full results.

---

## Features

**An allowlist, on by default and narrow.** With no configuration,
`passmcp_check` evaluates loopback addresses only: `localhost`, any
`*.localhost` name, and loopback IPs. `--allow` names more hosts, exactly
or as `.example.com` for every subdomain (not the apex). passmcp makes real
requests to the endpoint it is given, and an agent choosing that endpoint
from a prompt is exactly the situation in which a request should not go
anywhere the operator did not name. A URL that is not http or https, or
that carries credentials in it, is refused before passmcp runs.

**No credentials.** Every run is `passmcp check <endpoint> --auth none`, with
`PASSMCP_CONFIG` pointing at an empty configuration file, so no profile or
defaults block the operator wrote for their own use can add a credential,
switch on mutations or redirect the report. The operator's secrets do not
travel to wherever an agent points passmcp.

**Read-only, twice.** Every tool here is annotated `readOnlyHint: true`.
And no flag that allows a mutating tool is ever passed to passmcp, so passmcp
invokes only the evaluated server's tools that declare `readOnlyHint`.

**passmcp, not a copy of it.** passmcp-server runs the passmcp program and reads
the JSON report it prints; it does not link passmcp's engine. The version
`passmcp_version` reports is whatever `passmcp version` says, and every safety
property passmcp has holds, because the server can only ask for what passmcp's
own flags allow.

**Answers an agent can act on.** A failing check comes with its id,
severity, detail and a documentation link. A statement that does not
verify is a result saying why, not a tool error. A server that fails
everything returns the first 25 failures and says it truncated.

---

## Configuration

| Flag | Environment | Default | Meaning |
|---|---|---|---|
| `--allow` | `PASSMCP_SERVER_ALLOW` | empty: loopback only | Hosts `passmcp_check` may evaluate besides loopback, comma-separated; a leading dot allows subdomains |
| `--passmcp` | `PASSMCP_SERVER_PASSMCP` | `passmcp` on `PATH` | Path of the passmcp program |
| `--version` | — | — | Print the version and exit |
| `--completion` | — | — | Print a completion script for `bash`, `zsh` or `fish`, and exit |

A flag overrides its environment variable. Nothing else is read: passmcp
itself runs with an empty configuration file, whatever the operator's own
passmcp configuration says.

Shell completions come from the flag set, so they list every flag:

```sh
passmcp-server --completion bash > /etc/bash_completion.d/passmcp-server
passmcp-server --completion zsh > "${fpath[1]}/_passmcp-server"
passmcp-server --completion fish > ~/.config/fish/completions/passmcp-server.fish
```

---

## Examples

A `tools/call` for `passmcp_check`:

```json
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"passmcp_check","arguments":{"endpoint":"http://127.0.0.1:3000/mcp","phases":["handshake","catalog"]}}}
```

`arguments` for the other two:

```json
{"statement": "<the attestation's JSON text>", "endpoint": "https://mcp.example.com/mcp"}
```

```json
{}
```

The first is `passmcp_verify_attestation`, with `endpoint` optional; the
second is `passmcp_version`, which takes no arguments. An unknown argument
is refused rather than ignored, so a misspelling is reported.

To see the handshake and the tool list without a host:

```sh
make smoke
```

---

## When not to use passmcp-server

- **For a server that is a program.** `passmcp_check` takes a Streamable
  HTTP URL. An agent starting arbitrary programs is a different trust
  decision; run `passmcp check --stdio -- <command>` yourself.
- **For a server that needs credentials.** Every run is `--auth none`, so
  the evaluation is what an unauthenticated client sees. Run passmcp in a
  terminal with the credentials you were given for the rest.
- **As a CI gate.** Use [passmcp-action](https://github.com/sebastienrousseau/passmcp-action),
  which keeps the report and fails the job.
- **To check who signed an attestation.** `passmcp_verify_attestation`
  checks structure and integrity; the signature is the envelope's, and
  `cosign` or `gh attestation verify` checks it.
- **For many evaluations at once.** Requests are answered in turn; a host
  that wants two at once starts two servers.

---

## Development

```bash
make            # format, vet, lint, headers, tests, stdio smoke test
make test-race  # race detector, randomised order
make image      # the container image for this machine, without goreleaser
make family     # this repository's row in passmcp's family manifest
make lockstep   # the version is passmcp's latest release
make versions   # every version-bearing file, and go.mod, name that release
make coverage-json  # build/coverage.json, the document behind the badge
make install-smoke  # install and uninstall under a staged DESTDIR
```

Every gate CI runs has a local form; [DEVELOPMENT.md](DEVELOPMENT.md) maps
them. CI also builds passmcp at the lockstep version and has it evaluate this
server over stdio, failing below 90 or on any failing check but
`supply.provenance`.

---

## Security

The endpoint an agent passes is checked against the allowlist before passmcp
runs, and a refused endpoint never reaches passmcp. passmcp runs with
`--auth none` and an empty configuration file, and with no flag that
allows a mutating tool; each of those is a test in `internal/runner` and
`internal/server`. What reaches the agent is bounded: at most 25 failures
per result, and passmcp's error output cut to its last 400 characters. CI
runs `govulncheck` on every push, and releases are signed with cosign
keyless and carry SLSA build provenance.

Report vulnerabilities according to [`SECURITY.md`](SECURITY.md).

---

## Documentation

The four entry points, identical across every repo in the family:

- **[User Manual](https://satellion.com/passmcp/docs/)** — passmcp's rendered manual: the phases, the checks, the report
- **[API reference](https://pkg.go.dev/satellion.com/passmcp-server)** — this module's packages
- **[Developer docs](DEVELOPMENT.md)** — toolchain, task map, reproducing every CI gate locally
- **[Ecosystem map](https://github.com/sebastienrousseau/passmcp/blob/main/docs/ecosystem.md)** — the family, the published artefacts, the lockstep version rule

| Document | Covers |
|---|---|
| [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) | How a call flows, the packages, and the settings every run is fixed to |
| [`docs/adr/`](docs/adr/README.md) | Decision records for this repository |
| [`docs/COMPARISON.md`](docs/COMPARISON.md) | passmcp-server beside the other ways to run passmcp |
| [`docs/BENCHMARKS.md`](docs/BENCHMARKS.md) | What the server adds to a run, and how it was measured |
| [`docs/publishing.md`](docs/publishing.md) | Publishing the registry listing |
| [`docs/releases/`](docs/releases/v0.0.1.md) | Release highlights, one file per release |
| [`server.json`](server.json) | The MCP Registry listing |
| [`SECURITY.md`](SECURITY.md) | Disclosure policy, supported versions, what is guaranteed |
| [`CONTRIBUTING.md`](CONTRIBUTING.md) | Signed-commit and DCO policy, what a change needs |
| [`CHANGELOG.md`](CHANGELOG.md) | Per-release notes, and the lockstep version rule |
| [`SUPPORT.md`](SUPPORT.md) | Where to ask, and what to expect |

---

## Stability guarantees

passmcp-server is pre-1.0, carries passmcp's version, and follows SemVer with the
patch digit moving for everything until 1.0.

**The breaking axis is what an agent or a host relies on.** These are
breaking:

- Removing or renaming a tool, an argument or a structured result field
- Making an optional argument required
- Widening the default allowlist, sending a credential, or passing passmcp a
  flag that allows a mutating tool
- Changing a flag's or an environment variable's meaning

Added tools, optional arguments and result fields, and a new passmcp release
underneath are **not** breaking. What a run reports is passmcp's, and passmcp's
own stability rule governs it.

**Deprecation window.** A deprecated tool or argument keeps working for at
least one release after the release that announces it.

---

## License

Licensed under the **[GNU General Public License v3.0 only](LICENSE)**.

passmcp-server is GPL-3.0-only like [passmcp](https://github.com/sebastienrousseau/passmcp),
the program it runs. It uses the Apache-2.0
[passmcp-reporting](https://github.com/sebastienrousseau/passmcp-reporting)
verifier for attestations.

<p align="right"><a href="#contents">Back to Top</a></p>
