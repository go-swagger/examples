# CLAUDE.md — go-swagger examples

## Project Overview

Example projects illustrating how to use [go-swagger](https://github.com/go-swagger/go-swagger)
to generate server, client, and CLI code from OpenAPI 2.0 (Swagger) specs.

All generated code is committed and kept in sync with the latest `go-swagger` via
automated regeneration (see [CI pipelines](#ci-pipelines)).

### Example projects

| Directory | Purpose |
|-----------|---------|
| `2.0/` | Petstore-style examples (JSON, YAML) |
| `authentication/` | Basic and API-key auth |
| `auto-configure/` | Auto-configured server |
| `cli/` | CLI generation examples |
| `composed-auth/` | Multiple auth schemes |
| `contributed-templates/` | Custom template usage (stratoscale) |
| `external-types/` | External type references |
| `file-server/` | File upload/download |
| `flags/` | Custom CLI flags |
| `generated/` | Generated client/server showcase |
| `oauth2/` | OAuth2 flow examples |
| `stream-client/` / `stream-server/` | Streaming support |
| `task-tracker/` | Task tracker CRUD |
| `todo-list/` / `todo-list-errors/` / `todo-list-strict/` | Todo list variants |
| `tutorials/` | Step-by-step tutorials |

### Key dependencies

| Dependency | Role |
|------------|------|
| `go-openapi/runtime` | Runtime support for generated code |
| `go-openapi/errors` | Error types for generated code |
| `go-openapi/loads` | Spec loading |
| `go-openapi/validate` | Validation support |
| `go-openapi/strfmt` | String format types |
| `go-openapi/swag` | JSON/YAML utilities |
| `jessevdk/go-flags` | CLI flag parsing (generated servers) |
| `spf13/cobra` / `spf13/pflag` / `spf13/viper` | CLI generation (cobra-based) |

### Regeneration

Examples are regenerated from go-swagger specs using `hack/regen-samples.sh`.
This script invokes `swagger generate` for each example and runs `go test ./...` as a smoke test.

## Conventions

Coding conventions are found beneath `.claude/rules/` (also symlinked as `.github/copilot/`).

### Summary

- All `.go` files must have SPDX license headers (Apache-2.0).
- Commits require DCO sign-off (`git commit -s`).
- Go version policy: support the 2 latest stable Go minor versions.
- Linting: `golangci-lint run` — config in `.golangci.yml` (posture: `default: all` with explicit disables).
- Every `//nolint` directive **must** have an inline comment explaining why.
- Tests: `go test ./...` (single module).
- Generated code should not be hand-edited — changes must go through `hack/regen-samples.sh`.

See `.claude/rules/` for detailed rules on Go conventions, linting, testing, GitHub Actions, and contributions.

## CI pipelines

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `test.yml` | PR, push to master | Lint + build matrix (2 Go versions × 3 OS) |
| `regen.yml` | Weekly (Monday 06:00 UTC), manual | Scheduled regeneration from `swagger@master`, auto-merged via bot |
| `auto-merge.yml` | PR | Auto-approve/merge dependabot and scheduled bot PRs (reuses `go-openapi/ci-workflows`) |
| `contributors.yml` | Weekly (Saturday 04:18 UTC), manual | Update all-time contributors list (reuses `go-openapi/ci-workflows`) |
| `codeql.yml` | PR to master, push to master, weekly | GitHub CodeQL semantic analysis |
| `scanner.yml` | Push to master, weekly, branch protection | Trivy + govulncheck vulnerability scans |

Additionally, go-swagger PRs that touch codegen trigger a cross-repo regen pipeline
that creates PRs in this repo automatically (see go-swagger's `regen-examples.yaml`).
