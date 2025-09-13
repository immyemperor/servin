# Image Build Command

The `servin build` command allows you to build container images from a Buildfile (similar to a Dockerfile). This provides a familiar interface for creating custom images with layers, metadata, and configurations.

## Usage

```bash
servin build [OPTIONS] PATH
```

## Arguments

- `PATH`: Build context directory containing the Buildfile and source files

## Options

- `-t, --tag string`: Name and optionally a tag in the 'name:tag' format
- `-f, --file string`: Name of the Buildfile (default 'Buildfile')
- `--no-cache`: Do not use cache when building the image (placeholder for future caching)
- `-q, --quiet`: Suppress the build output and print image ID on success
- `--build-arg stringArray`: Set build-time variables
- `--label stringArray`: Set metadata for an image

## Buildfile Instructions

### FROM
Specifies the base image for the build.

```dockerfile
FROM alpine:latest
FROM scratch
```

### RUN
Executes commands during the build process.

```dockerfile
RUN apk add --no-cache curl
RUN apt-get update && apt-get install -y nginx
```

### COPY
Copies files from the build context to the image.

```dockerfile
COPY . /app
COPY src/ /app/src/
COPY config.json /etc/myapp/
```

### ADD
Similar to COPY but with additional features (URL support, auto-extraction).

```dockerfile
ADD . /app
ADD archive.tar.gz /opt/
```

### WORKDIR
Sets the working directory for subsequent instructions.

```dockerfile
WORKDIR /app
WORKDIR /var/log
```

### ENV
Sets environment variables.

```dockerfile
# key=value format
ENV NODE_ENV=production
ENV DEBUG=true
ENV PATH=/usr/local/bin:$PATH

# key value format
ENV NODE_ENV production
ENV API_URL https://api.example.com
```

### EXPOSE
Declares ports that the container will listen on.

```dockerfile
EXPOSE 80
EXPOSE 443
EXPOSE 8080/tcp
```

### CMD
Specifies the default command to run when the container starts.

```dockerfile
CMD ["nginx", "-g", "daemon off;"]
CMD ["node", "server.js"]
CMD echo "Hello World"
```

### ENTRYPOINT
Configures the container to run as an executable.

```dockerfile
ENTRYPOINT ["./docker-entrypoint.sh"]
ENTRYPOINT ["nginx"]
```

### LABEL
Adds metadata to the image.

```dockerfile
LABEL maintainer="user@example.com"
LABEL version="1.0.0"
LABEL description="My application"
```

### USER
Sets the username or UID for running subsequent instructions.

```dockerfile
USER nginx
USER 1001
USER appuser:appgroup
```

### VOLUME
Creates mount points for external volumes.

```dockerfile
VOLUME ["/data"]
VOLUME /var/log /var/cache
```

## Examples

### Basic Build

```bash
# Build from current directory with default Buildfile
servin build .

# Build with custom tag
servin build -t myapp:latest .

# Build with custom Buildfile name
servin build -f MyCustomBuildfile .
```

### Build Arguments

```bash
# Pass build-time variables
servin build --build-arg VERSION=1.0.0 --build-arg ENV=production .

# Use build args in Buildfile
# Buildfile:
# ARG VERSION=latest
# FROM alpine:$VERSION
# ENV APP_VERSION=$VERSION
```

### Labels and Metadata

```bash
# Add labels during build
servin build --label version=1.0.0 --label maintainer=user@example.com .
```

### Quiet Mode

```bash
# Suppress output, only show image ID
servin build -q .
# Output: 1757719823787012100
```

## Sample Buildfiles

### Simple Web Application

```dockerfile
# Simple web server
FROM alpine:latest

# Install web server
RUN apk add --no-cache nginx

# Copy configuration
COPY nginx.conf /etc/nginx/nginx.conf
COPY html/ /var/www/html/

# Set working directory
WORKDIR /var/www/html

# Create user
RUN adduser -D -s /bin/sh nginx

# Switch to non-root user
USER nginx

# Expose port
EXPOSE 80

# Start server
CMD ["nginx", "-g", "daemon off;"]
```

### Node.js Application

