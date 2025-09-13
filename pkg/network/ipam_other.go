//go:build !linux

package network

import (
	"fmt"
	"net"
)

// IPAddressManager manages IP address allocation for containers (stub for non-Linux)
type IPAddressManager struct{}

// NewIPAddressManager creates a new IP address manager (stub)
func NewIPAddressManager() *IPAddressManager {
	return &IPAddressManager{}
}

// AllocateIP allocates the next available IP address in the given subnet (stub)
func (ipam *IPAddressManager) AllocateIP(subnet *net.IPNet) (net.IP, error) {
	return nil, fmt.Errorf("IP address management is only supported on Linux")
}

// ReleaseIP releases an allocated IP address (stub)
func (ipam *IPAddressManager) ReleaseIP(subnet *net.IPNet, ip net.IP) {
	// No-op on non-Linux platforms
}

// IsIPAllocated checks if an IP is already allocated (stub)
func (ipam *IPAddressManager) IsIPAllocated(subnet *net.IPNet, ip net.IP) bool {
	return false
}

// GetAllocatedIPs returns all allocated IPs in a subnet (stub)
func (ipam *IPAddressManager) GetAllocatedIPs(subnet *net.IPNet) []net.IP {
	return nil
}

// GetAvailableIPCount returns the number of available IPs in a subnet (stub)
func (ipam *IPAddressManager) GetAvailableIPCount(subnet *net.IPNet) int {
	return 0
}
