package vfs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GenericVFS provides a fallback VFS implementation for unsupported platforms
type GenericVFS struct {
	containerRoots map[string]string
	workingDirs    map[string]string
}

// WindowsVFS implements VirtualFileSystem for Windows
type WindowsVFS struct {
	containerRoots map[string]string
	workingDirs    map[string]string
}

// Generic VFS Implementation
func (g *GenericVFS) Initialize(containerID string, imageRootfs string) error {
	if g.containerRoots == nil {
		g.containerRoots = make(map[string]string)
	}
	if g.workingDirs == nil {
		g.workingDirs = make(map[string]string)
	}

	// Create a basic container directory structure
	var containerDir string
	if runtime.GOOS == "windows" {
		// Use Windows-appropriate directory - try user temp first
		tempDir := os.Getenv("TEMP")
		if tempDir == "" {
			tempDir = os.Getenv("TMP")
		}
		if tempDir == "" {
			tempDir = "C:\\temp"
		}
		containerDir = filepath.Join(tempDir, "servin", "containers", containerID)
	} else {
		containerDir = filepath.Join("/tmp/servin/containers", containerID)
	}

	rootfsDir := filepath.Join(containerDir, "rootfs")

	// Use appropriate permissions for the platform
	var dirMode os.FileMode = 0755
	if runtime.GOOS == "windows" {
		dirMode = 0666 // Windows doesn't use Unix permissions the same way
	}

	if err := os.MkdirAll(rootfsDir, dirMode); err != nil {
		return fmt.Errorf("failed to create container directory: %w", err)
	}

	g.containerRoots[containerID] = rootfsDir
	g.workingDirs[containerID] = "/"

	return nil
}

func (g *GenericVFS) Mount(containerID string) error {
	if _, exists := g.containerRoots[containerID]; !exists {
		return fmt.Errorf("container %s not initialized", containerID)
	}
	return nil
}

func (g *GenericVFS) Unmount(containerID string) error {
	return nil
}

func (g *GenericVFS) List(containerID string, path string) ([]FileInfo, error) {
	hostPath, err := g.GetHostPath(containerID, path)
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
			Owner:       "unknown",
			Group:       "unknown",
		}

		files = append(files, fileInfo)
	}

	return files, nil
}

func (g *GenericVFS) Read(containerID string, path string) (io.ReadCloser, error) {
	hostPath, err := g.GetHostPath(containerID, path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(hostPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}

	return file, nil
}

func (g *GenericVFS) Write(containerID string, path string, data io.Reader) error {
	hostPath, err := g.GetHostPath(containerID, path)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(hostPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(hostPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer file.Close()

	_, err = io.Copy(file, data)
	return err
}

func (g *GenericVFS) Stat(containerID string, path string) (FileInfo, error) {
	hostPath, err := g.GetHostPath(containerID, path)
	if err != nil {
		return FileInfo{}, err
	}

	info, err := os.Stat(hostPath)
	if err != nil {
		return FileInfo{}, fmt.Errorf("failed to stat %s: %w", path, err)
	}

	return FileInfo{
		Name:        info.Name(),
		Size:        info.Size(),
		Mode:        info.Mode(),
		ModTime:     info.ModTime(),
		IsDir:       info.IsDir(),
		Permissions: info.Mode().String(),
		Owner:       "unknown",
		Group:       "unknown",
	}, nil
}

func (g *GenericVFS) MkDir(containerID string, path string, mode os.FileMode) error {
	hostPath, err := g.GetHostPath(containerID, path)
	if err != nil {
		return err
	}

	return os.MkdirAll(hostPath, mode)
}

func (g *GenericVFS) Remove(containerID string, path string) error {
	hostPath, err := g.GetHostPath(containerID, path)
	if err != nil {
		return err
	}

	return os.RemoveAll(hostPath)
}

func (g *GenericVFS) Copy(srcContainerID, srcPath, dstContainerID, dstPath string) error {
	// Basic copy implementation
	srcHostPath, err := g.GetHostPath(srcContainerID, srcPath)
	if err != nil {
		return err
	}

	dstHostPath, err := g.GetHostPath(dstContainerID, dstPath)
	if err != nil {
		return err
	}

	return g.copyPath(srcHostPath, dstHostPath)
}

func (g *GenericVFS) copyPath(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		return g.copyDir(src, dst)
	}
	return g.copyFile(src, dst, srcInfo.Mode())
}

func (g *GenericVFS) copyDir(src, dst string) error {
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

		return g.copyFile(path, dstPath, info.Mode())
	})
}

func (g *GenericVFS) copyFile(src, dst string, mode os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

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

func (g *GenericVFS) Move(containerID string, srcPath, dstPath string) error {
	srcHostPath, err := g.GetHostPath(containerID, srcPath)
	if err != nil {
		return err
	}

	dstHostPath, err := g.GetHostPath(containerID, dstPath)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dstHostPath), 0755); err != nil {
		return err
	}

	return os.Rename(srcHostPath, dstHostPath)
}

