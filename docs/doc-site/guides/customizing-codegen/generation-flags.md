---
title: "Generation flags"
weight: 3
description: "How the various generate flags materialize"
---

The generated `main.go` for a server isn't fixed — a couple of flags change how
it parses command-line options and whether it carries the spec inside the binary.
The **flags** example generates the *same* API six ways so you can compare the
results side by side.

{{% notice tip %}}
Source: [`flags/`](https://github.com/go-swagger/examples/tree/master/flags).
Six sub-packages (`pflag/`, `flag/`, `go-flags/` and their `x…` variants) are
each generated from one `swagger.yml` with a different flag combination.
{{% /notice %}}

## `--flag-strategy` — how the server parses CLI options

The flag strategy selects the library the generated `main.go` uses for its
command-line flags. The API is identical; only the flag plumbing differs:

| `--flag-strategy` | Library | Flag style |
|-------------------|---------|-----------|
| `go-flags` (default) | [`jessevdk/go-flags`](https://github.com/jessevdk/go-flags) | `--port=8080`, env-var bindings (`[$PORT]`), grouped options |
| `pflag` | [`spf13/pflag`](https://github.com/spf13/pflag) | GNU-style `--port 8080` |
| `flag` | stdlib `flag` | single-dash `-port 8080` |

```bash
(mkdir pflag && cd pflag && swagger generate server --spec=../swagger.yml --flag-strategy=pflag)
```

All three expose the same server options — listeners (`--scheme`), timeouts,
TLS settings, socket path — just rendered in each library's idiom. Pick the one
that matches the rest of your CLI.

## `--exclude-spec` — embedded vs. runtime spec

By default the spec is **embedded** in the generated binary (the `x…` variants
drop this):

- **default (embedded)** — the spec is baked in; the server is fully
  self-contained and serves its own `swagger.json`.
- **`--exclude-spec`** — the spec is *not* embedded. An extra `--spec` CLI flag
  appears so the server loads the document at startup instead.

Embed for a single shippable artifact; exclude when you want to swap the spec
without rebuilding, or to keep the binary small.

## Trying it

Build any variant's server and ask for help to see that strategy's flag layout:

```shellsession
$ go build -o srv ./flags/pflag/cmd/simple-to-do-list-api-server && ./srv -h
```

## Related

- [Custom templates](../custom-templates/) — change the generated code shape, not just its flags.
- [CLI client](../../clients-and-cli/cli-client/) — a different use of flags: a generated cobra command-line *client*.
