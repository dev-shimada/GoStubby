# SSL/TLS Configuration Guide

This guide covers the SSL/TLS implementation in GoStubby, including setup, configuration, and best practices for securing your mock server.

## Overview

GoStubby supports HTTPS through SSL/TLS, allowing you to:
- Run secure HTTPS endpoints
- Support both HTTP and HTTPS simultaneously
- Configure custom certificates
- Enforce modern security standards

## Quick Start

### Basic HTTPS Setup

1. Generate a self-signed certificate (for development):
```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout ./certs/server.key -out ./certs/server.crt
```

2. Start the server with HTTPS enabled:
```bash
go run main.go --cert ./certs/server.crt --key ./certs/server.key
```

## Command Line Options

### SSL/TLS Configuration Flags

| Flag | Long Form | Description | Default |
|------|-----------|-------------|---------|
| `-s` | `--https-port` | HTTPS port number | 8443 |
| `-t` | `--cert` | Path to SSL/TLS certificate file | - |
| `-k` | `--key` | Path to SSL/TLS private key file | - |

### Examples

```bash
# Run with custom HTTPS port
go run main.go --cert ./certs/server.crt --key ./certs/server.key --https-port 443

# Run both HTTP and HTTPS with custom ports
go run main.go --port 8080 --https-port 8443 \
  --cert ./certs/server.crt --key ./certs/server.key
```

## Certificate Management

### 1. Self-Signed Certificates

Generate a self-signed certificate for development:

```bash
# Generate private key and certificate
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout ./certs/server.key \
  -out ./certs/server.crt \
  -subj "/CN=localhost"

# Generate with custom details
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout ./certs/server.key \
  -out ./certs/server.crt \
  -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"
```

### 2. Using Let's Encrypt

For production environments, obtain a free certificate from Let's Encrypt:

1. Install certbot:
```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install certbot

# macOS
brew install certbot
```

2. Generate certificate:
```bash
sudo certbot certonly --standalone -d your-domain.com
```

3. Use the generated certificates:
```bash
go run main.go \
  --cert /etc/letsencrypt/live/your-domain.com/fullchain.pem \
  --key /etc/letsencrypt/live/your-domain.com/privkey.pem
```

### 3. Commercial Certificates

When using commercial SSL certificates:

1. Combine certificate files:
```bash
cat domain.crt intermediate.crt root.crt > fullchain.pem
```

2. Start server with combined certificate:
```bash
go run main.go \
  --cert ./certs/fullchain.pem \
  --key ./certs/private.key
```

## Security Configuration

### TLS Version Control

GoStubby enforces modern TLS standards:
- Minimum TLS version: 1.2
- Recommended cipher suites
- Perfect Forward Secrecy (PFS)

### Cipher Suite Configuration

Default cipher suites are configured for maximum security:

```go
[]uint16{
    tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
    tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
    tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
    tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
}
```

## Best Practices

### 1. Certificate Management
- Keep private keys secure
- Regularly rotate certificates
- Monitor certificate expiration
- Use appropriate key sizes (minimum 2048-bit RSA)

### 2. Security Headers

Configure security headers in your responses:

```json
{
  "response": {
    "headers": {
      "Strict-Transport-Security": "max-age=31536000; includeSubDomains",
      "X-Content-Type-Options": "nosniff",
      "X-Frame-Options": "DENY",
      "X-XSS-Protection": "1; mode=block"
    }
  }
}
```

### 3. Production Setup
- Use valid SSL certificates from trusted CAs
- Enable HTTP/2 support
- Configure proper CORS headers
- Implement rate limiting
- Monitor security logs

## Troubleshooting

### Common Issues

1. **Certificate Problems**
```
Error: tls: failed to load certificate
```
- Verify file paths
- Check file permissions
- Validate certificate format
- Ensure certificate matches private key

2. **Port Access Issues**
```
Error: listen tcp :443: bind: permission denied
```
- Use ports > 1024 for non-root users
- Configure proper system permissions
- Use port forwarding if needed

3. **Certificate Trust Issues**
```
Error: x509: certificate signed by unknown authority
```
- Add root certificate to trust store
- Use proper certificate chain
- Verify intermediate certificates

### Validation Tools

1. SSL Labs Test:
```bash
# Test your server
curl https://www.ssllabs.com/ssltest/analyze.html?d=your-domain.com
```

2. OpenSSL Verification:
```bash
# Verify certificate
openssl x509 -in server.crt -text -noout

# Test connection
openssl s_client -connect localhost:8443 -tls1_2
```

## Security Considerations

1. **Certificate Storage**
- Use secure storage for private keys
- Implement proper access controls
- Consider using HSM for production

2. **Updates and Maintenance**
- Keep TLS libraries updated
- Monitor security advisories
- Plan certificate renewals
- Regular security audits

3. **Logging and Monitoring**
- Log TLS handshake failures
- Monitor certificate expiration
- Track security header compliance
- Alert on security events
