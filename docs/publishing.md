<!-- SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com> -->
<!-- SPDX-License-Identifier: GPL-3.0-only -->

# Publishing the registry listing

[`server.json`](https://github.com/sebastienrousseau/passmcp-server/blob/main/server.json) lists passmcp-server in the
[official MCP Registry](https://registry.modelcontextprotocol.io) as
`com.sebastienrousseau/passmcp-server`, pointing at the container image.
It is published by hand, after the release, because it needs the
Maintainer's GitHub login.

## Before

- The release for the version in `server.json` is out, and
  `ghcr.io/sebastienrousseau/passmcp-server:X.Y.Z` is public.
- The image carries the ownership label the registry checks. It must equal
  the `name` in `server.json`:

  ```sh
  docker buildx imagetools inspect ghcr.io/sebastienrousseau/passmcp-server:X.Y.Z \
    --format '{{ json (index .Image "linux/amd64").Config.Labels }}'
  ```

  and the output must map `io.modelcontextprotocol.server.name` to
  `com.sebastienrousseau/passmcp-server`.
- `make server-json` passes: the file validates against the schema it
  names.

## Publish

The name `com.sebastienrousseau/passmcp-server` is in a domain namespace,
so the registry grants it to whoever proves control of
`sebastienrousseau.com`. GitHub login is not enough: it grants only
`io.github.sebastienrousseau/*`
([registry authentication](https://github.com/modelcontextprotocol/registry/blob/main/docs/modelcontextprotocol-io/authentication.mdx)).

Install `mcp-publisher` (`brew install mcp-publisher`, or the binary from
the [registry's releases](https://github.com/modelcontextprotocol/registry/releases)).

Once, generate the signing key and publish its public half. The key stays
with the Maintainer and never enters this repository. Ed25519 needs
OpenSSL 3, not macOS's LibreSSL (`brew install openssl@3`):

```sh
openssl=/opt/homebrew/opt/openssl@3/bin/openssl
"$openssl" genpkey -algorithm Ed25519 -out key.pem
PUBLIC_KEY="$("$openssl" pkey -in key.pem -pubout -outform DER | tail -c 32 | base64)"
echo "sebastienrousseau.com. IN TXT \"v=MCPv1; k=ed25519; p=${PUBLIC_KEY}\""
```

Add that TXT record at the **apex** of `sebastienrousseau.com`, not under
a selector, and wait for it to propagate. Then, for each release, from the
repository root:

```sh
openssl=/opt/homebrew/opt/openssl@3/bin/openssl
PRIVATE_KEY="$("$openssl" pkey -in key.pem -noout -text | grep -A3 "priv:" | tail -n +2 | tr -d ' :\n')"
mcp-publisher login dns --domain sebastienrousseau.com --private-key "${PRIVATE_KEY}"
mcp-publisher publish
```

## After

Read the listing back from the registry rather than trusting the command's
output:

```sh
curl -fsSL "https://registry.modelcontextprotocol.io/v0.1/servers?search=com.sebastienrousseau/passmcp-server" -o listing.json
jq '.servers[].server | {name, version}' listing.json
```

The version must be the one just released.
