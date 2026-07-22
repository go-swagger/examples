---
title: "Contributing"
weight: 3
description: "How to contribute examples"
---

Contributions are always welcome — and not just code. Reporting issues, improving
docs, triaging bugs, and adding test coverage all help. These guidelines are the
standard ones shared across every `go-openapi` and `go-swagger` repository; if
you've contributed to a Go project on GitHub before, you'll feel at home.

{{% notice tip %}}
Authoritative sources:
[`.github/CONTRIBUTING.md`](https://github.com/go-swagger/examples/blob/master/.github/CONTRIBUTING.md)
and [`docs/STYLE.md`](https://github.com/go-swagger/examples/blob/master/docs/STYLE.md).
This page summarizes the essentials.
{{% /notice %}}

## Git flow

Fork the repo, branch from `master`, and open a pull request from your fork.
Branch naming is not enforced (it's your fork), but the common convention is
`fix/XXX-something` or `feature/XXX-something`, where `XXX` is the issue number.

Keep pull requests **focused** — small, single-purpose PRs are reviewed faster and
are less likely to lose the thread than large ones.

## A special note for generated code

Most of this repository is **generated** and must not be hand-edited — a
[regeneration](../regeneration/) would silently overwrite your change. If your
contribution affects generated output:

- change the **spec** or the **generation command** (in `hack/tools/regen.go`), not
  the generated `.go` files;
- run `go run ./hack/tools regen` and commit the regenerated result;
- hand-written glue (the `configure_*.go` files, custom handlers, `main.go` in the
  custom-server example) *is* editable — that's the code these examples exist to
  illustrate.

## Tests

Submit unit tests for your changes and run the full suite before opening a PR:

```sh
go test ./...
```

CI measures patch coverage; aim for at least 80% of your change. It's an indicator
maintainers weigh, not a hard gate.

## Code style & linting

The project runs the `golangci-lint` meta-linter with a deliberate posture:
**`default: all`, then disable what doesn't earn its keep**. The disabled list in
[`.golangci.yml`](https://github.com/go-swagger/examples/blob/master/.golangci.yml)
is a design rationale, not technical debt. Two rules matter most when contributing:

- every `//nolint` directive **must** carry an inline comment explaining why;
- prefer disabling a linter globally over scattering `//nolint` — if a linter
  fights an intentional pattern, the linter goes, not the code.

Run it (and the formatter) before committing:

```sh
golangci-lint run
golangci-lint fmt
```

## Sign your work (DCO)

Every commit must be **signed off** under the
[Developer Certificate of Origin](https://developercertificate.org), using your
real name and email:

```
Signed-off-by: Joe Smith <joe@example.com>
```

Add it automatically with `git commit -s`. PGP-signed commits are appreciated but
not required. Squash your commits into logical units (`git rebase -i`) before
requesting review.

## AI agents

Agentic contributors are welcome, with a few rules:

1. Issues and PRs written or posted by an agent should mention the original
   **human** poster for reference.
2. PRs must **not** be attributed to an agent as author — no commits authored by
   `@claude.code` or similar. Agents and bots may be listed as `Co-Authored-By:`;
   the commit author must be the human sponsor.
3. Security reports produced by an agent must be filed **privately** (see the
   [security policy](https://github.com/go-swagger/examples/blob/master/SECURITY.md))
   and mention the human poster.
