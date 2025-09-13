package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"servin/pkg/errors"
	"servin/pkg/logger"
	"servin/pkg/volume"

	"github.com/spf13/cobra"
)

var volumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Manage volumes",
	Long:  "Manage container volumes including creating, listing, and removing volumes.",
}

var volumeLsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List volumes",
	RunE:    runVolumeList,
}

var volumeCreateCmd = &cobra.Command{
	Use:   "create [OPTIONS] VOLUME_NAME",
	Short: "Create a volume",
	Long: `Create a new volume that containers can mount.

Examples:
  servin volume create myvolume
  servin volume create --driver local --label env=prod datavolume`,
	Args: cobra.ExactArgs(1),
	RunE: runVolumeCreate,
}

var volumeRmCmd = &cobra.Command{
	Use:     "rm [OPTIONS] VOLUME [VOLUME...]",
	Aliases: []string{"remove"},
	Short:   "Remove one or more volumes",
	Long: `Remove one or more volumes. You cannot remove a volume that is in use by a container.

Examples:
  servin volume rm myvolume
  servin volume rm volume1 volume2
  servin volume rm --force myvolume`,
	Args: cobra.MinimumNArgs(1),
	RunE: runVolumeRemove,
}

var volumePruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove all unused local volumes",
	Long:  "Remove all unused local volumes. Unused volumes are those not referenced by any containers.",
	RunE:  runVolumePrune,
}

var volumeInspectCmd = &cobra.Command{
	Use:   "inspect VOLUME [VOLUME...]",
	Short: "Display detailed information on one or more volumes",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runVolumeInspect,
}

var volumeRmAllCmd = &cobra.Command{
	Use:   "rm-all",
	Short: "Remove all volumes",
	Long:  "Remove all volumes. Use with caution as this will remove all data in volumes.",
	RunE:  runVolumeRemoveAll,
}

// Volume create flags
var (
	volumeDriver string
	volumeLabels []string
	volumeOpts   []string
)

// Volume remove flags
var (
	volumeForce bool
)

func init() {
	// Add subcommands to volume command
	volumeCmd.AddCommand(volumeLsCmd)
	volumeCmd.AddCommand(volumeCreateCmd)
	volumeCmd.AddCommand(volumeRmCmd)
	volumeCmd.AddCommand(volumeRmAllCmd)
	volumeCmd.AddCommand(volumePruneCmd)
	volumeCmd.AddCommand(volumeInspectCmd)

	// Volume create flags
	volumeCreateCmd.Flags().StringVar(&volumeDriver, "driver", "local", "Specify volume driver name")
	volumeCreateCmd.Flags().StringSliceVarP(&volumeLabels, "label", "l", []string{}, "Set metadata for a volume")
	volumeCreateCmd.Flags().StringSliceVar(&volumeOpts, "opt", []string{}, "Set driver specific options")

	// Volume remove flags
	volumeRmCmd.Flags().BoolVarP(&volumeForce, "force", "f", false, "Force the removal of one or more volumes")

	// Add volume command to root
	rootCmd.AddCommand(volumeCmd)
}

func runVolumeList(cmd *cobra.Command, args []string) error {
	logger.Debug("Starting volume list operation")

	if err := checkRoot(); err != nil {
		return err
	}

	volManager := volume.NewManager()
	volumes, err := volManager.ListVolumes()
	if err != nil {
		logger.Error("Failed to list volumes: %v", err)
		return errors.WrapError(err, errors.ErrTypeVolume, "runVolumeList", "failed to retrieve volume list")
	}

	logger.Info("Found %d volumes", len(volumes))

	if len(volumes) == 0 {
		fmt.Println("No volumes found")
		return nil
	}

	// Create table output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "VOLUME NAME\tDRIVER\tMOUNTPOINT\tCREATED")

	for _, vol := range volumes {
		created := formatTimeVolume(vol.CreatedAt)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			vol.Name, vol.Driver, vol.Mountpoint, created)
		logger.Debug("Listed volume: %s (driver: %s)", vol.Name, vol.Driver)
	}

	logger.Debug("Volume list operation completed successfully")
	return nil
}

func runVolumeCreate(cmd *cobra.Command, args []string) error {
	volumeName := args[0]
	logger.Debug("Starting volume create operation for: %s", volumeName)

	if err := checkRoot(); err != nil {
		return err
	}

	// Validate volume name
	if strings.TrimSpace(volumeName) == "" {
		err := errors.NewValidationError("runVolumeCreate", "volume name cannot be empty")
		logger.Error("Volume creation failed: %v", err)
		return err
	}

	volManager := volume.NewManager()

	// Parse labels
	labels := make(map[string]string)
	for _, label := range volumeLabels {
		if strings.Contains(label, "=") {
			parts := strings.SplitN(label, "=", 2)
			labels[parts[0]] = parts[1]
			logger.Debug("Parsed label: %s=%s", parts[0], parts[1])
		} else {
			labels[label] = ""
			logger.Debug("Parsed label: %s (empty value)", label)
		}
	}

	// Parse options
	options := make(map[string]string)
	for _, opt := range volumeOpts {
		if strings.Contains(opt, "=") {
			parts := strings.SplitN(opt, "=", 2)
			options[parts[0]] = parts[1]
			logger.Debug("Parsed option: %s=%s", parts[0], parts[1])
		} else {
			options[opt] = ""
			logger.Debug("Parsed option: %s (empty value)", opt)
		}
	}

	// Create volume
	vol, err := volManager.CreateVolume(volumeName, volumeDriver, options, labels)
	if err != nil {
		logger.Error("Failed to create volume '%s': %v", volumeName, err)
		return errors.WrapError(err, errors.ErrTypeVolume, "runVolumeCreate",
			fmt.Sprintf("failed to create volume '%s'", volumeName)).
			WithContext("volume_name", volumeName).
			WithContext("driver", volumeDriver).
			WithContext("labels", labels).
			WithContext("options", options)
	}

	logger.Info("Volume '%s' created successfully at %s", vol.Name, vol.Mountpoint)
	fmt.Printf("Volume '%s' created successfully\n", vol.Name)
	fmt.Printf("Mountpoint: %s\n", vol.Mountpoint)
	return nil
}

