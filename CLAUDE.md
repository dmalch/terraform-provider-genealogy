# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Terraform provider for managing genealogical data on Geni.com. Uses the modern Terraform Plugin Framework (v1.17.0) with Protocol v6. Provider name is `geni`, registry address `registry.terraform.io/dmalch/genealogy`.

## Common Commands

```bash
make build              # Build binary to bin/
make build-local        # Build for local darwin_arm64 with versioned path
make test               # Run all tests (go test -v ./...)
make docs               # Generate provider docs via tfplugindocs
make clean              # Remove built binaries

# Run a single test
go test -v ./test/acceptance/ -run TestAccGeniProfile -timeout 10m

# Acceptance tests require TF_ACC=1
TF_ACC=1 go test -v ./test/acceptance/ -timeout 10m

# Lint (used in CI)
golangci-lint run
```

## Architecture

**Provider entry:** `main.go` → `internal/provider.go` (`GeniProvider`).

**Resources & Data Sources** live under `internal/resource/` and `internal/datasource/`, each in its own package:
- `profile/`, `union/`, `document/` — resources with CRUD split across separate files (`create.go`, `read.go`, `update.go`, `delete.go`)
- `project/` — data source
- Each resource package has `schema.go` (Terraform schema), `model.go` (data structs), `convert.go` (API↔Terraform conversion), and `resource.go` (interface impl)
- `event/` — shared event schema (birth/death/burial/baptism) used across resources
- `geniplanmodifier/` — custom plan modifiers

**API layer** (`github.com/dmalch/go-geni`): standalone Go library — HTTP client with rate limiting (1 req/sec), retry logic for transient errors, and dual environment support (production `geni.com` / sandbox `api.sandbox.geni.com`). The OAuth helper lives under `github.com/dmalch/go-geni/auth`. Previously vendored under `internal/geni/` and `internal/authn/`; extracted in provider `v0.21.0` / library `v0.1.0`.

**Batch processing** (`internal/genibatch/`): Async channel-based bulk processors for unions, profiles, and documents with request deduplication. Three background goroutines are spawned during provider configuration.

**Caching** (`internal/genicache/`): Optional bigcache-based layer for profiles and documents, enabled via provider config.

**Authentication** (`github.com/dmalch/go-geni/auth`): Browser-based OAuth2 implicit flow with local Echo server on :8080 for callback. Tokens cached in `~/.genealogy/`. Falls back to manual `access_token` provider attribute.

## Testing

- Acceptance tests are in `test/acceptance/` and require `TF_ACC=1`
- Tests use `terraform-plugin-testing` framework with `ProtoV6ProviderFactories`
- A valid access token must be set in `test/acceptance/const_test.go` to run acceptance tests
- Unit tests live alongside implementation (e.g., `internal/resource/profile/convert_test.go`)
- CI runs acceptance tests against Terraform 1.11.x with a 10-minute timeout

### Unit test conventions

- Use `gomega` with a dot-import: `. "github.com/onsi/gomega"` and `Expect(...).To(...)` matchers — not `testify` or bare `t.Fatal`.
- Name the top-level test after the function under test (e.g. `TestValueFrom`, `TestResolveDocumentImport`); use `t.Run("Behavioral scenario description", ...)` for each case.
- Call `RegisterTestingT(t)` at the start of every `t.Run` sub-test (Gomega needs it per goroutine).
- Use `t.Context()` rather than `context.Background()`.
- Tests live in the same package as the code under test (no `_test` package) so unexported domain functions are directly callable.

## Linting

Configured in `.golangci.yml`. Notable enabled linters: `errcheck`, `forcetypeassert`, `godot`, `staticcheck`, `unparam`, `unused`, `usetesting`. All issues reported (no per-linter limits).

## Releasing

Releases are triggered by pushing a `v*` tag (e.g., `v0.16.1`). GitHub Actions runs GoReleaser (`.github/workflows/release.yaml`) to produce cross-platform builds signed with GPG.

### Per-PR (before the release tag)

Every PR that ships user-visible behavior **must** land these two artifacts in the same merge so the release-cut step is a pure rename:

1. **Add a CHANGELOG entry under `## X.Y.Z (Unreleased)`** at the top of `CHANGELOG.md`. Choose the section by impact: `FEATURES` (new capabilities), `BUG FIXES`, `IMPROVEMENTS`, `BREAKING CHANGES`, `SECURITY`, `BEHAVIORAL CHANGES`. Reference the issue and PR numbers in parentheses, e.g. `(#82, #87)`.
2. **Regenerate provider docs**: `make docs` (which runs `tfplugindocs generate --provider-name geni`). Required when schemas, resources, data sources, or list resources change. Commit any `docs/**` deltas alongside the code.
   - `tfplugindocs` is not vendored. If `make docs` reports `No such file or directory`, install it once with `go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest` and ensure `$(go env GOPATH)/bin` is on `PATH`.

### Tagging the release

1. Promote `## X.Y.Z (Unreleased)` → `## X.Y.Z` in `CHANGELOG.md`, and add a new `## X.Y.Z+1 (Unreleased)` placeholder above it.
2. Confirm `docs/**` is in sync with the current schema (`make docs` should produce no diff).
3. Commit, tag, push: `git tag v0.X.Y && git push && git push --tags`.
4. Verify: `gh run list --workflow=release.yaml`.

## Go Version

Go 1.26 with vendored dependencies (`vendor/`).
