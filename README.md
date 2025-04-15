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
go run main.go
```

The server will start on port 8080 by default.

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