func runVolumeRemove(cmd *cobra.Command, args []string) error {
	logger.Debug("Starting volume remove operation for volumes: %v", args)

	if err := checkRoot(); err != nil {
		return err
	}

	volManager := volume.NewManager()
	var errorList []string

	for _, volumeName := range args {
		logger.Debug("Attempting to remove volume: %s", volumeName)

		if err := volManager.RemoveVolume(volumeName, volumeForce); err != nil {
			logger.Error("Failed to remove volume '%s': %v", volumeName, err)
			errorList = append(errorList, fmt.Sprintf("failed to remove volume '%s': %v", volumeName, err))
			continue
		}

		logger.Info("Volume '%s' removed successfully", volumeName)
		fmt.Printf("Volume '%s' removed successfully\n", volumeName)
	}

	if len(errorList) > 0 {
		logger.Error("Volume removal completed with errors: %d failed out of %d total", len(errorList), len(args))
		return errors.NewVolumeError("runVolumeRemove", fmt.Sprintf("errors occurred:\n%s", strings.Join(errorList, "\n"))).
			WithContext("failed_volumes", len(errorList)).
			WithContext("total_volumes", len(args))
	}

	logger.Debug("All volume removals completed successfully")
	return nil
}

func runVolumePrune(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	volManager := volume.NewManager()

	// Get confirmation from user
	fmt.Print("WARNING! This will remove all unused volumes.\nAre you sure you want to continue? [y/N] ")
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
		fmt.Println("Operation cancelled")
		return nil
	}

	prunedVolumes, err := volManager.PruneVolumes()
	if err != nil {
		return fmt.Errorf("failed to prune volumes: %v", err)
	}

	if len(prunedVolumes) == 0 {
		fmt.Println("No unused volumes found")
	} else {
		fmt.Printf("Removed volumes:\n")
		for _, volumeName := range prunedVolumes {
			fmt.Printf("  %s\n", volumeName)
		}
	}

	return nil
}

func runVolumeInspect(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	volManager := volume.NewManager()

	for i, volumeName := range args {
		if i > 0 {
			fmt.Println() // Add spacing between volumes
		}

		vol, err := volManager.GetVolume(volumeName)
		if err != nil {
			fmt.Printf("Error inspecting volume '%s': %v\n", volumeName, err)
			continue
		}

		fmt.Printf("Volume: %s\n", vol.Name)
		fmt.Printf("Driver: %s\n", vol.Driver)
		fmt.Printf("Mountpoint: %s\n", vol.Mountpoint)
		fmt.Printf("Created: %s\n", vol.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Scope: %s\n", vol.Scope)

		if len(vol.Labels) > 0 {
			fmt.Println("Labels:")
			for key, value := range vol.Labels {
				fmt.Printf("  %s=%s\n", key, value)
			}
		}

		if len(vol.Options) > 0 {
			fmt.Println("Options:")
			for key, value := range vol.Options {
				fmt.Printf("  %s=%s\n", key, value)
			}
		}

		if len(vol.Status) > 0 {
			fmt.Println("Status:")
			for key, value := range vol.Status {
				fmt.Printf("  %s=%s\n", key, value)
			}
		}
	}

	return nil
}

func formatTimeVolume(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func runVolumeRemoveAll(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	volManager := volume.NewManager()

	// Get confirmation from user
	fmt.Print("WARNING! This will remove ALL volumes and their data.\nAre you sure you want to continue? [y/N] ")
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
		fmt.Println("Operation cancelled")
		return nil
	}

	// Get list of volumes before removal
	volumes, err := volManager.ListVolumes()
	if err != nil {
		return fmt.Errorf("failed to list volumes: %v", err)
	}

	if len(volumes) == 0 {
		fmt.Println("No volumes found to remove")
		return nil
	}

	// Remove all volumes
	if err := volManager.RemoveAllVolumes(volumeForce); err != nil {
		return fmt.Errorf("failed to remove all volumes: %v", err)
	}

	fmt.Printf("Successfully removed %d volumes:\n", len(volumes))
	for _, vol := range volumes {
		fmt.Printf("  %s\n", vol.Name)
	}

	return nil
}
