package vfs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// MacOSVFS implements VirtualFileSystem for macOS using directory overlays and chroot simulation
// Since macOS doesn't have Linux namespaces, we simulate container filesystem isolation
// using directory structures and path mapping
type MacOSVFS struct {
	containerRoots map[string]string // containerID -> root path mapping
	workingDirs    map[string]string // containerID -> working directory mapping
}

// Initialize sets up the filesystem for a container on macOS
func (m *MacOSVFS) Initialize(containerID string, imageRootfs string) error {
	if m.containerRoots == nil {
		m.containerRoots = make(map[string]string)
	}
	if m.workingDirs == nil {
		m.workingDirs = make(map[string]string)
	}

	// Create container-specific directory structure
	containerDir := filepath.Join("/tmp/servin/containers", containerID)
	rootfsDir := filepath.Join(containerDir, "rootfs")

	// Create the directory structure
	if err := os.MkdirAll(rootfsDir, 0755); err != nil {
		return fmt.Errorf("failed to create container directory: %w", err)
	}

	// If we have an image rootfs, create a copy or overlay
	if imageRootfs != "" && imageRootfs != rootfsDir {
		if err := m.createOverlay(imageRootfs, rootfsDir); err != nil {
			return fmt.Errorf("failed to create filesystem overlay: %w", err)
		}
	} else {
		// Create a minimal filesystem structure
		if err := m.createMinimalFS(rootfsDir); err != nil {
			return fmt.Errorf("failed to create minimal filesystem: %w", err)
		}
	}

	m.containerRoots[containerID] = rootfsDir
	m.workingDirs[containerID] = "/"

	return nil
}

// createOverlay creates a filesystem overlay (copy-on-write simulation)
func (m *MacOSVFS) createOverlay(source, dest string) error {
	// For simplicity, we'll use a copy approach on macOS
	// In production, this could be optimized with hard links or other techniques

	// Check if source exists and is accessible
	if _, err := os.Stat(source); os.IsNotExist(err) {
		// Source doesn't exist, create minimal structure
		return m.createMinimalFS(dest)
	}

	// Copy the source directory structure
	return m.copyDir(source, dest)
}

// createMinimalFS creates a basic Unix filesystem structure
func (m *MacOSVFS) createMinimalFS(rootPath string) error {
	dirs := []string{
		"bin", "etc", "home", "lib", "lib64", "usr", "usr/bin", "usr/lib",
		"var", "var/log", "tmp", "dev", "proc", "sys", "root",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(rootPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
	}

	// Create some basic files
	files := map[string]string{
		"etc/hostname": "container",
		"etc/hosts":    "127.0.0.1 localhost\n",
		"etc/passwd":   "root:x:0:0:root:/root:/bin/sh\n",
		"etc/group":    "root:x:0:\n",
	}

	for file, content := range files {
		filePath := filepath.Join(rootPath, file)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

// copyDir recursively copies a directory
func (m *MacOSVFS) copyDir(src, dst string) error {
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

		// Copy file
		return m.copyFile(path, dstPath, info.Mode())
	})
}

// copyFile copies a single file
func (m *MacOSVFS) copyFile(src, dst string, mode os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// Mount activates the container filesystem (no-op on macOS)
func (m *MacOSVFS) Mount(containerID string) error {
	// On macOS, mounting is simulated through directory structure
	// Verify the container directory exists
	if rootPath, exists := m.containerRoots[containerID]; !exists {
		return fmt.Errorf("container %s not initialized", containerID)
	} else if _, err := os.Stat(rootPath); err != nil {
		return fmt.Errorf("container rootfs not accessible: %w", err)
	}
	return nil
}

// Unmount deactivates the container filesystem (no-op on macOS)
func (m *MacOSVFS) Unmount(containerID string) error {
	// No actual unmounting needed on macOS
	return nil
}

// List returns files and directories in a container path
func (m *MacOSVFS) List(containerID string, path string) ([]FileInfo, error) {
	hostPath, err := m.GetHostPath(containerID, path)
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

		// Get owner/group information if possible
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			fileInfo.Owner = fmt.Sprintf("%d", stat.Uid)
			fileInfo.Group = fmt.Sprintf("%d", stat.Gid)
		}

		files = append(files, fileInfo)
	}

	return files, nil
}

