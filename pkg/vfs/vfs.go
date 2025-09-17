package vfs

import (
	"context"
	"io"
	"os"
	"time"
)

// FileInfo represents file/directory information
type FileInfo struct {
	Name        string
	Size        int64
	Mode        os.FileMode
	ModTime     time.Time
	IsDir       bool
	Permissions string
	Owner       string
	Group       string
}

// VirtualFileSystem defines the interface for container filesystem operations
// This abstraction allows different implementations for different platforms
type VirtualFileSystem interface {
	// Initialize the filesystem for a container
	Initialize(containerID string, imageRootfs string) error

	// Mount/unmount the container filesystem
	Mount(containerID string) error
	Unmount(containerID string) error

	// File/directory operations
	List(containerID string, path string) ([]FileInfo, error)
	Read(containerID string, path string) (io.ReadCloser, error)
	Write(containerID string, path string, data io.Reader) error
	Stat(containerID string, path string) (FileInfo, error)

	// Directory operations
	MkDir(containerID string, path string, mode os.FileMode) error
	Remove(containerID string, path string) error

	// File operations
	Copy(srcContainerID, srcPath, dstContainerID, dstPath string) error
	Move(containerID string, srcPath, dstPath string) error
	Chmod(containerID string, path string, mode os.FileMode) error

	// Search operations
	Find(containerID string, basePath string, name string, recursive bool) ([]string, error)

	// Working directory operations
	GetWorkingDir(containerID string) (string, error)
	SetWorkingDir(containerID string, path string) error

	// Get the actual host path for a container path (for platform-specific operations)
	GetHostPath(containerID string, containerPath string) (string, error)

	// Cleanup resources
	Cleanup(containerID string) error
}

// VFSProvider creates platform-specific VFS implementations
type VFSProvider interface {
	CreateVFS() VirtualFileSystem
	SupportsNamespaces() bool
	GetPlatform() string
}

// FileOperation represents a file operation for batch processing
type FileOperation struct {
	Type        string // "copy", "move", "delete", "mkdir"
	Source      string
	Destination string
	Mode        os.FileMode
	ContainerID string
}

// VFSManager manages VFS instances and provides high-level operations
type VFSManager struct {
	vfs      VirtualFileSystem
	provider VFSProvider
}

// NewVFSManager creates a new VFS manager with the appropriate provider
func NewVFSManager() (*VFSManager, error) {
	provider := GetPlatformProvider()
	vfs := provider.CreateVFS()

	return &VFSManager{
		vfs:      vfs,
		provider: provider,
	}, nil
}

// GetVFS returns the underlying VFS implementation
func (m *VFSManager) GetVFS() VirtualFileSystem {
	return m.vfs
}

// GetProvider returns the platform provider
func (m *VFSManager) GetProvider() VFSProvider {
	return m.provider
}

// BatchOperations executes multiple file operations in a transaction-like manner
func (m *VFSManager) BatchOperations(ctx context.Context, operations []FileOperation) error {
	// TODO: Implement transaction-like batch operations
	for _, op := range operations {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		switch op.Type {
		case "copy":
			if err := m.vfs.Copy(op.ContainerID, op.Source, op.ContainerID, op.Destination); err != nil {
				return err
			}
		case "move":
			if err := m.vfs.Move(op.ContainerID, op.Source, op.Destination); err != nil {
				return err
			}
		case "delete":
			if err := m.vfs.Remove(op.ContainerID, op.Source); err != nil {
				return err
			}
		case "mkdir":
			if err := m.vfs.MkDir(op.ContainerID, op.Destination, op.Mode); err != nil {
				return err
			}
		}
	}

	return nil
}
