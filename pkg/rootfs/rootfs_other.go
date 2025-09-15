//go:build !linux

package rootfs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"servin/pkg/image"
)

// RootFS manages container root filesystem (cross-platform implementation)
type RootFS struct {
	ContainerID  string
	RootPath     string
	ImagePath    string
	ImageManager *image.Manager
}

// New creates a new RootFS manager (cross-platform)
func New(containerID, imagePath string) *RootFS {
	// Use platform-appropriate path
	var rootPath string

	switch runtime.GOOS {
	case "windows":
		// Windows: Use user home directory
		homeDir, _ := os.UserHomeDir()
		rootPath = filepath.Join(homeDir, ".servin", "containers", containerID, "rootfs")
	case "darwin":
		// macOS: Use user home directory (similar to Windows but Unix paths)
		homeDir, _ := os.UserHomeDir()
		rootPath = filepath.Join(homeDir, ".servin", "containers", containerID, "rootfs")
	default:
		// Other Unix-like systems: Use /var/lib (but not Linux since that has its own implementation)
		rootPath = fmt.Sprintf("/var/lib/servin/containers/%s/rootfs", containerID)
	}

	return &RootFS{
		ContainerID:  containerID,
		RootPath:     rootPath,
		ImagePath:    imagePath,
		ImageManager: image.NewManager(),
	}
}

// Create sets up the container's root filesystem (cross-platform)
func (r *RootFS) Create() error {
	platformInfo := ""
	switch runtime.GOOS {
	case "windows":
		platformInfo = "Windows"
	case "darwin":
		platformInfo = "macOS"
	default:
		platformInfo = runtime.GOOS
	}

	fmt.Printf("Creating cross-platform rootfs for %s at %s\n", platformInfo, r.RootPath)

	// Create the root directory
	if err := os.MkdirAll(r.RootPath, 0755); err != nil {
		return fmt.Errorf("failed to create rootfs directory: %v", err)
	}

	// Try to use image-based rootfs first
	if r.ImagePath != "" {
		if err := r.createFromImage(); err != nil {
			fmt.Printf("Warning: failed to create rootfs from image '%s': %v\n", r.ImagePath, err)
			fmt.Println("Falling back to basic rootfs creation...")
			return r.createBasicRootFS()
		}
		return nil
	}

	// Fall back to basic rootfs creation
	return r.createBasicRootFS()
}

// createFromImage creates rootfs from a managed image (cross-platform)
func (r *RootFS) createFromImage() error {
	// Get image from manager
	img, err := r.ImageManager.GetImage(r.ImagePath)
	if err != nil {
		return fmt.Errorf("image not found: %v", err)
	}

	// Copy image rootfs to container rootfs
	if err := r.copyDirectory(img.RootFSPath, r.RootPath); err != nil {
		return fmt.Errorf("failed to copy image rootfs: %v", err)
	}

	fmt.Printf("Created rootfs from image %s at %s\n", r.ImagePath, r.RootPath)
	return nil
}

// createBasicRootFS creates a minimal rootfs without images (cross-platform)
func (r *RootFS) createBasicRootFS() error {
	// Create essential directories in the container rootfs
	essentialDirs := []string{
		"bin", "etc", "proc", "sys", "tmp", "var", "usr/bin", "usr/lib",
		"dev", "home", "root", "mnt", "opt",
	}

	for _, dir := range essentialDirs {
		dirPath := filepath.Join(r.RootPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	// Copy basic files from host if they exist
	if err := r.copyEssentialFiles(); err != nil {
		fmt.Printf("Warning: failed to copy some essential files: %v\n", err)
	}

	fmt.Printf("Created basic rootfs at %s\n", r.RootPath)
	return nil
}

// Enter simulates entering the container's root filesystem (cross-platform)
func (r *RootFS) Enter() error {
	// On non-Linux platforms, we can't actually chroot, but we can simulate it
	platformInfo := ""
	switch runtime.GOOS {
	case "windows":
		platformInfo = "Windows - chroot not available"
	case "darwin":
		platformInfo = "macOS - chroot requires root privileges and SIP considerations"
	default:
		platformInfo = fmt.Sprintf("%s - chroot capabilities unknown", runtime.GOOS)
	}

	fmt.Printf("Simulating container rootfs entry (path: %s)\n", r.RootPath)
	fmt.Printf("Note: %s - container will use host filesystem\n", platformInfo)
	return nil
}

// EnterChroot is not available on non-Linux platforms
func (r *RootFS) EnterChroot() error {
	return fmt.Errorf("chroot not available on %s", runtime.GOOS)
}

// SetupMounts simulates filesystem mounts (cross-platform)
func (r *RootFS) SetupMounts() error {
	fmt.Printf("Simulating filesystem mounts for %s (not available on %s)\n", r.RootPath, runtime.GOOS)
	return nil
}

// Cleanup removes the container's filesystem (cross-platform)
func (r *RootFS) Cleanup() error {
	if r.RootPath != "" {
		return os.RemoveAll(filepath.Dir(r.RootPath))
	}
	return nil
}

// copyDirectory recursively copies a directory (cross-platform)
func (r *RootFS) copyDirectory(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return r.copyFile(path, destPath)
	})
}

// copyEssentialFiles copies essential files from host to container (cross-platform)
func (r *RootFS) copyEssentialFiles() error {
	// Create a minimal /etc/passwd
	passwdContent := `root:x:0:0:root:/root:/bin/sh
nobody:x:65534:65534:nobody:/nonexistent:/usr/sbin/nologin
`
	passwdPath := filepath.Join(r.RootPath, "etc", "passwd")
	if err := os.WriteFile(passwdPath, []byte(passwdContent), 0644); err != nil {
		return fmt.Errorf("failed to create /etc/passwd: %v", err)
	}

	// Create a minimal /etc/group
	groupContent := `root:x:0:
nogroup:x:65534:
`
	groupPath := filepath.Join(r.RootPath, "etc", "group")
	if err := os.WriteFile(groupPath, []byte(groupContent), 0644); err != nil {
		return fmt.Errorf("failed to create /etc/group: %v", err)
	}

	// Create a simple shell script for testing
	scriptContent := `#!/bin/sh
echo "Container rootfs created successfully"
echo "Platform: Cross-platform (non-Linux)"
echo "RootFS: %s"
`
	scriptPath := filepath.Join(r.RootPath, "bin", "test.sh")
	if err := os.WriteFile(scriptPath, []byte(fmt.Sprintf(scriptContent, r.RootPath)), 0755); err != nil {
		fmt.Printf("Warning: failed to create test script: %v\n", err)
	}

	return nil
}

// copyFile copies a file from src to dst (cross-platform)
func (r *RootFS) copyFile(src, dst string) error {
	sourceData, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	// Get source file info for permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, sourceData, srcInfo.Mode())
}
