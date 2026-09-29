<!-- SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com> -->
<!-- SPDX-License-Identifier: GPL-3.0-only -->

# Comparison

How passmcp-server sits beside the other ways to put an MCP server in
front of a diagnostic. The question it answers is narrow: can an agent
run the evaluation itself, safely?

| Approach | An agent can call it | Targets limited by the operator | Credentials the agent can send | Scored, with remediation links | Mutating tools on the target |
| :--- | :---: | :---: | :---: | :---: | :---: |
| **passmcp-server** | yes, as three MCP tools | yes — loopback unless `--allow` | none; every run is `--auth none` | yes, passmcp's score and guidance | never invoked |
| `passmcp check` in a terminal | no | the operator types the URL | whatever the operator configures | yes | only with the operator's explicit flags |
| [passmcp-action](https://github.com/sebastienrousseau/passmcp-action) in CI | no; a pipeline runs it | the pipeline names the URL | the `token` input, from a pipeline secret | yes, and `fail-on` fails the job | only if `args` passes passmcp's flags for them |
| [MCP Inspector](https://github.com/modelcontextprotocol/inspector) | no — an interactive tool for a person | the operator types the URL | whatever the operator enters | no | whatever the operator calls |

## Evidence

- **passmcp-server.** The tools and their annotations are in
  `internal/server/tools.go`; the allowlist is `internal/server/allow.go`
  ([ADR 0002](adr/0002-loopback-only-allowlist.md)); the fixed arguments
  and environment of every run are asserted by
  `internal/runner/runner_test.go`
  ([ADR 0003](adr/0003-no-operator-configuration-credentials-or-mutations.md)).
- **passmcp.** Its flags and configuration are documented in
  [passmcp's manual](https://satellion.com/passmcp/docs/); its read-only
  default is [passmcp's ADR 0004](https://github.com/sebastienrousseau/passmcp/blob/main/docs/adr/0004-read-only-by-default.md).
- **passmcp-action.** Its inputs (`endpoint`, `token`, `args`, `fail-on`
  and the rest) are declared in its `action.yml` and documented in its README.
- **MCP Inspector.** Its README describes it as a visual testing tool for
  MCP servers, driven by a person through a web UI or its CLI mode; it
  lists and calls a server's tools, and does not score the server.

## When the others are the better choice

- To evaluate a server that needs credentials, or one that is a program
  (`--stdio`): passmcp in a terminal.
- To gate a build: passmcp-action, which keeps the report and fails the job.
- To poke at a server's tools by hand while writing it: MCP Inspector.
