package vfs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LinuxVFS implements VirtualFileSystem for Linux with namespace support
type LinuxVFS struct {
	containerRoots map[string]string
	workingDirs    map[string]string
}

// Initialize sets up the filesystem for a container on Linux
func (l *LinuxVFS) Initialize(containerID string, imageRootfs string) error {
	if l.containerRoots == nil {
		l.containerRoots = make(map[string]string)
	}
	if l.workingDirs == nil {
		l.workingDirs = make(map[string]string)
	}

	// On Linux, we can use the actual container rootfs path
	// This would typically be managed by the container runtime
	containerDir := filepath.Join("/var/lib/servin/containers", containerID)
	rootfsDir := filepath.Join(containerDir, "rootfs")

	if imageRootfs != "" {
		rootfsDir = imageRootfs
	}

	l.containerRoots[containerID] = rootfsDir
	l.workingDirs[containerID] = "/"

	return nil
}

// Mount activates the container filesystem using Linux namespaces
func (l *LinuxVFS) Mount(containerID string) error {
	// On Linux, this would set up mount namespaces
	// For now, we'll use the simplified approach
	if rootPath, exists := l.containerRoots[containerID]; !exists {
		return fmt.Errorf("container %s not initialized", containerID)
	} else if _, err := os.Stat(rootPath); err != nil {
		return fmt.Errorf("container rootfs not accessible: %w", err)
	}
	return nil
}

// Unmount deactivates the container filesystem
func (l *LinuxVFS) Unmount(containerID string) error {
	// On Linux, this would clean up mount namespaces
	return nil
}

// List returns files and directories in a container path
func (l *LinuxVFS) List(containerID string, path string) ([]FileInfo, error) {
	hostPath, err := l.GetHostPath(containerID, path)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(hostPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", path, err)
	}

	var files []FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		fileInfo := FileInfo{
			Name:        entry.Name(),
			Size:        info.Size(),
			Mode:        info.Mode(),
			ModTime:     info.ModTime(),
			IsDir:       entry.IsDir(),
			Permissions: info.Mode().String(),
		}

		// Get owner/group information
		owner, group := getFileOwnerInfo(info)
		fileInfo.Owner = owner
		fileInfo.Group = group

		files = append(files, fileInfo)
	}

	return files, nil
}

// Read returns a reader for file contents
func (l *LinuxVFS) Read(containerID string, path string) (io.ReadCloser, error) {
	hostPath, err := l.GetHostPath(containerID, path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(hostPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}

	return file, nil
}

// Write writes data to a file in the container
func (l *LinuxVFS) Write(containerID string, path string, data io.Reader) error {
	hostPath, err := l.GetHostPath(containerID, path)
	if err != nil {
		return err
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(hostPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(hostPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer file.Close()

	_, err = io.Copy(file, data)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}

// Stat returns file information
func (l *LinuxVFS) Stat(containerID string, path string) (FileInfo, error) {
	hostPath, err := l.GetHostPath(containerID, path)
	if err != nil {
		return FileInfo{}, err
	}

	info, err := os.Stat(hostPath)
	if err != nil {
		return FileInfo{}, fmt.Errorf("failed to stat %s: %w", path, err)
	}

	fileInfo := FileInfo{
		Name:        info.Name(),
		Size:        info.Size(),
		Mode:        info.Mode(),
		ModTime:     info.ModTime(),
		IsDir:       info.IsDir(),
		Permissions: info.Mode().String(),
	}

	// Get owner/group information
	owner, group := getFileOwnerInfo(info)
	fileInfo.Owner = owner
	fileInfo.Group = group

	return fileInfo, nil
}

// MkDir creates a directory
func (l *LinuxVFS) MkDir(containerID string, path string, mode os.FileMode) error {
	hostPath, err := l.GetHostPath(containerID, path)
	if err != nil {
		return err
	}

	err = os.MkdirAll(hostPath, mode)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}

	return nil
}

// Remove deletes a file or directory
func (l *LinuxVFS) Remove(containerID string, path string) error {
	hostPath, err := l.GetHostPath(containerID, path)
	if err != nil {
		return err
	}

	err = os.RemoveAll(hostPath)
	if err != nil {
		return fmt.Errorf("failed to remove %s: %w", path, err)
	}

	return nil
}

// Copy copies files between containers or within a container
func (l *LinuxVFS) Copy(srcContainerID, srcPath, dstContainerID, dstPath string) error {
	srcHostPath, err := l.GetHostPath(srcContainerID, srcPath)
	if err != nil {
		return err
	}

	dstHostPath, err := l.GetHostPath(dstContainerID, dstPath)
	if err != nil {
		return err
	}

	// Use efficient copy methods available on Linux
	return l.copyPath(srcHostPath, dstHostPath)
}

// copyPath copies a file or directory
func (l *LinuxVFS) copyPath(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		return l.copyDir(src, dst)
	} else {
		return l.copyFile(src, dst, srcInfo.Mode())
	}
}

// copyDir recursively copies a directory
func (l *LinuxVFS) copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return l.copyFile(path, dstPath, info.Mode())
	})
}

