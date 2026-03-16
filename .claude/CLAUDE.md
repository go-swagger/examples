# CLAUDE.md — go-swagger

## Project Overview

Examples that illustrate how to use `go-swagger`

### Key packages

| Package | Purpose |
|---------|---------|

### Key dependencies

| Dependency | Role |
|------------|------|

## Conventions

Coding conventions are found beneath `.claude/rules/` (also symlinked as `.github/copilot/`).

### Summary

- All `.go` files must have SPDX license headers (Apache-2.0).
- Commits require DCO sign-off (`git commit -s`).
- Go version policy: support the 2 latest stable Go minor versions.
- Linting: `golangci-lint run` — config in `.golangci.yml` (posture: `default: all` with explicit disables).
- Every `//nolint` directive **must** have an inline comment explaining why.
- Tests: `go test ./...` (single module). CI runs on `{ubuntu, macos, windows} x {stable, oldstable}` with `-race`.
- Test framework: `github.com/go-openapi/testify/v2` (not `stretchr/testify`; `testifylint` does not work).

See `.claude/rules/` for detailed rules on Go conventions, linting, testing, GitHub Actions, and contributions.
See `.github/STYLE.md` for the linting posture rationale.

## CI / Release pipeline

| Workflow | Trigger | Purpose |
|----------|---------|---------|

