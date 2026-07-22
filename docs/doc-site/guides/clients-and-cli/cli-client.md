---
title: "CLI client"
weight: 2
description: "Generate a command-line client tool"
---

`swagger generate cli` produces a command-line tool that wraps the generated
client: it reads flags and arguments, builds the operation parameters, calls the
server, and prints the response. It's built on [cobra](https://github.com/spf13/cobra)
and [viper](https://github.com/spf13/viper), with shell-completion support.

{{% notice tip %}}
Source: [`cli/`](https://github.com/go-swagger/examples/tree/master/cli).
Generated with `swagger generate cli --spec ./swagger.yml --cli-app-name todoctl`.
It targets the same spec as [auto-configure](../../servers/auto-configure/), so you
can run that server and drive it from this CLI.
{{% /notice %}}

## Command layout

The generated command tree mirrors the spec:

- the **root** command holds global flags (`--hostname`, `--scheme`, auth tokens);
- each **tag** becomes a sub-command (an *operation group*);
- each **operationId** becomes a sub-command under its tag;
- each path/query parameter becomes a flag; the body becomes a `--body` JSON flag,
  with a flag per body field layered on top.

The tag → sub-command mapping is a `cobra.Command` per group, wiring in one child
per operation:

{{< code file="cli/cli/cli.go" lang="go" lines="209-240" >}}

## An operation command

Each operationId gets its own command with a `RunE` that calls the server, plus
generated flag registration for its parameters:

{{< code file="cli/cli/add_one_operation.go" lang="go" lines="17-29" >}}

Body parameters are handled two ways at once — a whole-body `--body` JSON string
as a base payload, and a generated flag per field (recursing into sub-definitions)
that overrides it. That's where `--item.description` comes from:

{{< code file="cli/cli/add_one_operation.go" lang="go" lines="67-84" >}}

## Running it

Drive the [auto-configure](../../servers/auto-configure/) server with the tool:

```shellsession
$ go run ./cli/cmd/todoctl/main.go --hostname localhost:12345 \
    --x-todolist-token "example token" \
    todos addOne --item.description "hi" --body "{}"
{"description":"hi"}
```

The path is `todoctl` → `todos` (the tag) → `addOne` (the operationId), with
`--item.description` setting a body field.

## Config files and completion

Common flags — `hostname`, `scheme`, `base_path`, auth tokens — can live in a
config file instead of the command line, loaded via viper from
`~/.config/<app>/config.json` (or `--config`, in JSON/YAML/env form):

```json
{
    "hostname": "localhost:12345",
    "scheme": "http",
    "x-todolist-token": "example token"
}
```

Shell completions (bash, zsh, fish, PowerShell) come for free from cobra:

```shellsession
$ source <(./todoctl completion bash)
```

{{% notice note %}}
The CLI generator is under active development. A few spec shapes aren't covered
yet — arrays/maps in a body, and enums in help text and completions.
{{% /notice %}}

## Related

- [Generated client SDK](../generated-client/) — the typed client this CLI wraps.
- [dockerctl](https://github.com/go-swagger/dockerctl) — a full CLI generated this way for the Docker Engine API.
