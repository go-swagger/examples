// SPDX-FileCopyrightText: Copyright 2015-2026 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// regenStep describes one example regeneration.
type regenStep struct {
	// dir is the example directory relative to the repo root.
	dir string

	// preserve lists files to save before cleaning and restore after generation.
	// Paths are relative to dir.
	preserve []string

	// clean lists directories to remove before regeneration.
	// Paths are relative to dir.
	clean []string

	// mkdirs lists directories to create before generation (after clean).
	mkdirs []string

	// commands lists swagger (or other) commands to run in dir.
	commands [][]string
}

var steps = []regenStep{
	{
		dir:   "generated",
		clean: []string{"cmd", "models", "restapi"},
		commands: [][]string{
			{"swagger", "generate", "server", "-f", "swagger-petstore.json", "-A", "Petstore"},
		},
	},
	{
		dir:   "todo-list",
		clean: []string{"client", "cmd", "models", "restapi"},
		commands: [][]string{
			{"swagger", "generate", "client", "-A", "TodoList", "-f", "./swagger.yml"},
			{"swagger", "generate", "server", "-A", "TodoList", "-f", "./swagger.yml", "--flag-strategy", "pflag"},
		},
	},
	{
		dir:   "authentication",
		clean: []string{"client", "cmd", "models", "restapi"},
		commands: [][]string{
			{"swagger", "generate", "client", "-A", "AuthSample", "-f", "./swagger.yml", "-P", "models.Principal"},
			{"swagger", "generate", "server", "-A", "AuthSample", "-f", "./swagger.yml", "-P", "models.Principal"},
		},
	},
	{
		dir:   "task-tracker",
		clean: []string{"client", "cmd", "models", "restapi"},
		commands: [][]string{
			{"swagger", "generate", "client", "-A", "TaskTracker", "-f", "./swagger.yml"},
			{"swagger", "generate", "server", "-A", "TaskTracker", "-f", "./swagger.yml"},
		},
	},
	{
		dir:      "stream-server",
		preserve: []string{"restapi/configure_countdown.go"},
		clean:    []string{"cmd", "models", "restapi"},
		commands: [][]string{
			{"swagger", "generate", "server", "-A", "Countdown", "-f", "./swagger.yml"},
			{"swagger", "generate", "client", "-f", "swagger.yml", "--skip-models"},
		},
	},
	{
		dir:      "oauth2",
		preserve: []string{"restapi/configure_oauth_sample.go", "restapi/implementation.go"},
		clean:    []string{"cmd", "models", "restapi"},
		commands: [][]string{
			{"swagger", "generate", "server", "-A", "oauthSample", "-P", "models.Principal", "-f", "./swagger.yml"},
		},
	},
	{
		dir:   "tutorials/todo-list/server-1",
		clean: []string{"cmd", "models", "restapi"},
		commands: [][]string{
			{"swagger", "generate", "server", "-A", "TodoList", "-f", "./swagger.yml"},
		},
	},
	{
		dir:   "tutorials/todo-list/server-2",
		clean: []string{"cmd", "models", "restapi"},
		commands: [][]string{
			{"swagger", "generate", "server", "-A", "TodoList", "-f", "./swagger.yml"},
		},
	},
	{
		dir:      "tutorials/todo-list/server-complete",
		preserve: []string{"restapi/configure_todo_list.go"},
		clean:    []string{"cmd", "models", "restapi"},
		commands: [][]string{
			{"swagger", "generate", "server", "-A", "TodoList", "-f", "./swagger.yml"},
		},
	},
	{
		dir:    "tutorials/custom-server",
		clean:  []string{"gen"},
		mkdirs: []string{"gen"},
		commands: [][]string{
			{"swagger", "generate", "server", "--exclude-main", "-A", "greeter", "-t", "gen", "-f", "./swagger/swagger.yml"},
		},
	},
	{
		dir:      "composed-auth",
		preserve: []string{"restapi/configure_multi_auth_example.go"},
		clean:    []string{"cmd", "models", "restapi"},
		commands: [][]string{
			{"swagger", "generate", "server", "-A", "multi-auth-example", "-P", "models.Principal", "-f", "./swagger.yml"},
		},
	},
	{
		dir:   "contributed-templates/stratoscale",
		clean: []string{"client", "cmd", "models", "restapi"},
		commands: [][]string{
			{"swagger", "generate", "client", "-A", "Petstore", "--template", "stratoscale"},
			{"swagger", "generate", "server", "-A", "Petstore", "--template", "stratoscale"},
		},
	},
	{
		dir:      "external-types",
		preserve: []string{"models/my_type.go"},
		clean:    []string{"cmd", "models", "restapi"},
		mkdirs:   []string{"models"},
		commands: [][]string{
			{"swagger", "generate", "server", "--skip-validation", "-f", "example-external-types.yaml", "-A", "external-types-demo"},
		},
	},
	{
		dir:   "stream-client",
		clean: []string{"client"},
		commands: [][]string{
			{"swagger", "generate", "client"},
		},
	},
	{
		dir:      "file-server",
		preserve: []string{"restapi/configure_file_upload.go"},
		clean:    []string{"client", "cmd", "restapi"},
		commands: [][]string{
			{"swagger", "generate", "server"},
			{"swagger", "generate", "client"},
		},
	},
	{
		dir:   "cli",
		clean: []string{"cli", "client", "cmd", "models"},
		commands: [][]string{
			{"swagger", "generate", "cli", "--spec=swagger.yml", "--cli-app-name", "todoctl"},
		},
	},
	{
		dir:   "auto-configure",
		clean: []string{"cmd", "models", "restapi"},
		commands: [][]string{
			{"swagger", "generate", "server", "--spec=swagger.yml",
				"--implementation-package=github.com/go-swagger/examples/auto-configure/implementation"},
		},
	},
	{
		dir:   "todo-list-strict",
		clean: []string{"cmd", "models", "restapi", "client"},
		commands: [][]string{
			{"swagger", "generate", "server", "--spec=swagger.yml", "--strict-responders"},
			{"swagger", "generate", "client", "--spec=swagger.yml", "--strict-responders"},
		},
	},
	{
		dir:   "flags",
		clean: []string{"pflag", "flag", "go-flags", "xpflag", "xflag", "xgo-flags"},
		mkdirs: []string{
			"pflag", "flag", "go-flags",
			"xpflag", "xflag", "xgo-flags",
		},
		commands: [][]string{
			{"swagger", "generate", "server", "--spec=swagger.yml", "--flag-strategy=pflag", "-t", "pflag"},
			{"swagger", "generate", "server", "--spec=swagger.yml", "--flag-strategy=flag", "-t", "flag"},
			{"swagger", "generate", "server", "--spec=swagger.yml", "--flag-strategy=go-flags", "-t", "go-flags"},
			{"swagger", "generate", "server", "--spec=swagger.yml", "--flag-strategy=pflag", "--exclude-spec", "-t", "xpflag"},
			{"swagger", "generate", "server", "--spec=swagger.yml", "--flag-strategy=flag", "--exclude-spec", "-t", "xflag"},
			{"swagger", "generate", "server", "--spec=swagger.yml", "--flag-strategy=go-flags", "--exclude-spec", "-t", "xgo-flags"},
		},
	},
	{
		dir:   "tutorials/client",
		clean: []string{"client", "models", "stratoscale-client"},
		commands: [][]string{
			{"swagger", "generate", "client", "-A", "TodoList", "--spec=swagger.yml", "--client-package=classic-client"},
			{"swagger", "generate", "client", "-A", "TodoList", "--spec=swagger.yml", "--template", "stratoscale",
				"--existing-models=github.com/go-swagger/examples/tutorials/client/models",
				"--client-package=stratoscale-client"},
		},
	},
}

