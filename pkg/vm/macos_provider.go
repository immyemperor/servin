//go:build darwin

package vm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// VirtualizationFrameworkProvider implements VM operations using macOS Virtualization.framework
type VirtualizationFrameworkProvider struct {
	config  *VMConfig
	vmPath  string
	sshPort int
	running bool
}

// NewVirtualizationFrameworkProvider creates a new Virtualization.framework provider
func NewVirtualizationFrameworkProvider(config *VMConfig) (VMProvider, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	vmPath := filepath.Join(homeDir, ".servin", "vms", config.Name)

	return &VirtualizationFrameworkProvider{
		config:  config,
		vmPath:  vmPath,
		sshPort: config.SSHPort,
		running: false,
	}, nil
}

// Create creates a new VM using QEMU with proper Alpine Linux setup
func (p *VirtualizationFrameworkProvider) Create(config *VMConfig) error {
	// Ensure VM directory exists
	if err := os.MkdirAll(p.vmPath, 0755); err != nil {
		return fmt.Errorf("failed to create VM directory: %v", err)
	}

	return p.createQEMUVM(config)
}

// Start starts the VM using QEMU
func (p *VirtualizationFrameworkProvider) Start() error {
	if p.running {
		return nil
	}

	// Ensure VM is created with Alpine Linux and cloud-init
	diskPath := filepath.Join(p.vmPath, "alpine.qcow2")
	if err := p.createBootableAlpineImage(diskPath); err != nil {
		return fmt.Errorf("failed to create bootable Alpine image: %v", err)
	}

	// Create cloud-init ISO for automated setup
	if err := p.createCloudInitISO(); err != nil {
		return fmt.Errorf("failed to create cloud-init ISO: %v", err)
	}

	return p.startQEMUVM()
}

func (p *VirtualizationFrameworkProvider) createCloudInitISO() error {
	vmDir := p.vmPath
	tempDir := filepath.Join(vmDir, "cloud-init-temp")
	isoPath := filepath.Join(vmDir, "cloud-init.iso")

	// Create temporary directory for cloud-init files
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create comprehensive auto-setup script for Alpine Linux
	autoSetupScript := `#!/bin/ash
# Automated SSH setup script for Alpine Linux VM
# This script runs automatically when the VM boots

echo "Starting automated SSH setup for Servin VM..."

# Wait for system to be ready
sleep 3

# Update package repository
apk update

# Install essential packages
apk add openssh sudo curl bash

# Configure SSH for remote access
echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config
echo 'PasswordAuthentication yes' >> /etc/ssh/sshd_config
echo 'ClientAliveInterval 60' >> /etc/ssh/sshd_config
echo 'ClientAliveCountMax 3' >> /etc/ssh/sshd_config

# Set root password
echo 'root:servin123' | chpasswd

# Create servin user for container management
adduser -D -s /bin/ash servin
echo 'servin:servin123' | chpasswd
echo 'servin ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

# Enable SSH service to start on boot
rc-update add sshd default
rc-service sshd start

# Set up container environment
echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf
sysctl -p

# Create directories for Servin
mkdir -p /var/lib/servin /usr/local/bin

# Create completion marker
echo "SSH setup completed at $(date)" > /var/log/servin-setup.log
echo "✅ SSH is now available on port 22"
echo "✅ Servin VM setup complete!"

# Make this script run automatically on boot
if ! grep -q "autosetup.sh" /etc/local.d/*.start 2>/dev/null; then
    echo '#!/bin/ash' > /etc/local.d/servin-setup.start
    echo 'if [ ! -f /var/log/servin-setup.log ]; then' >> /etc/local.d/servin-setup.start
    echo '    /tmp/autosetup.sh &' >> /etc/local.d/servin-setup.start
    echo 'fi' >> /etc/local.d/servin-setup.start
    chmod +x /etc/local.d/servin-setup.start
    rc-update add local default
fi
`

	// Create meta-data
	metaData := `instance-id: alpine-servin-vm
local-hostname: alpine-servin
`

	// Write files to temp directory
	if err := os.WriteFile(filepath.Join(tempDir, "autosetup.sh"), []byte(autoSetupScript), 0755); err != nil {
		return fmt.Errorf("failed to write auto-setup script: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tempDir, "meta-data"), []byte(metaData), 0644); err != nil {
		return fmt.Errorf("failed to write meta-data: %v", err)
	}

	// Create ISO using hdiutil (macOS native)
	fmt.Println("Creating auto-setup ISO using hdiutil...")

	// Remove existing ISO if it exists
	if p.fileExists(isoPath) {
		os.Remove(isoPath)
	}

	cmd := exec.Command("hdiutil", "makehybrid", "-iso", "-joliet", "-o", isoPath, tempDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create auto-setup ISO: %v, output: %s", err, string(output))
	}

	fmt.Printf("✅ Auto-setup ISO created at: %s\n", isoPath)
	return nil
}

// Stop stops the VM
func (p *VirtualizationFrameworkProvider) Stop() error {
	if !p.running {
		return nil
	}

	// Send shutdown signal via SSH
	cmd := exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"root@localhost",
		"shutdown -h now")

	cmd.Run() // Ignore errors as VM might shutdown before SSH responds

	// Wait for VM to stop
	time.Sleep(5 * time.Second)
	p.running = false

	return nil
}

