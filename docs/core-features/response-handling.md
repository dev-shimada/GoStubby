# Response Handling

GoStubby provides flexible response handling capabilities that allow you to create both static and dynamic mock responses. This document covers all aspects of response configuration and customization.

## Response Configuration

### Basic Structure

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

## Response Types

### 1. Direct Response Body

Use the `body` field to specify the response content directly in the configuration:

```json
{
  "response": {
    "status": 200,
    "body": "{\"message\": \"Hello, World!\"}",
    "headers": {
      "Content-Type": "application/json"
    }
  }
}
```

### 2. File-Based Response

Use the `bodyFileName` field to load response content from a file:

```json
{
  "response": {
    "status": 200,
    "bodyFileName": "responses/user-profile.json",
    "headers": {
      "Content-Type": "application/json"
    }
  }
}
```

## Template-Based Responses

GoStubby supports dynamic response generation using templates. Templates can access request parameters and generate customized responses.

### Available Template Variables

1. **Path Parameters**
   ```json
   {
     "response": {
       "body": "{\"userId\": \"{{.Path.id}}\", \"message\": \"User {{.Path.id}} details\"}"
     }
   }
   ```

2. **Query Parameters**
   ```json
   {
     "response": {
       "body": "{\"search\": \"{{.Query.q}}\", \"page\": {{.Query.page}}}"
     }
   }
   ```

### Template Syntax

- Path Parameters: `{{.Path.paramName}}`
- Query Parameters: `{{.Query.paramName}}`
- HTTP Method: `{{.Request.Method}}`
- Request Headers: `{{.Request.Header.headerName}}`

## Status Codes

Configure appropriate HTTP status codes for different scenarios:

```json
{
  "response": {
    "status": 201,  // Created
    "body": "{\"message\": \"Resource created successfully\"}"
  }
}
```

Common status codes:
- 200: OK (Success)
- 201: Created
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 500: Internal Server Error

## Response Headers

### Setting Custom Headers

```json
{
  "response": {
    "headers": {
      "Content-Type": "application/json",
      "Cache-Control": "no-cache",
      "X-Custom-Header": "custom-value"
    }
  }
}
```

### Common Headers
- Content-Type
- Cache-Control
- Access-Control-Allow-Origin
- Authorization
- X-Rate-Limit

## Advanced Response Features

### 1. Conditional Responses

Use different templates based on request parameters:

```json
{
  "response": {
    "body": "{% if eq .Query.type \"premium\" %}
      {\"message\": \"Premium content\"}
    {% else %}
      {\"message\": \"Basic content\"}
    {% endif %}"
  }
}
```

### 2. Dynamic File Loading

Load different response files based on parameters:

```json
{
  "response": {
    "bodyFileName": "responses/{{.Path.type}}.json"
  }
}
```

### 3. Custom Error Responses

```json
{
  "response": {
    "status": 400,
    "body": "{\"error\": \"Invalid request\", \"details\": \"Missing required field: {{.Path.field}}\"}"
  }
}
```

## Best Practices

1. **Response Organization**
   - Group related responses in directories
   - Use meaningful file names
   - Maintain consistent file structure

2. **Template Usage**
   - Keep templates simple and readable
   - Validate template syntax
   - Handle missing parameters gracefully

3. **Error Handling**
   - Use appropriate status codes
   - Provide meaningful error messages
   - Include relevant error details

4. **Performance**
   - Cache file-based responses when possible
   - Minimize template complexity
   - Use appropriate content compression

## Examples

### 1. Basic JSON Response

```json
{
  "response": {
    "status": 200,
    "headers": {
      "Content-Type": "application/json"
    },
    "body": "{\"id\": 1, \"name\": \"Example\"}"
  }
}
```

### 2. Dynamic Response with Templates

```json
{
  "response": {
    "status": 200,
    "headers": {
      "Content-Type": "application/json"
    },
    "body": "{
      \"userId\": \"{{.Path.id}}\",
      \"query\": \"{{.Query.q}}\",
      \"timestamp\": \"{{.Request.Time}}\"
    }"
  }
}
```

### 3. File-Based Response with Custom Headers

```json
{
  "response": {
    "status": 200,
    "headers": {
      "Content-Type": "application/json",
      "Cache-Control": "max-age=3600",
      "ETag": "\"123456\""
    },
    "bodyFileName": "responses/large-response.json"
  }
}
```

## Troubleshooting

1. **Template Errors**
   - Check template syntax
   - Verify parameter names
   - Ensure parameters are available

2. **File Loading Issues**
   - Verify file paths
   - Check file permissions
   - Validate file content

3. **Content Type Mismatches**
   - Ensure Content-Type header matches body
   - Verify JSON syntax
   - Check character encoding

4. **Performance Issues**
   - Monitor file sizes
   - Optimize template processing
   - Consider response caching
