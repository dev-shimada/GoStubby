# Getting Started with GoStubby

This guide will help you get up and running with GoStubby quickly. Follow these steps to install, configure, and start using the mock server for your development needs.

## Installation

There are two ways to install and use GoStubby:

### 1. Go Package Installation

Install GoStubby using Go's package management:

```bash
go install github.com/dev-shimada/gostubby@latest
```

### 2. Docker Image

Pull and use the official Docker image:

```bash
# Pull the image
docker pull ghcr.io/dev-shimada/gostubby:latest

# Run the container
docker run -p 8080:8080 -v $(pwd)/configs:/app/configs ghcr.io/dev-shimada/gostubby:latest
```

Using Docker provides several advantages:
- No need to install Go
- Consistent environment across platforms
- Easy deployment in containerized environments
- Automatic updates by pulling the latest image

## Quick Start Tutorial

### 1. Create Your First Configuration File

Create a new file named `configs/config.json` with a basic endpoint configuration:

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

Start the server with default settings:

```bash
go run main.go
```

The server will be available at:
- HTTP: `http://localhost:8080`
- HTTPS (if configured): `https://localhost:8443`

### 3. Test the Mock Endpoint

Use curl or any HTTP client to test the endpoint:

```bash
curl -s -X GET "http://localhost:8080/example/file?param1=123%3F"
```

Expected response:
```json
{
    "message": "This is a stub response", 
    "description": "Test",
    "path1": "file",
    "param1": "123?"
}
```

```bash
curl -s -X GET -H "Accept: application/json" -d '{"key": "value"}' "http://localhost:8080/example/123/aZ0/@@@/abc/acd?param1=false&param2=aZ0&param3=000&param4=abc&param5=acd" | jq
```
Expected response:
```json
{
  "message": "This is a stub response",
  "param1": "false",
  "param2": "aZ0",
  "param3": "000",
  "param4": "abc",
  "param5": "acd",
  "path1": "123",
  "path2": "aZ0",
  "path3": "@@@",
  "path4": "abc",
  "path5": "acd"
}
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
    "body": "{\"userId\": \"{{.Path.id}}\", \"name\": \"User{{.Path.id}}\"}"
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

GoStubby supports the following command-line options:

### HTTP Settings
- Port number: `-p` or `--port` (default: 8080)

### HTTPS Settings
- HTTPS port number: `-s` or `--https-port` (default: 8443)
- SSL/TLS certificate: `-t` or `--cert` (path to SSL/TLS certificate file)
- SSL/TLS private key: `-k` or `--key` (path to SSL/TLS private key file)

### General Settings
- Configuration file: `-c` or `--config` (default: "./configs")

Example with custom settings:
```bash
go run main.go --port 3000 --config ./my-configs
```

## HTTPS Configuration

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

- Learn about advanced request matching in [Request Matching](core-features/request-matching.md)
- Explore complex responses in [Response Handling](core-features/response-handling.md)
- Set up secure endpoints with [SSL/TLS](security/ssl-tls.md)
- Check detailed configuration options in [Configuration Guide](configuration/format.md)

## Common Issues and Solutions

1. **Port Already in Use**
   ```bash
   # Change the port
   go run main.go --port 3000
   ```

2. **Configuration File Not Found**
   - Verify the configuration file exists at the specified path
   - Use absolute paths or correct relative paths
   - Check file permissions

3. **Invalid Configuration Format**
   - Verify JSON syntax
   - Ensure all required fields are present
   - Verify matching patterns are correct

4. **Certificate Issues**
   - Check certificate and key file paths
   - Ensure files are readable
   - Verify certificate validity and expiration
