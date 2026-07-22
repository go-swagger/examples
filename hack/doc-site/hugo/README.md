# Hugo documentation site

This directory holds the Hugo configuration that builds
<https://goswagger.io/examples/>.

## Layout

```text
hugo/
├── hugo.yaml                  # Static Hugo config
├── examples.yaml.template      # Build-time config template (version info)
├── examples.yaml               # Generated from the template (git-ignored output)
├── gendoc.go                  # Local development helper (`go run gendoc.go`)
├── themes/
│   ├── hugo-relearn/          # Relearn theme (downloaded by CI / dev script)
│   ├── examples-assets/        # Custom logo / SCSS
│   └── examples-static/        # Static branding (favicon, …)
└── layouts/
    ├── shortcodes/            # Custom Hugo shortcodes
    └── partials/              # Custom partial templates
```

## Content

Markdown content is mounted from `../../../docs/doc-site/` via the
`module.mounts` block in `hugo.yaml`. Editing those files (or adding new ones)
is enough — no codegen or generator is involved.

## Local preview

```sh
go run gendoc.go
```

The script:

1. Extracts version info from git tags and the root `go.mod`
2. Renders `examples.yaml` from `examples.yaml.template`
3. Starts `hugo server` on <http://localhost:1313/examples/> with live reload

Requires `hugo` (extended, ≥ v0.150) and `git` on `PATH`.

## Configuration

Two-layer config, mirroring the pattern used by other go-openapi doc sites:

1. **`hugo.yaml`** — static configuration (theme, mounts, menu, params)
2. **`examples.yaml`** — dynamic configuration (Go version, latest release tag,
   build timestamp), generated from `examples.yaml.template`

Both files are passed together via `--config hugo.yaml,examples.yaml`. The
dynamic values land under `params.examples.*` and are referenced from the
markdown content.

## Deployment

GitHub Actions workflow `.github/workflows/update-doc.yml`:

- Builds on every push to `master` and on tags `v*` that touch `docs/**`,
  `hack/doc-site/**`, or the workflow itself
- Publishes the rendered site to GitHub Pages
  (<https://goswagger.io/examples/>)
