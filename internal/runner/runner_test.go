// SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com>
// SPDX-License-Identifier: GPL-3.0-only

//go:build !windows

package runner

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fakePassmcp writes a stand-in passmcp program that records how it was called
// and answers the way the given body says.
func fakePassmcp(t *testing.T, body string) (bin, record string) {
	t.Helper()
	dir := t.TempDir()
	record = filepath.Join(dir, "record")
	bin = filepath.Join(dir, "passmcp")
	script := "#!/bin/sh\n" +
		"printf '%s\\n' \"$*\" > " + record + "\n" +
		"printf 'config=%s\\n' \"$(cat \"$PASSMCP_CONFIG\" 2>/dev/null)\" >> " + record + "\n" +
		body + "\n"
	if err := os.WriteFile(bin, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	return bin, record
}

func TestCheckFixesTheSafetySettings(t *testing.T) {
	bin, record := fakePassmcp(t, `echo '{"target":{"endpoint":"http://127.0.0.1/mcp"},"score":{"total":50,"grade":"D"},"phases":[{"name":"x","findings":[{"id":"a.b","status":"fail"}]}]}'; exit 2`)
	r, err := Exec{Path: bin}.Check(context.Background(), "http://127.0.0.1/mcp", []string{"net", "auth"})
	if err != nil {
		t.Fatal(err)
	}
	if r.Score.Total != 50 || len(r.Failures()) != 1 {
		t.Errorf("report = %+v", r)
	}
	b, _ := os.ReadFile(record)
	got := string(b)
	for _, want := range []string{"check http://127.0.0.1/mcp", "--auth none", "--output json", "--phases net,auth", "config={}"} {
		if !strings.Contains(got, want) {
			t.Errorf("passmcp was not called with %q:\n%s", want, got)
		}
	}
	for _, never := range []string{"--allow-mutations", "--allow-destructive", "--token"} {
		if strings.Contains(got, never) {
			t.Errorf("passmcp was called with %s", never)
		}
	}
}

func TestCheckWithoutAReport(t *testing.T) {
	bin, _ := fakePassmcp(t, `echo "fatal: bad flag" >&2; exit 1`)
	if _, err := (Exec{Path: bin}).Check(context.Background(), "http://127.0.0.1/mcp", nil); !errors.Is(err, ErrNoReport) || !strings.Contains(err.Error(), "bad flag") {
		t.Errorf("err = %v", err)
	}
	bin, _ = fakePassmcp(t, `echo "not json"; exit 2`)
	if _, err := (Exec{Path: bin}).Check(context.Background(), "http://127.0.0.1/mcp", nil); !errors.Is(err, ErrNoReport) {
		t.Errorf("err = %v", err)
	}
	if _, err := (Exec{Path: filepath.Join(t.TempDir(), "absent")}).Check(context.Background(), "http://127.0.0.1/mcp", nil); !errors.Is(err, ErrNoReport) {
		t.Errorf("err = %v", err)
	}
}

func TestVersion(t *testing.T) {
	bin, _ := fakePassmcp(t, `echo "passmcp 0.0.5"`)
	if v, err := (Exec{Path: bin}).Version(context.Background()); err != nil || v != "passmcp 0.0.5" {
		t.Errorf("%q %v", v, err)
	}
	if _, err := (Exec{Path: filepath.Join(t.TempDir(), "absent")}).Version(context.Background()); err == nil {
		t.Error("a missing passmcp reported a version")
	}
	if (Exec{}).bin() != "passmcp" {
		t.Error("the default is not passmcp on PATH")
	}
}
