<!-- SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com> -->
<!-- SPDX-License-Identifier: GPL-3.0-only -->

# Architecture Decision Records

Decisions about this repository that will be questioned later, with the
reasoning that produced them. Records are immutable once merged; a
decision that changes gets a new record superseding the old one.

Decisions about the checks, the score and the attestation format are
passmcp's, in [passmcp's ADRs](https://github.com/sebastienrousseau/passmcp/blob/main/docs/adr/README.md).
This directory records what is decided here. The four records below were
made with the first commit and written down afterwards, from the package
documentation and tests that enforce them.

| # | Decision | Status |
|---|---|---|
| [0001](0001-run-passmcp-as-a-program.md) | Run the passmcp program; do not link its engine | Accepted |
| [0002](0002-loopback-only-allowlist.md) | The allowlist is on by default and allows loopback only | Accepted |
| [0003](0003-no-operator-configuration-credentials-or-mutations.md) | No operator configuration, no credentials, no mutations | Accepted |
| [0004](0004-version-in-lockstep-with-passmcp.md) | The version is passmcp's latest release, and the image is passmcp's | Accepted |
