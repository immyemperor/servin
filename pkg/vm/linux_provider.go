//go:build linux

package vm

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// KVMProvider implements VM operations using Linux KVM/QEMU
type KVMProvider struct {
	config  *VMConfig
	vmPath  string
	sshPort int
	running bool
	qemuCmd *exec.Cmd
	qemuPid int
}

// NewKVMProvider creates a new KVM provider
func NewKVMProvider(config *VMConfig) (VMProvider, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	vmPath := filepath.Join(homeDir, ".servin", "vms", config.Name)

	// Find available SSH port
	sshPort := config.SSHPort
	if sshPort == 0 {
		sshPort = 2222
	}
	if !isPortAvailable(sshPort) {
		for port := 2222; port <= 2299; port++ {
			if isPortAvailable(port) {
				sshPort = port
				break
			}
		}
	}

	return &KVMProvider{
		config:  config,
		vmPath:  vmPath,
		sshPort: sshPort,
		running: false,
	}, nil
}

// isPortAvailable checks if a port is available for use
func isPortAvailable(port int) bool {
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// Create creates a new VM using KVM/QEMU with automated Alpine Linux setup
func (p *KVMProvider) Create(config *VMConfig) error {
	// Check if KVM is available
	if !p.isKVMAvailable() {
		return fmt.Errorf("KVM is not available on this system. Enable virtualization in BIOS and load kvm modules")
	}

	// Ensure VM directory exists
	if err := os.MkdirAll(p.vmPath, 0755); err != nil {
		return fmt.Errorf("failed to create VM directory: %v", err)
	}

	return p.createKVMVM(config)
}

// createKVMVM creates a KVM VM with Alpine Linux and automated SSH setup
func (p *KVMProvider) createKVMVM(config *VMConfig) error {
	fmt.Println("Setting up KVM VM with Alpine Linux...")

	// Download Alpine Linux kernel and initramfs if not present
	if err := p.downloadAlpineKernel(); err != nil {
		return fmt.Errorf("failed to download Alpine kernel: %v", err)
	}

	// Create cloud-init ISO with automated SSH setup
	if err := p.createCloudInitISO(); err != nil {
		return fmt.Errorf("failed to create cloud-init ISO: %v", err)
	}

	// Create disk image
	if err := p.createDiskImage(); err != nil {
		return fmt.Errorf("failed to create disk image: %v", err)
	}

	fmt.Println("✅ KVM VM created successfully")
	return nil
}

// downloadAlpineKernel downloads Alpine Linux kernel and initramfs for KVM
func (p *KVMProvider) downloadAlpineKernel() error {
	kernelPath := filepath.Join(p.vmPath, "vmlinuz-virt")
	initramfsPath := filepath.Join(p.vmPath, "initramfs-virt")

	// Check if files already exist
	if _, err := os.Stat(kernelPath); err == nil {
		if _, err := os.Stat(initramfsPath); err == nil {
			fmt.Println("Alpine kernel files already exist")
			return nil
		}
	}

	fmt.Println("Downloading Alpine Linux kernel...")

	// Determine architecture
	arch := "x86_64"
	if strings.Contains(os.Getenv("GOARCH"), "arm") {
		arch = "aarch64"
	}

	baseURL := fmt.Sprintf("https://dl-cdn.alpinelinux.org/alpine/v3.19/releases/%s/netboot-3.19.1", arch)

	// Download kernel
	kernelURL := fmt.Sprintf("%s/vmlinuz-virt", baseURL)
	if err := downloadFile(kernelURL, kernelPath); err != nil {
		return fmt.Errorf("failed to download kernel: %v", err)
	}

	// Download initramfs
	initramfsURL := fmt.Sprintf("%s/initramfs-virt", baseURL)
	if err := downloadFile(initramfsURL, initramfsPath); err != nil {
		return fmt.Errorf("failed to download initramfs: %v", err)
	}

	fmt.Println("✅ Alpine kernel downloaded")
	return nil
}

// createCloudInitISO creates a cloud-init ISO with automated SSH setup
func (p *KVMProvider) createCloudInitISO() error {
	isoPath := filepath.Join(p.vmPath, "cloud-init.iso")
	tempDir := filepath.Join(p.vmPath, "cloud-init-temp")

	// Clean up any existing temp directory
	os.RemoveAll(tempDir)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create autosetup script
	setupScript := `#!/bin/sh
# Automated SSH setup for Servin VM on Linux KVM
echo "Starting automated SSH setup..."

# Update package index
apk update

# Install required packages
apk add openssh sudo curl wget bash

# Configure SSH
echo 'root:servin123' | chpasswd
echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config
echo 'PasswordAuthentication yes' >> /etc/ssh/sshd_config
echo 'PubkeyAuthentication yes' >> /etc/ssh/sshd_config

# Generate SSH host keys
ssh-keygen -A

# Enable and start SSH service
rc-update add sshd default
rc-service sshd start

# Setup container environment
modprobe overlay 2>/dev/null || true
modprobe bridge 2>/dev/null || true
echo overlay >> /etc/modules
echo bridge >> /etc/modules

# Enable IP forwarding
echo 'net.ipv4.ip_forward = 1' >> /etc/sysctl.conf
sysctl -p

# Create servin user
adduser -D servin
echo 'servin:servin123' | chpasswd
addgroup servin wheel

# Success indicator
echo "SSH setup completed successfully" > /tmp/ssh-setup-complete
echo "✅ SSH setup completed - ready for connections"
`

	scriptPath := filepath.Join(tempDir, "autosetup.sh")
	if err := os.WriteFile(scriptPath, []byte(setupScript), 0755); err != nil {
		return fmt.Errorf("failed to write setup script: %v", err)
	}

	// Create cloud-init user-data
	userData := `#cloud-config
users:
  - name: root
    password: servin123
    chpasswd: { expire: False }
    sudo: ['ALL=(ALL) NOPASSWD:ALL']

runcmd:
  - /autosetup.sh

bootcmd:
  - echo "Starting Servin VM setup..."

final_message: "Servin VM setup completed"
`

	userDataPath := filepath.Join(tempDir, "user-data")
	if err := os.WriteFile(userDataPath, []byte(userData), 0644); err != nil {
		return fmt.Errorf("failed to write user-data: %v", err)
	}

	// Create meta-data
	metaData := `instance-id: servin-vm-001
local-hostname: servin-vm
`

	metaDataPath := filepath.Join(tempDir, "meta-data")
	if err := os.WriteFile(metaDataPath, []byte(metaData), 0644); err != nil {
		return fmt.Errorf("failed to write meta-data: %v", err)
	}

	// Create ISO using genisoimage or mkisofs
	var cmd *exec.Cmd
	if _, err := exec.LookPath("genisoimage"); err == nil {
		cmd = exec.Command("genisoimage", "-output", isoPath, "-volid", "cidata", "-joliet", "-rock", tempDir)
	} else if _, err := exec.LookPath("mkisofs"); err == nil {
		cmd = exec.Command("mkisofs", "-o", isoPath, "-V", "cidata", "-J", "-R", tempDir)
	} else {
		return fmt.Errorf("neither genisoimage nor mkisofs found. Install cdrtools or genisoimage package")
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create ISO: %v", err)
	}

	fmt.Println("✅ Cloud-init ISO created")
	return nil
}

// createDiskImage creates a disk image for the VM
func (p *KVMProvider) createDiskImage() error {
	diskPath := filepath.Join(p.vmPath, "disk.qcow2")

	// Check if disk already exists
	if _, err := os.Stat(diskPath); err == nil {
		fmt.Println("Disk image already exists")
		return nil
	}

	// Create qcow2 disk image
	cmd := exec.Command("qemu-img", "create", "-f", "qcow2", diskPath, fmt.Sprintf("%dG", p.config.DiskSize))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create disk image: %v", err)
	}

	fmt.Println("✅ Disk image created")
	return nil
}

