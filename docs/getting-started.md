# Getting Started with GoStubby

This guide will help you get up and running with GoStubby quickly. Follow these steps to install, configure, and start using the mock server for your development needs.

## Installation

Install GoStubby using Go's package management:

```bash
go get github.com/dev-shimada/gostubby
```

## Quick Start Tutorial

### 1. Create Your First Configuration

Create a new file `configs/config.json` with a basic endpoint configuration:

```json
[
  {
    "request": {
      "urlPathTemplate": "/hello/{name}",
      "method": "GET",
      "pathParameters": {
        "name": {
          "matches": "^[a-zA-Z]+$"
        }
      }
    },
    "response": {
      "status": 200,
      "headers": {
        "Content-Type": "application/json"
      },
      "body": "{\"message\": \"Hello, {{.Path.name}}!\"}"
    }
  }
]
```

### 2. Start the Server

Run the server with default settings:

```bash
go run main.go
```

The server will start on:
- HTTP: `http://localhost:8080`
- HTTPS (if configured): `https://localhost:8443`

### 3. Test Your Mock Endpoint

Use curl or any HTTP client to test the endpoint:

```bash
curl http://localhost:8080/hello/world
```

Expected response:
```json
{"message": "Hello, world!"}
```

## Basic Configuration Examples

### 1. Static Response

```json
{
  "request": {
    "urlPathTemplate": "/api/status",
    "method": "GET"
  },
  "response": {
    "status": 200,
    "body": "{\"status\": \"operational\"}"
  }
}
```

### 2. Dynamic Response with Path Parameters

```json
{
  "request": {
    "urlPathTemplate": "/users/{id}/profile",
    "method": "GET",
    "pathParameters": {
      "id": {
        "matches": "^[0-9]+$"
      }
    }
  },
  "response": {
    "status": 200,
    "body": "{\"userId\": \"{{.Path.id}}\", \"name\": \"User {{.Path.id}}\"}"
  }
}
```

### 3. Request Body Validation

```json
{
  "request": {
    "urlPathTemplate": "/api/users",
    "method": "POST",
    "body": {
      "matches": ".*\"email\":\\s*\"[^@]+@[^@]+\\.[^@]+\".*"
    }
  },
  "response": {
    "status": 201,
    "body": "{\"message\": \"User created successfully\"}"
  }
}
```

## Command Line Options

GoStubby supports various command line options for configuration:

### HTTP Configuration
- Port: `-p` or `--port` (default: 8080)

### HTTPS Configuration
- HTTPS Port: `-s` or `--https-port` (default: 8443)
- Certificate: `-t` or `--cert` (path to SSL/TLS certificate file)
- Private Key: `-k` or `--key` (path to SSL/TLS private key file)

### General Configuration
- Configuration: `-c` or `--config` (default: "./configs")

Example with custom settings:
```bash
go run main.go --port 3000 --config ./my-configs
```

## Running with HTTPS

1. Generate a self-signed certificate (for development):
```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout ./certs/server.key -out ./certs/server.crt
```

2. Start the server with HTTPS:
```bash
go run main.go --cert ./certs/server.crt --key ./certs/server.key
```

## Next Steps

- Explore [Request Matching](core-features/request-matching.md) for advanced request matching
- Learn about [Response Handling](core-features/response-handling.md) for complex responses
- Configure [SSL/TLS](security/ssl-tls.md) for secure endpoints
- Check [Configuration Guide](configuration/format.md) for detailed configuration options

## Common Issues and Solutions

1. **Port Already in Use**
   ```bash
   # Change the port
   go run main.go --port 3000
   ```

2. **Configuration File Not Found**
   - Ensure the config file exists in the specified path
   - Use absolute paths or correct relative paths
   - Check file permissions

3. **Invalid Configuration Format**
   - Validate JSON syntax
   - Ensure all required fields are present
   - Check for proper matching patterns

4. **Certificate Issues**
   - Verify certificate and key file paths
   - Ensure files are readable
   - Check certificate expiration and validity
