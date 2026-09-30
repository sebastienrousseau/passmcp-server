<!-- SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com> -->
<!-- SPDX-License-Identifier: GPL-3.0-only -->

# Benchmarks

passmcp-server adds a process start and a JSON round trip to a passmcp
run. The run itself is passmcp's, and `passmcp_check` bounds it at five
minutes. These numbers say what the server costs; they say nothing about
how long a target server takes to evaluate, which passmcp reports per
phase.

## Results

| Scenario | Mean ± σ | Min … max | Runs |
| :--- | ---: | ---: | ---: |
| Start, `initialize`, `tools/list`, `passmcp_version`, exit | 27.2 ms ± 22.2 ms | 6.2 ms … 115.9 ms | 50 |
| passmcp's full stdio evaluation of passmcp-server | 242.3 ms ± 80.5 ms | 148.5 ms … 438.1 ms | 30 |

Environment: Apple A18 Pro, macOS, Go 1.27.1, passmcp 0.0.1, measured
2026-09-29 with a load average of 19 from other builds on the same
machine. The spread is wide for that reason; the minimum is the better
guide to the cost on an idle machine. The same evaluation scored
passmcp-server 95 (A), with `supply.provenance` the only failing check,
as expected for a local build.

## Method

Both binaries are built from source, then timed with
[hyperfine](https://github.com/sharkdp/hyperfine):

```sh
GOBIN="$PWD/build/bench" go install satellion.com/passmcp/cmd/passmcp@v0.0.3
CGO_ENABLED=0 go build -trimpath -o build/bench/passmcp-server ./cmd/passmcp-server

printf '%s\n' \
  '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"bench","version":"0"}}}' \
  '{"jsonrpc":"2.0","method":"notifications/initialized"}' \
  '{"jsonrpc":"2.0","id":2,"method":"tools/list"}' \
  '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"passmcp_version","arguments":{}}}' \
  > build/bench/session.jsonl

hyperfine --warmup 5 --runs 50 \
  "build/bench/passmcp-server --passmcp build/bench/passmcp < build/bench/session.jsonl"

hyperfine --warmup 3 --runs 30 \
  "build/bench/passmcp check --stdio --rps 0 --no-color --output json -- build/bench/passmcp-server --passmcp build/bench/passmcp > /dev/null; true"
```

The second command ends in `; true` because passmcp exits 2 when any
check fails, and a local build always fails `supply.provenance`.
