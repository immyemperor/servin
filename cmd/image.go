package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"servin/pkg/image"

	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Manage images",
	Long:  "Manage container images including importing, listing, and removing images.",
}

var imageLsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List images",
	RunE:    runImageList,
}

var imageImportCmd = &cobra.Command{
	Use:   "import TARBALL NAME:TAG",
	Short: "Import an image from a tarball",
	Long: `Import a container image from a tarball file.
The tarball should contain a complete filesystem that can be used as a container rootfs.

Examples:
  servin image import alpine.tar.gz alpine:latest
  servin image import ubuntu-base.tgz ubuntu:20.04`,
	Args: cobra.ExactArgs(2),
	RunE: runImageImport,
}

var imageRmCmd = &cobra.Command{
	Use:     "rm IMAGE [IMAGE...]",
	Aliases: []string{"remove"},
	Short:   "Remove one or more images",
	Args:    cobra.MinimumNArgs(1),
	RunE:    runImageRemove,
}

var imagePullCmd = &cobra.Command{
	Use:   "pull IMAGE",
	Short: "Pull an image from a registry",
	Long: `Pull an image from a container registry.
This is a placeholder for future registry support.`,
	Args: cobra.ExactArgs(1),
	RunE: runImagePull,
}

var imageInspectCmd = &cobra.Command{
	Use:   "inspect IMAGE",
	Short: "Display detailed information about an image",
	Args:  cobra.ExactArgs(1),
	RunE:  runImageInspect,
}

var imageTagCmd = &cobra.Command{
	Use:   "tag SOURCE_IMAGE[:TAG] TARGET_IMAGE[:TAG]",
	Short: "Create a tag TARGET_IMAGE that refers to SOURCE_IMAGE",
	Long: `Create a tag that refers to an existing image.

Examples:
  servin image tag alpine:latest alpine:v1.0
  servin image tag 45b0a36b30b7 myapp:latest
  servin image tag ubuntu ubuntu:backup`,
	Args: cobra.ExactArgs(2),
	RunE: runImageTag,
}

func init() {
	// Add subcommands to image command
	imageCmd.AddCommand(imageLsCmd)
	imageCmd.AddCommand(imageImportCmd)
	imageCmd.AddCommand(imageRmCmd)
	imageCmd.AddCommand(imagePullCmd)
	imageCmd.AddCommand(imageInspectCmd)
	imageCmd.AddCommand(imageTagCmd)

	// Add image command to root
	rootCmd.AddCommand(imageCmd)
}

func runImageList(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	imgManager := image.NewManager()
	images, err := imgManager.ListImages()
	if err != nil {
		return fmt.Errorf("failed to list images: %v", err)
	}

	if len(images) == 0 {
		fmt.Println("No images found")
		return nil
	}

	// Create table output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "REPOSITORY\tTAG\tIMAGE ID\tCREATED\tSIZE")

	for _, img := range images {
		for _, repoTag := range img.RepoTags {
			parts := strings.Split(repoTag, ":")
			repo := parts[0]
			tag := "latest"
			if len(parts) > 1 {
				tag = parts[1]
			}

			created := formatTimeImage(img.Created)
			size := formatSize(img.Size)

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				repo, tag, img.ID[:12], created, size)
		}
	}

	return nil
}

func runImageImport(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	tarballPath := args[0]
	imageRef := args[1]

	// Check if tarball exists
	if _, err := os.Stat(tarballPath); os.IsNotExist(err) {
		return fmt.Errorf("tarball file not found: %s", tarballPath)
	}

	// Parse image reference
	name, tag := parseImageReference(imageRef)

	fmt.Printf("Importing image %s:%s from %s...\n", name, tag, tarballPath)

	imgManager := image.NewManager()
	img, err := imgManager.CreateImageFromTarball(tarballPath, name, tag)
	if err != nil {
		return fmt.Errorf("failed to import image: %v", err)
	}

	fmt.Printf("Successfully imported image %s:%s (ID: %s)\n", name, tag, img.ID[:12])
	return nil
}

func runImageRemove(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	imgManager := image.NewManager()

	for _, imageRef := range args {
		fmt.Printf("Removing image %s...\n", imageRef)

		if err := imgManager.RemoveImage(imageRef); err != nil {
			fmt.Printf("Error removing image %s: %v\n", imageRef, err)
			continue
		}

		fmt.Printf("Successfully removed image %s\n", imageRef)
	}

	return nil
}

func runImagePull(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	imageRef := args[0]
	fmt.Printf("Pulling image %s...\n", imageRef)

	// For now, this is a placeholder
	return fmt.Errorf("registry pulling not yet implemented. Please use 'servin image import' with a tarball")
}

func runImageInspect(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	imageRef := args[0]

	imgManager := image.NewManager()
	img, err := imgManager.GetImage(imageRef)
	if err != nil {
		return fmt.Errorf("failed to get image: %v", err)
	}

	fmt.Printf("Image: %s\n", imageRef)
	fmt.Printf("ID: %s\n", img.ID)
	fmt.Printf("Created: %s\n", img.Created.Format(time.RFC3339))
	fmt.Printf("Size: %s\n", formatSize(img.Size))
	fmt.Printf("RootFS Type: %s\n", img.RootFSType)
	fmt.Printf("RootFS Path: %s\n", img.RootFSPath)

	if len(img.RepoTags) > 0 {
		fmt.Printf("Repo Tags: %s\n", strings.Join(img.RepoTags, ", "))
	}

	if len(img.Config.Env) > 0 {
		fmt.Printf("Environment:\n")
		for _, env := range img.Config.Env {
			fmt.Printf("  %s\n", env)
		}
	}

	if len(img.Config.Cmd) > 0 {
		fmt.Printf("Default Command: %s\n", strings.Join(img.Config.Cmd, " "))
	}

	if img.Config.WorkingDir != "" {
		fmt.Printf("Working Directory: %s\n", img.Config.WorkingDir)
	}

	if img.Config.User != "" {
		fmt.Printf("User: %s\n", img.Config.User)
	}

	if len(img.Metadata) > 0 {
		fmt.Printf("Metadata:\n")
		for key, value := range img.Metadata {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	return nil
}

// Helper functions
func parseImageReference(ref string) (string, string) {
	parts := strings.Split(ref, ":")
	if len(parts) == 1 {
		return parts[0], "latest"
	}
	return strings.Join(parts[:len(parts)-1], ":"), parts[len(parts)-1]
}

func formatTimeImage(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "Less than a minute ago"
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		return fmt.Sprintf("%d hours ago", hours)
	} else {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)
	}
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func runImageTag(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	sourceRef := args[0]
	targetTag := args[1]

	imgManager := image.NewManager()

	// Tag the image
	if err := imgManager.TagImage(sourceRef, targetTag); err != nil {
		return fmt.Errorf("failed to tag image: %v", err)
	}

	fmt.Printf("Successfully tagged %s as %s\n", sourceRef, targetTag)
	return nil
}
