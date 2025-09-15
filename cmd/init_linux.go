//go:build linux

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/unix"
)

// setupContainerEnvironment sets up the container's isolated environment
func setupContainerEnvironment() error {
	// Change to rootfs if available (passed via environment variable)
	if rootfsPath := os.Getenv("SERVIN_ROOTFS"); rootfsPath != "" {
		if err := changeRoot(rootfsPath); err != nil {
			return fmt.Errorf("failed to change root: %v", err)
		}

		// After changing root, set up the container filesystem
		if err := setupContainerFilesystem(); err != nil {
			return fmt.Errorf("failed to setup container filesystem: %v", err)
		}
	}

	// Set hostname if provided
	if hostname := os.Getenv("HOSTNAME"); hostname != "" {
		if err := setHostname(hostname); err != nil {
			fmt.Printf("Warning: failed to set hostname: %v\n", err)
		}
	}

	return nil
}

// setupContainerFilesystem sets up the container's internal filesystem
func setupContainerFilesystem() error {
	// Mount proc filesystem
	if err := mountProc(); err != nil {
		fmt.Printf("Warning: failed to mount proc: %v\n", err)
	}

	// Mount sys filesystem
	if err := mountSys(); err != nil {
		fmt.Printf("Warning: failed to mount sys: %v\n", err)
	}

	// Setup basic devices
	if err := setupDevices(); err != nil {
		fmt.Printf("Warning: failed to setup devices: %v\n", err)
	}

	return nil
}

// setHostname sets the container hostname
func setHostname(hostname string) error {
	return unix.Sethostname([]byte(hostname))
}

// mountProc mounts the proc filesystem
func mountProc() error {
	if err := os.MkdirAll("/proc", 0755); err != nil {
		return err
	}

	// Mount proc filesystem
	return unix.Mount("proc", "/proc", "proc", 0, "")
}

// mountSys mounts the sys filesystem
func mountSys() error {
	if err := os.MkdirAll("/sys", 0755); err != nil {
		return err
	}

	// Mount sys filesystem
	return unix.Mount("sysfs", "/sys", "sysfs", 0, "")
}

// setupDevices creates basic device nodes
func setupDevices() error {
	// Create /dev directory
	if err := os.MkdirAll("/dev", 0755); err != nil {
		return err
	}

	// Create basic device nodes
	devices := map[string]struct {
		major, minor int
		mode         uint32
	}{
		"/dev/null":    {1, 3, syscall.S_IFCHR | 0666},
		"/dev/zero":    {1, 5, syscall.S_IFCHR | 0666},
		"/dev/random":  {1, 8, syscall.S_IFCHR | 0444},
		"/dev/urandom": {1, 9, syscall.S_IFCHR | 0444},
	}

	for path, dev := range devices {
		if err := createDeviceNode(path, dev.major, dev.minor, dev.mode); err != nil {
			fmt.Printf("Warning: failed to create device %s: %v\n", path, err)
		}
	}

	return nil
}

// createDeviceNode creates a device node
func createDeviceNode(path string, major, minor int, mode uint32) error {
	dev := unix.Mkdev(uint32(major), uint32(minor))
	return unix.Mknod(path, mode, int(dev))
}

// changeRoot changes to the container's rootfs using pivot_root
func changeRoot(rootfsPath string) error {
	// Create old_root directory inside the new root
	oldRoot := filepath.Join(rootfsPath, ".old_root")
	if err := os.MkdirAll(oldRoot, 0755); err != nil {
		return fmt.Errorf("failed to create old_root directory: %v", err)
	}

	// Use pivot_root to change the root filesystem
	if err := unix.PivotRoot(rootfsPath, oldRoot); err != nil {
		// If pivot_root fails, fall back to chroot
		fmt.Printf("Warning: pivot_root failed, falling back to chroot: %v\n", err)
		return fallbackChroot(rootfsPath)
	}

	// Change to the new root
	if err := unix.Chdir("/"); err != nil {
		return fmt.Errorf("failed to chdir to new root: %v", err)
	}

	// Unmount the old root
	if err := unix.Unmount("/.old_root", unix.MNT_DETACH); err != nil {
		fmt.Printf("Warning: failed to unmount old root: %v\n", err)
	}

	// Remove the old root directory
	if err := os.RemoveAll("/.old_root"); err != nil {
		fmt.Printf("Warning: failed to remove old root directory: %v\n", err)
	}

	fmt.Printf("Successfully changed root to: %s\n", rootfsPath)
	return nil
}

// fallbackChroot performs a simple chroot as fallback
func fallbackChroot(rootfsPath string) error {
	if err := unix.Chroot(rootfsPath); err != nil {
		return fmt.Errorf("chroot failed: %v", err)
	}

	if err := unix.Chdir("/"); err != nil {
		return fmt.Errorf("failed to chdir after chroot: %v", err)
	}

	fmt.Printf("Changed root using chroot to: %s\n", rootfsPath)
	return nil
}
