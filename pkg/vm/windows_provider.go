//go:build windows

package vm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"net"
)

// HyperVProvider implements VM operations using Windows Hyper-V or VirtualBox
type HyperVProvider struct {
	config     *VMConfig
	vmPath     string
	sshPort    int
	running    bool
	vmBackend  string // "hyperv" or "virtualbox" or "wsl2"
}

// NewHyperVProvider creates a new Hyper-V provider
func NewHyperVProvider(config *VMConfig) (VMProvider, error) {
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

	provider := &HyperVProvider{
		config:  config,
		vmPath:  vmPath,
		sshPort: sshPort,
		running: false,
	}

	// Determine the best backend
	if provider.isHyperVAvailable() {
		provider.vmBackend = "hyperv"
	} else if provider.isWSL2Available() {
		provider.vmBackend = "wsl2"
	} else if provider.isVirtualBoxAvailable() {
		provider.vmBackend = "virtualbox"
	} else {
		return nil, fmt.Errorf("no supported virtualization backend found (Hyper-V, WSL2, or VirtualBox)")
	}

	return provider, nil
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

// Create creates a new VM using the best available backend
func (p *HyperVProvider) Create(config *VMConfig) error {
	// Ensure VM directory exists
	if err := os.MkdirAll(p.vmPath, 0755); err != nil {
		return fmt.Errorf("failed to create VM directory: %v", err)
	}

	switch p.vmBackend {
	case "hyperv":
		return p.createHyperVVM(config)
	case "wsl2":
		return p.createWSL2VM(config)
	case "virtualbox":
		return p.createVirtualBoxVM(config)
	default:
		return fmt.Errorf("unsupported backend: %s", p.vmBackend)
	}
}

// createWSL2VM creates a VM using WSL2 with automated setup
func (p *HyperVProvider) createWSL2VM(config *VMConfig) error {
	fmt.Println("Setting up WSL2 VM with Alpine Linux...")

	// Check if WSL2 is properly configured
	if err := p.ensureWSL2Setup(); err != nil {
		return fmt.Errorf("failed to setup WSL2: %v", err)
	}

	// Create or import Alpine Linux distribution
	distroName := fmt.Sprintf("servin-%s", config.Name)
	
	// Download Alpine Linux rootfs
	if err := p.downloadAlpineRootFS(); err != nil {
		return fmt.Errorf("failed to download Alpine rootfs: %v", err)
	}

	// Import into WSL2
	rootfsPath := filepath.Join(p.vmPath, "alpine-rootfs.tar.gz")
	cmd := exec.Command("wsl", "--import", distroName, p.vmPath, rootfsPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to import WSL2 distribution: %v", err)
	}

	// Configure the distribution
	if err := p.configureWSL2Distribution(distroName); err != nil {
		return fmt.Errorf("failed to configure WSL2 distribution: %v", err)
	}

	fmt.Println("✅ WSL2 VM created successfully")
	return nil
}

// createHyperVVM creates a VM using Hyper-V with Alpine Linux
func (p *HyperVProvider) createHyperVVM(config *VMConfig) error {
	fmt.Println("Setting up Hyper-V VM with Alpine Linux...")

	// Download Alpine Linux ISO
	if err := p.downloadAlpineISO(); err != nil {
		return fmt.Errorf("failed to download Alpine ISO: %v", err)
	}

	// Create Hyper-V VM
	vmName := config.Name
	vhdPath := filepath.Join(p.vmPath, "disk.vhdx")

	// Create VHD
	createVHDCmd := fmt.Sprintf(`
New-VHD -Path '%s' -SizeBytes %dGB -Dynamic
`, vhdPath, config.DiskSize)

	cmd := exec.Command("powershell", "-Command", createVHDCmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create VHD: %v", err)
	}

	// Create VM
	createVMCmd := fmt.Sprintf(`
New-VM -Name '%s' -MemoryStartupBytes %dMB -VHDPath '%s' -Generation 2
Set-VMProcessor -VMName '%s' -Count %d
Set-VMMemory -VMName '%s' -DynamicMemoryEnabled $false
`, vmName, config.Memory, vhdPath, vmName, config.CPUs, vmName)

	cmd = exec.Command("powershell", "-Command", createVMCmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create Hyper-V VM: %v", err)
	}

	// Attach Alpine ISO
	isoPath := filepath.Join(p.vmPath, "alpine.iso")
	attachISOCmd := fmt.Sprintf(`
Add-VMDvdDrive -VMName '%s' -Path '%s'
`, vmName, isoPath)

	cmd = exec.Command("powershell", "-Command", attachISOCmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to attach ISO: %v", err)
	}

	fmt.Println("✅ Hyper-V VM created successfully")
	return nil
}

// createVirtualBoxVM creates a VM using VirtualBox
func (p *HyperVProvider) createVirtualBoxVM(config *VMConfig) error {
	fmt.Println("Setting up VirtualBox VM with Alpine Linux...")

	vmName := config.Name
	vdiPath := filepath.Join(p.vmPath, "disk.vdi")

	// Create VM
	cmd := exec.Command("VBoxManage", "createvm", "--name", vmName, "--register")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create VirtualBox VM: %v", err)
	}

	// Configure VM
	configCmds := [][]string{
		{"VBoxManage", "modifyvm", vmName, "--memory", strconv.Itoa(config.Memory)},
		{"VBoxManage", "modifyvm", vmName, "--cpus", strconv.Itoa(config.CPUs)},
		{"VBoxManage", "modifyvm", vmName, "--ostype", "Linux26_64"},
		{"VBoxManage", "modifyvm", vmName, "--natpf1", fmt.Sprintf("ssh,tcp,,%d,,22", p.sshPort)},
	}

	for _, cmdArgs := range configCmds {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to configure VirtualBox VM: %v", err)
		}
	}

	// Create disk
	cmd = exec.Command("VBoxManage", "createhd", "--filename", vdiPath, "--size", strconv.Itoa(config.DiskSize*1024))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create VirtualBox disk: %v", err)
	}

	// Attach disk
	cmd = exec.Command("VBoxManage", "storagectl", vmName, "--name", "SATA", "--add", "sata")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add storage controller: %v", err)
	}

	cmd = exec.Command("VBoxManage", "storageattach", vmName, "--storagectl", "SATA", "--port", "0", "--device", "0", "--type", "hdd", "--medium", vdiPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to attach disk: %v", err)
	}

	// Download and attach Alpine ISO
	if err := p.downloadAlpineISO(); err != nil {
		return fmt.Errorf("failed to download Alpine ISO: %v", err)
	}

	isoPath := filepath.Join(p.vmPath, "alpine.iso")
	cmd = exec.Command("VBoxManage", "storageattach", vmName, "--storagectl", "SATA", "--port", "1", "--device", "0", "--type", "dvddrive", "--medium", isoPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to attach ISO: %v", err)
	}

	fmt.Println("✅ VirtualBox VM created successfully")
	return nil
}

