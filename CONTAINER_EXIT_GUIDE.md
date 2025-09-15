# Servin Container Exit Guide

## How to Exit from Running Containers

### Interactive Containers (Foreground)

When you run a container in the foreground (default behavior), you can exit using:

**Method 1: Ctrl+C (Recommended)**
```bash
# Start container
./servin run alpine sh

# Exit with Ctrl+C
# Press Ctrl+C to send SIGINT signal
```

**Method 2: Type 'exit' (for shell containers)**
```bash
# If running a shell, you can type exit
./servin run alpine sh
/ # exit
```

**Method 3: Use another terminal**
```bash
# In another terminal, list containers
./servin ls

# Stop the container
./servin stop <container-id-or-name>
```

### Background Containers (Detached)

For containers running in the background:

```bash
# Start container in background
./servin run -d --name mycontainer alpine sleep 60

# Check running containers
./servin ls

# Stop the container
./servin stop mycontainer
```

## Signal Handling

Servin now properly handles signals:

1. **SIGINT (Ctrl+C)**: Graceful termination attempt
2. **SIGTERM**: Graceful termination (used by stop command)
3. **SIGKILL**: Force termination (after 5-second timeout)

## Exit Sequence

When you press Ctrl+C or run stop command:

1. **SIGTERM sent** to container process
2. **5-second grace period** for graceful shutdown
3. **SIGKILL sent** if process doesn't exit
4. **Container cleanup** (filesystem, networks, cgroups)

## Container Run Options

### Foreground (Default)
```bash
# Runs in foreground, shows output, Ctrl+C to exit
./servin run alpine echo "Hello World"
```

### Background/Detached
```bash
# Runs in background, prints container ID
./servin run -d alpine sleep 300
```

### Interactive Shell
```bash
# Interactive shell - exit with 'exit' command or Ctrl+C
./servin run alpine sh
```

## Troubleshooting Exit Issues

If container won't exit:

1. **Check if container is actually running**
   ```bash
   ./servin ls
   ```

2. **Force stop with different terminal**
   ```bash
   ./servin stop <container-id>
   ```

3. **Check for processes in namespace** (Linux only)
   ```bash
   # This shows processes in container
   ps aux | grep <container-command>
   ```

4. **Kill the container process** (emergency)
   ```bash
   # Find the main servin process
   ps aux | grep servin
   # Kill it with SIGKILL
   kill -9 <pid>
   ```

## Examples

### Quick Test Container
```bash
# Run a simple command (exits automatically)
./servin run alpine echo "Hello from container"
```

### Long-Running Container
```bash
# Run container in background
./servin run -d --name web-server alpine sleep 3600

# Check it's running
./servin ls

# Stop it when done
./servin stop web-server
```

### Interactive Shell
```bash
# Start interactive shell
./servin run alpine sh

# Inside container:
/ # ls
/ # ps aux
/ # exit  # or press Ctrl+C
```

### Container with Specific Settings
```bash
# Run with resource limits and exit instruction
./servin run \
  --name myapp \
  --memory 512m \
  --env APP_ENV=test \
  alpine sh

# Exit with Ctrl+C when done
```

## Best Practices

1. **Use detached mode** for long-running services
2. **Use descriptive names** for easier container management
3. **Clean up stopped containers** regularly
4. **Use Ctrl+C** for quick exit from interactive containers
5. **Use proper signals** for graceful application shutdown

The improved signal handling ensures that containers can be exited cleanly without hanging or requiring force termination.