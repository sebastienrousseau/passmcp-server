#!/usr/bin/env bash
# SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com>
# SPDX-License-Identifier: GPL-3.0-only
#
# Fail unless every place that names the version being released names it:
# the CHANGELOG heading, the install snippets in the README and docs, the
# version the README's ecosystem section states, CITATION.cff, the registry
# listing in server.json, the passmcp release the Dockerfile builds on, and
# the passmcp-reporting release go.mod requires. The family moves in lockstep,
# so every one of them is the same version.
#
#   scripts/verify-release-versions.sh v0.0.1
set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."
tag="${1:-${GITHUB_REF_NAME:-}}"
[ -n "$tag" ] || { echo "usage: $0 vX.Y.Z" >&2; exit 2; }
ver="${tag#v}"
fail=0

grep -Eq "^## \[$ver\]" CHANGELOG.md || { echo "CHANGELOG.md has no '## [$ver]' heading" >&2; fail=1; }

# Both install lines, passmcp-server's and the passmcp it runs; at least
# one must be there, so a path that stops matching cannot pass silently.
installs=$(grep -Eoh 'satellion\.com/passmcp(-server)?/cmd/passmcp(-server)?@v[0-9]+\.[0-9]+\.[0-9]+' README.md docs/*.md || true)
if [ -z "$installs" ]; then
  echo "README.md names no go install line to check" >&2; fail=1
elif grep -v "@v$ver\$" <<<"$installs"; then
  echo "the docs pin a go install version other than $ver" >&2; fail=1
fi

# The ecosystem section states the family's one version.
if ! grep -Fq "Every component is released at **$ver**" README.md; then
  echo "README.md's ecosystem section does not state $ver" >&2; fail=1
fi

if ! grep -Eq "^version: \"?$ver\"?\$" CITATION.cff; then
  echo "CITATION.cff does not say version $ver" >&2; fail=1
fi

# The attestation verifier comes from passmcp-reporting at the same release.
if ! grep -Eq "satellion\.com/passmcp-reporting v$ver\$" go.mod; then
  echo "go.mod does not require satellion.com/passmcp-reporting v$ver" >&2; fail=1
fi
if grep -Eo 'ghcr\.io/sebastienrousseau/passmcp-server:[0-9]+\.[0-9]+\.[0-9]+' README.md docs/*.md | grep -v ":$ver\$"; then
  echo "the docs pin an image version other than $ver" >&2; fail=1
fi

# server.json: the listing's version and the image it points at.
listed=$(python3 - "$ver" <<'PYEOF'
import json, sys
ver = sys.argv[1]
s = json.load(open("server.json"))
bad = []
if s.get("version") != ver:
    bad.append("version is %r" % s.get("version"))
for p in s.get("packages", []):
    if p.get("registryType") == "oci" and not p["identifier"].endswith(":" + ver):
        bad.append("OCI identifier is %r" % p["identifier"])
print("; ".join(bad))
PYEOF
)
if [ -n "$listed" ]; then
  echo "server.json disagrees with $ver: $listed" >&2; fail=1
fi

# The image is built on passmcp's image for the same release.
base=$(sed -nE 's/^ARG PASSMCP_VERSION=([0-9.]+)$/\1/p' Dockerfile)
if [ "$base" != "$ver" ]; then
  echo "Dockerfile builds on passmcp ${base:-<none>}, not $ver; move the FROM digest and PASSMCP_VERSION together" >&2; fail=1
fi

[ "$fail" -eq 0 ] && echo "release versions agree on $ver"
exit "$fail"
