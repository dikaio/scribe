# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands
- Build: `go build -o scribe ./cmd/scribe`
- Install: `go install ./cmd/scribe`
- Run: `go run ./cmd/scribe/main.go`

## Test Commands
- Run all tests: `go test ./...`
- Run specific test: `go test ./internal/[package] -run TestName`
- Example: `go test ./internal/config -run TestDefaultConfig`

## Code Style Guidelines
- **Imports**: Group standard library first, then third-party, sorted alphabetically
- **Naming**: PascalCase for exported, camelCase for unexported, package names lowercase
- **Types**: Document all exported types/functions with comments
- **Errors**: Return as last value, propagate with context, use descriptive messages
- **Testing**: Use table-driven tests, clean up resources
- **Formatting**: Use standard Go formatting (gofmt)
- **Design**: Functions perform single tasks, methods for operations on types
- **Package Structure**: Organize by functionality in cmd/, internal/, pkg/

The project emphasizes simplicity with minimal dependencies and follows idiomatic Go practices.