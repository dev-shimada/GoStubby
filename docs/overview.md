# GoStubby Overview

GoStubby is a flexible and powerful mock server implementation in Go that allows you to define mock endpoints with advanced request matching capabilities and templated responses. GoStubby can be installed either as a Go package or run as a Docker container, making it accessible to both Go and non-Go developers.

## Installation Options

1. **Go Package**:
   ```bash
   go install github.com/dev-shimada/gostubby@latest
   ```

2. **Docker Container**:
   ```bash
   docker pull ghcr.io/dev-shimada/gostubby:latest
   ```

## Features

### Flexible Request Matching
- URL path matching with templates (e.g., `/users/{id}`)
- Regular expression pattern matching
- Query parameter validation
- Request body validation
- Multiple matching patterns: `equalTo`, `matches`, `doesNotMatch`, `contains`, `doesNotContain`

[Rest of the content remains the same...]