func (g *GenericVFS) Chmod(containerID string, path string, mode os.FileMode) error {
	hostPath, err := g.GetHostPath(containerID, path)
	if err != nil {
		return err
	}

	return os.Chmod(hostPath, mode)
}

func (g *GenericVFS) Find(containerID string, basePath string, name string, recursive bool) ([]string, error) {
	hostBasePath, err := g.GetHostPath(containerID, basePath)
	if err != nil {
		return nil, err
	}

	var matches []string

	walkFunc := func(hostPath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		relPath, err := filepath.Rel(hostBasePath, hostPath)
		if err != nil {
			return nil
		}

		containerPath := filepath.Join(basePath, relPath)
		if basePath == "/" {
			containerPath = "/" + relPath
		}

		if name == "" || strings.Contains(info.Name(), name) {
			matches = append(matches, containerPath)
		}

		if !recursive && info.IsDir() && hostPath != hostBasePath {
			return filepath.SkipDir
		}

		return nil
	}

	err = filepath.Walk(hostBasePath, walkFunc)
	return matches, err
}

func (g *GenericVFS) GetWorkingDir(containerID string) (string, error) {
	if wd, exists := g.workingDirs[containerID]; exists {
		return wd, nil
	}
	return "/", nil
}

func (g *GenericVFS) SetWorkingDir(containerID string, path string) error {
	if _, err := g.Stat(containerID, path); err != nil {
		return fmt.Errorf("directory does not exist: %s", path)
	}

	if g.workingDirs == nil {
		g.workingDirs = make(map[string]string)
	}

	g.workingDirs[containerID] = path
	return nil
}

func (g *GenericVFS) GetHostPath(containerID string, containerPath string) (string, error) {
	rootPath, exists := g.containerRoots[containerID]
	if !exists {
		return "", fmt.Errorf("container %s not found", containerID)
	}

	cleanPath := filepath.Clean(containerPath)
	if !filepath.IsAbs(cleanPath) {
		cleanPath = "/" + cleanPath
	}

	// Handle platform-specific path conversion
	if runtime.GOOS == "windows" {
		// Convert Unix-style container paths to Windows paths
		cleanPath = strings.ReplaceAll(cleanPath, "/", string(filepath.Separator))
		cleanPath = strings.TrimPrefix(cleanPath, string(filepath.Separator))
		hostPath := filepath.Join(rootPath, cleanPath)
		return hostPath, nil
	}

	// Unix-style systems
	hostPath := filepath.Join(rootPath, strings.TrimPrefix(cleanPath, "/"))
	return hostPath, nil
}

func (g *GenericVFS) Cleanup(containerID string) error {
	if _, exists := g.containerRoots[containerID]; exists {
		delete(g.containerRoots, containerID)
	}

	if g.workingDirs != nil {
		delete(g.workingDirs, containerID)
	}

	return nil
}

// Windows VFS Implementation (inherits from Generic but with Windows-specific optimizations)
func (w *WindowsVFS) Initialize(containerID string, imageRootfs string) error {
	if w.containerRoots == nil {
		w.containerRoots = make(map[string]string)
	}
	if w.workingDirs == nil {
		w.workingDirs = make(map[string]string)
	}

	// Use Windows-appropriate directory structure
	tempDir := os.Getenv("TEMP")
	if tempDir == "" {
		tempDir = os.Getenv("TMP")
	}
	if tempDir == "" {
		tempDir = "C:\\temp"
	}

	containerDir := filepath.Join(tempDir, "servin", "containers", containerID)
	rootfsDir := filepath.Join(containerDir, "rootfs")

	// Create the directory structure with Windows-appropriate permissions
	if err := os.MkdirAll(rootfsDir, 0666); err != nil {
		return fmt.Errorf("failed to create container directory: %w", err)
	}

	// Create Windows-compatible filesystem structure
	if err := w.createWindowsMinimalFS(rootfsDir); err != nil {
		return fmt.Errorf("failed to create minimal filesystem: %w", err)
	}

	w.containerRoots[containerID] = rootfsDir
	w.workingDirs[containerID] = "/"

	return nil
}

