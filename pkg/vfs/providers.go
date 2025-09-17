package vfs

import (
	"runtime"
)

// GetPlatformProvider returns the appropriate VFS provider for the current platform
func GetPlatformProvider() VFSProvider {
	switch runtime.GOOS {
	case "linux":
		return &LinuxProvider{}
	case "darwin":
		return &MacOSProvider{}
	case "windows":
		return &WindowsProvider{}
	default:
		// Fallback to generic provider for unsupported platforms
		return &GenericProvider{}
	}
}

// LinuxProvider provides VFS for Linux with full namespace support
type LinuxProvider struct{}

func (p *LinuxProvider) CreateVFS() VirtualFileSystem {
	return &LinuxVFS{}
}

func (p *LinuxProvider) SupportsNamespaces() bool {
	return true
}

func (p *LinuxProvider) GetPlatform() string {
	return "linux"
}

// MacOSProvider provides VFS for macOS with chroot and overlay simulation
type MacOSProvider struct{}

func (p *MacOSProvider) CreateVFS() VirtualFileSystem {
	return &MacOSVFS{}
}

func (p *MacOSProvider) SupportsNamespaces() bool {
	return false
}

func (p *MacOSProvider) GetPlatform() string {
	return "darwin"
}

// WindowsProvider provides VFS for Windows
type WindowsProvider struct{}

func (p *WindowsProvider) CreateVFS() VirtualFileSystem {
	return &WindowsVFS{}
}

func (p *WindowsProvider) SupportsNamespaces() bool {
	return false
}

func (p *WindowsProvider) GetPlatform() string {
	return "windows"
}

// GenericProvider provides a fallback VFS implementation
type GenericProvider struct{}

func (p *GenericProvider) CreateVFS() VirtualFileSystem {
	return &GenericVFS{}
}

func (p *GenericProvider) SupportsNamespaces() bool {
	return false
}

func (p *GenericProvider) GetPlatform() string {
	return "generic"
}