```dockerfile
# Node.js app
FROM scratch

# Set working directory
WORKDIR /app

# Set environment
ENV NODE_ENV=production
ENV PORT=3000

# Copy application files
COPY package.json .
COPY src/ ./src/

# Set labels
LABEL maintainer="dev@example.com"
LABEL version="1.0.0"
LABEL description="Node.js web application"

# Create user
USER node

# Expose port
EXPOSE 3000

# Start application
CMD ["node", "src/server.js"]
```

### Multi-stage Build Simulation

```dockerfile
# Build stage
FROM alpine:latest
WORKDIR /build
COPY src/ .
RUN apk add --no-cache build-essential
RUN make compile

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=build /build/app .
USER appuser
EXPOSE 8080
CMD ["./app"]
```

## Build Process

The build process follows these steps:

1. **Context Resolution**: Resolve and validate the build context path
2. **Buildfile Parsing**: Parse the Buildfile and validate syntax
3. **Instruction Processing**: Execute each instruction in order
4. **Image Creation**: Create the final image with all layers and metadata
5. **Tagging**: Apply the specified tag to the built image
6. **Storage**: Save the image to the local image store

## Build Context

The build context is the directory containing the Buildfile and source files:

```
my-project/
├── Buildfile
├── src/
│   ├── main.go
│   └── config.json
├── static/
│   └── index.html
└── scripts/
    └── entrypoint.sh
```

All files in the build context are available for COPY and ADD instructions.

## Build Arguments

Build arguments allow parameterization of builds:

```bash
# Define in Buildfile
ARG VERSION=latest
ARG ENVIRONMENT=development
FROM alpine:$VERSION
ENV APP_ENV=$ENVIRONMENT

# Pass at build time
servin build --build-arg VERSION=3.14 --build-arg ENVIRONMENT=production .
```

## Error Handling

The build command provides detailed error messages:

```bash
# Missing Buildfile
servin build .
# Error: [NOT_FOUND] build: Buildfile 'Buildfile' not found

# Invalid instruction
# Error: step 3 failed: INVALID instruction requires an argument

# Build context not found
servin build /nonexistent
# Error: [NOT_FOUND] build: build context '/nonexistent' not found
```

## Platform Support

### Linux
Full functionality with complete instruction support.

### Windows/macOS  
Full build functionality available. Built images can be used for development workflows and cross-platform testing.

## Image Storage

Built images are stored in the local image repository:

- **Linux**: `/var/lib/servin/images/`
- **Windows**: `%USERPROFILE%\.servin\images\`
- **macOS**: `~/.servin/images/`

## Integration with Other Commands

The build command integrates seamlessly with other servin commands:

```bash
# Build and run
servin build -t myapp:latest .
servin run myapp:latest

# Build and inspect
servin build -t myapp:latest .
servin image inspect myapp:latest

# Build and tag
servin build -t myapp:latest .
servin image tag myapp:latest myapp:v1.0.0
```

## Performance Considerations

- **Build Context Size**: Minimize the build context by using .dockerignore-style filtering
- **Layer Optimization**: Combine RUN instructions to reduce layer count
- **Caching**: Future versions will support build caching for faster rebuilds

## Differences from Docker

While similar to Docker's build process, servin build has some differences:

- **Simplified Implementation**: Focus on core functionality
- **Cross-platform Development**: Emphasis on development workflows
- **Metadata Rich**: Extensive build metadata storage
- **Educational Focus**: Clear, understandable build process

## Future Enhancements

Planned improvements for the build command:

1. **Build Caching**: Layer caching for faster rebuilds
2. **Multi-stage Builds**: Support for complex build workflows
3. **BuildKit Integration**: Advanced build features
4. **Squash Layers**: Optimize final image size
5. **Build Secrets**: Secure handling of sensitive build data
6. **Parallel Builds**: Concurrent instruction execution

## Troubleshooting

### Common Issues

**Buildfile not found**
- Ensure the Buildfile exists in the build context
- Check the `-f` flag for custom Buildfile names

**Context path errors**
- Use absolute paths or relative paths from current directory
- Ensure build context directory exists

**Instruction errors**
- Check instruction syntax against supported formats
- Verify file paths in COPY/ADD instructions exist

**Permission errors**
- Ensure read permissions on build context files
- Check write permissions for image storage directory

## Related Commands

- `servin image ls`: List built images
- `servin image inspect IMAGE`: View image details and build metadata
- `servin image tag SOURCE TARGET`: Tag built images
- `servin run IMAGE`: Run containers from built images
- `servin image rm IMAGE`: Remove built images
