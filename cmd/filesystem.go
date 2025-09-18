package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"servin/pkg/state"
	"servin/pkg/vfs"

	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "fs-ls [FLAGS] CONTAINER [PATH]",
	Short: "List files and directories in a container",
	Long:  "List files and directories in the specified container filesystem path",
	Args:  cobra.MinimumNArgs(1),
	RunE:  listFiles,
}

var catCmd = &cobra.Command{
	Use:   "cat CONTAINER FILE",
	Short: "Display file contents in a container",
	Long:  "Display the contents of a file in the specified container",
	Args:  cobra.ExactArgs(2),
	RunE:  displayFile,
}

var statCmd = &cobra.Command{
	Use:   "stat CONTAINER PATH",
	Short: "Display file status in a container",
	Long:  "Display detailed file/directory status information in a container",
	Args:  cobra.ExactArgs(2),
	RunE:  statFile,
}

var findCmd = &cobra.Command{
	Use:   "find CONTAINER PATH [NAME]",
	Short: "Find files and directories in a container",
	Long:  "Search for files and directories matching criteria in a container",
	Args:  cobra.MinimumNArgs(2),
	RunE:  findFiles,
}

var pwdCmd = &cobra.Command{
	Use:   "pwd CONTAINER",
	Short: "Show working directory in a container",
	Long:  "Show the current working directory in the specified container",
	Args:  cobra.ExactArgs(1),
	RunE:  showWorkingDir,
}

func init() {
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(catCmd)
	rootCmd.AddCommand(statCmd)
	rootCmd.AddCommand(findCmd)
	rootCmd.AddCommand(pwdCmd)

	// Add flags for ls command
	lsCmd.Flags().BoolP("long", "l", false, "Use long listing format")
	lsCmd.Flags().BoolP("all", "a", false, "Show hidden files")
	lsCmd.Flags().Bool("human", false, "Human readable sizes") // Remove -h shorthand to avoid conflict
	lsCmd.Flags().BoolP("recursive", "R", false, "List subdirectories recursively")

	// Add flags for find command
	findCmd.Flags().StringP("name", "n", "", "Find files by name pattern")
	findCmd.Flags().StringP("type", "t", "", "Find by type (f=file, d=directory)")
	findCmd.Flags().IntP("maxdepth", "d", -1, "Maximum search depth")
}

func listFiles(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	containerID := args[0]
	path := "/"
	if len(args) > 1 {
		path = args[1]
	}

	// Create VFS manager
	vfsManager, err := vfs.NewVFSManager()
	if err != nil {
		return fmt.Errorf("failed to create VFS manager: %v", err)
	}

	// Initialize container filesystem if not already done
	if err := initializeContainerVFS(vfsManager, containerID); err != nil {
		return err
	}

	// Get flags
	longFormat, _ := cmd.Flags().GetBool("long")
	showAll, _ := cmd.Flags().GetBool("all")
	humanReadable, _ := cmd.Flags().GetBool("human")
	recursive, _ := cmd.Flags().GetBool("recursive")

	return listDirectoryVFS(vfsManager, containerID, path, longFormat, showAll, humanReadable, recursive)
}

func displayFile(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	containerID := args[0]
	filePath := args[1]

	// Create VFS manager
	vfsManager, err := vfs.NewVFSManager()
	if err != nil {
		return fmt.Errorf("failed to create VFS manager: %v", err)
	}

	// Initialize container filesystem
	if err := initializeContainerVFS(vfsManager, containerID); err != nil {
		return err
	}

	vfsSystem := vfsManager.GetVFS()

	// Read file contents
	reader, err := vfsSystem.Read(containerID, filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	defer reader.Close()

	// Copy contents to stdout
	content, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read file contents: %v", err)
	}

	fmt.Print(string(content))
	return nil
}

