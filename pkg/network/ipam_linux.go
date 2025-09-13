//go:build linux

package network

import (
	"fmt"
	"net"
	"sync"
)

// IPAddressManager manages IP address allocation for containers
type IPAddressManager struct {
	allocated map[string]map[string]bool // subnet -> ip -> allocated
	mutex     sync.RWMutex
}

// NewIPAddressManager creates a new IP address manager
func NewIPAddressManager() *IPAddressManager {
	return &IPAddressManager{
		allocated: make(map[string]map[string]bool),
	}
}

// AllocateIP allocates the next available IP address in the given subnet
func (ipam *IPAddressManager) AllocateIP(subnet *net.IPNet) (net.IP, error) {
	ipam.mutex.Lock()
	defer ipam.mutex.Unlock()

	subnetKey := subnet.String()

	// Initialize subnet tracking if not exists
	if ipam.allocated[subnetKey] == nil {
		ipam.allocated[subnetKey] = make(map[string]bool)
	}

	// Get the network address and find the next available IP
	networkIP := subnet.IP.Mask(subnet.Mask)

	// Start from .2 (skip .0 for network and .1 for gateway)
	ip := make(net.IP, len(networkIP))
	copy(ip, networkIP)

	// Increment to .2
	incrementIP(ip)
	incrementIP(ip)

	// Find next available IP
	for subnet.Contains(ip) {
		ipStr := ip.String()

		// Skip broadcast address
		if isBroadcast(ip, subnet) {
			break
		}

		// Check if IP is available
		if !ipam.allocated[subnetKey][ipStr] {
			ipam.allocated[subnetKey][ipStr] = true
			return dupIP(ip), nil
		}

		incrementIP(ip)
	}

	return nil, fmt.Errorf("no available IP addresses in subnet %s", subnet.String())
}

// ReleaseIP releases an allocated IP address
func (ipam *IPAddressManager) ReleaseIP(subnet *net.IPNet, ip net.IP) {
	ipam.mutex.Lock()
	defer ipam.mutex.Unlock()

	subnetKey := subnet.String()
	ipStr := ip.String()

	if ipam.allocated[subnetKey] != nil {
		delete(ipam.allocated[subnetKey], ipStr)
	}
}

// IsIPAllocated checks if an IP is already allocated
func (ipam *IPAddressManager) IsIPAllocated(subnet *net.IPNet, ip net.IP) bool {
	ipam.mutex.RLock()
	defer ipam.mutex.RUnlock()

	subnetKey := subnet.String()
	ipStr := ip.String()

	if ipam.allocated[subnetKey] == nil {
		return false
	}

	return ipam.allocated[subnetKey][ipStr]
}

// GetAllocatedIPs returns all allocated IPs in a subnet
func (ipam *IPAddressManager) GetAllocatedIPs(subnet *net.IPNet) []net.IP {
	ipam.mutex.RLock()
	defer ipam.mutex.RUnlock()

	subnetKey := subnet.String()
	var ips []net.IP

	if ipam.allocated[subnetKey] != nil {
		for ipStr := range ipam.allocated[subnetKey] {
			if ip := net.ParseIP(ipStr); ip != nil {
				ips = append(ips, ip)
			}
		}
	}

	return ips
}

// GetAvailableIPCount returns the number of available IPs in a subnet
func (ipam *IPAddressManager) GetAvailableIPCount(subnet *net.IPNet) int {
	ipam.mutex.RLock()
	defer ipam.mutex.RUnlock()

	// Calculate total IPs in subnet
	ones, bits := subnet.Mask.Size()
	totalIPs := 1 << uint(bits-ones)

	// Subtract network, gateway, and broadcast
	usableIPs := totalIPs - 3

	// Subtract allocated IPs
	subnetKey := subnet.String()
	allocatedCount := 0
	if ipam.allocated[subnetKey] != nil {
		allocatedCount = len(ipam.allocated[subnetKey])
	}

	available := usableIPs - allocatedCount
	if available < 0 {
		available = 0
	}

	return available
}

// Helper functions

// incrementIP increments an IP address by one
func incrementIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
}

// dupIP duplicates an IP address
func dupIP(ip net.IP) net.IP {
	dup := make(net.IP, len(ip))
	copy(dup, ip)
	return dup
}

// isBroadcast checks if an IP is the broadcast address for the subnet
func isBroadcast(ip net.IP, subnet *net.IPNet) bool {
	// Calculate broadcast address
	broadcast := make(net.IP, len(subnet.IP))
	copy(broadcast, subnet.IP)

	// Set host bits to 1
	for i := 0; i < len(broadcast); i++ {
		broadcast[i] |= ^subnet.Mask[i]
	}

	return ip.Equal(broadcast)
}
