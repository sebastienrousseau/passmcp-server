// SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com>
// SPDX-License-Identifier: GPL-3.0-only

package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"testing"

	"satellion.com/passmcp-server/internal/runner"
)

// recordingRunner remembers every endpoint the server hands passmcp.
type recordingRunner struct{ endpoints []string }

func (r *recordingRunner) Check(_ context.Context, endpoint string, _ []string) (runner.Report, error) {
	r.endpoints = append(r.endpoints, endpoint)
	return runner.Report{}, nil
}

func (r *recordingRunner) Version(context.Context) (string, error) { return "passmcp 0.0.0", nil }

// FuzzServe feeds arbitrary bytes to the stdio loop, standing in for
// whatever a client or a prompt-driven agent writes. Whatever arrives:
// Serve does not panic or fail, stdout carries only JSON-RPC 2.0
// responses (one result or one error each, never more answers than
// lines), and passmcp is never pointed at an endpoint the default
// allowlist refuses.
func FuzzServe(f *testing.F) {
	for _, seed := range []string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18"}}`,
		`{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n" + `{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
		`{"jsonrpc":"2.0","id":"a","method":"server/discover"}`,
		call("passmcp_check", map[string]any{"endpoint": "http://127.0.0.1:3000/mcp", "phases": []string{"net"}}),
		call("passmcp_check", map[string]any{"endpoint": "https://user:pw@evil.example/mcp"}),
		call("passmcp_check", map[string]any{"endpoint": "http://localhost.evil.example/mcp"}),
		call("passmcp_verify_attestation", map[string]any{"statement": signedStatement(f), "endpoint": "https://mcp.example.com/mcp"}),
		call("passmcp_verify_attestation", map[string]any{"statement": `{"_type":1}`}),
		call("passmcp_version", map[string]any{"extra": true}),
		`{"jsonrpc":"2.0","id":null,"method":"ping"}`,
		`not json` + "\n\n" + `{"jsonrpc":"1.0","id":3,"method":"ping"}`,
	} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		rec := &recordingRunner{}
		s := &Server{Version: "0.0.0", Runner: rec}
		var out bytes.Buffer
		if err := s.Serve(context.Background(), bytes.NewReader(data), &out); err != nil {
			t.Fatalf("Serve: %v", err)
		}
		checkResponses(t, data, out.Bytes())
		for _, e := range rec.endpoints {
			if err := (Allowlist{}).Check(e); err != nil {
				t.Fatalf("passmcp was pointed at a refused endpoint %q: %v", e, err)
			}
		}
	})
}

// checkResponses asserts that every line of out is one JSON-RPC 2.0
// response, and that there are no more of them than lines of in.
func checkResponses(t *testing.T, in, out []byte) {
	t.Helper()
	lines := bytes.Split(bytes.TrimSuffix(out, []byte("\n")), []byte("\n"))
	if len(out) == 0 {
		lines = nil
	}
	if most := bytes.Count(in, []byte("\n")) + 1; len(lines) > most {
		t.Fatalf("%d responses to %d lines", len(lines), most)
	}
	for _, l := range lines {
		var r struct {
			JSONRPC string          `json:"jsonrpc"`
			ID      json.RawMessage `json:"id"`
			Result  json.RawMessage `json:"result"`
			Error   json.RawMessage `json:"error"`
		}
		if err := json.Unmarshal(l, &r); err != nil {
			t.Fatalf("stdout carried a line that is not JSON: %q", l)
		}
		if r.JSONRPC != "2.0" || len(r.ID) == 0 || (len(r.Result) == 0) == (len(r.Error) == 0) {
			t.Fatalf("not a JSON-RPC 2.0 response with one result or error: %s", l)
		}
	}
}

// FuzzAllowlist checks the allowlist against arbitrary entries and
// endpoints. ParseAllowlist keeps only trimmed, lower-case, non-empty
// hosts; an accepted endpoint is always an http or https URL with a host
// and no credentials; widening the list never refuses what the default
// accepts; and whatever the widened list alone accepts is named by one of
// its entries, exactly or as a dot-suffix.
func FuzzAllowlist(f *testing.F) {
	for _, seed := range [][2]string{
		{"", "http://127.0.0.1:3000/mcp"},
		{"", "http://[::1]/mcp"},
		{"", "http://app.localhost/mcp"},
		{"mcp.example.com", "https://mcp.example.com/mcp"},
		{".example.com", "https://a.b.example.com/mcp"},
		{".example.com", "https://example.com/mcp"},
		{" Example.COM , ,x", "https://EXAMPLE.com/mcp"},
		{"", "https://user:pw@127.0.0.1/mcp"},
		{"", "ftp://127.0.0.1/"},
		{"evil.example", "http://localhost.evil.example/"},
	} {
		f.Add(seed[0], seed[1])
	}
	f.Fuzz(func(t *testing.T, allow, endpoint string) {
		a := parsedAllowlist(t, allow)
		byDefault := (Allowlist{}).Check(endpoint) == nil
		widened := a.Check(endpoint) == nil
		if byDefault && !widened {
			t.Fatalf("%q: accepted by default, refused by %v", endpoint, a.Hosts)
		}
		if !widened {
			return
		}
		host := acceptedHost(t, endpoint)
		if !byDefault && !named(a.Hosts, host) {
			t.Fatalf("%q accepted, but %q is not named by %v", endpoint, host, a.Hosts)
		}
	})
}

// parsedAllowlist parses allow and fails t if any host it kept is empty,
// untrimmed, not lower case, or still holds a separator.
func parsedAllowlist(t *testing.T, allow string) Allowlist {
	t.Helper()
	a := ParseAllowlist(allow)
	for _, h := range a.Hosts {
		if h == "" || h != strings.ToLower(strings.TrimSpace(h)) || strings.Contains(h, ",") {
			t.Fatalf("ParseAllowlist(%q) kept %q", allow, h)
		}
	}
	return a
}

// acceptedHost fails t unless endpoint is an http or https URL with a
// host and no credentials, and returns its host in lower case.
func acceptedHost(t *testing.T, endpoint string) string {
	t.Helper()
	u, err := url.Parse(endpoint)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" || u.User != nil {
		t.Fatalf("accepted %q, which is not a credential-free http(s) URL with a host", endpoint)
	}
	return strings.ToLower(u.Hostname())
}

func named(hosts []string, host string) bool {
	for _, h := range hosts {
		if host == h || (strings.HasPrefix(h, ".") && strings.HasSuffix(host, h)) {
			return true
		}
	}
	return false
}