// Start starts the VM using KVM/QEMU
func (p *KVMProvider) Start() error {
	if p.running {
		return nil
	}

	return p.startKVMVM()
}

// startKVMVM starts the KVM VM with proper acceleration and SSH automation
func (p *KVMProvider) startKVMVM() error {
	kernelPath := filepath.Join(p.vmPath, "vmlinuz-virt")
	initramfsPath := filepath.Join(p.vmPath, "initramfs-virt")
	diskPath := filepath.Join(p.vmPath, "disk.qcow2")
	isoPath := filepath.Join(p.vmPath, "cloud-init.iso")

	// Verify required files exist
	for _, path := range []string{kernelPath, initramfsPath, diskPath, isoPath} {
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("required file not found: %s", path)
		}
	}

	// Determine QEMU binary
	qemuBinary := "qemu-system-x86_64"
	if strings.Contains(os.Getenv("GOARCH"), "arm") {
		qemuBinary = "qemu-system-aarch64"
	}

	// Check if QEMU is available
	if _, err := exec.LookPath(qemuBinary); err != nil {
		return fmt.Errorf("QEMU not found: %s. Install qemu-kvm package", qemuBinary)
	}

	// Build QEMU command
	qemuArgs := []string{
		"-enable-kvm", // Enable KVM acceleration
		"-m", strconv.Itoa(p.config.Memory),
		"-smp", strconv.Itoa(p.config.CPUs),
		"-kernel", kernelPath,
		"-initrd", initramfsPath,
		"-append", "console=ttyS0 ip=dhcp ssh=1 SERVIN_AUTO_SETUP=1",
		"-drive", fmt.Sprintf("file=%s,format=qcow2", diskPath),
		"-drive", fmt.Sprintf("file=%s,media=cdrom", isoPath),
		"-netdev", fmt.Sprintf("user,id=net0,hostfwd=tcp::%d-:22", p.sshPort),
		"-device", "virtio-net,netdev=net0",
		"-nographic",
		"-serial", "stdio",
		"-daemonize",
	}

	// Add CPU features for better performance
	qemuArgs = append(qemuArgs, "-cpu", "host")

	fmt.Printf("Starting KVM VM with SSH on port %d...\n", p.sshPort)
	fmt.Println("VM will boot Alpine Linux with automated SSH setup")

	p.qemuCmd = exec.Command(qemuBinary, qemuArgs...)
	if err := p.qemuCmd.Start(); err != nil {
		return fmt.Errorf("failed to start QEMU: %v", err)
	}

	p.qemuPid = p.qemuCmd.Process.Pid
	p.running = true

	fmt.Printf("✅ KVM VM started (PID: %d)\n", p.qemuPid)
	fmt.Println("⏳ Waiting for SSH setup to complete...")

	// Monitor SSH connectivity
	go p.monitorSSHAndDeploy()

	return nil
}

