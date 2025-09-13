//go:build linux

package network

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

// NetworkMode represents different network modes for containers
type NetworkMode string

const (
	// BridgeMode uses a bridge network for container communication
	BridgeMode NetworkMode = "bridge"
	// HostMode uses the host network directly
	HostMode NetworkMode = "host"
	// NoneMode disables networking
	NoneMode NetworkMode = "none"
)

// Network represents a container network
type Network struct {
	Name       string      `json:"name"`
	Mode       NetworkMode `json:"mode"`
	Bridge     string      `json:"bridge"`
	Subnet     *net.IPNet  `json:"subnet"`
	Gateway    net.IP      `json:"gateway"`
	IPAMDriver string      `json:"ipam_driver"`
}

// ContainerNetwork represents network configuration for a specific container
type ContainerNetwork struct {
	ContainerID   string            `json:"container_id"`
	NetworkName   string            `json:"network_name"`
	IP            net.IP            `json:"ip"`
	MAC           string            `json:"mac"`
	VethHost      string            `json:"veth_host"`      // Host-side veth interface
	VethContainer string            `json:"veth_container"` // Container-side veth interface
	PortMappings  []PortMapping     `json:"port_mappings"`
	ExtraHosts    map[string]string `json:"extra_hosts"`
}

// PortMapping represents port forwarding from host to container
type PortMapping struct {
	HostPort      int    `json:"host_port"`
	ContainerPort int    `json:"container_port"`
	Protocol      string `json:"protocol"` // tcp, udp
	HostIP        string `json:"host_ip"`  // bind to specific host IP (optional)
}

// NetworkManager manages container networks
type NetworkManager struct {
	networks map[string]*Network
	ipam     *IPAddressManager
}

// NewNetworkManager creates a new network manager
func NewNetworkManager() *NetworkManager {
	nm := &NetworkManager{
		networks: make(map[string]*Network),
		ipam:     NewIPAddressManager(),
	}

	// Create default bridge network
	if err := nm.CreateDefaultBridge(); err != nil {
		fmt.Printf("Warning: failed to create default bridge: %v\n", err)
	}

	return nm
}

// CreateDefaultBridge creates the default servin bridge network
func (nm *NetworkManager) CreateDefaultBridge() error {
	// Default subnet for servin containers: 172.17.0.0/16
	_, subnet, err := net.ParseCIDR("172.17.0.0/16")
	if err != nil {
		return fmt.Errorf("failed to parse default subnet: %v", err)
	}

	gateway := net.ParseIP("172.17.0.1")
	if gateway == nil {
		return fmt.Errorf("failed to parse gateway IP")
	}

	network := &Network{
		Name:       "servin0",
		Mode:       BridgeMode,
		Bridge:     "servin0",
		Subnet:     subnet,
		Gateway:    gateway,
		IPAMDriver: "default",
	}

	return nm.CreateBridge(network)
}

// CreateBridge creates a bridge network
func (nm *NetworkManager) CreateBridge(network *Network) error {
	bridgeName := network.Bridge

	// Check if bridge already exists
	if nm.bridgeExists(bridgeName) {
		fmt.Printf("Bridge %s already exists\n", bridgeName)
		nm.networks[network.Name] = network
		return nil
	}

	// Create bridge interface
	if err := nm.runCommand("ip", "link", "add", "name", bridgeName, "type", "bridge"); err != nil {
		return fmt.Errorf("failed to create bridge %s: %v", bridgeName, err)
	}

	// Set bridge IP address
	cidr := fmt.Sprintf("%s/%d", network.Gateway.String(),
		func() int { ones, _ := network.Subnet.Mask.Size(); return ones }())

	if err := nm.runCommand("ip", "addr", "add", cidr, "dev", bridgeName); err != nil {
		return fmt.Errorf("failed to set bridge IP: %v", err)
	}

	// Bring bridge up
	if err := nm.runCommand("ip", "link", "set", bridgeName, "up"); err != nil {
		return fmt.Errorf("failed to bring bridge up: %v", err)
	}

	// Enable IP forwarding
	if err := nm.runCommand("sysctl", "-w", "net.ipv4.ip_forward=1"); err != nil {
		fmt.Printf("Warning: failed to enable IP forwarding: %v\n", err)
	}

	// Add iptables rules for NAT (outbound traffic)
	if err := nm.setupNATRules(network); err != nil {
		fmt.Printf("Warning: failed to setup NAT rules: %v\n", err)
	}

	nm.networks[network.Name] = network
	fmt.Printf("Created bridge network %s (%s)\n", network.Name, bridgeName)
	return nil
}