// createWindowsMinimalFS creates a Windows-compatible filesystem structure
func (w *WindowsVFS) createWindowsMinimalFS(rootPath string) error {
	// Create both Unix-style and Windows-style directory structure
	dirs := []string{
		// Unix-style directories for compatibility
		"bin", "etc", "home", "tmp", "usr", "usr\\bin", "var", "var\\log",
		// Windows-style directories
		"Program Files", "Windows", "Windows\\System32", "Users",
		"ProgramData", "System32",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(rootPath, dir)
		if err := os.MkdirAll(dirPath, 0666); err != nil {
			return err
		}
	}

	// Create some basic files with Windows line endings
	files := map[string]string{
		"etc\\hostname":    "container\r\n",
		"etc\\hosts":       "127.0.0.1 localhost\r\n::1 localhost\r\n",
		"etc\\passwd":      "root:x:0:0:root:/root:/bin/sh\r\n",
		"etc\\group":       "root:x:0:\r\n",
		"Windows\\win.ini": "; Windows configuration\r\n[fonts]\r\n[extensions]\r\n",
		"autoexec.bat":     "@echo off\r\nrem Container startup\r\n",
		"config.sys":       "rem Container configuration\r\n",
	}

	for file, content := range files {
		filePath := filepath.Join(rootPath, file)
		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(filePath), 0666); err != nil {
			return err
		}
		if err := os.WriteFile(filePath, []byte(content), 0666); err != nil {
			return err
		}
	}

	return nil
}

func (w *WindowsVFS) Mount(containerID string) error {
	g := (*GenericVFS)(w)
	return g.Mount(containerID)
}

func (w *WindowsVFS) Unmount(containerID string) error {
	g := (*GenericVFS)(w)
	return g.Unmount(containerID)
}

func (w *WindowsVFS) List(containerID string, path string) ([]FileInfo, error) {
	g := (*GenericVFS)(w)
	return g.List(containerID, path)
}

func (w *WindowsVFS) Read(containerID string, path string) (io.ReadCloser, error) {
	g := (*GenericVFS)(w)
	return g.Read(containerID, path)
}

func (w *WindowsVFS) Write(containerID string, path string, data io.Reader) error {
	g := (*GenericVFS)(w)
	return g.Write(containerID, path, data)
}

func (w *WindowsVFS) Stat(containerID string, path string) (FileInfo, error) {
	g := (*GenericVFS)(w)
	return g.Stat(containerID, path)
}

func (w *WindowsVFS) MkDir(containerID string, path string, mode os.FileMode) error {
	g := (*GenericVFS)(w)
	return g.MkDir(containerID, path, mode)
}

func (w *WindowsVFS) Remove(containerID string, path string) error {
	g := (*GenericVFS)(w)
	return g.Remove(containerID, path)
}

func (w *WindowsVFS) Copy(srcContainerID, srcPath, dstContainerID, dstPath string) error {
	g := (*GenericVFS)(w)
	return g.Copy(srcContainerID, srcPath, dstContainerID, dstPath)
}

func (w *WindowsVFS) Move(containerID string, srcPath, dstPath string) error {
	g := (*GenericVFS)(w)
	return g.Move(containerID, srcPath, dstPath)
}

func (w *WindowsVFS) Chmod(containerID string, path string, mode os.FileMode) error {
	g := (*GenericVFS)(w)
	return g.Chmod(containerID, path, mode)
}

func (w *WindowsVFS) Find(containerID string, basePath string, name string, recursive bool) ([]string, error) {
	g := (*GenericVFS)(w)

	// Call the generic find but with Windows-specific case-insensitive handling
	hostBasePath, err := g.GetHostPath(containerID, basePath)
	if err != nil {
		return nil, err
	}

	var matches []string

	walkFunc := func(hostPath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		relPath, err := filepath.Rel(hostBasePath, hostPath)
		if err != nil {
			return nil
		}

		// Convert Windows path separators to Unix style for container paths
		containerPath := strings.ReplaceAll(relPath, "\\", "/")
		if basePath != "/" {
			containerPath = basePath + "/" + containerPath
		} else if containerPath != "." {
			containerPath = "/" + containerPath
		} else {
			containerPath = "/"
		}

		// Windows case-insensitive name matching
		if name == "" || strings.Contains(strings.ToLower(info.Name()), strings.ToLower(name)) {
			matches = append(matches, containerPath)
		}

		if !recursive && info.IsDir() && hostPath != hostBasePath {
			return filepath.SkipDir
		}

		return nil
	}

	err = filepath.Walk(hostBasePath, walkFunc)
	return matches, err
}

func (w *WindowsVFS) GetWorkingDir(containerID string) (string, error) {
	g := (*GenericVFS)(w)
	return g.GetWorkingDir(containerID)
}

func (w *WindowsVFS) SetWorkingDir(containerID string, path string) error {
	g := (*GenericVFS)(w)
	return g.SetWorkingDir(containerID, path)
}

func (w *WindowsVFS) GetHostPath(containerID string, containerPath string) (string, error) {
	g := (*GenericVFS)(w)
	return g.GetHostPath(containerID, containerPath)
}

func (w *WindowsVFS) Cleanup(containerID string) error {
	g := (*GenericVFS)(w)
	return g.Cleanup(containerID)
}