func runRegen() error {
	// Ensure swagger is available
	swaggerPath, err := exec.LookPath("swagger")
	if err != nil {
		fmt.Println("go-swagger is not present: installing from master")
		cmd := exec.Command("go", "install", "github.com/go-swagger/go-swagger/cmd/swagger@master")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("installing swagger: %w", err)
		}
	} else {
		fmt.Printf("Using swagger at %s\n", swaggerPath)
		run(".", "swagger", "version")
	}

	root, err := repoRoot()
	if err != nil {
		return err
	}

	for _, step := range steps {
		dir := filepath.Join(root, step.dir)
		fmt.Printf("\n=== %s ===\n", step.dir)

		// Save preserved files
		saved := make(map[string][]byte)
		for _, p := range step.preserve {
			data, readErr := os.ReadFile(filepath.Join(dir, p))
			if readErr != nil {
				return fmt.Errorf("preserving %s/%s: %w", step.dir, p, readErr)
			}
			saved[p] = data
		}

		// Clean
		for _, d := range step.clean {
			os.RemoveAll(filepath.Join(dir, d))
		}

		// Create directories
		for _, d := range step.mkdirs {
			os.MkdirAll(filepath.Join(dir, d), 0o755)
		}

		// Run commands
		for _, args := range step.commands {
			if err := run(dir, args[0], args[1:]...); err != nil {
				return fmt.Errorf("%s: %s: %w", step.dir, strings.Join(args, " "), err)
			}
		}

		// Restore preserved files
		for p, data := range saved {
			dest := filepath.Join(dir, p)
			os.MkdirAll(filepath.Dir(dest), 0o755)
			if err := os.WriteFile(dest, data, 0o644); err != nil {
				return fmt.Errorf("restoring %s/%s: %w", step.dir, p, err)
			}
		}
	}

	// Final build + test
	fmt.Println("\n=== go test ./... ===")

	return run(root, "go", "test", "-v", "./...")
}

func run(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func repoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("finding repo root: %w", err)
	}

	return strings.TrimSpace(string(out)), nil
}
