<!-- SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com> -->
<!-- SPDX-License-Identifier: GPL-3.0-only -->

# 0002 — The allowlist is on by default and allows loopback only

**Status:** Accepted · **Decided:** 2026-09-28, with the first commit ·
**Recorded:** 2026-09-29

## Context

passmcp makes real requests to the endpoint it is given. Through
`passmcp_check`, that endpoint comes from an agent, and the agent takes it
from a prompt. A prompt can be written by anyone whose text reaches the
agent. Without a limit, the server would send requests from the
operator's network to any host a prompt names, including hosts inside
that network.

## Decision

`passmcp_check` refuses an endpoint before passmcp starts unless it is:

- an `http` or `https` URL with a host and no credentials in it, and
- a loopback address (`localhost`, any `*.localhost` name, a loopback
  IP), or a host the operator named with `--allow` or
  `PASSMCP_SERVER_ALLOW`: exactly, or as `.example.com` for every
  subdomain (not the apex).

With no configuration, only loopback is allowed.

## Consequences

- A server on the operator's own machine can be evaluated with no setup;
  anything else needs the operator to name it once.
- Inside the container, loopback is the container itself; the README says
  how to reach the host.
- Widening the default is a breaking change under the README's Stability
  guarantees, and needs the week-long tracking issue CONTRIBUTING.md
  requires.

## Where it is enforced

`internal/server/allow.go`, and the allowlist cases in
`internal/server/server_test.go`.