// monitorSSHAndDeploy monitors SSH connectivity and deploys Servin when ready
func (p *KVMProvider) monitorSSHAndDeploy() {
	maxWait := 90 * time.Second
	start := time.Now()

	for time.Since(start) < maxWait {
		if p.testSSHConnectivity() {
			fmt.Println("✅ SSH is ready!")
			time.Sleep(2 * time.Second) // Let SSH fully stabilize

			if err := p.deployServinToVM(); err != nil {
				fmt.Printf("⚠️ Failed to deploy Servin to VM: %v\n", err)
			}
			return
		}
		time.Sleep(2 * time.Second)
	}

	fmt.Println("⚠️ SSH setup timeout - manual configuration may be needed")
}

// testSSHConnectivity tests if SSH is working
func (p *KVMProvider) testSSHConnectivity() bool {
	cmd := exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "ConnectTimeout=2",
		"-o", "BatchMode=yes",
		"root@localhost",
		"echo 'SSH_WORKING'")

	output, err := cmd.Output()
	return err == nil && strings.Contains(string(output), "SSH_WORKING")
}

// deployServinToVM deploys the Servin binary to the VM
func (p *KVMProvider) deployServinToVM() error {
	// Find the servin binary
	servinBinary := "./servin"
	if _, err := os.Stat(servinBinary); err != nil {
		// Try alternative locations
		for _, path := range []string{"servin", "/usr/local/bin/servin", "build/servin"} {
			if _, err := os.Stat(path); err == nil {
				servinBinary = path
				break
			}
		}
	}

	if _, err := os.Stat(servinBinary); err != nil {
		return fmt.Errorf("servin binary not found")
	}

	fmt.Println("📦 Deploying Servin to VM...")

	// Copy binary to VM
	cmd := exec.Command("scp",
		"-P", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		servinBinary,
		"root@localhost:/usr/local/bin/servin")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to copy binary: %v", err)
	}

	// Make it executable
	cmd = exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"root@localhost",
		"chmod +x /usr/local/bin/servin")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to make binary executable: %v", err)
	}

	fmt.Println("✅ Servin deployed to VM")
	return nil
}