// Destroy removes the VM completely
func (p *VirtualizationFrameworkProvider) Destroy() error {
	if p.running {
		p.Stop()
	}

	return os.RemoveAll(p.vmPath)
}

// IsRunning checks if the VM is currently running
func (p *VirtualizationFrameworkProvider) IsRunning() bool {
	// Check for running QEMU process
	cmd := exec.Command("pgrep", "-f", "qemu-system-aarch64.*"+p.vmPath)
	err := cmd.Run()

	if err == nil {
		p.running = true
		return true
	}

	// Fallback: Check if we can connect via SSH (if SSH is configured)
	sshCmd := exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "ConnectTimeout=1",
		"root@localhost",
		"echo 'alive'")

	sshErr := sshCmd.Run()
	p.running = (sshErr == nil)
	return p.running
}

// GetInfo returns VM information
func (p *VirtualizationFrameworkProvider) GetInfo() (*VMInfo, error) {
	return &VMInfo{
		Name:       p.config.Name,
		Status:     p.getStatus(),
		Platform:   "macOS",
		Provider:   "Virtualization.framework",
		CPUs:       p.config.CPUs,
		Memory:     p.config.Memory,
		IPAddress:  "127.0.0.1",
		SSHPort:    p.sshPort,
		DockerPort: p.config.DockerPort,
		Capabilities: map[string]bool{
			"containers":   true,
			"networking":   true,
			"volumes":      true,
			"port_forward": true,
		},
	}, nil
}

// RunContainer runs a container inside the VM using Servin's native containerization
func (p *VirtualizationFrameworkProvider) RunContainer(config *ContainerConfig) (*ContainerResult, error) {
	if !p.IsRunning() {
		return nil, fmt.Errorf("VM is not running")
	}

	// Build Servin container command (not Docker!)
	servinCmd := p.buildServinCommand(config)

	// Execute via SSH to run Servin container natively in Linux VM
	cmd := exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "ConnectTimeout=5",
		"root@localhost",
		servinCmd)

	output, err := cmd.CombinedOutput()

	result := &ContainerResult{
		Name:   config.Name,
		Output: string(output),
	}

	if err != nil {
		result.Error = err.Error()
		result.ExitCode = 1
	} else {
		result.Status = "running"
		result.ExitCode = 0
		// Extract container ID from Servin output
		if lines := strings.Split(strings.TrimSpace(string(output)), "\n"); len(lines) > 0 {
			result.ID = strings.TrimSpace(lines[len(lines)-1])
		}
	}

	return result, nil
}

// ListContainers lists containers in the VM
func (p *VirtualizationFrameworkProvider) ListContainers() ([]*ContainerInfo, error) {
	if !p.IsRunning() {
		return nil, fmt.Errorf("VM is not running")
	}

	cmd := exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"root@localhost",
		"docker ps -a --format 'table {{.ID}}\\t{{.Names}}\\t{{.Image}}\\t{{.Status}}\\t{{.CreatedAt}}\\t{{.Command}}'")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %v", err)
	}

	return p.parseContainerList(string(output)), nil
}

// StopContainer stops a container in the VM
func (p *VirtualizationFrameworkProvider) StopContainer(id string) error {
	return p.executeDockerCommand(fmt.Sprintf("docker stop %s", id))
}

// RemoveContainer removes a container in the VM
func (p *VirtualizationFrameworkProvider) RemoveContainer(id string) error {
	return p.executeDockerCommand(fmt.Sprintf("docker rm %s", id))
}

