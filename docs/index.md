<!-- SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com> -->
<!-- SPDX-License-Identifier: GPL-3.0-only -->

# passmcp-server documentation

passmcp, the MCP server diagnostic, as read-only MCP tools over stdio.

| Document | Covers |
|---|---|
| [../README.md](https://github.com/sebastienrousseau/passmcp-server/blob/main/README.md) | Install, host configuration, the three tools, the allowlist |
| [ARCHITECTURE.md](ARCHITECTURE.md) | How a call flows, the packages, and the settings every run is fixed to |
| [adr/](adr/README.md) | Decision records for this repository |
| [COMPARISON.md](COMPARISON.md) | passmcp-server beside the other ways to run passmcp |
| [BENCHMARKS.md](BENCHMARKS.md) | What the server adds to a run, and how it was measured |
| [publishing.md](publishing.md) | Publishing the MCP Registry listing after a release |
| [Release 0.0.4](releases/v0.0.4.md) | The highlights of the 0.0.4 release |
| [Release 0.0.3](releases/v0.0.3.md) | The highlights of the 0.0.3 release |
| [Release 0.0.2](releases/v0.0.2.md) | The highlights of the 0.0.2 release |
| [Release 0.0.1](releases/v0.0.1.md) | The highlights of the first release |
| [../server.json](https://github.com/sebastienrousseau/passmcp-server/blob/main/server.json) | The registry listing itself |

What each check means, and how to fix it, is passmcp's manual:
<https://satellion.com/passmcp/docs/>.

The README's coverage badge reads
[coverage.json](https://sebastienrousseau.com/passmcp-server/coverage.json), published
with this site by the Manual workflow from the statement coverage CI
measures on `main`.