// Start starts the VM using the configured backend
func (p *HyperVProvider) Start() error {
	if p.running {
		return nil
	}

	switch p.vmBackend {
	case "hyperv":
		return p.startHyperVVM()
	case "wsl2":
		return p.startWSL2VM()
	case "virtualbox":
		return p.startVirtualBoxVM()
	default:
		return fmt.Errorf("unsupported backend: %s", p.vmBackend)
	}
}

// startWSL2VM starts the WSL2 VM with automated SSH setup
func (p *HyperVProvider) startWSL2VM() error {
	distroName := fmt.Sprintf("servin-%s", p.config.Name)

	fmt.Printf("Starting WSL2 VM: %s\n", distroName)

	// Start the distribution and run setup script
	setupScript := `#!/bin/sh
# Setup SSH in WSL2
apk update
apk add openssh sudo curl wget bash

# Configure SSH
echo 'root:servin123' | chpasswd
echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config
echo 'PasswordAuthentication yes' >> /etc/ssh/sshd_config

# Generate SSH host keys
ssh-keygen -A

# Start SSH (WSL2 doesn't use systemd)
/usr/sbin/sshd

# Setup container environment
modprobe overlay 2>/dev/null || true
modprobe bridge 2>/dev/null || true

echo "WSL2 setup completed" > /tmp/setup-complete
`

	// Write setup script to VM
	scriptPath := filepath.Join(p.vmPath, "setup.sh")
	if err := os.WriteFile(scriptPath, []byte(setupScript), 0755); err != nil {
		return fmt.Errorf("failed to write setup script: %v", err)
	}

	// Copy script to WSL2 and execute
	cmd := exec.Command("wsl", "-d", distroName, "--", "sh", "-c", "cat > /tmp/setup.sh && chmod +x /tmp/setup.sh && /tmp/setup.sh")
	cmd.Stdin = strings.NewReader(setupScript)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run setup script: %v", err)
	}

	// Set up port forwarding for SSH
	cmd = exec.Command("netsh", "interface", "portproxy", "add", "v4tov4",
		fmt.Sprintf("listenport=%d", p.sshPort), "listenaddress=127.0.0.1",
		"connectport=22", "connectaddress=127.0.0.1")
	cmd.Run() // Ignore errors as rule might already exist

	p.running = true
	fmt.Printf("✅ WSL2 VM started with SSH on port %d\n", p.sshPort)

	// Monitor SSH connectivity
	go p.monitorSSHAndDeploy()

	return nil
}

