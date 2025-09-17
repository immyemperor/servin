//go:build linux

package vfs

import (
	"fmt"
	"os"
	"syscall"
)

// getFileOwnerInfo extracts owner and group information from file info on Linux
func getFileOwnerInfo(info os.FileInfo) (string, string) {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		return fmt.Sprintf("%d", stat.Uid), fmt.Sprintf("%d", stat.Gid)
	}
	return "0", "0"
}
