//go:build windows

package vfs

import (
	"os"
)

// getFileOwnerInfo extracts owner and group information from file info on Windows
// Windows doesn't use the same UID/GID system, so we return default values
func getFileOwnerInfo(info os.FileInfo) (string, string) {
	// Windows file ownership is more complex and uses SIDs
	// For simplicity in containerization, we return default values
	return "0", "0"
}
