# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

@PRD.md
@README.md

## Build & Run Commands
- Build: `go build -o bin/scribe ./cmd/scribe`
- Install: `go install ./cmd/scribe`
- Run: `go run ./cmd/scribe/main.go`

## Test Commands
- Run all tests: `go test ./...`
- Run specific test: `go test ./internal/[package] -run TestName`
- Example: `go test ./internal/config -run TestDefaultConfig`
- Run with coverage: `go test ./... -coverprofile=build/coverage/coverage.out`
- View coverage report: `go tool cover -html=build/coverage/coverage.out -o build/coverage/coverage.html`
- Check coverage stats: `go tool cover -func=build/coverage/coverage.out`

## Task Commands
- Build: `task build`
- Install: `task install`
- Create a new site: `task new -- site-name`
- Serve a site: `task serve`
- Run tests: `task test`
- Create a release: `task release:patch` (or `minor`/`major`)

## Release Process
When changes are ready for a new release:
1. Make sure all changes are committed and pushed to GitHub
2. Run `task release:patch` (or `minor`/`major` depending on the type of changes)
3. The script will update version numbers, generate changelog, create a tag, and push to GitHub

## Code Style Guidelines
- **Imports**: Group standard library first, then third-party, sorted alphabetically
- **Naming**: PascalCase for exported, camelCase for unexported, package names lowercase
- **Types**: Document all exported types/functions with comments
- **Errors**: Return as last value, propagate with context, use descriptive messages
- **Testing**: Use table-driven tests, clean up resources
- **Formatting**: Use standard Go formatting (gofmt)
- **Design**: Functions perform single tasks, methods for operations on types
- **Package Structure**: Organize by functionality in cmd/, internal/, pkg/

## Commit Message Conventions
- `feat:` or `feature:` for new features
- `fix:` for bug fixes
- `refactor:` for code refactoring
- `docs:` for documentation changes
- `test:` for test-related changes
- `chore:` for maintenance tasks

The project emphasizes simplicity with minimal dependencies and follows idiomatic Go practices.