# Container Logs Command

The `servin logs` command allows you to fetch and display the logs of running or stopped containers. This command retrieves stdout and stderr output from containers that have been run with log capturing enabled.

## Usage

```bash
servin logs [OPTIONS] CONTAINER
```

## Arguments

- `CONTAINER`: Container ID or name to fetch logs from

## Options

- `-f, --follow`: Follow log output (for running containers)
- `-t, --timestamps`: Show timestamps for each log line
- `--tail string`: Number of lines to show from the end of the logs (default "all")
- `--since string`: Show logs since timestamp (e.g. 2013-01-02T13:23:37Z) or relative (e.g. 42m for 42 minutes)
- `--until string`: Show logs before a timestamp (e.g. 2013-01-02T13:23:37Z) or relative (e.g. 42m for 42 minutes)

## Examples

### Basic Usage

```bash
# Show all logs for a container
servin logs my-container

# Show logs with timestamps
servin logs --timestamps my-container

# Show only the last 10 lines
servin logs --tail 10 my-container

# Follow logs in real-time (for running containers)
servin logs --follow my-container
```

### Time Filtering

```bash
# Show logs from the last hour
servin logs --since 1h my-container

# Show logs from the last 30 minutes
servin logs --since 30m my-container

# Show logs since a specific timestamp
servin logs --since 2025-09-13T04:00:00Z my-container

# Show logs between two timestamps
servin logs --since 2025-09-13T04:00:00Z --until 2025-09-13T05:00:00Z my-container
```

### Combining Options

```bash
# Show last 20 lines with timestamps from the last 2 hours
servin logs --tail 20 --timestamps --since 2h my-container

# Follow logs with timestamps
servin logs --follow --timestamps my-container
```

## Log Storage

Container logs are stored in the servin directory structure:

- **Linux**: `/var/lib/servin/logs/<container-id>/`
- **Windows**: `C:\Users\<username>\.servin\logs\<container-id>\`
- **macOS**: `~/.servin/logs/<container-id>/`

Each container has two log files:
- `stdout.log`: Standard output from the container
- `stderr.log`: Standard error output from the container

## Log Format

Log entries are automatically timestamped when captured from containers:

```
2025-09-13T04:47:00.123456789Z This is a log message
```

When displayed with `--timestamps`, additional metadata is shown:

```
2025-09-13T04:47:00+05:30 [stdout] This is a log message from stdout
2025-09-13T04:47:01+05:30 [stderr] This is an error message from stderr
```

## Platform Support

### Linux
Full functionality with automatic log capture during container execution.

### Windows/macOS
Limited functionality due to namespace limitations. Logs are captured when possible, but full containerization features are not available.

## Error Handling

The logs command provides detailed error messages for common scenarios:

```bash
# Container not found
servin logs non-existent-container
# Error: [NOT_FOUND] container: no container found with this ID or name

# No logs available
servin logs container-without-logs
# No logs available for container container-without-logs
```

## Real-time Following

The `--follow` option enables real-time log streaming for running containers:

```bash
servin logs --follow my-running-container
```

This will:
1. Display existing logs
2. Continue streaming new log entries as they appear
3. Show both stdout and stderr in chronological order
4. Run until interrupted with Ctrl+C

## Performance Considerations

- Large log files may take time to process with time filtering
- The `--tail` option improves performance by limiting output
- Real-time following uses periodic polling (1-second intervals)
- For production use, consider log rotation and external log management

## Integration with Other Commands

The logs command works seamlessly with other servin commands:

```bash
# Run a container and immediately follow its logs
servin run --name web-server nginx &
servin logs --follow web-server

# Check logs after container stops
servin run --name test-app alpine echo "Hello World"
servin logs test-app

# Debug failed containers
servin run --name debug-container alpine /bin/sh -c "echo 'Starting'; exit 1"
servin logs debug-container
```

## Troubleshooting

### No logs available
- Ensure the container was run after logs functionality was implemented
- Check that the container actually produced output
- Verify container exists: `servin ls`

### Permission errors
- Ensure proper permissions to access log directory
- On Linux, may require root access depending on installation

### Encoding issues
- Log files use UTF-8 encoding
- Special characters in container output may display differently based on terminal settings

## Future Enhancements

Planned improvements for the logs command:

1. **Log rotation**: Automatic log file rotation to prevent disk space issues
2. **Remote logging**: Support for sending logs to external logging systems
3. **Structured logging**: JSON format option for machine-readable logs
4. **Log compression**: Automatic compression of old log files
5. **Efficient following**: Use of file system events instead of polling
6. **Log aggregation**: Combined logging across multiple containers

## Related Commands

- `servin run`: Create and run containers (with log capture)
- `servin ls`: List containers to find container IDs/names
- `servin stop`: Stop running containers
- `servin rm`: Remove containers (logs are preserved separately)
