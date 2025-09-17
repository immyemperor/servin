package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var cpCmd = &cobra.Command{
	Use:   "cp SOURCE_CONTAINER:SOURCE_PATH DEST_CONTAINER:DEST_PATH",
	Short: "Copy files between containers or between host and container",
	Long:  "Copy files/directories between containers or between host filesystem and container",
	Args:  cobra.ExactArgs(2),
	RunE:  copyFiles,
}

var mvCmd = &cobra.Command{
	Use:   "mv CONTAINER OLD_PATH NEW_PATH",
	Short: "Move/rename files in container",
	Long:  "Move or rename files and directories within a container filesystem",
	Args:  cobra.ExactArgs(3),
	RunE:  moveFiles,
}

var mkdirCmd = &cobra.Command{
	Use:   "mkdir CONTAINER DIRECTORY",
	Short: "Create directory in container",
	Long:  "Create one or more directories in the container filesystem",
	Args:  cobra.MinimumNArgs(2),
	RunE:  makeDirectory,
}

var rmdirCmd = &cobra.Command{
	Use:   "rmdir CONTAINER DIRECTORY",
	Short: "Remove empty directory from container",
	Long:  "Remove empty directories from the container filesystem",
	Args:  cobra.MinimumNArgs(2),
	RunE:  removeDirectory,
}

var rmCmd = &cobra.Command{
	Use:   "rm CONTAINER FILE",
	Short: "Remove files from container",
	Long:  "Remove files and directories from the container filesystem",
	Args:  cobra.MinimumNArgs(2),
	RunE:  removeFiles,
}

var chmodCmd = &cobra.Command{
	Use:   "chmod CONTAINER MODE FILE",
	Short: "Change file permissions in container",
	Long:  "Change file mode/permissions for files in the container filesystem",
	Args:  cobra.ExactArgs(3),
	RunE:  changeMode,
}

func init() {
	rootCmd.AddCommand(cpCmd)
	rootCmd.AddCommand(mvCmd)
	rootCmd.AddCommand(mkdirCmd)
	rootCmd.AddCommand(rmdirCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(chmodCmd)

	// Add flags
	cpCmd.Flags().BoolP("recursive", "r", false, "Copy directories recursively")
	mkdirCmd.Flags().BoolP("parents", "p", false, "Create parent directories as needed")
	rmCmd.Flags().BoolP("recursive", "r", false, "Remove directories recursively")
	rmCmd.Flags().BoolP("force", "f", false, "Force removal without prompting")
}

func copyFiles(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	source := args[0]
	dest := args[1]

	// Parse source and destination
	sourceContainer, sourcePath := parseContainerPath(source)
	destContainer, destPath := parseContainerPath(dest)

	recursive, _ := cmd.Flags().GetBool("recursive")

	// Handle different copy scenarios
	if sourceContainer != "" && destContainer != "" {
		// Container to container copy
		return copyContainerToContainer(sourceContainer, sourcePath, destContainer, destPath, recursive)
	} else if sourceContainer != "" {
		// Container to host copy
		return copyContainerToHost(sourceContainer, sourcePath, destPath, recursive)
	} else if destContainer != "" {
		// Host to container copy
		return copyHostToContainer(sourcePath, destContainer, destPath, recursive)
	} else {
		return fmt.Errorf("invalid copy arguments - at least one path must specify a container")
	}
}

func moveFiles(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	containerID := args[0]
	oldPath := args[1]
	newPath := args[2]

	// Get container rootfs path
	rootfsPath, err := getContainerRootFS(containerID)
	if err != nil {
		return err
	}

	fullOldPath := filepath.Join(rootfsPath, oldPath)
	fullNewPath := filepath.Join(rootfsPath, newPath)

	// Check if source exists
	if _, err := os.Stat(fullOldPath); os.IsNotExist(err) {
		return fmt.Errorf("source path does not exist: %s", oldPath)
	}

	// Move the file/directory
	err = os.Rename(fullOldPath, fullNewPath)
	if err != nil {
		return fmt.Errorf("failed to move %s to %s: %v", oldPath, newPath, err)
	}

	fmt.Printf("Moved %s to %s in container %s\n", oldPath, newPath, containerID[:12])
	return nil
}

func makeDirectory(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	containerID := args[0]
	directories := args[1:]

	// Get container rootfs path
	rootfsPath, err := getContainerRootFS(containerID)
	if err != nil {
		return err
	}

	parents, _ := cmd.Flags().GetBool("parents")

	for _, dir := range directories {
		fullPath := filepath.Join(rootfsPath, dir)

		if parents {
			err = os.MkdirAll(fullPath, 0755)
		} else {
			err = os.Mkdir(fullPath, 0755)
		}

		if err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}

		fmt.Printf("Created directory %s in container %s\n", dir, containerID[:12])
	}

	return nil
}