// CopyToVM copies files from host to VM
func (p *VirtualizationFrameworkProvider) CopyToVM(hostPath, vmPath string) error {
	cmd := exec.Command("scp",
		"-P", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		hostPath,
		fmt.Sprintf("root@localhost:%s", vmPath))

	return cmd.Run()
}

// CopyFromVM copies files from VM to host
func (p *VirtualizationFrameworkProvider) CopyFromVM(vmPath, hostPath string) error {
	cmd := exec.Command("scp",
		"-P", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		fmt.Sprintf("root@localhost:%s", vmPath),
		hostPath)

	return cmd.Run()
}

// ForwardPort forwards a port from host to VM
func (p *VirtualizationFrameworkProvider) ForwardPort(hostPort, vmPort int) error {
	// Port forwarding is typically configured during VM creation
	// For dynamic forwarding, we'd need to use SSH tunneling
	return fmt.Errorf("dynamic port forwarding not implemented")
}

// RemovePortForward removes a port forward
func (p *VirtualizationFrameworkProvider) RemovePortForward(hostPort int) error {
	return fmt.Errorf("dynamic port forwarding not implemented")
}

// Helper methods

func (p *VirtualizationFrameworkProvider) isUTMAvailable() bool {
	cmd := exec.Command("which", "utm")
	return cmd.Run() == nil
}

func (p *VirtualizationFrameworkProvider) createUTMVM(config *VMConfig) error {
	// UTM VM creation would require UTM CLI or API
	// For now, fall back to QEMU
	return p.createQEMUVM(config)
}

func (p *VirtualizationFrameworkProvider) createQEMUVM(config *VMConfig) error {
	fmt.Printf("Creating QEMU VM at path: %s\n", p.vmPath)
	diskPath := filepath.Join(p.vmPath, "alpine.qcow2")

	// Check if we already have a VM disk
	if _, err := os.Stat(diskPath); os.IsNotExist(err) {
		fmt.Println("Creating Alpine Linux VM disk...")

		// Create a bootable Alpine Linux image using cloud-init
		if err := p.createBootableAlpineImage(diskPath); err != nil {
			return fmt.Errorf("failed to create bootable Alpine image: %v", err)
		}
	} else {
		fmt.Printf("VM disk already exists at: %s\n", diskPath)
	}

	fmt.Println("QEMU VM disk ready")
	return nil
}

func (p *VirtualizationFrameworkProvider) startUTMVM() error {
	// Start UTM VM
	return fmt.Errorf("UTM integration not implemented")
}

