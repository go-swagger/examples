// SPDX-FileCopyrightText: Copyright 2015-2026 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

// Command tools provides development utilities for the go-swagger examples repository.
//
// Usage:
//
//	go run ./hack/tools <command> [flags]
//
// Commands:
//
//	regen       Regenerate all examples from their specs
//	gen-certs   Generate self-signed TLS certificates for todo-list-errors
//	gen-tokens  Generate RSA keypair and JWT tokens for composed-auth
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	var err error

	switch os.Args[1] {
	case "regen":
		err = runRegen()
	case "gen-certs":
		err = runGenCerts()
	case "gen-tokens":
		err = runGenTokens()
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		usage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `Usage: go run ./hack/tools <command> [flags]

Commands:
  regen       Regenerate all examples from their specs
  gen-certs   Generate self-signed TLS certificates for todo-list-errors
  gen-tokens  Generate RSA keypair and JWT tokens for composed-auth`)
}
