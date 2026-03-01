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

**API layer** (`internal/geni/`): HTTP client with rate limiting (1 req/sec), retry logic for transient errors, and dual environment support (production `geni.com` / sandbox `api.sandbox.geni.com`).

**Batch processing** (`internal/genibatch/`): Async channel-based bulk processors for unions, profiles, and documents with request deduplication. Three background goroutines are spawned during provider configuration.

**Caching** (`internal/genicache/`): Optional bigcache-based layer for profiles and documents, enabled via provider config.

**Authentication** (`internal/authn/`): Browser-based OAuth2 implicit flow with local Echo server on :8080 for callback. Tokens cached in `~/.genealogy/`. Falls back to manual `access_token` provider attribute.

## Testing

- Acceptance tests are in `test/acceptance/` and require `TF_ACC=1`
- Tests use `terraform-plugin-testing` framework with `ProtoV6ProviderFactories`
- A valid access token must be set in `test/acceptance/const_test.go` to run acceptance tests
- Unit tests for conversion logic live alongside implementation (e.g., `internal/resource/profile/convert_test.go`)
- CI runs acceptance tests against Terraform 1.11.x with a 10-minute timeout

## Linting

Configured in `.golangci.yml`. Notable enabled linters: `errcheck`, `forcetypeassert`, `godot`, `staticcheck`, `unparam`, `unused`, `usetesting`. All issues reported (no per-linter limits).

## Releasing

Releases are triggered by pushing a `v*` tag (e.g., `v0.16.1`). GitHub Actions runs GoReleaser (`.github/workflows/release.yaml`) to produce cross-platform builds signed with GPG.

Steps:
1. Update `CHANGELOG.md`: remove "(Unreleased)" from the current section, add entries, and add a new `## X.Y.Z (Unreleased)` header above it.
2. Commit and tag: `git tag v0.X.Y`
3. Push commit and tag: `git push && git push --tags`
4. Verify: `gh run list --workflow=release.yaml`

## Go Version

Go 1.24.1 with vendored dependencies (`vendor/`).
