# Image Management Enhancement

## Overview

Enhanced the Servin container runtime with a comprehensive image management system that supports proper container image storage, import/export, and rootfs creation from images.

## Key Features Implemented

### 1. Image Storage & Metadata Management
- **Local Image Storage**: Images stored in `~/.servin/images` (Windows) or `/var/lib/servin/images` (Linux)
- **JSON Index**: Centralized image metadata storage with `index.json`
- **Image Metadata**: Comprehensive image information including:
  - Unique image ID (SHA256-based)
  - Repository tags (name:tag format)
  - Creation timestamp and size
  - Layer information
  - Configuration (environment, commands, working directory)
  - Custom metadata and labels

### 2. Image Commands
- **`servin image ls`**: List all available images with repository, tag, ID, creation time, and size
- **`servin image import TARBALL NAME:TAG`**: Import container images from tarball files
- **`servin image rm IMAGE`**: Remove images by name:tag or ID
- **`servin image inspect IMAGE`**: Display detailed image information
- **`servin image pull IMAGE`**: Placeholder for future registry support

### 3. Enhanced RootFS Creation
- **Image-based RootFS**: Containers can now be created from imported images
- **Automatic Fallback**: If image is not found, falls back to basic rootfs creation
- **Efficient Copying**: Recursive directory copying from image to container rootfs
- **Cross-platform Support**: Windows stubs with Linux implementation

### 4. Tarball Import/Export System
- **Compressed Tarball Support**: Handles .tar.gz and .tgz files
- **Security**: Path traversal protection during extraction
- **File Permissions**: Preserves original file permissions and ownership
- **Symlink Support**: Handles symbolic links with warnings for failures

## Architecture

### Image Manager (`pkg/image/`)
```
pkg/image/
├── image.go      # Core image management, storage, and metadata
└── utils.go      # Tarball handling, ID generation, utilities
```

### Enhanced RootFS (`pkg/rootfs/`)
- Updated to integrate with image manager
- Automatic image resolution and rootfs creation
- Backward compatibility with existing basic rootfs creation

### CLI Commands (`cmd/image.go`)
- Complete image management CLI with subcommands
- Consistent interface with other Servin commands
- Proper error handling and user feedback

## Usage Examples

### Import an Image
```bash
# Import from tarball
servin image import alpine-base.tar.gz alpine:latest

# Import with custom tag
servin image import ubuntu-20.04.tgz ubuntu:20.04
```

### List Images
```bash
servin image ls
# Output:
# REPOSITORY  TAG     IMAGE ID      CREATED        SIZE
# alpine      latest  a1b2c3d4e5f6  2 hours ago    5.6 MB
# ubuntu      20.04   f6e5d4c3b2a1  1 day ago      72.8 MB
```

### Inspect Image
```bash
servin image inspect alpine:latest
# Shows detailed image information including metadata, environment, commands
```

### Create Container from Image
```bash
# Use specific image for container
servin run --name web-server alpine:latest /bin/sh

# Falls back to basic rootfs if image not found
servin run --name test-container unknown:image /bin/bash
```

### Remove Images
```bash
# Remove by tag
servin image rm alpine:latest

# Remove by ID
servin image rm a1b2c3d4e5f6

# Remove multiple images
servin image rm alpine:latest ubuntu:20.04
```

## Implementation Details

### Image ID Generation
- Uses SHA256 hash of name:tag combination plus process ID for uniqueness
- Truncated to 16 characters for display (full ID stored internally)

### Storage Structure
```
~/.servin/images/           # Image storage directory
├── index.json             # Image metadata index
├── a1b2c3d4e5f6/          # Image directory (by ID)
│   ├── bin/               # Extracted filesystem
│   ├── etc/
│   └── ...
└── f6e5d4c3b2a1/          # Another image
    └── ...
```

### Metadata Format
```json
{
  "id": "a1b2c3d4e5f6",
  "repo_tags": ["alpine:latest"],
  "created": "2025-09-13T03:54:30Z",
  "size": 5872640,
  "layers": ["a1b2c3d4e5f6"],
  "rootfs_type": "tarball",
  "rootfs_path": "/var/lib/servin/images/a1b2c3d4e5f6",
  "config": {
    "env": ["PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"],
    "cmd": ["/bin/sh"],
    "working_dir": "/",
    "user": "root"
  },
  "metadata": {
    "source": "tarball",
    "original_path": "alpine-base.tar.gz"
  }
}
```

## Benefits

1. **Proper Image Management**: No longer requires manual rootfs directory setup
2. **Reusable Images**: One image can be used for multiple containers
3. **Efficient Storage**: Images are stored once and reused
4. **Metadata Tracking**: Rich metadata for image inspection and management
5. **Standard Format**: Compatible with future OCI image format support
6. **Cross-platform**: Works on both Windows (dev) and Linux (production)

## Future Enhancements

1. **Registry Support**: Pull images from Docker Hub and other registries
2. **Layer Management**: Support for image layers and layer caching
3. **Image Building**: `servin image build` from Dockerfile
4. **Image Pushing**: Push images to registries
5. **OCI Compatibility**: Full OCI image format support
6. **Image Signing**: Security features for image verification

## Testing

The image management system has been tested with:
- ✅ Image import from compressed tarballs
- ✅ Image listing with proper formatting
- ✅ Image inspection with detailed metadata
- ✅ Image removal with cleanup
- ✅ Container creation from images
- ✅ Fallback to basic rootfs when image not found
- ✅ Cross-platform compatibility (Windows dev mode)

This enhancement significantly improves the container runtime's usability and brings it closer to production-ready container management systems.