// startHyperVVM starts the Hyper-V VM
func (p *HyperVProvider) startHyperVVM() error {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Start-VM -Name '%s'", p.config.Name))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start Hyper-V VM: %v", err)
	}

	p.running = true
	fmt.Printf("✅ Hyper-V VM started\n")
	return nil
}

// startVirtualBoxVM starts the VirtualBox VM
func (p *HyperVProvider) startVirtualBoxVM() error {
	cmd := exec.Command("VBoxManage", "startvm", p.config.Name, "--type", "headless")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start VirtualBox VM: %v", err)
	}

	p.running = true
	fmt.Printf("✅ VirtualBox VM started with SSH on port %d\n", p.sshPort)

	// Monitor SSH connectivity
	go p.monitorSSHAndDeploy()

	return nil
}

// Stop stops the VM
func (p *HyperVProvider) Stop() error {
	if !p.running {
		return nil
	}

	switch p.vmBackend {
	case "hyperv":
		return p.stopHyperVVM()
	case "wsl2":
		return p.stopWSL2VM()
	case "virtualbox":
		return p.stopVirtualBoxVM()
	default:
		return fmt.Errorf("unsupported backend: %s", p.vmBackend)
	}
}

// stopWSL2VM stops the WSL2 VM
func (p *HyperVProvider) stopWSL2VM() error {
	distroName := fmt.Sprintf("servin-%s", p.config.Name)
	cmd := exec.Command("wsl", "--terminate", distroName)
	cmd.Run() // Ignore errors

	// Remove port forwarding
	cmd = exec.Command("netsh", "interface", "portproxy", "delete", "v4tov4",
		fmt.Sprintf("listenport=%d", p.sshPort), "listenaddress=127.0.0.1")
	cmd.Run()

	p.running = false
	fmt.Println("✅ WSL2 VM stopped")
	return nil
}

// stopHyperVVM stops the Hyper-V VM
func (p *HyperVProvider) stopHyperVVM() error {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Stop-VM -Name '%s' -Force", p.config.Name))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop Hyper-V VM: %v", err)
	}

	p.running = false
	fmt.Println("✅ Hyper-V VM stopped")
	return nil
}

// stopVirtualBoxVM stops the VirtualBox VM
func (p *HyperVProvider) stopVirtualBoxVM() error {
	cmd := exec.Command("VBoxManage", "controlvm", p.config.Name, "poweroff")
	cmd.Run() // Ignore errors

	p.running = false
	fmt.Println("✅ VirtualBox VM stopped")
	return nil
}

