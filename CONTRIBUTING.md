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
4. Wait for GitHub Actions to run tests on your PR
5. Address any feedback from reviewers

## Release Process

The release process is largely automated through GitHub Actions and scripts.

### Testing a Potential Release

Before creating an official release, you can test the release process:

1. Go to GitHub Actions > Pre-Release Testing
2. Click "Run workflow"
3. Select the version type (patch, minor, major)
4. This will verify that the code builds and tests pass across all platforms without creating an actual release

### Creating an Official Release

Maintainers can create new releases using either of these methods:

#### Option 1: Using the Release Script (Recommended)

```bash
# For a patch release (0.1.0 -> 0.1.1)
./scripts/release.sh

# For a minor release (0.1.0 -> 0.2.0)
./scripts/release.sh --type minor

# For a major release (0.1.0 -> 1.0.0)
./scripts/release.sh --type major
```

The script will:
- Verify your working directory is clean
- Run tests
- Calculate the new version
- Generate a changelog from commit messages
- Update version in code
- Commit changes
- Create and push a Git tag
- Trigger GitHub Actions to build and publish the release

#### Option 2: Manual Tag Creation

Experienced maintainers can also create a release by:
1. Updating the version in `pkg/cli/cli.go`
2. Updating the CHANGELOG.md
3. Committing these changes
4. Creating and pushing a tag:
   ```bash
   git tag -a "vX.Y.Z" -m "Release vX.Y.Z"
   git push origin vX.Y.Z
   ```

This will trigger the GitHub Actions release workflow automatically.

## Automated Release Process

When a new tag is pushed, GitHub Actions automatically:

1. Builds the binaries for all supported platforms
2. Creates a GitHub release with these binaries
3. Updates the Homebrew formula

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