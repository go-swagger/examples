---
title: "Task tracker"
weight: 6
description: "A CRUD task-tracker API"
---

The **task-tracker** example is a larger, realistic CRUD API — an issue tracker
with tasks, comments and file attachments. Its spec is the one go-swagger uses to
exercise code generation, so it deliberately packs in almost every construct:
nested resources, composition, arrays, file uploads, and multiple security
schemes. It's the best example to see what the generator produces for a
non-trivial contract.

{{% notice tip %}}
Source: [`task-tracker/`](https://github.com/go-swagger/examples/tree/master/task-tracker).
Generated with `swagger generate server --name TaskTracker --spec ./swagger.yml --principal any`.
{{% /notice %}}

## The API surface

The spec defines a full CRUD surface across three nested resources:

| Path | Operations |
|------|-----------|
| `/tasks` | `listTasks`, `createTask` |
| `/tasks/{id}` | `getTaskDetails`, `updateTask`, `deleteTask` |
| `/tasks/{id}/comments` | `addCommentToTask`, `getTaskComments` |
| `/tasks/{id}/files` | `uploadTaskFile` |

Each becomes a typed handler on the generated API, with parameters (path, query,
body, multipart) bound and validated for you.

## Model composition

The spec builds `Task` on top of a shared `TaskCard`, and the generator preserves
that with Go embedding — plus read-only fields, maps of attachments, and slices of
related models:

{{< code file="task-tracker/models/task.go" lang="go" lines="23-53" >}}

Read-only fields (`Comments`, `LastUpdated`) are populated by the server on
responses and ignored on input, exactly as the spec declares.

## Two API-key schemes

The spec declares two `apiKey` security definitions — one carried as a query
parameter, one as a header:

{{< code file="task-tracker/swagger.yml" lang="yaml" region="security" >}}

The generator turns each into an authenticator hook you implement. The scaffold
leaves them returning `NotImplemented`; you fill in the token validation:

{{< code file="task-tracker/restapi/configure_task_tracker.go" lang="go" region="auth" >}}

For worked authentication examples (basic, api-key, composed, OAuth2), see the
[Authentication guides](../../authentication/).

## Related

- [Todo list server](../todo-list/) — the minimal CRUD server to start from.
- [File server](../file-server/) — the `type: file` upload mechanics used by `/tasks/{id}/files`.
- [Petstore](../petstore/) — another full spec, generated end to end.
