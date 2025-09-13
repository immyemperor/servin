//go:build !linux

package network

import (
	"fmt"
	"net"
)

// NetworkMode represents different network modes for containers
type NetworkMode string

const (
	BridgeMode NetworkMode = "bridge"
	HostMode   NetworkMode = "host"
	NoneMode   NetworkMode = "none"
)

// Network represents a container network (stub for non-Linux)
type Network struct {
	Name       string      `json:"name"`
	Mode       NetworkMode `json:"mode"`
	Bridge     string      `json:"bridge"`
	Subnet     *net.IPNet  `json:"subnet"`
	Gateway    net.IP      `json:"gateway"`
	IPAMDriver string      `json:"ipam_driver"`
}

// ContainerNetwork represents network configuration for a specific container (stub)
type ContainerNetwork struct {
	ContainerID   string            `json:"container_id"`
	NetworkName   string            `json:"network_name"`
	IP            net.IP            `json:"ip"`
	MAC           string            `json:"mac"`
	VethHost      string            `json:"veth_host"`
	VethContainer string            `json:"veth_container"`
	PortMappings  []PortMapping     `json:"port_mappings"`
	ExtraHosts    map[string]string `json:"extra_hosts"`
}

// PortMapping represents port forwarding from host to container (stub)
type PortMapping struct {
	HostPort      int    `json:"host_port"`
	ContainerPort int    `json:"container_port"`
	Protocol      string `json:"protocol"`
	HostIP        string `json:"host_ip"`
}

// NetworkManager manages container networks (stub for non-Linux)
type NetworkManager struct{}

// NewNetworkManager creates a new network manager (stub)
func NewNetworkManager() *NetworkManager {
	return &NetworkManager{}
}

// CreateDefaultBridge creates the default servin bridge network (stub)
func (nm *NetworkManager) CreateDefaultBridge() error {
	return fmt.Errorf("networking is only supported on Linux")
}

// CreateBridge creates a bridge network (stub)
func (nm *NetworkManager) CreateBridge(network *Network) error {
	return fmt.Errorf("networking is only supported on Linux")
}

// CreateVethPair creates a virtual ethernet pair for container networking (stub)
func (nm *NetworkManager) CreateVethPair(containerID string) (*ContainerNetwork, error) {
	return nil, fmt.Errorf("networking is only supported on Linux")
}

// AttachContainerToNetwork attaches a container to the bridge network (stub)
func (nm *NetworkManager) AttachContainerToNetwork(containerNet *ContainerNetwork, netNS string) error {
	return fmt.Errorf("networking is only supported on Linux")
}

// DetachContainerFromNetwork removes container from network (stub)
func (nm *NetworkManager) DetachContainerFromNetwork(containerNet *ContainerNetwork) error {
	return fmt.Errorf("networking is only supported on Linux")
}

// SetupPortMapping configures port forwarding from host to container (stub)
func (nm *NetworkManager) SetupPortMapping(containerNet *ContainerNetwork, mapping PortMapping) error {
	return fmt.Errorf("networking is only supported on Linux")
}

// Cleanup removes all servin network resources (stub)
func (nm *NetworkManager) Cleanup() error {
	return nil
}