// CreateVethPair creates a virtual ethernet pair for container networking
func (nm *NetworkManager) CreateVethPair(containerID string) (*ContainerNetwork, error) {
	// Generate unique interface names
	vethHost := fmt.Sprintf("veth%s", containerID[:8])
	vethContainer := fmt.Sprintf("veth%s_c", containerID[:8])

	// Create veth pair
	if err := nm.runCommand("ip", "link", "add", vethHost, "type", "veth", "peer", "name", vethContainer); err != nil {
		return nil, fmt.Errorf("failed to create veth pair: %v", err)
	}

	// Get default network
	defaultNetwork := nm.networks["servin0"]
	if defaultNetwork == nil {
		return nil, fmt.Errorf("default network not found")
	}

	// Allocate IP address for container
	containerIP, err := nm.ipam.AllocateIP(defaultNetwork.Subnet)
	if err != nil {
		nm.runCommand("ip", "link", "del", vethHost) // Cleanup on failure
		return nil, fmt.Errorf("failed to allocate IP: %v", err)
	}

	// Generate MAC address
	mac := generateMAC()

	containerNet := &ContainerNetwork{
		ContainerID:   containerID,
		NetworkName:   "servin0",
		IP:            containerIP,
		MAC:           mac,
		VethHost:      vethHost,
		VethContainer: vethContainer,
		PortMappings:  []PortMapping{},
		ExtraHosts:    make(map[string]string),
	}

	return containerNet, nil
}

// AttachContainerToNetwork attaches a container to the bridge network
func (nm *NetworkManager) AttachContainerToNetwork(containerNet *ContainerNetwork, netNS string) error {
	vethHost := containerNet.VethHost
	vethContainer := containerNet.VethContainer
	bridgeName := "servin0"

	// Attach host-side veth to bridge
	if err := nm.runCommand("ip", "link", "set", vethHost, "master", bridgeName); err != nil {
		return fmt.Errorf("failed to attach veth to bridge: %v", err)
	}

	// Bring host-side veth up
	if err := nm.runCommand("ip", "link", "set", vethHost, "up"); err != nil {
		return fmt.Errorf("failed to bring host veth up: %v", err)
	}

	// Move container-side veth to container network namespace
	if netNS != "" {
		if err := nm.runCommand("ip", "link", "set", vethContainer, "netns", netNS); err != nil {
			return fmt.Errorf("failed to move veth to netns: %v", err)
		}
	}

	// Configure container-side interface (this would typically be done inside the container)
	// For now, we'll set up what we can from the host side
	if netNS != "" {
		// Set container interface IP and bring it up
		cidr := fmt.Sprintf("%s/16", containerNet.IP.String())
		if err := nm.runInNetNS(netNS, "ip", "addr", "add", cidr, "dev", vethContainer); err != nil {
			return fmt.Errorf("failed to set container IP: %v", err)
		}

		if err := nm.runInNetNS(netNS, "ip", "link", "set", vethContainer, "up"); err != nil {
			return fmt.Errorf("failed to bring container veth up: %v", err)
		}

		// Set default route
		gateway := "172.17.0.1"
		if err := nm.runInNetNS(netNS, "ip", "route", "add", "default", "via", gateway); err != nil {
			fmt.Printf("Warning: failed to set default route: %v\n", err)
		}
	}

	fmt.Printf("Attached container %s to network (IP: %s)\n",
		containerNet.ContainerID[:12], containerNet.IP.String())
	return nil
}

