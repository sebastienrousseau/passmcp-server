// SPDX-FileCopyrightText: 2026 Sebastien Rousseau <sebastian.rousseau@gmail.com>
// SPDX-License-Identifier: GPL-3.0-only

// Command passmcp-server is an MCP server over stdio that exposes passmcp's
// diagnostics as tools. An MCP host starts it as a child process:
//
//	{"command": "passmcp-server", "args": ["--allow", ".internal.example.com"]}
//
// It evaluates only loopback endpoints unless --allow (or PASSMCP_SERVER_ALLOW)
// names more hosts, sends no credentials, and runs the passmcp program found
// on PATH or at --passmcp.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"satellion.com/passmcp-server/internal/runner"
	"satellion.com/passmcp-server/internal/server"
)

// Version is stamped by the release build; a local build says dev.
var Version = "dev"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "passmcp-server:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, in io.Reader, out, errw io.Writer) error {
	fs := flag.NewFlagSet("passmcp-server", flag.ContinueOnError)
	fs.SetOutput(errw)
	allow := fs.String("allow", os.Getenv("PASSMCP_SERVER_ALLOW"), "hosts passmcp may be pointed at besides loopback, comma-separated; a leading dot allows subdomains")
	passmcpPath := fs.String("passmcp", os.Getenv("PASSMCP_SERVER_PASSMCP"), "path of the passmcp program; empty means passmcp on PATH")
	version := fs.Bool("version", false, "print the version and exit")
	completion := fs.String("completion", "", "print a completion script for bash, zsh or fish, and exit")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *completion != "" {
		return writeCompletion(out, fs, *completion)
	}
	if *version {
		_, err := fmt.Fprintln(out, "passmcp-server", Version)
		return err
	}
	s := &server.Server{
		Version: Version,
		Runner:  runner.Exec{Path: *passmcpPath},
		Allow:   server.ParseAllowlist(*allow),
	}
	return s.Serve(ctx, in, out)
}
