# Contributing to CF ECS Log Plugin

Thank you for your interest in contributing to the CF ECS Log Plugin! This document provides guidelines and instructions for contributing.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/cf-ecslogs-plugin.git`
3. Create a new branch: `git checkout -b feature/your-feature-name`

## Prerequisites

- Go 1.22 or later
- Cloud Foundry CLI v6.7.0 or later
- Make

## Development Setup

```bash
# Install dependencies
make deps

# Build the plugin
make build

# Install the plugin locally for testing
make install
```

## Making Changes

1. Make your changes in a feature branch
2. Follow Go best practices and conventions
3. Add tests for new functionality
4. Ensure your code builds: `make build`
5. Format your code: `make fmt`
6. Run the linter: `make vet`

## Code Style

- Follow standard Go formatting (use `gofmt` or `make fmt`)
- Write clear, descriptive commit messages
- Keep functions small and focused
- Add comments for exported functions and types
- Use meaningful variable and function names

## Testing

While there's a known issue with Go test binaries on macOS (dyld UUID error), you should:

1. Verify the code compiles: `make build`
2. Test manually by installing and running against a CF environment
3. Add unit tests for new functionality in `*_test.go` files

## Pull Request Process

1. Update the README.md with details of changes if applicable
2. Update the CHANGELOG if there is one
3. Ensure your branch is up to date with main
4. Submit a pull request with a clear description of:
   - What changes were made
   - Why they were made
   - How to test the changes

## Reporting Issues

When reporting issues, please include:

- Go version: `go version`
- CF CLI version: `cf version`
- Operating system and version
- Steps to reproduce the issue
- Expected vs actual behavior
- Any relevant log output or error messages

## Questions?

Feel free to open an issue for any questions about contributing.

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.
