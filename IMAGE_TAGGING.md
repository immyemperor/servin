# Image Tagging Examples

This document demonstrates the image tagging functionality in Servin.

## Basic Tagging

### Tag an existing image with a new name:
```bash
# Tag alpine:test as alpine:v1.0
servin image tag alpine:test alpine:v1.0

# Tag using image ID
servin image tag 45b0a36b30b7 alpine:latest

# Tag with a different repository name
servin image tag alpine:test myapp:production
```

### Automatic :latest tag
```bash
# When no tag is specified, :latest is automatically appended
servin image tag alpine:test myapp
# This creates myapp:latest
```

## Use Cases

### 1. Version Management
```bash
# Create version tags for releases
servin image tag myapp:latest myapp:v1.0.0
servin image tag myapp:latest myapp:v1.0
servin image tag myapp:latest myapp:stable
```

### 2. Environment Tagging
```bash
# Tag images for different environments
servin image tag myapp:latest myapp:development
servin image tag myapp:latest myapp:staging
servin image tag myapp:latest myapp:production
```

### 3. Backup and Archive
```bash
# Create backup tags before updates
servin image tag myapp:latest myapp:backup-$(date +%Y%m%d)
servin image tag myapp:latest myapp:archive
```

### 4. Multi-Repository Tagging
```bash
# Tag the same image under different repository names
servin image tag alpine:test webapp:latest
servin image tag alpine:test microservice:v1
servin image tag alpine:test api-server:production
```

## Error Handling

### Tag already exists
```bash
$ servin image tag alpine:test alpine:latest
Error: failed to tag image: tag 'alpine:latest' already exists on image 45b0a36b30b7
```

### Source image not found
```bash
$ servin image tag nonexistent:tag newimage:tag
Error: failed to tag image: source image not found: image 'nonexistent:tag' not found
```

## Viewing Tagged Images

```bash
# List all images to see tags
servin image ls

# Example output:
REPOSITORY  TAG         IMAGE ID      CREATED         SIZE
alpine      test        45b0a36b30b7  27 minutes ago  279 B
alpine      v1.0        45b0a36b30b7  27 minutes ago  279 B
alpine      latest      45b0a36b30b7  27 minutes ago  279 B
myapp       production  45b0a36b30b7  27 minutes ago  279 B
myapp       latest      45b0a36b30b7  27 minutes ago  279 B
```

## Notes

- All tags pointing to the same image share the same Image ID
- Tags are stored in the image metadata and persist across restarts
- Tagging does not create a copy of the image data, only adds a reference
- Images can be referenced by any of their tags in container run commands
- Removing an image removes all its tags