// Destroy removes the VM completely
func (p *HyperVProvider) Destroy() error {
	if p.running {
		p.Stop()
	}

	switch p.vmBackend {
	case "hyperv":
		cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Remove-VM -Name '%s' -Force", p.config.Name))
		cmd.Run()
	case "wsl2":
		distroName := fmt.Sprintf("servin-%s", p.config.Name)
		cmd := exec.Command("wsl", "--unregister", distroName)
		cmd.Run()
	case "virtualbox":
		cmd := exec.Command("VBoxManage", "unregistervm", p.config.Name, "--delete")
		cmd.Run()
	}

	return os.RemoveAll(p.vmPath)
}

// IsRunning checks if the VM is currently running
func (p *HyperVProvider) IsRunning() bool {
	switch p.vmBackend {
	case "wsl2":
		distroName := fmt.Sprintf("servin-%s", p.config.Name)
		cmd := exec.Command("wsl", "-l", "--running")
		output, err := cmd.Output()
		if err != nil {
			p.running = false
			return false
		}
		p.running = strings.Contains(string(output), distroName)
		return p.running
	case "hyperv":
		cmd := exec.Command("powershell", "-Command", fmt.Sprintf("(Get-VM -Name '%s').State", p.config.Name))
		output, err := cmd.Output()
		if err != nil {
			p.running = false
			return false
		}
		p.running = strings.Contains(string(output), "Running")
		return p.running
	case "virtualbox":
		cmd := exec.Command("VBoxManage", "showvminfo", p.config.Name, "--machinereadable")
		output, err := cmd.Output()
		if err != nil {
			p.running = false
			return false
		}
		p.running = strings.Contains(string(output), "VMState=\"running\"")
		return p.running
	}
	return false
}

// GetInfo returns VM information
func (p *HyperVProvider) GetInfo() (*VMInfo, error) {
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
		Platform:   "Windows",
		Provider:   fmt.Sprintf("%s/%s", "Windows", strings.ToUpper(p.vmBackend)),
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
			"nested_virt":  p.vmBackend == "hyperv",
			"ssh_access":   p.testSSHConnectivity(),
		},
	}, nil
}

// Helper methods for checking backend availability
func (p *HyperVProvider) isHyperVAvailable() bool {
	cmd := exec.Command("powershell", "-Command", "Get-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "State") && strings.Contains(string(output), "Enabled")
}

func (p *HyperVProvider) isWSL2Available() bool {
	cmd := exec.Command("wsl", "--status")
	err := cmd.Run()
	return err == nil
}

func (p *HyperVProvider) isVirtualBoxAvailable() bool {
	_, err := exec.LookPath("VBoxManage")
	return err == nil
}

// Download helpers
func (p *HyperVProvider) downloadAlpineISO() error {
	isoPath := filepath.Join(p.vmPath, "alpine.iso")
	if _, err := os.Stat(isoPath); err == nil {
		return nil // Already exists
	}

	fmt.Println("Downloading Alpine Linux ISO...")
	url := "https://dl-cdn.alpinelinux.org/alpine/v3.19/releases/x86_64/alpine-standard-3.19.1-x86_64.iso"
	
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`
Invoke-WebRequest -Uri '%s' -OutFile '%s'
`, url, isoPath))
	
	return cmd.Run()
}

func (p *HyperVProvider) downloadAlpineRootFS() error {
	rootfsPath := filepath.Join(p.vmPath, "alpine-rootfs.tar.gz")
	if _, err := os.Stat(rootfsPath); err == nil {
		return nil // Already exists
	}

	fmt.Println("Downloading Alpine Linux rootfs...")
	url := "https://dl-cdn.alpinelinux.org/alpine/v3.19/releases/x86_64/alpine-minirootfs-3.19.1-x86_64.tar.gz"
	
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`
Invoke-WebRequest -Uri '%s' -OutFile '%s'
`, url, rootfsPath))
	
	return cmd.Run()
}

// WSL2 specific helpers
func (p *HyperVProvider) ensureWSL2Setup() error {
	// Check if WSL2 is the default version
	cmd := exec.Command("wsl", "--set-default-version", "2")
	return cmd.Run()
}