func (p *VirtualizationFrameworkProvider) startQEMUVM() error {
	fmt.Printf("Starting QEMU VM from path: %s\n", p.vmPath)
	diskPath := filepath.Join(p.vmPath, "alpine.qcow2")

	// Check if required files exist
	if !p.fileExists(diskPath) {
		return fmt.Errorf("VM disk not found at %s - VM needs to be created first", diskPath)
	}

	// QEMU command with auto-setup support
	args := []string{
		"-M", "virt,accel=hvf", // Use Hypervisor.framework
		"-cpu", "host",
		"-smp", strconv.Itoa(p.config.CPUs),
		"-m", strconv.Itoa(p.config.Memory),
		"-drive", fmt.Sprintf("file=%s,if=virtio,format=qcow2", diskPath),
		"-netdev", fmt.Sprintf("user,id=net0,hostfwd=tcp::%d-:22", p.sshPort),
		"-device", "virtio-net-pci,netdev=net0",
		"-nographic",
	}

	// Check if we have netboot kernel files and auto-setup ISO
	kernelPath := filepath.Join(p.vmPath, "vmlinuz-virt")
	initrdPath := filepath.Join(p.vmPath, "initramfs-virt")
	autoSetupISO := filepath.Join(p.vmPath, "cloud-init.iso")

	if p.fileExists(kernelPath) && p.fileExists(initrdPath) {
		fmt.Println("Using netboot kernel with automated SSH setup...")
		args = append(args, "-kernel", kernelPath)
		args = append(args, "-initrd", initrdPath)

		// Enhanced kernel parameters for auto-setup
		appendCmd := "console=ttyS0 ip=dhcp modules=loop,squashfs,sd-mod,usb-storage alpine_repo=http://dl-cdn.alpinelinux.org/alpine/v3.19/main"

		// Add auto-setup ISO if available
		if p.fileExists(autoSetupISO) {
			args = append(args, "-drive", fmt.Sprintf("file=%s,if=virtio,media=cdrom", autoSetupISO))
			appendCmd += " autosetup=cdrom"
			fmt.Println("Auto-setup ISO attached for automated SSH configuration")
		}

		args = append(args, "-append", appendCmd)
	} else {
		// Use BIOS boot (slower but more reliable)
		fmt.Println("Using BIOS boot mode...")
	}

	fmt.Printf("Starting QEMU with command: qemu-system-aarch64 %s\n", strings.Join(args, " "))
	cmd := exec.Command("qemu-system-aarch64", args...)

	// Capture stderr to see any QEMU errors
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start QEMU VM: %v", err)
	}

	fmt.Printf("QEMU VM started with PID: %d\n", cmd.Process.Pid)
	fmt.Println("🚀 VM is starting with automated SSH setup...")

	// Wait for VM to boot and auto-configure SSH
	fmt.Println("Waiting for Alpine Linux to boot and configure SSH automatically...")
	for i := 0; i < 60; i++ {
		// Check if SSH is available
		if p.testSSHConnectivity() {
			p.running = true
			fmt.Println("✅ VM is now running with SSH configured automatically!")
			fmt.Printf("SSH access: ssh root@localhost -p %d (password: servin123)\n", p.sshPort)

			// Deploy Servin binary to VM
			if err := p.deployServinToVM(); err != nil {
				fmt.Printf("Warning: Failed to deploy Servin to VM: %v\n", err)
			} else {
				fmt.Println("✅ Servin binary deployed to VM successfully!")
			}

			return nil
		}

		// Show progress
		if i%5 == 0 {
			fmt.Printf("Waiting for SSH auto-setup... (%d/60 seconds)\n", i)
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Println("⚠️  SSH auto-setup taking longer than expected")
	fmt.Printf("Manual setup may be required. Connect to VM console and run:\n")
	fmt.Printf("  mount /dev/sr0 /mnt && /mnt/autosetup.sh\n")
	fmt.Printf("SSH will be available at: ssh root@localhost -p %d\n", p.sshPort)

	p.running = true
	return nil
}

func (p *VirtualizationFrameworkProvider) getStatus() string {
	if p.IsRunning() {
		return "running"
	}
	return "stopped"
}

func (p *VirtualizationFrameworkProvider) buildDockerCommand(config *ContainerConfig) string {
	cmd := []string{"docker", "run"}

	if config.Detached {
		cmd = append(cmd, "-d")
	}

	if config.Name != "" {
		cmd = append(cmd, "--name", config.Name)
	}

	// Add environment variables
	for key, value := range config.Environment {
		cmd = append(cmd, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	// Add port mappings
	for hostPort, containerPort := range config.Ports {
		cmd = append(cmd, "-p", fmt.Sprintf("%s:%s", hostPort, containerPort))
	}

	// Add volumes
	for hostPath, containerPath := range config.Volumes {
		cmd = append(cmd, "-v", fmt.Sprintf("%s:%s", hostPath, containerPath))
	}

	// Add working directory
	if config.WorkDir != "" {
		cmd = append(cmd, "-w", config.WorkDir)
	}

	// Add image
	cmd = append(cmd, config.Image)

	// Add command
	cmd = append(cmd, config.Command...)

	return strings.Join(cmd, " ")
}

func (p *VirtualizationFrameworkProvider) executeDockerCommand(dockerCmd string) error {
	cmd := exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"root@localhost",
		dockerCmd)

	return cmd.Run()
}

func (p *VirtualizationFrameworkProvider) parseContainerList(output string) []*ContainerInfo {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) <= 1 { // Skip header line
		return []*ContainerInfo{}
	}

	containers := make([]*ContainerInfo, 0, len(lines)-1)

	for _, line := range lines[1:] { // Skip header
		fields := strings.Split(line, "\t")
		if len(fields) >= 6 {
			containers = append(containers, &ContainerInfo{
				ID:      fields[0],
				Name:    fields[1],
				Image:   fields[2],
				Status:  fields[3],
				Created: fields[4],
				Command: fields[5],
			})
		}
	}

	return containers
}

// downloadAlpineISO downloads Alpine Linux ISO for VM setup
func (p *VirtualizationFrameworkProvider) downloadAlpineISO(isoPath string) error {
	// Use a lightweight Alpine Linux ISO
	url := "https://dl-cdn.alpinelinux.org/alpine/v3.19/releases/aarch64/alpine-virt-3.19.1-aarch64.iso"

	// Simple download implementation
	cmd := exec.Command("curl", "-L", "-o", isoPath, url)
	return cmd.Run()
}

// createCloudInitConfig creates cloud-init configuration for automated VM setup
func (p *VirtualizationFrameworkProvider) createCloudInitConfig() error {
	userDataPath := filepath.Join(p.vmPath, "user-data")
	metaDataPath := filepath.Join(p.vmPath, "meta-data")

	// Create user-data for cloud-init (automated setup)
	userData := `#cloud-config
users:
  - name: servin
    sudo: ALL=(ALL) NOPASSWD:ALL
    shell: /bin/bash
    ssh_authorized_keys:
      - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC... # We'll generate this
packages:
  - bash
  - curl
  - iproute2
  - iptables
runcmd:
  - mkdir -p /var/lib/servin
  - chown servin:servin /var/lib/servin
  - echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf
  - sysctl -p
ssh_pwauth: true
password: servin
chpasswd:
  expire: false
`

	if err := os.WriteFile(userDataPath, []byte(userData), 0644); err != nil {
		return err
	}

	// Create minimal meta-data
	metaData := `instance-id: servin-vm
local-hostname: servin
`

	if err := os.WriteFile(metaDataPath, []byte(metaData), 0644); err != nil {
		return err
	}

	// In production, we'd create an ISO with cloud-init data
	return nil
}

// createBootableAlpineImage creates a bootable Alpine Linux image with cloud-init
func (p *VirtualizationFrameworkProvider) createBootableAlpineImage(diskPath string) error {
	fmt.Println("Creating bootable Alpine Linux disk image with SSH and Servin support...")

	// Ensure VM directory exists
	if err := os.MkdirAll(p.vmPath, 0755); err != nil {
		return fmt.Errorf("failed to create VM directory: %v", err)
	}

	// Step 1: Create empty disk
	fmt.Println("Creating empty disk image...")
	cmd := exec.Command("qemu-img", "create", "-f", "qcow2", diskPath, "8G")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create disk image: %v, output: %s", err, string(output))
	}

	// Step 2: Download Alpine netboot kernel and initrd
	kernelPath := filepath.Join(p.vmPath, "vmlinuz-virt")
	initrdPath := filepath.Join(p.vmPath, "initramfs-virt")

	if !p.fileExists(kernelPath) {
		fmt.Println("Downloading Alpine kernel...")
		kernelURL := "https://dl-cdn.alpinelinux.org/alpine/v3.19/releases/aarch64/netboot-3.19.1/vmlinuz-virt"
		cmd = exec.Command("curl", "-L", "-o", kernelPath, kernelURL)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to download kernel: %v", err)
		}
	}

	if !p.fileExists(initrdPath) {
		fmt.Println("Downloading Alpine initramfs...")
		initrdURL := "https://dl-cdn.alpinelinux.org/alpine/v3.19/releases/aarch64/netboot-3.19.1/initramfs-virt"
		cmd = exec.Command("curl", "-L", "-o", initrdPath, initrdURL)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to download initramfs: %v", err)
		}
	}

	// Step 3: Create cloud-init configuration for automated setup
	if err := p.createCloudInitISO(); err != nil {
		return fmt.Errorf("failed to create cloud-init ISO: %v", err)
	}

	// Step 4: Set up SSH keys
	if err := p.setupSSHKeys(); err != nil {
		fmt.Printf("Warning: SSH key setup failed: %v\n", err)
	}

	fmt.Println("Alpine Linux VM components ready")
	fmt.Println("VM will auto-configure SSH and install Servin on first boot")

	return nil
}

