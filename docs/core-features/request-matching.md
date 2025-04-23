# Request Matching

GoStubby provides powerful request matching capabilities that allow you to define precise conditions for when a mock response should be returned. This document details the various request matching options available.

## URL Path Templates

GoStubby uses a template-based system for matching URL paths. Templates can contain fixed segments and variable parameters.

### Basic Template Structure

```
/fixed/path/{variable}/segments/{another_variable}
```

### Examples

```json
{
  "urlPathTemplate": "/users/{id}/posts/{postId}",
  "method": "GET"
}
```

### Variable Naming Rules
- Use alphanumeric characters and underscores
- Case-sensitive
- Must be unique within the template
- Cannot start with numbers

## Matching Types

GoStubby supports several matching patterns that can be applied to path parameters, query parameters, and request bodies.

### 1. Exact Match (`equalTo`)

Requires the value to exactly match the specified string.

```json
{
  "pathParameters": {
    "id": {
      "equalTo": "12345"
    }
  }
}
```

### 2. Regular Expression (`matches`)

Matches the value against a regular expression pattern.

```json
{
  "pathParameters": {
    "id": {
      "matches": "^[0-9]+$"
    }
  }
}
```

### 3. Negative Regular Expression (`doesNotMatch`)

Ensures the value does not match a regular expression pattern.

```json
{
  "pathParameters": {
    "username": {
      "doesNotMatch": "[0-9]+"
    }
  }
}
```

### 4. Contains (`contains`)

Checks if the value contains a specified substring.

```json
{
  "queryParameters": {
    "tags": {
      "contains": "important"
    }
  }
}
```

### 5. Does Not Contain (`doesNotContain`)

Ensures the value does not contain a specified substring.

```json
{
  "queryParameters": {
    "status": {
      "doesNotContain": "deleted"
    }
  }
}
```

## Parameter Types

### Path Parameters

Define validation rules for URL path variables.

```json
{
  "urlPathTemplate": "/users/{id}",
  "pathParameters": {
    "id": {
      "matches": "^[0-9]+$"
    }
  }
}
```

### Query Parameters

Validate query string parameters.

```json
{
  "queryParameters": {
    "page": {
      "matches": "^[0-9]+$"
    },
    "limit": {
      "matches": "^[0-9]+$"
    },
    "sort": {
      "equalTo": "desc"
    }
  }
}
```

### Request Body

Apply matching patterns to the request body content.

```json
{
  "body": {
    "matches": ".*\"email\":\\s*\"[^@]+@[^@]+\\.[^@]+\".*"
  }
}
```

## Multiple Conditions

You can combine multiple matching conditions for more precise control:

```json
{
  "request": {
    "urlPathTemplate": "/api/users/{id}",
    "method": "PUT",
    "pathParameters": {
      "id": {
        "matches": "^[0-9]+$"
      }
    },
    "queryParameters": {
      "version": {
        "equalTo": "2.0"
      }
    },
    "body": {
      "matches": ".*\"status\":\\s*\"active\".*"
    }
  }
}
```

## Best Practices

1. **Pattern Specificity**
   - Use specific patterns to avoid unintended matches
   - Consider edge cases in your patterns
   - Test patterns with various inputs

2. **Regular Expressions**
   - Keep patterns simple and readable
   - Test regular expressions thoroughly
   - Consider performance implications of complex patterns

3. **Error Handling**
   - Provide appropriate error responses for non-matching requests
   - Use clear validation patterns for better error messages
   - Consider adding custom error responses for validation failures

4. **Maintenance**
   - Document complex patterns
   - Use consistent naming conventions
   - Group related endpoints together

## Examples

### 1. Basic API Endpoint

```json
{
  "request": {
    "urlPathTemplate": "/api/v1/products/{id}",
    "method": "GET",
    "pathParameters": {
      "id": {
        "matches": "^[0-9]+$"
      }
    }
  }
}
```

### 2. Search Endpoint with Query Parameters

```json
{
  "request": {
    "urlPathTemplate": "/api/v1/search",
    "method": "GET",
    "queryParameters": {
      "q": {
        "matches": ".{3,}"
      },
      "category": {
        "matches": "^(electronics|books|clothing)$"
      }
    }
  }
}
```

### 3. User Creation with Body Validation

```json
{
  "request": {
    "urlPathTemplate": "/api/v1/users",
    "method": "POST",
    "body": {
      "matches": ".*\"email\":\\s*\"[^@]+@[^@]+\\.[^@]+\".*",
      "contains": "\"age\":",
      "doesNotContain": "\"password\":"
    }
  }
}
```

## Troubleshooting

1. **Pattern Not Matching**
   - Verify the pattern syntax
   - Test with simpler patterns first
   - Use regex testing tools for complex patterns

2. **Multiple Matches**
   - Review pattern specificity
   - Check for conflicting patterns
   - Consider pattern order

3. **Performance Issues**
   - Simplify complex patterns
   - Avoid excessive use of complex regex
   - Consider caching frequently used patterns
