# Copilot Instructions

## Project Overview

Example projects illustrating how to use [go-swagger](https://github.com/go-swagger/go-swagger)
to generate server, client, and CLI code from OpenAPI 2.0 (Swagger) specs.
All generated code is committed and kept in sync via automated regeneration (`hack/regen-samples.sh`).

### Key dependencies

- `go-openapi/runtime` — runtime support for generated code
- `go-openapi/errors` — error types for generated code
- `go-openapi/loads` — spec loading
- `go-openapi/validate` — validation support
- `go-openapi/strfmt` — string format types
- `go-openapi/swag` — JSON/YAML utilities
- `jessevdk/go-flags` — CLI flag parsing (generated servers)
- `spf13/cobra` / `spf13/pflag` / `spf13/viper` — CLI generation (cobra-based)

## Conventions

Coding conventions are found beneath `.github/copilot/`

### Summary

- All `.go` files must have SPDX license headers (Apache-2.0).
- Commits require DCO sign-off (`git commit -s`).
- Linting: `golangci-lint run` — config in `.golangci.yml` (posture: `default: all` with explicit disables).
- Every `//nolint` directive **must** have an inline comment explaining why.
- Tests: `go test ./...` (single module).
- Generated code should not be hand-edited — changes must go through `hack/regen-samples.sh`.

See `.github/copilot/` (symlinked to `.claude/rules/`) for detailed rules on Go conventions, linting, testing, and contributions.