// DetachContainerFromNetwork removes container from network
func (nm *NetworkManager) DetachContainerFromNetwork(containerNet *ContainerNetwork) error {
	// Delete veth pair (this automatically removes both ends)
	if err := nm.runCommand("ip", "link", "del", containerNet.VethHost); err != nil {
		// Don't return error if interface doesn't exist
		if !strings.Contains(err.Error(), "Cannot find device") {
			return fmt.Errorf("failed to delete veth pair: %v", err)
		}
	}

	// Release IP address
	network := nm.networks[containerNet.NetworkName]
	if network != nil {
		nm.ipam.ReleaseIP(network.Subnet, containerNet.IP)
	}

	fmt.Printf("Detached container %s from network\n", containerNet.ContainerID[:12])
	return nil
}

// SetupPortMapping configures port forwarding from host to container
func (nm *NetworkManager) SetupPortMapping(containerNet *ContainerNetwork, mapping PortMapping) error {
	hostIP := mapping.HostIP
	if hostIP == "" {
		hostIP = "0.0.0.0"
	}

	protocol := strings.ToLower(mapping.Protocol)
	if protocol == "" {
		protocol = "tcp"
	}

	// Add iptables DNAT rule for port forwarding
	rule := []string{
		"-t", "nat",
		"-A", "PREROUTING",
		"-p", protocol,
		"-d", hostIP,
		"--dport", strconv.Itoa(mapping.HostPort),
		"-j", "DNAT",
		"--to-destination", fmt.Sprintf("%s:%d", containerNet.IP.String(), mapping.ContainerPort),
	}

	if err := nm.runCommand("iptables", rule...); err != nil {
		return fmt.Errorf("failed to add port mapping rule: %v", err)
	}

	containerNet.PortMappings = append(containerNet.PortMappings, mapping)
	fmt.Printf("Port mapping: %s:%d -> %s:%d (%s)\n",
		hostIP, mapping.HostPort, containerNet.IP.String(), mapping.ContainerPort, protocol)

	return nil
}

// Helper methods

func (nm *NetworkManager) bridgeExists(bridgeName string) bool {
	err := nm.runCommand("ip", "link", "show", bridgeName)
	return err == nil
}

func (nm *NetworkManager) runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command '%s %v' failed: %v, output: %s",
			name, args, err, string(output))
	}
	return nil
}

func (nm *NetworkManager) runInNetNS(netns, name string, args ...string) error {
	// Prepend 'ip netns exec <namespace>' to the command
	fullArgs := append([]string{"netns", "exec", netns, name}, args...)
	return nm.runCommand("ip", fullArgs...)
}

func (nm *NetworkManager) setupNATRules(network *Network) error {
	subnet := network.Subnet.String()

	// Enable masquerading for outbound traffic from containers
	rules := [][]string{
		{"-t", "nat", "-A", "POSTROUTING", "-s", subnet, "!", "-o", network.Bridge, "-j", "MASQUERADE"},
		{"-A", "FORWARD", "-o", network.Bridge, "-j", "ACCEPT"},
		{"-A", "FORWARD", "-i", network.Bridge, "!", "-o", network.Bridge, "-j", "ACCEPT"},
		{"-A", "FORWARD", "-i", network.Bridge, "-o", network.Bridge, "-j", "ACCEPT"},
	}

	for _, rule := range rules {
		if err := nm.runCommand("iptables", rule...); err != nil {
			return fmt.Errorf("failed to add iptables rule %v: %v", rule, err)
		}
	}

	return nil
}

// generateMAC generates a random MAC address for container interface
func generateMAC() string {
	// Use a locally administered MAC address (second bit of first octet set)
	return fmt.Sprintf("02:42:%02x:%02x:%02x:%02x",
		byte(0), byte(0), byte(0), byte(1)) // Simplified for demo
}

// Cleanup removes all servin network resources
func (nm *NetworkManager) Cleanup() error {
	for _, network := range nm.networks {
		if network.Mode == BridgeMode {
			// Remove bridge
			nm.runCommand("ip", "link", "del", network.Bridge)
		}
	}

	// Remove iptables rules (would need more sophisticated tracking)
	// For now, just print a warning
	fmt.Println("Note: Manual iptables cleanup may be required")

	return nil
}
