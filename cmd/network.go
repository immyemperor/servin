package cmd

import (
	"fmt"

	"servin/pkg/network"

	"github.com/spf13/cobra"
)

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage container networks",
	Long:  "Create, list, and manage container networks",
}

var networkLsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List networks",
	Long:    "List all container networks",
	RunE:    listNetworks,
}

var networkInspectCmd = &cobra.Command{
	Use:   "inspect [NETWORK]",
	Short: "Display detailed information about a network",
	Args:  cobra.MaximumNArgs(1),
	RunE:  inspectNetwork,
}

func init() {
	rootCmd.AddCommand(networkCmd)
	networkCmd.AddCommand(networkLsCmd)
	networkCmd.AddCommand(networkInspectCmd)
}

func listNetworks(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	nm := network.NewNetworkManager()

	fmt.Printf("%-15s %-10s %-15s %-20s %-10s\n",
		"NETWORK ID", "NAME", "DRIVER", "SCOPE", "IPAM")

	// For now, just show the default network
	// In a full implementation, you'd have network persistence
	fmt.Printf("%-15s %-10s %-15s %-20s %-10s\n",
		"servin0", "servin0", "bridge", "local", "default")

	_ = nm // Avoid unused variable warning
	return nil
}

func inspectNetwork(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	networkName := "servin0"
	if len(args) > 0 {
		networkName = args[0]
	}

	fmt.Printf("Network: %s\n", networkName)
	fmt.Printf("Driver: bridge\n")
	fmt.Printf("Subnet: 172.17.0.0/16\n")
	fmt.Printf("Gateway: 172.17.0.1\n")
	fmt.Printf("Containers: (network inspection not fully implemented)\n")

	return nil
}
