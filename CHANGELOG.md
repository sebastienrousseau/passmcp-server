<!-- SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com> -->
<!-- SPDX-License-Identifier: GPL-3.0-only -->

# Changelog

All notable changes to passmcp-server are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and versions are
[Semantic Versioning](https://semver.org/) shaped.

**This repository carries passmcp's version.** It is in lockstep with
[passmcp](https://github.com/sebastienrousseau/passmcp): it wraps passmcp's
release — the container image is built on passmcp's image for the same
version — so its version is passmcp's latest release, exactly, and a release
here with nothing in it is the version rule working.

## [0.0.1] — Unreleased

The first release.

### Added

- **passmcp as MCP tools.** An agent can evaluate an MCP server, or check an
  attestation about one, from inside the editor. The tools are read-only,
  the targets allowlisted, and no credential is ever sent.
