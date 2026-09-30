# syntax=docker/dockerfile:1.7@sha256:a57df69d0ea827fb7266491f2813635de6f17269be881f696fbfdf2d83dda33e
# SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com>
# SPDX-License-Identifier: GPL-3.0-only
#
# Runtime image for passmcp-server.
#
# The binary is built by goreleaser, not here: the image ships the same
# artefact the release archive does. dockers_v2 lays the build context out
# by platform (linux/amd64/passmcp-server, linux/arm64/passmcp-server) and buildx
# sets TARGETPLATFORM for each architecture.
#
# The base is passmcp's own release image, pinned by digest: distroless
# nonroot with passmcp at /usr/local/bin/passmcp. passmcp_check runs that program,
# so the image carries the passmcp release this one is in lockstep with and
# needs nothing else. PASSMCP_VERSION names that release, and
# scripts/verify-release-versions.sh refuses a release where it and the
# tag disagree; move it and the digest together. To resolve the digest:
#   docker buildx imagetools inspect ghcr.io/sebastienrousseau/passmcp:<version>
ARG PASSMCP_VERSION=0.0.4
FROM ghcr.io/sebastienrousseau/passmcp@sha256:a619e709cddc568394e7bf3b700f9ca44d925c658087c882d5a077ad65cdfdc7

ARG PASSMCP_VERSION
ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/passmcp-server /usr/local/bin/passmcp-server

# io.modelcontextprotocol.server.name is how the MCP Registry verifies that
# whoever publishes server.json owns this image; it must equal server.json's
# name exactly.
LABEL io.modelcontextprotocol.server.name="com.sebastienrousseau/passmcp-server" \
      io.github.sebastienrousseau.passmcp.version="${PASSMCP_VERSION}" \
      org.opencontainers.image.source="https://github.com/sebastienrousseau/passmcp-server" \
      org.opencontainers.image.url="https://github.com/sebastienrousseau/passmcp-server" \
      org.opencontainers.image.description="passmcp-server: passmcp's MCP server diagnostics as read-only MCP tools, over stdio" \
      org.opencontainers.image.licenses="GPL-3.0-only" \
      org.opencontainers.image.title="passmcp-server" \
      org.opencontainers.image.vendor="Sebastien Rousseau"

# distroless nonroot is uid/gid 65532. Numeric, so Kubernetes runAsNonRoot
# can verify it.
USER 65532:65532
WORKDIR /home/nonroot

# An MCP host runs this with `docker run -i --rm`; stdin and stdout are the
# transport. The path of passmcp is fixed so PATH inside the image does not
# matter.
ENTRYPOINT ["/usr/local/bin/passmcp-server", "--passmcp", "/usr/local/bin/passmcp"]
CMD []