func statFile(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	containerID := args[0]
	filePath := args[1]

	// Get container rootfs path
	rootfsPath, err := getContainerRootFS(containerID)
	if err != nil {
		return err
	}

	// Construct full path
	fullPath := filepath.Join(rootfsPath, filePath)

	// Get file info
	info, err := os.Stat(fullPath)
	if err != nil {
		return fmt.Errorf("failed to stat file: %v", err)
	}

	// Display file information
	fmt.Printf("  File: %s\n", filePath)
	fmt.Printf("  Size: %-12d Blocks: %-8d IO Block: 4096\n", info.Size(), (info.Size()+4095)/4096)

	fileType := "regular file"
	if info.IsDir() {
		fileType = "directory"
	}
	fmt.Printf("Device: unknown     Inode: unknown      Links: 1      Type: %s\n", fileType)
	fmt.Printf("Access: %s\n", info.Mode().String())
	fmt.Printf("Modify: %s\n", info.ModTime().Format("2006-01-02 15:04:05.000000000 -0700"))

	return nil
}

func findFiles(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	containerID := args[0]
	searchPath := args[1]

	// Get container rootfs path
	rootfsPath, err := getContainerRootFS(containerID)
	if err != nil {
		return err
	}

	// Construct full search path
	fullSearchPath := filepath.Join(rootfsPath, searchPath)

	// Get search criteria
	namePattern, _ := cmd.Flags().GetString("name")
	fileType, _ := cmd.Flags().GetString("type")
	maxDepth, _ := cmd.Flags().GetInt("maxdepth")

	// If name provided as argument, use it
	if len(args) > 2 {
		namePattern = args[2]
	}

	return findInDirectory(fullSearchPath, searchPath, namePattern, fileType, maxDepth, 0)
}

func showWorkingDir(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	containerID := args[0]

	// For now, just return root as the working directory
	// In a real implementation, this would check the container's current working directory
	fmt.Println("/")

	// TODO: Implement proper working directory tracking for containers
	_ = containerID
	return nil
}

// Helper function to get container rootfs path
func getContainerRootFS(containerID string) (string, error) {
	sm := state.NewStateManager()

	// Load container state
	container, err := sm.LoadContainer(containerID)
	if err != nil {
		return "", fmt.Errorf("container not found: %s", containerID)
	}

	// Try different possible rootfs paths
	possiblePaths := []string{
		fmt.Sprintf("/var/lib/servin/containers/%s/rootfs", container.ID),
		fmt.Sprintf("/var/lib/servin/containers/%s/rootfs", containerID),
		fmt.Sprintf("/tmp/servin/containers/%s/rootfs", container.ID),
		fmt.Sprintf("/tmp/servin/containers/%s/rootfs", containerID),
		fmt.Sprintf("/Users/%s/.servin/containers/%s/rootfs", os.Getenv("USER"), container.ID),
		fmt.Sprintf("/Users/%s/.servin/containers/%s/rootfs", os.Getenv("USER"), containerID),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// For testing purposes, return a fallback path that indicates no rootfs
	// The calling function should handle this case
	return fmt.Sprintf("/nonexistent/rootfs/%s", containerID), nil
}

// Helper function to list directory contents
func listDirectory(fullPath, containerPath string, longFormat, showAll, humanReadable, recursive bool) error {
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	for _, entry := range entries {
		// Skip hidden files unless -a flag is used
		if !showAll && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		if longFormat {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			// Format like 'ls -l'
			mode := info.Mode().String()
			size := info.Size()
			modTime := info.ModTime().Format("Jan 02 15:04")

			sizeStr := strconv.FormatInt(size, 10)
			if humanReadable {
				sizeStr = formatFileSize(size)
			}

			fmt.Printf("%s %8s %s %s\n", mode, sizeStr, modTime, entry.Name())
		} else {
			fmt.Println(entry.Name())
		}

		// Recursive listing
		if recursive && entry.IsDir() {
			subPath := filepath.Join(fullPath, entry.Name())
			subContainerPath := filepath.Join(containerPath, entry.Name())
			fmt.Printf("\n%s:\n", subContainerPath)
			listDirectory(subPath, subContainerPath, longFormat, showAll, humanReadable, false)
		}
	}

	return nil
}

// Helper function to find files
func findInDirectory(fullPath, containerPath, namePattern, fileType string, maxDepth, currentDepth int) error {
	if maxDepth >= 0 && currentDepth > maxDepth {
		return nil
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil // Skip directories we can't read
	}

	for _, entry := range entries {
		entryPath := filepath.Join(containerPath, entry.Name())

		// Check type filter
		if fileType != "" {
			if fileType == "f" && entry.IsDir() {
				continue
			}
			if fileType == "d" && !entry.IsDir() {
				continue
			}
		}

		// Check name pattern
		matched := true
		if namePattern != "" {
			matched, _ = filepath.Match(namePattern, entry.Name())
		}

		if matched {
			fmt.Println(entryPath)
		}

		// Recurse into directories
		if entry.IsDir() {
			subFullPath := filepath.Join(fullPath, entry.Name())
			findInDirectory(subFullPath, entryPath, namePattern, fileType, maxDepth, currentDepth+1)
		}
	}

	return nil
}

// Helper function to format file sizes
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%c", float64(size)/float64(div), "KMGTPE"[exp])
}

