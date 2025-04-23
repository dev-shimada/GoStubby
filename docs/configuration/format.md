# Configuration Format Guide

This document provides a comprehensive guide to GoStubby's configuration format, including all available options and examples.

## Configuration Structure

GoStubby uses JSON format for configuration files. Each configuration file contains an array of stub mappings:

```json
[
  {
    "request": {
      // Request matching configuration
    },
    "response": {
      // Response configuration
    }
  }
]
```

## Full Configuration Schema

```json
{
  "request": {
    "urlPathTemplate": string,          // URL template with path parameters
    "method": string,                   // HTTP method (GET, POST, etc.)
    "pathParameters": {                 // Path parameter validation rules
      "paramName": {
        "equalTo": string,             // Exact match
        "matches": string,             // Regex pattern
        "doesNotMatch": string,        // Negative regex pattern
        "contains": string,            // String contains
        "doesNotContain": string       // String does not contain
      }
    },
    "queryParameters": {               // Query parameter validation
      "paramName": {
        // Same rules as pathParameters
      }
    },
    "body": {                         // Request body validation
      // Same rules as parameters
    }
  },
  "response": {
    "status": number,                 // HTTP status code
    "body": string,                   // Direct response content
    "bodyFileName": string,           // File-based response
    "headers": {                      // Response headers
      "headerName": string
    }
  }
}
```

## Request Configuration

### URL Path Templates

Templates support fixed segments and variable parameters:

```json
{
  "urlPathTemplate": "/api/v1/users/{id}/posts/{postId}"
}
```

### HTTP Methods

Supported HTTP methods:
- GET
- POST
- PUT
- DELETE
- PATCH
- HEAD
- OPTIONS

```json
{
  "method": "POST"
}
```

### Parameter Validation

#### Path Parameters

```json
{
  "pathParameters": {
    "id": {
      "matches": "^[0-9]+$"
    },
    "category": {
      "equalTo": "electronics"
    }
  }
}
```

#### Query Parameters

```json
{
  "queryParameters": {
    "page": {
      "matches": "^[0-9]+$"
    },
    "sort": {
      "equalTo": "desc"
    }
  }
}
```

#### Request Body

```json
{
  "body": {
    "matches": ".*\"email\":\\s*\"[^@]+@[^@]+\\.[^@]+\".*",
    "contains": "\"active\":true"
  }
}
```

## Response Configuration

### Status Codes

```json
{
  "status": 200  // Any valid HTTP status code
}
```

### Direct Response Body

```json
{
  "body": "{\"message\": \"Success\"}"
}
```

### File-Based Response

```json
{
  "bodyFileName": "responses/user-profile.json"
}
```

### Custom Headers

```json
{
  "headers": {
    "Content-Type": "application/json",
    "Cache-Control": "no-cache",
    "X-Custom-Header": "custom-value"
  }
}
```

## Template Variables

Available variables in response templates:

```json
{
  "body": {
    "path": "{{.Path.paramName}}",      // Path parameters
    "query": "{{.Query.paramName}}",     // Query parameters
    "method": "{{.Request.Method}}",     // HTTP method
    "header": "{{.Request.Header.name}}" // Request headers
  }
}
```

## Configuration Management

### File Organization

Recommended directory structure:
```
configs/
├── api/
│   ├── users.json
│   └── products.json
├── mock/
│   └── test-data.json
└── config.json
```

### Multiple Configuration Files

When using multiple files:
1. Each file must contain a valid JSON array
2. Files are loaded in alphabetical order
3. Later definitions override earlier ones

## Examples

### 1. Basic REST Endpoint

```json
{
  "request": {
    "urlPathTemplate": "/api/users/{id}",
    "method": "GET",
    "pathParameters": {
      "id": {
        "matches": "^[0-9]+$"
      }
    }
  },
  "response": {
    "status": 200,
    "headers": {
      "Content-Type": "application/json"
    },
    "body": "{\"id\": {{.Path.id}}, \"name\": \"User {{.Path.id}}\"}"
  }
}
```

### 2. Complex Validation

```json
{
  "request": {
    "urlPathTemplate": "/api/products",
    "method": "POST",
    "queryParameters": {
      "version": {
        "equalTo": "2.0"
      }
    },
    "body": {
      "matches": ".*\"price\":\\s*[0-9]+(\\.?[0-9]*)?.*",
      "contains": "\"category\":"
    }
  },
  "response": {
    "status": 201,
    "headers": {
      "Content-Type": "application/json",
      "Location": "/api/products/{{.Response.id}}"
    },
    "body": "{\"message\": \"Product created\", \"id\": \"12345\"}"
  }
}
```

### 3. File-Based Configuration

```json
{
  "request": {
    "urlPathTemplate": "/api/data/{type}",
    "method": "GET",
    "pathParameters": {
      "type": {
        "matches": "^(users|products|orders)$"
      }
    }
  },
  "response": {
    "status": 200,
    "headers": {
      "Content-Type": "application/json"
    },
    "bodyFileName": "responses/{{.Path.type}}.json"
  }
}
```

## Best Practices

1. **File Organization**
   - Use meaningful file names
   - Group related stubs together
   - Maintain consistent structure

2. **Validation Rules**
   - Keep regex patterns simple
   - Use specific matching rules
   - Document complex patterns

3. **Response Management**
   - Use templates for dynamic content
   - Organize response files logically
   - Maintain consistent response formats

4. **Version Control**
   - Version your configurations
   - Document changes
   - Use meaningful commit messages

## Troubleshooting

1. **Invalid JSON**
   - Use JSON validators
   - Check syntax errors
   - Verify file encoding

2. **Pattern Matching Issues**
   - Test regex patterns
   - Verify URL templates
   - Check parameter names

3. **File Loading**
   - Verify file paths
   - Check permissions
   - Validate JSON structure