func removeDirectory(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	containerID := args[0]
	directories := args[1:]

	// Get container rootfs path
	rootfsPath, err := getContainerRootFS(containerID)
	if err != nil {
		return err
	}

	for _, dir := range directories {
		fullPath := filepath.Join(rootfsPath, dir)

		// Check if directory exists and is empty
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", dir)
		}

		err = os.Remove(fullPath)
		if err != nil {
			return fmt.Errorf("failed to remove directory %s: %v", dir, err)
		}

		fmt.Printf("Removed directory %s from container %s\n", dir, containerID[:12])
	}

	return nil
}

func removeFiles(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	containerID := args[0]
	files := args[1:]

	// Get container rootfs path
	rootfsPath, err := getContainerRootFS(containerID)
	if err != nil {
		return err
	}

	recursive, _ := cmd.Flags().GetBool("recursive")
	force, _ := cmd.Flags().GetBool("force")

	for _, file := range files {
		fullPath := filepath.Join(rootfsPath, file)

		// Check if file exists
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			if !force {
				return fmt.Errorf("file does not exist: %s", file)
			}
			continue
		}

		if recursive {
			err = os.RemoveAll(fullPath)
		} else {
			err = os.Remove(fullPath)
		}

		if err != nil {
			if !force {
				return fmt.Errorf("failed to remove %s: %v", file, err)
			}
			continue
		}

		fmt.Printf("Removed %s from container %s\n", file, containerID[:12])
	}

	return nil
}

func changeMode(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	containerID := args[0]
	modeStr := args[1]
	filePath := args[2]

	// Get container rootfs path
	rootfsPath, err := getContainerRootFS(containerID)
	if err != nil {
		return err
	}

	fullPath := filepath.Join(rootfsPath, filePath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	// Parse mode string (simplified - assumes octal notation)
	var mode os.FileMode
	if len(modeStr) == 3 {
		// Parse octal mode like "755"
		for i, char := range modeStr {
			digit := int(char - '0')
			if digit < 0 || digit > 7 {
				return fmt.Errorf("invalid mode: %s", modeStr)
			}
			mode |= os.FileMode(digit) << uint(3*(2-i))
		}
	} else {
		return fmt.Errorf("unsupported mode format: %s (use octal notation like '755')", modeStr)
	}

	err = os.Chmod(fullPath, mode)
	if err != nil {
		return fmt.Errorf("failed to change mode of %s: %v", filePath, err)
	}

	fmt.Printf("Changed mode of %s to %s in container %s\n", filePath, modeStr, containerID[:12])
	return nil
}

// Helper functions

func parseContainerPath(path string) (container, filePath string) {
	parts := strings.SplitN(path, ":", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", path
}

func copyContainerToContainer(sourceContainer, sourcePath, destContainer, destPath string, recursive bool) error {
	// Get source container rootfs
	sourceRootfs, err := getContainerRootFS(sourceContainer)
	if err != nil {
		return fmt.Errorf("source container error: %v", err)
	}

	// Get destination container rootfs
	destRootfs, err := getContainerRootFS(destContainer)
	if err != nil {
		return fmt.Errorf("destination container error: %v", err)
	}

	sourceFullPath := filepath.Join(sourceRootfs, sourcePath)
	destFullPath := filepath.Join(destRootfs, destPath)

	return copyPath(sourceFullPath, destFullPath, recursive)
}

func copyContainerToHost(containerID, containerPath, hostPath string, recursive bool) error {
	// Get container rootfs
	rootfsPath, err := getContainerRootFS(containerID)
	if err != nil {
		return err
	}

	sourceFullPath := filepath.Join(rootfsPath, containerPath)
	return copyPath(sourceFullPath, hostPath, recursive)
}

func copyHostToContainer(hostPath, containerID, containerPath string, recursive bool) error {
	// Get container rootfs
	rootfsPath, err := getContainerRootFS(containerID)
	if err != nil {
		return err
	}

	destFullPath := filepath.Join(rootfsPath, containerPath)
	return copyPath(hostPath, destFullPath, recursive)
}

func copyPath(src, dst string, recursive bool) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("source path error: %v", err)
	}

	if srcInfo.IsDir() {
		if !recursive {
			return fmt.Errorf("source is a directory, use -r flag for recursive copy")
		}
		return copyDir(src, dst)
	}

	return copyFile(src, dst)
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination directory if needed
	destDir := filepath.Dir(dst)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Copy file permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create destination directory
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}