// Initialize container VFS if not already done
func initializeContainerVFS(vfsManager *vfs.VFSManager, containerID string) error {
	vfsSystem := vfsManager.GetVFS()

	// Try to mount first (in case it's already initialized)
	if err := vfsSystem.Mount(containerID); err == nil {
		return nil // Already initialized and mounted
	}

	// Get container information to find rootfs
	rootfsPath, err := getContainerRootFS(containerID)
	if err != nil {
		// If traditional rootfs lookup fails, initialize with empty rootfs
		if initErr := vfsSystem.Initialize(containerID, ""); initErr != nil {
			return fmt.Errorf("failed to initialize container VFS: %v", initErr)
		}
	} else {
		// Initialize with existing rootfs
		if initErr := vfsSystem.Initialize(containerID, rootfsPath); initErr != nil {
			return fmt.Errorf("failed to initialize container VFS: %v", initErr)
		}
	}

	// Mount the filesystem
	if err := vfsSystem.Mount(containerID); err != nil {
		return fmt.Errorf("failed to mount container VFS: %v", err)
	}

	return nil
}

// List directory using VFS
func listDirectoryVFS(vfsManager *vfs.VFSManager, containerID, path string, longFormat, showAll, humanReadable, recursive bool) error {
	vfsSystem := vfsManager.GetVFS()

	files, err := vfsSystem.List(containerID, path)
	if err != nil {
		return fmt.Errorf("failed to list directory: %v", err)
	}

	// Filter hidden files if not showing all
	if !showAll {
		var filteredFiles []vfs.FileInfo
		for _, file := range files {
			if !strings.HasPrefix(file.Name, ".") {
				filteredFiles = append(filteredFiles, file)
			}
		}
		files = filteredFiles
	}

	// Print files
	for _, file := range files {
		if longFormat {
			// Long format: permissions size date name
			modeStr := file.Permissions
			sizeStr := fmt.Sprintf("%d", file.Size)
			if humanReadable {
				sizeStr = formatFileSize(file.Size)
			}
			timeStr := file.ModTime.Format("Jan 02 15:04")

			fmt.Printf("%s %8s %s %s\n", modeStr, sizeStr, timeStr, file.Name)
		} else {
			// Simple format: just names
			fmt.Println(file.Name)
		}

		// Recursive listing
		if recursive && file.IsDir && !strings.HasPrefix(file.Name, ".") {
			subPath := path
			if subPath == "/" {
				subPath = "/" + file.Name
			} else {
				subPath = path + "/" + file.Name
			}

			fmt.Printf("\n%s:\n", subPath)
			if err := listDirectoryVFS(vfsManager, containerID, subPath, longFormat, showAll, humanReadable, false); err != nil {
				fmt.Printf("Error listing %s: %v\n", subPath, err)
			}
		}
	}

	return nil
}