// Read returns a reader for file contents
func (m *MacOSVFS) Read(containerID string, path string) (io.ReadCloser, error) {
	hostPath, err := m.GetHostPath(containerID, path)
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
func (m *MacOSVFS) Write(containerID string, path string, data io.Reader) error {
	hostPath, err := m.GetHostPath(containerID, path)
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
func (m *MacOSVFS) Stat(containerID string, path string) (FileInfo, error) {
	hostPath, err := m.GetHostPath(containerID, path)
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

	// Get owner/group information if possible
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		fileInfo.Owner = fmt.Sprintf("%d", stat.Uid)
		fileInfo.Group = fmt.Sprintf("%d", stat.Gid)
	}

	return fileInfo, nil
}

// MkDir creates a directory
func (m *MacOSVFS) MkDir(containerID string, path string, mode os.FileMode) error {
	hostPath, err := m.GetHostPath(containerID, path)
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
func (m *MacOSVFS) Remove(containerID string, path string) error {
	hostPath, err := m.GetHostPath(containerID, path)
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
func (m *MacOSVFS) Copy(srcContainerID, srcPath, dstContainerID, dstPath string) error {
	srcHostPath, err := m.GetHostPath(srcContainerID, srcPath)
	if err != nil {
		return err
	}

	dstHostPath, err := m.GetHostPath(dstContainerID, dstPath)
	if err != nil {
		return err
	}

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(dstHostPath), 0755); err != nil {
		return err
	}

	// Get source info
	srcInfo, err := os.Stat(srcHostPath)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		return m.copyDir(srcHostPath, dstHostPath)
	} else {
		return m.copyFile(srcHostPath, dstHostPath, srcInfo.Mode())
	}
}

// Move moves/renames files within a container
func (m *MacOSVFS) Move(containerID string, srcPath, dstPath string) error {
	srcHostPath, err := m.GetHostPath(containerID, srcPath)
	if err != nil {
		return err
	}

	dstHostPath, err := m.GetHostPath(containerID, dstPath)
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
func (m *MacOSVFS) Chmod(containerID string, path string, mode os.FileMode) error {
	hostPath, err := m.GetHostPath(containerID, path)
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
func (m *MacOSVFS) Find(containerID string, basePath string, name string, recursive bool) ([]string, error) {
	hostBasePath, err := m.GetHostPath(containerID, basePath)
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
func (m *MacOSVFS) GetWorkingDir(containerID string) (string, error) {
	if wd, exists := m.workingDirs[containerID]; exists {
		return wd, nil
	}
	return "/", nil
}

// SetWorkingDir sets the working directory for a container
func (m *MacOSVFS) SetWorkingDir(containerID string, path string) error {
	// Validate the path exists
	if _, err := m.Stat(containerID, path); err != nil {
		return fmt.Errorf("directory does not exist: %s", path)
	}

	if m.workingDirs == nil {
		m.workingDirs = make(map[string]string)
	}

	m.workingDirs[containerID] = path
	return nil
}

// GetHostPath converts a container path to the corresponding host path
func (m *MacOSVFS) GetHostPath(containerID string, containerPath string) (string, error) {
	rootPath, exists := m.containerRoots[containerID]
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
func (m *MacOSVFS) Cleanup(containerID string) error {
	if rootPath, exists := m.containerRoots[containerID]; exists {
		// Remove the container directory
		if err := os.RemoveAll(filepath.Dir(rootPath)); err != nil {
			return fmt.Errorf("failed to cleanup container directory: %w", err)
		}

		delete(m.containerRoots, containerID)
	}

	if m.workingDirs != nil {
		delete(m.workingDirs, containerID)
	}

	return nil
}
