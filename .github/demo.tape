# SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com>
# SPDX-License-Identifier: GPL-3.0-only
#
# The README demo: passmcp-server's read-only tools answering over stdio,
# called the way an MCP host would, with `passmcp call`. passmcp_check
# evaluates passmcp's example toxic-pair server on loopback. Rendered by
# `make demo`, which installs passmcp and the example server at the
# lockstep release, builds passmcp-server, and puts all three on PATH.
Output .github/demo.gif
Require passmcp
Require passmcp-server
Require servers
Set Shell "bash"
Set FontSize 15
Set Width 1200
Set Height 280
Set Padding 20
Set Theme "Catppuccin Mocha"
Set TypingSpeed 40ms

Hide
# The example server stops itself after 90s: a hidden kill would still
# print a line.
Type "PS1='$ '; timeout 90 servers -flaw toxic-pair -addr 127.0.0.1:17802 >/dev/null 2>&1 &" Enter
# Let the set-up finish and clear it before recording starts.
Sleep 2s
Type "clear" Enter
Sleep 500ms
Show

Type "passmcp call --stdio passmcp_version -- passmcp-server" Sleep 400ms Enter
Sleep 3s
Type "passmcp call --stdio passmcp_check --arg endpoint=http://127.0.0.1:17802/mcp -- passmcp-server" Sleep 400ms Enter
# A full evaluation takes about 18s, most of it the performance phase.
Sleep 23s
