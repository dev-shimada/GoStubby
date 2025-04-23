# Development Guide

This guide provides information for developers who want to contribute to GoStubby. It covers setting up a development environment, running tests, and contributing guidelines.

## Development Environment Setup

### Prerequisites

1. Go 1.16 or later
2. Git
3. Make (optional, but recommended)
4. OpenSSL (for SSL/TLS certificate generation)

### Getting Started

1. Clone the repository:
```bash
git clone https://github.com/dev-shimada/GoStubby.git
cd GoStubby
```

2. Install dependencies:
```bash
go mod download
```

3. Build the project:
```bash
go build
```

## Project Structure

```
GoStubby/
├── .github/            # GitHub Actions workflows
├── configs/            # Configuration examples
├── docs/              # Documentation
├── testdata/          # Test fixtures
├── body/              # Response body templates
├── main.go            # Application entry point
├── main_test.go       # Main package tests
├── go.mod             # Go module file
└── go.sum             # Go module checksum
```

## Testing

### Running Tests

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Generate coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Categories

1. **Unit Tests**
   - Located alongside source files
   - Name pattern: `*_test.go`
   - Focus on individual components

2. **Integration Tests**
   - Located in `test` directory
   - Test component interactions
   - Use test fixtures

3. **Performance Tests**
   - Benchmark critical operations
   - Located in `*_test.go` files
   - Use `testing.B` benchmarks

### Writing Tests

Example test structure:
```go
func TestFeature(t *testing.T) {
    // Test setup
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "valid input",
            input:    "test",
            expected: "result",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test execution
            result := Feature(tt.input)
            
            // Assertions
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## Code Style

### Go Guidelines

1. Follow standard Go formatting:
```bash
gofmt -s -w .
```

2. Use go lint:
```bash
golint ./...
```

3. Run go vet:
```bash
go vet ./...
```

### Code Organization

1. **Package Structure**
   - One package per directory
   - Clear package responsibilities
   - Minimal public interfaces

2. **File Organization**
   - Related functionality together
   - Clear file naming
   - Logical grouping

3. **Code Documentation**
   - Document all exported symbols
   - Include examples
   - Clear and concise comments

## Contributing

### Development Workflow

1. **Create Issue**
   - Describe the problem/feature
   - Add relevant labels
   - Link related issues

2. **Branch Creation**
   - Use descriptive names
   - Include issue number
   - Example: `feature/123-add-ssl-support`

3. **Development**
   - Write tests first
   - Implement changes
   - Update documentation

4. **Code Review**
   - Submit pull request
   - Address review comments
   - Update as needed

### Pull Request Guidelines

1. **Preparation**
   - Rebase on main
   - Run all tests
   - Update documentation

2. **PR Description**
   - Clear description
   - Reference issues
   - List changes

3. **Code Quality**
   - Pass all tests
   - Follow style guide
   - Include documentation

### Commit Messages

Follow conventional commits format:
```
type(scope): description

[optional body]

[optional footer]
```

Types:
- feat: New feature
- fix: Bug fix
- docs: Documentation
- style: Formatting
- refactor: Code restructuring
- test: Adding tests
- chore: Maintenance

Example:
```
feat(ssl): add SSL/TLS support

- Add certificate loading
- Implement HTTPS server
- Update documentation

Closes #123
```

## Debugging

### VSCode Configuration

1. Create `.vscode/launch.json`:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug GoStubby",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}",
            "args": ["--config", "./configs/debug.json"]
        }
    ]
}
```

2. Add debug points
3. Use VSCode debugging features

### Logging

1. Development logging:
```go
log.Printf("Debug: %v", value)
```

2. Production logging:
```go
log.Printf("Error: %v", err)
```

## Performance Profiling

### CPU Profiling

```bash
go test -cpuprofile cpu.prof -bench .
go tool pprof cpu.prof
```

### Memory Profiling

```bash
go test -memprofile mem.prof -bench .
go tool pprof mem.prof
```

## Release Process

1. **Version Update**
   - Update version numbers
   - Update CHANGELOG.md
   - Update documentation

2. **Testing**
   - Run all tests
   - Perform integration testing
   - Check documentation

3. **Release**
   - Create release branch
   - Tag version
   - Push to repository

4. **Post-Release**
   - Update main branch
   - Clean up branches
   - Update documentation

## Support

- GitHub Issues: Bug reports and feature requests
- Discussions: General questions and discussions
- Pull Requests: Code contributions

Remember to:
- Search existing issues
- Provide clear descriptions
- Include minimal examples
- Be respectful and constructive
