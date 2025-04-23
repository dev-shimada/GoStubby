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

[Rest of the content remains the same...]
