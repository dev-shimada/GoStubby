# GoStubby
[![Go Report Card](https://goreportcard.com/badge/github.com/dev-shimada/GoStubby)](https://goreportcard.com/report/github.com/dev-shimada/GoStubby)
[![Coverage Status](https://coveralls.io/repos/github/dev-shimada/GoStubby/badge.svg?branch=main)](https://coveralls.io/github/dev-shimada/GoStubby?branch=main)
[![CI](https://github.com/dev-shimada/GoStubby/actions/workflows/CI.yaml/badge.svg)](https://github.com/dev-shimada/GoStubby/actions/workflows/CI.yaml)
[![build](https://github.com/dev-shimada/GoStubby/actions/workflows/build-docker-image.yaml/badge.svg)](https://github.com/dev-shimada/GoStubby/actions/workflows/build-docker-image.yaml)
[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/dev-shimada/GoStubby/blob/master/LICENSE)

A flexible and powerful mock server implementation in Go that allows you to define mock endpoints with advanced request matching capabilities and templated responses.

## Features

- **Flexible Request Matching**:
  - URL path matching with templates (e.g., `/users/{id}`)
  - Regular expression pattern matching
  - Query parameter validation
  - Request body validation
  - Multiple matching patterns: `equalTo`, `matches`, `doesNotMatch`, `contains`, `doesNotContain`

- **Powerful Response Handling**:
  - Template-based response bodies with access to request parameters
  - File-based response bodies
  - Custom HTTP status codes
  - Custom response headers

## Installation

```bash
go get github.com/dev-shimada/gostubby
```

## Usage

1. Create a configuration file in `./configs/config.json`:

```json
[
  {
    "request": {
      "urlPathTemplate": "/users/{id}",
      "method": "GET",
      "pathParameters": {
        "id": {
          "matches": "^[0-9]+$"
        }
      }
    },
    "response": {
      "status": 200,
      "body": "{\"id\": \"{{.Path.id}}\", \"name\": \"User {{.Path.id}}\"}"
    }
  }
]
```

2. Run the server:

```bash
# Default port (8080) and config directory
go run main.go

# Custom port using short option
go run main.go -p 3000

# Custom port using long option
go run main.go --port 3000

# Custom config path using short option
go run main.go -c ./path/to/config.json

# Custom config path using long option
go run main.go --config ./path/to/configs
```

The server supports the following command-line options:

HTTP Configuration:
- Port: `-p` or `--port` (default: 8080)

HTTPS Configuration:
- HTTPS Port: `-s` or `--https-port` (default: 8443)
- Certificate: `-t` or `--cert` (path to SSL/TLS certificate file)
- Private Key: `-k` or `--key` (path to SSL/TLS private key file)

General Configuration:
- Configuration: `-c` or `--config` (default: "./configs")

You can specify either a single JSON configuration file or a directory containing multiple JSON configuration files. When a directory is specified, all JSON files in that directory will be loaded.

### SSL/TLS Support

The server supports running in HTTPS mode when SSL/TLS certificates are provided. You can run the server with both HTTP and HTTPS enabled simultaneously.

To enable HTTPS:
1. Obtain SSL/TLS certificate and private key files
2. Run the server with the certificate and key file paths:

```bash
# Run with both HTTP and HTTPS
go run main.go --cert ./certs/server.crt --key ./certs/server.key

# Custom ports for both HTTP and HTTPS
go run main.go --port 8080 --https-port 8443 --cert ./certs/server.crt --key ./certs/server.key
```

For development and testing, you can generate a self-signed certificate:
```bash
# Generate private key and self-signed certificate
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout ./certs/server.key -out ./certs/server.crt
```

Note: The server enforces TLS 1.2 or higher for security.

## Configuration Format

### Request Matching

```json
{
  "request": {
    "urlPathTemplate": "/example/{param}",  // URL template with path parameters
    "method": "GET",                        // HTTP method
    "pathParameters": {                     // Path parameter validation rules
      "param": {
        "equalTo": "value",                 // Exact match
        "matches": "^[0-9]+$",             // Regex pattern match
        "doesNotMatch": "[a-z]+",          // Negative regex pattern match
        "contains": "substring",            // String contains
        "doesNotContain": "substring"       // String does not contain
      }
    },
    "queryParameters": {                    // Query parameter validation
      "param": {
        // Same matching rules as pathParameters
      }
    },
    "body": {                              // Request body validation
      // Same matching rules as parameters
    }
  }
}
```

### Response Configuration

```json
{
  "response": {
    "status": 200,                         // HTTP status code
    "body": "Response content",            // Direct response content
    "bodyFileName": "response.json",       // OR file-based response
    "headers": {                           // Custom response headers
      "Content-Type": "application/json"
    }
  }
}
```

### Template Variables

In response bodies, you can use the following template variables:
- Path parameters: `{{.Path.paramName}}`
- Query parameters: `{{.Query.paramName}}`

## Example Configurations

1. Basic endpoint with path parameter:
```json
{
  "request": {
    "urlPathTemplate": "/users/{id}",
    "method": "GET",
    "pathParameters": {
      "id": {
        "matches": "^[0-9]+$"
      }
    }
  },
  "response": {
    "status": 200,
    "body": "{\"id\": \"{{.Path.id}}\", \"name\": \"User {{.Path.id}}\"}"
  }
}
```

2. Endpoint with file-based response:
```json
{
  "request": {
    "urlPathTemplate": "/data/{type}",
    "method": "GET"
  },
  "response": {
    "status": 200,
    "bodyFileName": "responses/data.json"
  }
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

*Read this in other languages: [日本語](README.ja.md)*
