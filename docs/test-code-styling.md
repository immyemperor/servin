---
layout: default
title: Test Code Styling
permalink: /test-code-styling/
---

# Test Code Styling

This is a test page to verify that code blocks are properly styled.

## Inline Code
Here is some `inline code` that should have proper background and text colors.

## Bash Commands
```bash
servin containers create ubuntu:latest
servin containers list
servin containers start my-container
```

## Docker Commands
```dockerfile
FROM ubuntu:latest
RUN apt-get update && apt-get install -y curl
COPY . /app
WORKDIR /app
CMD ["./start.sh"]
```

## YAML Configuration
```yaml
version: '3.8'
services:
  web:
    image: nginx:latest
    ports:
      - "80:80"
    volumes:
      - ./html:/usr/share/nginx/html
```

## JSON Configuration
```json
{
  "name": "servin-container",
  "version": "1.0.0",
  "runtime": {
    "type": "runc",
    "options": {
      "systemdCgroup": true
    }
  }
}
```

## Go Code
```go
package main

import (
    "fmt"
    "log"
)

func main() {
    fmt.Println("Hello, Servin!")
    if err := startContainer(); err != nil {
        log.Fatal(err)
    }
}
```

## JavaScript
```javascript
const container = {
    name: 'web-server',
    image: 'nginx:latest',
    ports: ['80:80']
};

async function createContainer(config) {
    try {
        const result = await servin.containers.create(config);
        console.log('Container created:', result.id);
        return result;
    } catch (error) {
        console.error('Failed to create container:', error);
        throw error;
    }
}
```

## Plain Text Code Block
```
This is a plain text code block
without syntax highlighting.
It should still have proper
background and text styling.
```

## Code in Headers

### Running `servin --help` Command

The `servin` CLI provides the following commands:

#### Container Commands with `servin containers`

You can use `servin containers list` to see all containers.