// Stop stops the VM
func (p *KVMProvider) Stop() error {
	if !p.running {
		return nil
	}

	// Try graceful shutdown via SSH first
	if p.testSSHConnectivity() {
		cmd := exec.Command("ssh",
			"-p", strconv.Itoa(p.sshPort),
			"-o", "StrictHostKeyChecking=no",
			"-o", "UserKnownHostsFile=/dev/null",
			"root@localhost",
			"shutdown -h now")
		cmd.Run() // Ignore errors as VM might shutdown before SSH responds

		// Wait for graceful shutdown
		time.Sleep(10 * time.Second)
	}

	// Force kill QEMU process if still running
	if p.qemuPid > 0 {
		if process, err := os.FindProcess(p.qemuPid); err == nil {
			process.Signal(syscall.SIGTERM)
			time.Sleep(5 * time.Second)
			process.Signal(syscall.SIGKILL)
		}
	}

	p.running = false
	p.qemuPid = 0
	fmt.Println("✅ VM stopped")
	return nil
}

// Destroy removes the VM completely
func (p *KVMProvider) Destroy() error {
	if p.running {
		p.Stop()
	}

	return os.RemoveAll(p.vmPath)
}

// IsRunning checks if the VM is currently running
func (p *KVMProvider) IsRunning() bool {
	// First check if we think it's running
	if !p.running {
		return false
	}

	// Check if QEMU process is still alive
	if p.qemuPid > 0 {
		if process, err := os.FindProcess(p.qemuPid); err == nil {
			// Send signal 0 to check if process exists
			if err := process.Signal(syscall.Signal(0)); err == nil {
				// Process exists, now check SSH connectivity
				p.running = p.testSSHConnectivity()
				return p.running
			}
		}
	}

	// Process not found or SSH not working
	p.running = false
	p.qemuPid = 0
	return false
}

// GetInfo returns VM information
func (p *KVMProvider) GetInfo() (*VMInfo, error) {
	status := "stopped"
	if p.IsRunning() {
		status = "running"
	}

	uptime := ""
	if p.running && p.testSSHConnectivity() {
		// Get uptime from VM
		cmd := exec.Command("ssh",
			"-p", strconv.Itoa(p.sshPort),
			"-o", "StrictHostKeyChecking=no",
			"-o", "UserKnownHostsFile=/dev/null",
			"-o", "ConnectTimeout=2",
			"root@localhost",
			"uptime -p")
		if output, err := cmd.Output(); err == nil {
			uptime = strings.TrimSpace(string(output))
		}
	}

	return &VMInfo{
		Name:       p.config.Name,
		Status:     status,
		Platform:   "Linux",
		Provider:   "KVM/QEMU",
		CPUs:       p.config.CPUs,
		Memory:     p.config.Memory,
		IPAddress:  "127.0.0.1",
		SSHPort:    p.sshPort,
		DockerPort: p.config.DockerPort,
		Uptime:     uptime,
		Capabilities: map[string]bool{
			"containers":   true,
			"networking":   true,
			"volumes":      true,
			"port_forward": true,
			"nested_virt":  true, // KVM supports nested virtualization
			"ssh_access":   p.testSSHConnectivity(),
		},
	}, nil
}

// isKVMAvailable checks if KVM acceleration is available
func (p *KVMProvider) isKVMAvailable() bool {
	// Check if KVM device exists
	if _, err := os.Stat("/dev/kvm"); err != nil {
		return false
	}

	// Check if we can access it
	cmd := exec.Command("test", "-r", "/dev/kvm", "-a", "-w", "/dev/kvm")
	return cmd.Run() == nil
}

// downloadFile downloads a file from URL to destination
func downloadFile(url, dest string) error {
	cmd := exec.Command("curl", "-L", "-o", dest, url)
	return cmd.Run()
}

// RunContainer runs a container inside the VM using native Servin runtime
func (p *KVMProvider) RunContainer(config *ContainerConfig) (*ContainerResult, error) {
	if !p.IsRunning() {
		return nil, fmt.Errorf("VM is not running")
	}

	// Build servin run command
	servinCmd := p.buildServinCommand(config)

	// Execute via SSH
	cmd := exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
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
		// Extract container ID from output if available
		if lines := strings.Split(strings.TrimSpace(string(output)), "\n"); len(lines) > 0 {
			result.ID = strings.TrimSpace(lines[len(lines)-1])
		}
	}

	return result, nil
}

