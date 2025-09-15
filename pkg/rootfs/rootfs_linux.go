//go:build linux

package rootfs

import (
	"fmt"
	"os"
	"path/filepath"

	"servin/pkg/image"

	"golang.org/x/sys/unix"
)

// RootFS manages container root filesystem
type RootFS struct {
	ContainerID  string
	RootPath     string
	ImagePath    string
	ImageManager *image.Manager
}

// New creates a new RootFS manager with image support
func New(containerID, imageRef string) *RootFS {
	rootPath := fmt.Sprintf("/var/lib/servin/containers/%s/rootfs", containerID)
	return &RootFS{
		ContainerID:  containerID,
		RootPath:     rootPath,
		ImagePath:    imageRef,
		ImageManager: image.NewManager(),
	}
}

// Create sets up the container's root filesystem from an image
func (r *RootFS) Create() error {
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

// createFromImage creates rootfs from a managed image
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

// createBasicRootFS creates a minimal rootfs without images
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

// Enter enters the container's root filesystem and sets up the environment
func (r *RootFS) Enter() error {
	// Set environment variable for the init process to use
	if err := os.Setenv("SERVIN_ROOTFS", r.RootPath); err != nil {
		return fmt.Errorf("failed to set SERVIN_ROOTFS environment: %v", err)
	}

	fmt.Printf("Prepared rootfs environment at %s\n", r.RootPath)
	return nil
}

// EnterChroot performs actual chroot operation (used by init process)
func (r *RootFS) EnterChroot() error {
	// Change to the container's root filesystem
	if err := unix.Chroot(r.RootPath); err != nil {
		return fmt.Errorf("failed to chroot to %s: %v", r.RootPath, err)
	}

	// Change working directory to root inside the chroot
	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("failed to change to root directory: %v", err)
	}

	fmt.Printf("Entered container rootfs (chroot to %s)\n", r.RootPath)
	return nil
}

// SetupMounts sets up necessary filesystems inside the container
func (r *RootFS) SetupMounts() error {
	mounts := []struct {
		source string
		target string
		fstype string
		flags  uintptr
		data   string
	}{
		{"proc", "/proc", "proc", 0, ""},
		{"sysfs", "/sys", "sysfs", 0, ""},
		{"tmpfs", "/tmp", "tmpfs", 0, ""},
		{"devtmpfs", "/dev", "devtmpfs", 0, ""},
	}

	for _, mount := range mounts {
		targetPath := filepath.Join(r.RootPath, mount.target)
		if err := os.MkdirAll(targetPath, 0755); err != nil {
			fmt.Printf("Warning: failed to create mount point %s: %v\n", targetPath, err)
			continue
		}

		if err := unix.Mount(mount.source, targetPath, mount.fstype, mount.flags, mount.data); err != nil {
			fmt.Printf("Warning: failed to mount %s: %v\n", mount.target, err)
		} else {
			fmt.Printf("Mounted %s at %s\n", mount.source, targetPath)
		}
	}

	return nil
}

// Cleanup removes the container's filesystem
func (r *RootFS) Cleanup() error {
	return os.RemoveAll(filepath.Dir(r.RootPath))
}

// copyDirectory recursively copies a directory
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

// copyEssentialFiles copies essential files from host to container
func (r *RootFS) copyEssentialFiles() error {
	// Copy essential binaries
	binaries := map[string]string{
		"/bin/sh":   "bin/sh",
		"/bin/bash": "bin/bash",
		"/bin/ls":   "bin/ls",
		"/bin/cat":  "bin/cat",
		"/bin/echo": "bin/echo",
	}

	for hostPath, containerPath := range binaries {
		if _, err := os.Stat(hostPath); err == nil {
			targetPath := filepath.Join(r.RootPath, containerPath)
			if err := r.copyFile(hostPath, targetPath); err != nil {
				fmt.Printf("Warning: failed to copy %s: %v\n", hostPath, err)
			}
		}
	}

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

	return nil
}

// copyFile copies a file from src to dst
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
