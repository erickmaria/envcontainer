# AGENTS.md - Guidance for AI coding agents

Purpose
- Provide concise, actionable guidance for AI coding agents working on this repository.

Quick orientation
- Entry point: `main.go` invokes the Cobra CLI implemented under `cmd/`.
- Primary areas of work: `cmd/` (CLI), `internal/runtime/` (runtime abstractions), `internal/runtime/docker/` (Docker implementation), and `internal/` for helper packages.

Build & run (developer commands)
- Build binary: `make build` (invokes `go build` producing `envcontainer`).
- Run locally: `make run` (runs `go run cmd/envcontainer/*.go`).
- Create compact release: `make compact/linux` (build + zip).
- Bump version: `make bump-version/major|minor|patch` (uses `bump2version`).
- Release: `make release` (push tags to `main`).

Testing & static checks
- Run tests if present: `go test ./...`.
- Format and vet: `gofmt -w .` and `go vet ./...`.
- Maintain modules: `go mod tidy`.

Container & runtime notes
- The project manages containers via an internal runtime interface (`internal/runtime/runtime.go`).
- Docker-specific implementation lives in `internal/runtime/docker/` and uses the Docker SDK.
- The README documents `.envcontainer.yaml` configuration and recommends using the `docker` CLI when `devcontainer` lacks features.

> Agent behavior and conventions
- Keep changes minimal and targeted to the user's request. Prefer small, testable edits.
- When modifying behavior, update `README.md` and add short notes in `CHANGELOG.md` when relevant.
- Respect existing coding style: follow idiomatic Go, run `gofmt`, and keep exported APIs stable.
- For CLI work, treat `cmd/` as the canonical place for user-facing changes (Cobra commands).

When to use Docker vs. Go runtime code
- If the task modifies container lifecycle, prefer changes under `internal/runtime/` and `internal/runtime/docker/`.
- For build, packaging, or release tasks modify `Makefile` and release tooling only after checking existing targets.

Recommended quick checks before edits
- Read `README.md` and `Makefile` to capture expected UX and commands.
- Locate CLI commands in `cmd/envcontainer` to understand flags and behavior.
- Search for related types under `internal/types/` and `internal/pkg/` for data shapes.

Suggested next steps for humans (after agent edits)
- Run `gofmt -w .` and `go test ./...` locally.
- Build the binary with `make build` to verify compile-time changes.
- If releasing or bumping versions, use `make bump-version/*` and `make release` as appropriate.

Contact / context
- This document is intended to be concise; add repository-specific edge-cases here when discovered.
