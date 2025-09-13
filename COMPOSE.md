# Servin Compose - Multi-Container Application Management

## Overview

Servin Compose is a tool for defining and running multi-container applications using Servin. With Compose, you can use a YAML file to configure your application's services, networks, and volumes, then create and start all services with a single command.

## Key Features

- **Docker Compose Compatible**: Uses familiar `servin-compose.yml` syntax
- **Service Orchestration**: Automatic dependency resolution and startup ordering
- **Network Management**: Shared networking between services
- **Volume Management**: Named volumes and bind mounts
- **Environment Configuration**: Service-specific environment variables
- **Build Integration**: Build images as part of the compose workflow
- **Log Aggregation**: View logs from all services or specific services
- **Interactive Execution**: Execute commands in running service containers

## Getting Started

### 1. Create a servin-compose.yml file

```yaml
version: '3.8'

services:
  web:
    build: .
    ports:
      - "8080:80"
    environment:
      - NODE_ENV=production
    volumes:
      - .:/app
    depends_on:
      - db

  db:
    image: postgres:13
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    volumes:
      - db_data:/var/lib/postgresql/data

volumes:
  db_data:
```

### 2. Start your application

```bash
servin compose up
```

### 3. Stop your application

```bash
servin compose down
```

## Command Reference

### Core Commands

#### `servin compose up [OPTIONS]`
Create and start services defined in the compose file.

**Options:**
- `-d, --detach`: Run in detached mode (background)
- `-f, --file FILE`: Specify alternate compose file (default: servin-compose.yml)
- `-p, --project-name NAME`: Specify project name (default: directory name)

**Examples:**
```bash
servin compose up                    # Start all services
servin compose up -d                 # Start in background
servin compose up web                # Start only 'web' service
servin compose -f prod-compose.yml up # Use different compose file
```

#### `servin compose down [OPTIONS]`
Stop and remove containers, networks created by up.

**Options:**
- `-v, --volumes`: Remove named volumes
- `-f, --file FILE`: Specify alternate compose file
- `-p, --project-name NAME`: Specify project name

**Examples:**
```bash
servin compose down                  # Stop and remove containers
servin compose down --volumes        # Also remove volumes
```

#### `servin compose ps [OPTIONS]`
List containers for services defined in the compose file.

**Options:**
- `-a, --all`: Show all containers (default shows just running)
- `-f, --file FILE`: Specify alternate compose file
- `-p, --project-name NAME`: Specify project name

#### `servin compose logs [OPTIONS] [SERVICE...]`
View output from services.

**Options:**
- `-f, --follow`: Follow log output
- `-t, --timestamps`: Show timestamps
- `--tail NUM`: Number of lines to show from end of logs

**Examples:**
```bash
servin compose logs                  # Show logs from all services
servin compose logs web              # Show logs from 'web' service
servin compose logs -f web api       # Follow logs from multiple services
servin compose logs --tail 50 web    # Show last 50 lines
```

#### `servin compose exec [OPTIONS] SERVICE COMMAND [ARG...]`
Execute a command in a running service container.

**Options:**
- `-i, --interactive`: Keep STDIN open

**Examples:**
```bash
servin compose exec web sh           # Open shell in web service
servin compose exec db psql -U postgres # Run psql in db service
servin compose exec web npm test     # Run tests in web service
```

## Compose File Reference

### File Structure

```yaml
version: '3.8'          # Compose file format version

services:               # Define services (containers)
  service_name:
    # Service configuration

networks:               # Define custom networks (optional)
  network_name:
    # Network configuration

volumes:                # Define named volumes (optional)
  volume_name:
    # Volume configuration
```

### Service Configuration

#### Image and Build

```yaml
services:
  app:
    image: nginx:alpine           # Use existing image
    
  custom:
    build:                        # Build from source
      context: .                  # Build context
      dockerfile: Buildfile       # Custom buildfile name
      args:                       # Build arguments
        VERSION: "1.0"
```

#### Ports

```yaml
services:
  web:
    ports:
      - "8080:80"                 # Host:Container
      - "443:443"
      - "9000"                    # Random host port
```

#### Volumes

```yaml
services:
  app:
    volumes:
      - ./src:/app/src            # Bind mount
      - app_data:/app/data        # Named volume
      - /tmp:/tmp:ro              # Read-only mount
```

#### Environment Variables

```yaml
services:
  app:
    environment:
      - NODE_ENV=production       # Array format
      - DEBUG=1
    # OR
    environment:                  # Map format
      NODE_ENV: production
      DEBUG: 1
```

#### Networks

```yaml
services:
  web:
    networks:
      - frontend                  # Connect to custom network
      - backend
```

#### Dependencies

```yaml
services:
  web:
    depends_on:
      - api                       # Start api before web
      - db
  api:
    depends_on:
      - db
```

#### Advanced Configuration