func (p *HyperVProvider) configureWSL2Distribution(distroName string) error {
	// Set default user and configure
	cmd := exec.Command("wsl", "-d", distroName, "--", "sh", "-c", `
# Setup root user
echo 'root:servin123' | chpasswd

# Setup package manager
echo "http://dl-cdn.alpinelinux.org/alpine/v3.19/main" > /etc/apk/repositories
echo "http://dl-cdn.alpinelinux.org/alpine/v3.19/community" >> /etc/apk/repositories
apk update
`)
	return cmd.Run()
}

// SSH and monitoring helpers
func (p *HyperVProvider) testSSHConnectivity() bool {
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

func (p *HyperVProvider) monitorSSHAndDeploy() {
	maxWait := 90 * time.Second
	start := time.Now()

	for time.Since(start) < maxWait {
		if p.testSSHConnectivity() {
			fmt.Println("✅ SSH is ready!")
			time.Sleep(2 * time.Second)
			
			if err := p.deployServinToVM(); err != nil {
				fmt.Printf("⚠️ Failed to deploy Servin to VM: %v\n", err)
			}
			return
		}
		time.Sleep(2 * time.Second)
	}

	fmt.Println("⚠️ SSH setup timeout - manual configuration may be needed")
}

func (p *HyperVProvider) deployServinToVM() error {
	// Find the servin binary
	servinBinary := "./servin.exe"
	if _, err := os.Stat(servinBinary); err != nil {
		for _, path := range []string{"servin.exe", "build/servin.exe"} {
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

	// Copy binary to VM using scp (if available) or WSL copy
	if p.vmBackend == "wsl2" {
		distroName := fmt.Sprintf("servin-%s", p.config.Name)
		cmd := exec.Command("wsl", "-d", distroName, "--", "mkdir", "-p", "/usr/local/bin")
		cmd.Run()

		// Copy using WSL
		cmd = exec.Command("cmd", "/C", fmt.Sprintf(`copy "%s" "\\\\wsl$\\%s\\usr\\local\\bin\\servin"`, servinBinary, distroName))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to copy binary: %v", err)
		}
	} else {
		// Use SCP for Hyper-V and VirtualBox
		cmd := exec.Command("scp",
			"-P", strconv.Itoa(p.sshPort),
			"-o", "StrictHostKeyChecking=no",
			"-o", "UserKnownHostsFile=/dev/null",
			servinBinary,
			"root@localhost:/usr/local/bin/servin")

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to copy binary: %v", err)
		}
	}

	// Make it executable
	cmd := exec.Command("ssh",
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

// Container operations (using Servin runtime)
func (p *HyperVProvider) RunContainer(config *ContainerConfig) (*ContainerResult, error) {
	if !p.IsRunning() {
		return nil, fmt.Errorf("VM is not running")
	}

	// Build servin run command
	servinCmd := p.buildServinCommand(config)

	// Execute via SSH or WSL
	var cmd *exec.Cmd
	if p.vmBackend == "wsl2" {
		distroName := fmt.Sprintf("servin-%s", p.config.Name)
		cmd = exec.Command("wsl", "-d", distroName, "--", "sh", "-c", servinCmd)
	} else {
		cmd = exec.Command("ssh",
			"-p", strconv.Itoa(p.sshPort),
			"-o", "StrictHostKeyChecking=no",
			"-o", "UserKnownHostsFile=/dev/null",
			"root@localhost",
			servinCmd)
	}

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
		if lines := strings.Split(strings.TrimSpace(string(output)), "\n"); len(lines) > 0 {
			result.ID = strings.TrimSpace(lines[len(lines)-1])
		}
	}

	return result, nil
}

func (p *HyperVProvider) buildServinCommand(config *ContainerConfig) string {
	parts := []string{"/usr/local/bin/servin", "run"}

	if config.Name != "" {
		parts = append(parts, "--name", config.Name)
	}

	for hostPort, containerPort := range config.Ports {
		parts = append(parts, "-p", fmt.Sprintf("%s:%s", hostPort, containerPort))
	}

	for hostPath, containerPath := range config.Volumes {
		parts = append(parts, "-v", fmt.Sprintf("%s:%s", hostPath, containerPath))
	}

	for key, value := range config.Environment {
		parts = append(parts, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	if config.WorkDir != "" {
		parts = append(parts, "-w", config.WorkDir)
	}

	if config.Detached {
		parts = append(parts, "-d")
	}

	parts = append(parts, config.Image)

	if len(config.Command) > 0 {
		parts = append(parts, config.Command...)
	}

	return strings.Join(parts, " ")
}

func (p *HyperVProvider) ListContainers() ([]*ContainerInfo, error) {
	if !p.IsRunning() {
		return nil, fmt.Errorf("VM is not running")
	}

	var cmd *exec.Cmd
	if p.vmBackend == "wsl2" {
		distroName := fmt.Sprintf("servin-%s", p.config.Name)
		cmd = exec.Command("wsl", "-d", distroName, "--", "/usr/local/bin/servin", "list")
	} else {
		cmd = exec.Command("ssh",
			"-p", strconv.Itoa(p.sshPort),
			"-o", "StrictHostKeyChecking=no",
			"-o", "UserKnownHostsFile=/dev/null",
			"root@localhost",
			"/usr/local/bin/servin list")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %v", err)
	}

	return p.parseServinContainerList(string(output)), nil
}

func (p *HyperVProvider) parseServinContainerList(output string) []*ContainerInfo {
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

func (p *HyperVProvider) StopContainer(id string) error {
	return p.executeServinCommand(fmt.Sprintf("/usr/local/bin/servin stop %s", id))
}

func (p *HyperVProvider) RemoveContainer(id string) error {
	return p.executeServinCommand(fmt.Sprintf("/usr/local/bin/servin remove %s", id))
}

func (p *HyperVProvider) executeServinCommand(command string) error {
	var cmd *exec.Cmd
	if p.vmBackend == "wsl2" {
		distroName := fmt.Sprintf("servin-%s", p.config.Name)
		cmd = exec.Command("wsl", "-d", distroName, "--", "sh", "-c", command)
	} else {
		cmd = exec.Command("ssh",
			"-p", strconv.Itoa(p.sshPort),
			"-o", "StrictHostKeyChecking=no",
			"-o", "UserKnownHostsFile=/dev/null",
			"root@localhost",
			command)
	}

	return cmd.Run()
}

func (p *HyperVProvider) CopyToVM(hostPath, vmPath string) error {
	if !p.IsRunning() {
		return fmt.Errorf("VM is not running")
	}

	if p.vmBackend == "wsl2" {
		distroName := fmt.Sprintf("servin-%s", p.config.Name)
		cmd := exec.Command("cmd", "/C", fmt.Sprintf(`copy "%s" "\\\\wsl$\\%s\\%s"`, hostPath, distroName, vmPath))
		return cmd.Run()
	} else {
		cmd := exec.Command("scp",
			"-P", strconv.Itoa(p.sshPort),
			"-o", "StrictHostKeyChecking=no",
			"-o", "UserKnownHostsFile=/dev/null",
			hostPath,
			fmt.Sprintf("root@localhost:%s", vmPath))
		return cmd.Run()
	}
}

func (p *HyperVProvider) CopyFromVM(vmPath, hostPath string) error {
	if !p.IsRunning() {
		return fmt.Errorf("VM is not running")
	}

	if p.vmBackend == "wsl2" {
		distroName := fmt.Sprintf("servin-%s", p.config.Name)
		cmd := exec.Command("cmd", "/C", fmt.Sprintf(`copy "\\\\wsl$\\%s\\%s" "%s"`, distroName, vmPath, hostPath))
		return cmd.Run()
	} else {
		cmd := exec.Command("scp",
			"-P", strconv.Itoa(p.sshPort),
			"-o", "StrictHostKeyChecking=no",
			"-o", "UserKnownHostsFile=/dev/null",
			fmt.Sprintf("root@localhost:%s", vmPath),
			hostPath)
		return cmd.Run()
	}
}

func (p *HyperVProvider) ForwardPort(hostPort, vmPort int) error {
	return fmt.Errorf("dynamic port forwarding not implemented, configure during VM start")
}

func (p *HyperVProvider) RemovePortForward(hostPort int) error {
	return fmt.Errorf("dynamic port forwarding not implemented")
}