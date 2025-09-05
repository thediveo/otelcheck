# `otelcheck`

[![PkgGoDev](https://img.shields.io/badge/-reference-blue?logo=go&logoColor=white&labelColor=505050)](https://pkg.go.dev/github.com/thediveo/otelcheck)
[![License](https://img.shields.io/github/license/thediveo/otelcheck)](https://img.shields.io/github/license/thediveo/otelcheck)
![Coverage](https://img.shields.io/badge/Coverage-95.3%25-brightgreen)

`otelcheck` (voiced with a proper [H
aspiré](https://en.wikipedia.org/wiki/Aspirated_h) and tongue in cheek) provides
[Gomega matchers](https://onsi.github.io/gomega/) for reasoning about
[OpenTelemetry signals](https://opentelemetry.io/docs/concepts/signals/) in unit
tests. This package is opinionated in that its OTel-specific DSL hides the
complex and error-prone “structure field and method chasing” of the OTel
SDK-internal data models so that test spec writers are able to focus on their
(flat) mental model again.

That is, `otelcheck` somewhat _undoes_ OTel's observability models _overdoing_.

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md).

## DevContainer

> [!CAUTION]
>
> Do **not** use VSCode's "~~Dev Containers: Clone Repository in Container
> Volume~~" command, as it is utterly broken by design, ignoring
> `.devcontainer/devcontainer.json`.

1. `git clone https://github.com/thediveo/otelcheck`
2. in VSCode: Ctrl+Shift+P, "Dev Containers: Open Workspace in Container..."
3. select `otelcheck.code-workspace` and off you go...

## Go Version Support

`otelcheck` supports versions of Go that are noted by the Go release policy, that
is, major versions _N_ and _N_-1 (where _N_ is the current major version).

## Copyright and License

`pyrotest` is Copyright 2025 Harald Albrecht, and licensed under the Apache
License, Version 2.0.