// buildServinCommand builds a servin command for container operations
func (p *KVMProvider) buildServinCommand(config *ContainerConfig) string {
	parts := []string{"/usr/local/bin/servin", "run"}

	// Add name if specified
	if config.Name != "" {
		parts = append(parts, "--name", config.Name)
	}

	// Add ports
	for hostPort, containerPort := range config.Ports {
		parts = append(parts, "-p", fmt.Sprintf("%s:%s", hostPort, containerPort))
	}

	// Add volumes
	for hostPath, containerPath := range config.Volumes {
		parts = append(parts, "-v", fmt.Sprintf("%s:%s", hostPath, containerPath))
	}

	// Add environment variables
	for key, value := range config.Environment {
		parts = append(parts, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	// Add working directory
	if config.WorkDir != "" {
		parts = append(parts, "-w", config.WorkDir)
	}

	// Add detached mode
	if config.Detached {
		parts = append(parts, "-d")
	}

	// Add image
	parts = append(parts, config.Image)

	// Add command
	if len(config.Command) > 0 {
		parts = append(parts, config.Command...)
	}

	return strings.Join(parts, " ")
}

// ListContainers lists containers in the VM using Servin
func (p *KVMProvider) ListContainers() ([]*ContainerInfo, error) {
	if !p.IsRunning() {
		return nil, fmt.Errorf("VM is not running")
	}

	cmd := exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"root@localhost",
		"/usr/local/bin/servin list")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %v", err)
	}

	return p.parseServinContainerList(string(output)), nil
}

// parseServinContainerList parses Servin container list output
func (p *KVMProvider) parseServinContainerList(output string) []*ContainerInfo {
	var containers []*ContainerInfo
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if strings.Contains(line, "CONTAINER") || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 4 {
			container := &ContainerInfo{
				ID:     fields[0],
				Name:   fields[1],
				Image:  fields[2],
				Status: fields[3],
			}
			if len(fields) >= 5 {
				container.Created = fields[4]
			}
			containers = append(containers, container)
		}
	}

	return containers
}

// StopContainer stops a container in the VM
func (p *KVMProvider) StopContainer(id string) error {
	return p.executeServinCommand(fmt.Sprintf("/usr/local/bin/servin stop %s", id))
}

// RemoveContainer removes a container in the VM
func (p *KVMProvider) RemoveContainer(id string) error {
	return p.executeServinCommand(fmt.Sprintf("/usr/local/bin/servin remove %s", id))
}

// executeServinCommand executes a Servin command in the VM
func (p *KVMProvider) executeServinCommand(command string) error {
	cmd := exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"root@localhost",
		command)

	return cmd.Run()
}

// CopyToVM copies a file from host to VM
func (p *KVMProvider) CopyToVM(hostPath, vmPath string) error {
	if !p.IsRunning() {
		return fmt.Errorf("VM is not running")
	}

	cmd := exec.Command("scp",
		"-P", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		hostPath,
		fmt.Sprintf("root@localhost:%s", vmPath))

	return cmd.Run()
}

// CopyFromVM copies a file from VM to host
func (p *KVMProvider) CopyFromVM(vmPath, hostPath string) error {
	if !p.IsRunning() {
		return fmt.Errorf("VM is not running")
	}

	cmd := exec.Command("scp",
		"-P", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		fmt.Sprintf("root@localhost:%s", vmPath),
		hostPath)

	return cmd.Run()
}

// ForwardPort forwards a port from host to VM
func (p *KVMProvider) ForwardPort(hostPort, vmPort int) error {
	// Port forwarding is set up during VM start
	// This is a placeholder for dynamic port forwarding
	return fmt.Errorf("dynamic port forwarding not implemented, configure during VM start")
}

// RemovePortForward removes a port forward
func (p *KVMProvider) RemovePortForward(hostPort int) error {
	// Port forwarding is managed during VM start
	// This is a placeholder for dynamic port forwarding
	return fmt.Errorf("dynamic port forwarding not implemented")
}
