# Improved Error Handling and Logging

## Overview

Servin now includes comprehensive error handling and logging capabilities to provide better debugging, monitoring, and troubleshooting support across all platforms.

## Logging System

### Features
- **Structured Logging**: Timestamped logs with severity levels and caller information
- **Multiple Log Levels**: DEBUG, INFO, WARN, ERROR, FATAL
- **Cross-Platform File Logging**: Platform-appropriate log file locations
- **Console and File Output**: Simultaneous logging to console and file
- **Verbose Mode**: Enhanced logging with caller information

### Log Levels

| Level | Description | Use Case |
|-------|-------------|----------|
| **DEBUG** | Detailed diagnostic information | Development and troubleshooting |
| **INFO** | General informational messages | Normal operation tracking |
| **WARN** | Warning conditions that don't prevent operation | Potential issues |
| **ERROR** | Error conditions that affect operation | Failed operations |
| **FATAL** | Critical errors that cause program termination | Unrecoverable failures |

### Global Logging Flags

```bash
# Set log level
servin --log-level debug volume create myvolume

# Enable verbose output with caller information
servin --verbose volume ls

# Specify custom log file
servin --log-file /custom/path/servin.log image ls

# Combine flags
servin --verbose --log-level debug --log-file ./debug.log run alpine echo test
```

### Log File Locations

| Platform | Default Log Path |
|----------|------------------|
| **Linux** | `/var/log/servin/servin.log` |
| **Windows** | `%USERPROFILE%\.servin\logs\servin.log` |
| **macOS** | `~/.servin/logs/servin.log` |

### Example Log Output

```
2025-09-13 04:36:47 [DEBUG] [volume.go:123] Creating volume: test-volume (driver: local)
2025-09-13 04:36:47 [INFO] [volume.go:156] Volume 'test-volume' created successfully at C:\Users\immye\.servin\volumes\test-volume
2025-09-13 04:36:57 [WARN] [volume.go:134] Attempted to create volume that already exists: test-volume
2025-09-13 04:36:57 [ERROR] [volume.go:201] Failed to create volume 'test-volume': [CONFLICT] CreateVolume: volume 'test-volume' already exists
```

## Error Handling System

### Structured Errors

Servin now uses structured errors with:
- **Error Types**: Categorized error classification
- **Operation Context**: Which operation failed
- **Underlying Cause**: Root cause error chaining
- **Contextual Data**: Additional metadata for debugging
- **Stack Traces**: Full call stack for debugging

### Error Types

| Type | Description | Examples |
|------|-------------|----------|
| **SYSTEM** | System-level errors | OS permission issues, kernel features |
| **IO** | Input/output errors | File read/write, directory creation |
| **NETWORK** | Network-related errors | Bridge creation, iptables configuration |
| **PERMISSION** | Permission/privilege errors | Root access, file permissions |
| **CONTAINER** | Container lifecycle errors | Process creation, namespace setup |
| **IMAGE** | Image management errors | Tarball extraction, metadata parsing |
| **VOLUME** | Volume management errors | Directory creation, volume conflicts |
| **VALIDATION** | Input validation errors | Invalid names, malformed arguments |
| **NOT_FOUND** | Resource not found errors | Missing containers, images, volumes |
| **CONFLICT** | Resource conflict errors | Duplicate names, resource in use |
| **CONFIG** | Configuration errors | Invalid settings, missing configuration |

### Error Examples

#### Validation Error
```bash
$ servin volume create ""
Error: [VALIDATION] CreateVolume: volume name cannot be empty
```

#### Conflict Error
```bash
$ servin volume create existing-volume
Error: [VOLUME] runVolumeCreate: failed to create volume 'existing-volume' (caused by: [CONFLICT] CreateVolume: volume 'existing-volume' already exists)
```

#### Permission Error
```bash
$ servin run alpine echo test  # without sudo on Linux
Error: [PERMISSION] checkRoot: root privileges required on Linux
```

#### Not Found Error
```bash
$ servin volume rm nonexistent
Error: [VOLUME] runVolumeRemove: errors occurred:
failed to remove volume 'nonexistent': volume 'nonexistent' not found
```

## Implementation Details

### Logger Package (`pkg/logger/`)

- **Global Logger**: Default logger with sensible defaults
- **Custom Loggers**: Create loggers with specific configurations
- **File Rotation**: Automatic log file management
- **Multi-Writer**: Simultaneous console and file output
- **Platform Awareness**: OS-specific log paths

### Error Package (`pkg/errors/`)

- **ServinError Struct**: Rich error information with context
- **Error Constructors**: Type-specific error creation functions
- **Error Wrapping**: Chain errors with additional context
- **Context Addition**: Add metadata to errors
- **Stack Traces**: Capture and format call stacks

### CLI Integration

- **Consistent Error Display**: Structured error messages in CLI
- **Contextual Help**: Show relevant help on errors
- **Logging Flags**: Global flags for logging configuration
- **Verbose Mode**: Enhanced output for debugging

## Usage Examples

### Development Debugging

```bash
# Enable full debug logging
servin --verbose --log-level debug volume create debug-volume

# Monitor logs in real-time
tail -f ~/.servin/logs/servin.log

# Custom log file for specific operations
servin --log-file ./container-debug.log run --verbose alpine ls
```

### Production Monitoring

```bash
# Set appropriate log level for production
servin --log-level warn container-command

# Parse structured logs
grep "ERROR" /var/log/servin/servin.log | tail -20

# Monitor specific operations
grep "volume" /var/log/servin/servin.log
```

### Error Investigation

```bash
# Get detailed error information
servin --verbose --log-level debug failing-command 2>&1 | tee error-investigation.log

# Check error types
grep "\[ERROR\]" ~/.servin/logs/servin.log | grep "PERMISSION"
```

## Benefits

### For Developers
- **Enhanced Debugging**: Detailed logs and stack traces
- **Context Awareness**: Rich error information with metadata
- **Development Mode**: Verbose logging and debug information

### For Operations
- **Monitoring**: Structured logs for operational monitoring
- **Troubleshooting**: Clear error categorization and context
- **Audit Trail**: Complete operation logging

### For Users
- **Clear Error Messages**: Understandable error descriptions
- **Helpful Context**: Relevant information for resolving issues
- **Consistent Experience**: Uniform error handling across commands

## Future Enhancements

### Planned Features
- **Log Rotation**: Automatic log file rotation and cleanup
- **Remote Logging**: Support for remote log aggregation
- **Metrics Integration**: Export metrics from logs
- **Error Recovery**: Automatic retry mechanisms
- **Performance Monitoring**: Operation timing and performance logs
