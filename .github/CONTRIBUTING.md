# Contributing to ggraph

First off, thank you for considering contributing to ggraph! It's people like you that make ggraph such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the issue list as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible using our bug report template.

**Good Bug Reports** include:

- A quick summary and/or background
- Steps to reproduce
  - Be specific!
  - Give sample code if you can
- What you expected would happen
- What actually happens
- Notes (possibly including why you think this might be happening, or stuff you tried that didn't work)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

- Use a clear and descriptive title
- Provide a step-by-step description of the suggested enhancement
- Provide specific examples to demonstrate the steps
- Describe the current behavior and explain which behavior you expected to see instead
- Explain why this enhancement would be useful

### Pull Requests

1. Fork the repo and create your branch from `main`
2. If you've added code that should be tested, add tests
3. If you've changed APIs, update the documentation
4. Ensure the test suite passes with `make test`
5. Make sure your code lints with `make lint`
6. Issue that pull request!

## Development Setup

### Prerequisites

- Go 1.24 or higher
- Git

### Setting Up Your Environment

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/ggraph.git
   cd ggraph
   ```
3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/morphy76/ggraph.git
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```

### Running Tests

```bash
# Run all tests
make test

# Run tests with benchmarks
make test-bench

# Run linting
make lint
```

### Project Structure

```
ggraph/
├── internal/          # Internal packages (not exported)
│   ├── agent/        # Agent implementation details
│   └── graph/        # Graph implementation details
├── pkg/              # Public packages
│   ├── agent/        # Agent API
│   ├── builders/     # Builder utilities
│   └── graph/        # Graph API
└── examples/         # Example implementations
```

## Style Guidelines

### Go Code Style

- Follow standard Go conventions and idioms
- Use `gofmt` to format your code
- Use meaningful variable and function names
- Write comments for exported functions and types
- Keep functions focused and concise

### Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line
- Consider starting the commit message with an applicable prefix:
  - `feat:` - A new feature
  - `fix:` - A bug fix
  - `docs:` - Documentation only changes
  - `style:` - Changes that do not affect the meaning of the code
  - `refactor:` - A code change that neither fixes a bug nor adds a feature
  - `perf:` - A code change that improves performance
  - `test:` - Adding missing tests or correcting existing tests
  - `chore:` - Changes to the build process or auxiliary tools

### Documentation

- Update the README.md if you change functionality
- Comment all exported functions, types, and packages
- Provide examples for new features
- Keep documentation up to date with code changes

## Testing Guidelines

- Write unit tests for all new code
- Aim for high test coverage (>80%)
- Include both positive and negative test cases
- Use table-driven tests where appropriate
- Test edge cases and error conditions

Example test structure:

```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name    string
        input   interface{}
        want    interface{}
        wantErr bool
    }{
        // test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

## Release Process

Releases are managed by the maintainers. Version numbers follow [Semantic Versioning](https://semver.org/).

## Questions?

Feel free to open an issue with your question or reach out to the maintainers.

## Recognition

Contributors will be recognized in the project's README and release notes.

Thank you for contributing! 🎉
