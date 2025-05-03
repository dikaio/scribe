# Contributing to Scribe

Thank you for your interest in contributing to Scribe! This document outlines the process for contributing to the project and making changes.

## Development Process

### 1. Setting Up Development Environment

1. Fork the repository on GitHub
2. Clone your fork to your local machine
3. Set up the development environment:
   ```bash
   go build -o bin/scribe ./cmd/scribe
   ```

### 2. Making Changes

When implementing new features or fixing bugs, follow these steps:

1. Create a new branch for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```
   
2. Implement your changes following Go best practices
3. Run tests to ensure your changes don't break existing functionality:
   ```bash
   go test ./...
   ```

### 3. Committing Changes

When you're ready to commit your changes:

1. Stage your changes:
   ```bash
   git add .
   ```
   
2. Commit with a descriptive message following conventional commit standards:
   ```bash
   git commit -m "type: brief description of changes"
   ```
   
   Common types include:
   - `feat`: A new feature
   - `fix`: A bug fix
   - `refactor`: Code changes that neither fix a bug nor add a feature
   - `docs`: Documentation changes
   - `test`: Adding or modifying tests
   - `chore`: Changes to the build process, tools, etc.

3. Push your changes to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

### 4. Creating a Pull Request

1. Create a pull request from your fork to the main repository
2. Provide a clear description of the changes and their purpose
3. Link any related issues in the pull request description

### 5. Releases

Project maintainers will handle releases following this process:

1. Determine the appropriate release type:
   - **patch**: For backwards-compatible bug fixes
   - **minor**: For new features that are backwards-compatible
   - **major**: For changes that break backward compatibility

2. Run the release script with the appropriate release type:
   ```bash
   ./scripts/release.sh --type [patch|minor|major]
   ```

   This script will:
   - Update the version in the code
   - Update the CHANGELOG.md
   - Create a git tag
   - Push changes to GitHub

3. Monitor the GitHub Actions workflow to ensure the release is properly published

## Code Style Guidelines

- Follow standard Go code style (use `gofmt`)
- Document all exported functions and types
- Write clear commit messages
- Include tests for new functionality
- Use descriptive variable and function names

## Testing

All changes should include appropriate tests:

```bash
go test ./...
```

## Review Process

1. All pull requests require at least one review from a maintainer
2. CI checks must pass before a PR can be merged
3. Address any feedback from reviews

Thank you for contributing to Scribe!