// setupSSHKeys sets up SSH keys for passwordless VM access
func (p *VirtualizationFrameworkProvider) setupSSHKeys() error {
	sshDir := filepath.Join(os.Getenv("HOME"), ".ssh")
	keyPath := filepath.Join(sshDir, "servin_vm_key")

	// Generate SSH key if it doesn't exist
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		fmt.Println("Generating SSH key for VM access...")
		cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "2048", "-f", keyPath, "-N", "")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to generate SSH key: %v", err)
		}
	}

	return nil
}

// needsInstallation checks if the VM disk needs OS installation
func (p *VirtualizationFrameworkProvider) needsInstallation(diskPath string) bool {
	stat, err := os.Stat(diskPath)
	if err != nil {
		return true
	}
	// If disk is very small (just created), it needs installation
	return stat.Size() < 10*1024*1024 // Less than 10MB means empty/needs installation
}

// fileExists checks if a file exists
func (p *VirtualizationFrameworkProvider) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// buildServinCommand builds a Servin container command for execution in the VM
func (p *VirtualizationFrameworkProvider) buildServinCommand(config *ContainerConfig) string {
	cmd := []string{"/usr/local/bin/servin", "run"}

	// Add image
	cmd = append(cmd, config.Image)

	// Add command if specified
	if len(config.Command) > 0 {
		cmd = append(cmd, config.Command...)
	}

	// Add container name
	if config.Name != "" {
		cmd = append(cmd, "--name", config.Name)
	}

	// Add environment variables
	for key, value := range config.Environment {
		cmd = append(cmd, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	// Add ports
	for hostPort, containerPort := range config.Ports {
		cmd = append(cmd, "-p", fmt.Sprintf("%s:%s", hostPort, containerPort))
	}

	// Add volumes
	for hostPath, containerPath := range config.Volumes {
		cmd = append(cmd, "-v", fmt.Sprintf("%s:%s", hostPath, containerPath))
	}

	// Add working directory
	if config.WorkDir != "" {
		cmd = append(cmd, "-w", config.WorkDir)
	}

	// Add detached mode if specified
	if config.Detached {
		cmd = append(cmd, "-d")
	}

	return strings.Join(cmd, " ")
}

// createAutoSetupScript creates a script for automated Alpine setup
func (p *VirtualizationFrameworkProvider) createAutoSetupScript() error {
	setupPath := filepath.Join(p.vmPath, "setup.sh")

	setupScript := `#!/bin/sh
# Alpine Linux auto-setup script for Servin VM
set -e

# Install to disk if not already done
if [ ! -f /mnt/servin-installed ]; then
    echo "Setting up Alpine Linux for Servin..."
    
    # Setup repositories
    echo "http://dl-cdn.alpinelinux.org/alpine/v3.19/main" > /etc/apk/repositories
    echo "http://dl-cdn.alpinelinux.org/alpine/v3.19/community" >> /etc/apk/repositories
    apk update
    
    # Install essential packages
    apk add openssh bash curl
    
    # Enable SSH
    rc-update add sshd default
    
    # Create servin user
    adduser -D -s /bin/bash servin
    echo "servin:servin" | chpasswd
    
    # Allow SSH with password (temporary)
    sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config
    sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
    
    # Start SSH
    service sshd start
    
    # Create marker file
    mkdir -p /mnt
    touch /mnt/servin-installed
    
    echo "Servin VM setup complete!"
fi
`

	return os.WriteFile(setupPath, []byte(setupScript), 0755)
}

// commandExists checks if a command is available in PATH
func (p *VirtualizationFrameworkProvider) commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// testSSHConnectivity tests if SSH is available and working
func (p *VirtualizationFrameworkProvider) testSSHConnectivity() bool {
	cmd := exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "ConnectTimeout=2",
		"-o", "BatchMode=yes",
		"root@localhost",
		"echo 'SSH_READY'")

	output, err := cmd.CombinedOutput()
	return err == nil && strings.Contains(string(output), "SSH_READY")
}

// deployServinToVM copies the Servin binary to the VM and makes it executable
func (p *VirtualizationFrameworkProvider) deployServinToVM() error {
	// Get the current Servin binary path
	servinPath := "./servin"
	if !p.fileExists(servinPath) {
		// Try to find servin binary
		if wd, err := os.Getwd(); err == nil {
			servinPath = filepath.Join(wd, "servin")
		}
		if !p.fileExists(servinPath) {
			return fmt.Errorf("servin binary not found")
		}
	}

	// Copy Servin binary to VM
	cmd := exec.Command("scp",
		"-P", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "BatchMode=yes",
		servinPath,
		"root@localhost:/usr/local/bin/servin")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to copy servin binary: %v", err)
	}

	// Make it executable
	cmd = exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "BatchMode=yes",
		"root@localhost",
		"chmod +x /usr/local/bin/servin")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to make servin executable: %v", err)
	}

	return nil
}