// copyFile copies a single file
func (l *LinuxVFS) copyFile(src, dst string, mode os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// Move moves/renames files within a container
func (l *LinuxVFS) Move(containerID string, srcPath, dstPath string) error {
	srcHostPath, err := l.GetHostPath(containerID, srcPath)
	if err != nil {
		return err
	}

	dstHostPath, err := l.GetHostPath(containerID, dstPath)
	if err != nil {
		return err
	}

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(dstHostPath), 0755); err != nil {
		return err
	}

	err = os.Rename(srcHostPath, dstHostPath)
	if err != nil {
		return fmt.Errorf("failed to move %s to %s: %w", srcPath, dstPath, err)
	}

	return nil
}

// Chmod changes file permissions
func (l *LinuxVFS) Chmod(containerID string, path string, mode os.FileMode) error {
	hostPath, err := l.GetHostPath(containerID, path)
	if err != nil {
		return err
	}

	err = os.Chmod(hostPath, mode)
	if err != nil {
		return fmt.Errorf("failed to chmod %s: %w", path, err)
	}

	return nil
}

// Find searches for files matching criteria
func (l *LinuxVFS) Find(containerID string, basePath string, name string, recursive bool) ([]string, error) {
	hostBasePath, err := l.GetHostPath(containerID, basePath)
	if err != nil {
		return nil, err
	}

	var matches []string

	walkFunc := func(hostPath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue on errors
		}

		// Convert back to container path
		relPath, err := filepath.Rel(hostBasePath, hostPath)
		if err != nil {
			return nil
		}

		containerPath := filepath.Join(basePath, relPath)
		if basePath == "/" {
			containerPath = "/" + relPath
		}

		// Check if name matches
		if name == "" || strings.Contains(info.Name(), name) {
			matches = append(matches, containerPath)
		}

		// Control recursion
		if !recursive && info.IsDir() && hostPath != hostBasePath {
			return filepath.SkipDir
		}

		return nil
	}

	err = filepath.Walk(hostBasePath, walkFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return matches, nil
}

// GetWorkingDir returns the current working directory for a container
func (l *LinuxVFS) GetWorkingDir(containerID string) (string, error) {
	if wd, exists := l.workingDirs[containerID]; exists {
		return wd, nil
	}
	return "/", nil
}

// SetWorkingDir sets the working directory for a container
func (l *LinuxVFS) SetWorkingDir(containerID string, path string) error {
	// Validate the path exists
	if _, err := l.Stat(containerID, path); err != nil {
		return fmt.Errorf("directory does not exist: %s", path)
	}

	if l.workingDirs == nil {
		l.workingDirs = make(map[string]string)
	}

	l.workingDirs[containerID] = path
	return nil
}

// GetHostPath converts a container path to the corresponding host path
func (l *LinuxVFS) GetHostPath(containerID string, containerPath string) (string, error) {
	rootPath, exists := l.containerRoots[containerID]
	if !exists {
		return "", fmt.Errorf("container %s not found", containerID)
	}

	// Clean and normalize the container path
	cleanPath := filepath.Clean(containerPath)
	if !filepath.IsAbs(cleanPath) {
		cleanPath = "/" + cleanPath
	}

	// Convert to host path
	hostPath := filepath.Join(rootPath, strings.TrimPrefix(cleanPath, "/"))

	return hostPath, nil
}

// Cleanup removes all resources for a container
func (l *LinuxVFS) Cleanup(containerID string) error {
	if _, exists := l.containerRoots[containerID]; exists {
		// On Linux, we might need to unmount namespaces here
		delete(l.containerRoots, containerID)
	}

	if l.workingDirs != nil {
		delete(l.workingDirs, containerID)
	}

	return nil
}
