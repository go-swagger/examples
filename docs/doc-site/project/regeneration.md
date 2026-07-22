---
title: "Regeneration"
weight: 4
description: "How examples stay in sync with go-swagger"
---

The generated code in this repository is **committed**, yet it always reflects the
*current* go-swagger generator. That's the whole point of the examples: what you
read here is what `swagger generate` emits from `master` today, not a frozen
snapshot. A single tool keeps them in sync.

## Regenerating everything

```sh
go run ./hack/tools regen
```

The `regen` command drives the whole repository from one place —
[`hack/tools/regen.go`](https://github.com/go-swagger/examples/blob/master/hack/tools/regen.go)
holds a table with one entry per example: which directories to clean, and the exact
`swagger generate` command(s) to run. For each example it:

1. **ensures the generator is present** — if `swagger` isn't on your `PATH`, it
   installs it from source (`go install github.com/go-swagger/go-swagger/cmd/swagger@master`),
   pinning regeneration to the latest generator;
2. **cleans** the generated sub-directories (`models`, `restapi`, `client`, `cmd`, …)
   so nothing stale survives;
3. **runs** that example's generate command(s) — server, client, or both, with the
   flags and templates each example needs;
4. **restores** any preserved hand-written files that live inside a cleaned tree;
5. finally, once every example is regenerated, runs `go test ./...` across the whole
   module as a smoke test.

Because the command list is data, adding or changing an example's generation is a
one-line edit to that table — not a shell script to maintain.

## Auth material for runnable examples

A few examples need secrets to actually *run* (TLS certificates, JWT signing keys).
These are deliberately **not committed**. Generate them locally with the same tool:

```sh
go run ./hack/tools gen-certs    # self-signed TLS certs (todo-list-errors)
go run ./hack/tools gen-tokens   # RSA keypair + JWT tokens (composed-auth)
```

## Automated regeneration in CI

Regeneration also runs unattended, so the committed code never drifts from the
generator:

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `regen.yml` | Weekly (Mon 06:00 UTC) + manual | Regenerate from `swagger@master`; open an auto-merged PR if output changed |
| `go-test.yml` | PR, push to `master` | Lint + build matrix (2 Go versions × 3 OS) |
| `auto-merge.yml` | PR | Auto-approve/merge bot PRs (dependabot, scheduled regen) |
| `codeql.yml` | PR, push, weekly | CodeQL semantic analysis |
| `scanner.yml` | Push, weekly | Trivy + govulncheck vulnerability scans |
| `contributors.yml` | Weekly + manual | Refresh the all-time contributors list |

On top of the weekly job, go-swagger itself triggers a **cross-repo** pipeline:
PRs to go-swagger that touch code generation open a regeneration PR *here*
automatically, so a generator change and its effect on the examples are reviewed
together.

## Related

- [Contributing](../contributing/) — why you edit the spec or the regen table, never the generated files.
- [Repository README](../readme/) — how the examples are organized.