```yaml
services:
  app:
    command: ["npm", "start"]     # Override default command
    entrypoint: ["/entrypoint.sh"] # Override entrypoint
    working_dir: /app             # Set working directory
    user: "1000:1000"             # Run as specific user
    hostname: myapp               # Set container hostname
    restart: unless-stopped       # Restart policy
    labels:                       # Container labels
      traefik.enable: "true"
    expose:                       # Expose ports (internal only)
      - "3000"
```

### Networks

```yaml
networks:
  frontend:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
          gateway: 172.20.0.1
  backend:
    driver: bridge
```

### Volumes

```yaml
volumes:
  db_data:                        # Basic named volume
  app_logs:
    driver: local                 # Specify driver
    labels:
      description: "Application logs"
```

## Use Cases and Patterns

### Web Application Stack

```yaml
version: '3.8'

services:
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - app

  app:
    build: .
    environment:
      DATABASE_URL: postgresql://postgres:password@db:5432/myapp
    volumes:
      - .:/app
    depends_on:
      - db

  db:
    image: postgres:13
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    volumes:
      - db_data:/var/lib/postgresql/data

volumes:
  db_data:
```

### Microservices Architecture

```yaml
version: '3.8'

services:
  api-gateway:
    build: ./gateway
    ports:
      - "80:8080"
    depends_on:
      - user-service
      - product-service

  user-service:
    build: ./user-service
    environment:
      DATABASE_URL: postgresql://postgres:password@user-db:5432/users
    depends_on:
      - user-db

  product-service:
    build: ./product-service
    environment:
      DATABASE_URL: postgresql://postgres:password@product-db:5432/products
    depends_on:
      - product-db

  user-db:
    image: postgres:13
    environment:
      POSTGRES_DB: users
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    volumes:
      - user_data:/var/lib/postgresql/data

  product-db:
    image: postgres:13
    environment:
      POSTGRES_DB: products
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    volumes:
      - product_data:/var/lib/postgresql/data

volumes:
  user_data:
  product_data:
```

### Development Environment

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Buildfile.dev
    ports:
      - "3000:3000"
    environment:
      NODE_ENV: development
      HOT_RELOAD: "true"
    volumes:
      - .:/app
      - /app/node_modules
    depends_on:
      - db
      - redis

  db:
    image: postgres:13
    environment:
      POSTGRES_DB: myapp_dev
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - db_data:/var/lib/postgresql/data

  redis:
    image: redis:6-alpine
    ports:
      - "6379:6379"

volumes:
  db_data:
```

## Best Practices

### File Organization

- Keep compose files in project root
- Use descriptive service names
- Group related services logically
- Use `.env` files for environment-specific values

### Service Design

- One process per service
- Use health checks for service readiness
- Implement graceful shutdown handling
- Use specific image tags (avoid `latest`)

### Data Management

- Use named volumes for persistent data
- Use bind mounts for development
- Back up important data volumes
- Set appropriate file permissions

### Security

- Don't hardcode sensitive data in compose files
- Use environment variables for secrets
- Run services as non-root users when possible
- Use read-only mounts where appropriate

### Development Workflow

```bash
# Start development environment
servin compose up -d

# View logs
servin compose logs -f

# Run tests
servin compose exec app npm test

# Access database
servin compose exec db psql -U postgres

# Rebuild and restart service
servin compose up --build app

# Clean up
servin compose down --volumes
```

## Platform Considerations

### Linux
- Full networking support with service-to-service communication
- Complete volume mounting capabilities
- Production-ready deployment

### Windows/macOS
- Limited networking (development mode)
- Volume mounting with platform-specific paths
- Ideal for development and testing

### Cross-Platform Development
1. Develop and test with Compose on any platform
2. Use relative paths in volume mounts
3. Deploy to Linux for production

## Troubleshooting

### Common Issues

**Services won't start:**
- Check image availability: `servin image ls`
- Verify Buildfile syntax
- Check port conflicts

**Network connectivity issues:**
- Ensure services are on same network
- Use service names, not IP addresses
- Check firewall settings (Linux)

**Volume mount problems:**
- Verify file permissions
- Check path syntax for your platform
- Ensure directories exist

**Build failures:**
- Check Buildfile syntax
- Verify build context
- Review build arguments

### Debug Commands

```bash
# Check service status
servin compose ps

# View service logs
servin compose logs service_name

# Inspect service configuration
servin compose config

# Execute commands in containers
servin compose exec service_name sh
```

## Migration from Docker Compose

Servin Compose is designed to be compatible with Docker Compose files. Most compose files should work with minimal or no changes:

1. Rename `docker-compose.yml` to `servin-compose.yml`
2. Replace `Dockerfile` references with `Buildfile`
3. Update any Docker-specific features to Servin equivalents
4. Test with `servin compose up`

## Examples

See the `examples/` directory for complete compose file examples:

- `examples/servin-compose.yml` - Full-featured web application
- `examples/simple-compose.yml` - Basic two-service setup

## Limitations

- No support for Docker Swarm features
- Limited networking on Windows/macOS
- No built-in secret management
- No support for multiple compose files per project

## Future Enhancements

- Secret management integration
- Multi-file compose support
- Service scaling capabilities
- Enhanced networking features
- Registry integration for image pulling
