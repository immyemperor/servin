package vm

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

// TestCrossPlatformVMCreation tests VM creation on all platforms
func TestCrossPlatformVMCreation(t *testing.T) {
	config := DefaultVMConfig("test-vm")
	config.Memory = 1024 // Use less memory for testing
	config.CPUs = 1

	provider, err := GetVMProvider(config)
	if err != nil {
		t.Fatalf("Failed to get VM provider for %s: %v", runtime.GOOS, err)
	}

	// Test VM creation
	if err := provider.Create(config); err != nil {
		t.Fatalf("Failed to create VM on %s: %v", runtime.GOOS, err)
	}

	// Clean up
	defer func() {
		if err := provider.Destroy(); err != nil {
			t.Logf("Failed to destroy test VM: %v", err)
		}
	}()

	// Test VM info
	info, err := provider.GetInfo()
	if err != nil {
		t.Fatalf("Failed to get VM info: %v", err)
	}

	if info.Name != config.Name {
		t.Errorf("Expected VM name %s, got %s", config.Name, info.Name)
	}

	if info.Platform != runtime.GOOS {
		t.Errorf("Expected platform %s, got %s", runtime.GOOS, info.Platform)
	}

	t.Logf("✅ VM creation test passed on %s with provider %s", runtime.GOOS, info.Provider)
}

// TestCrossPlatformVMLifecycle tests full VM lifecycle
func TestCrossPlatformVMLifecycle(t *testing.T) {
	config := DefaultVMConfig("lifecycle-test-vm")
	config.Memory = 1024
	config.CPUs = 1

	provider, err := GetVMProvider(config)
	if err != nil {
		t.Fatalf("Failed to get VM provider: %v", err)
	}

	// Create VM
	if err := provider.Create(config); err != nil {
		t.Fatalf("Failed to create VM: %v", err)
	}

	defer func() {
		provider.Destroy()
	}()

	// Test initial state
	if provider.IsRunning() {
		t.Error("VM should not be running after creation")
	}

	// Start VM
	if err := provider.Start(); err != nil {
		t.Fatalf("Failed to start VM: %v", err)
	}

	// Wait a moment for VM to start
	time.Sleep(5 * time.Second)

	// Test running state
	if !provider.IsRunning() {
		t.Error("VM should be running after start")
	}

	// Test VM info while running
	info, err := provider.GetInfo()
	if err != nil {
		t.Fatalf("Failed to get VM info while running: %v", err)
	}

	if info.Status != "running" {
		t.Errorf("Expected status 'running', got '%s'", info.Status)
	}

	// Stop VM
	if err := provider.Stop(); err != nil {
		t.Fatalf("Failed to stop VM: %v", err)
	}

	// Wait for VM to stop
	time.Sleep(5 * time.Second)

	// Test stopped state
	if provider.IsRunning() {
		t.Error("VM should not be running after stop")
	}

	t.Logf("✅ VM lifecycle test passed on %s", runtime.GOOS)
}

// TestCrossPlatformVMCapabilities tests platform-specific capabilities
func TestCrossPlatformVMCapabilities(t *testing.T) {
	config := DefaultVMConfig("capabilities-test-vm")
	provider, err := GetVMProvider(config)
	if err != nil {
		t.Fatalf("Failed to get VM provider: %v", err)
	}

	if err := provider.Create(config); err != nil {
		t.Fatalf("Failed to create VM: %v", err)
	}

	defer provider.Destroy()

	info, err := provider.GetInfo()
	if err != nil {
		t.Fatalf("Failed to get VM info: %v", err)
	}

	// Test expected capabilities by platform
	expectedCapabilities := map[string]map[string]bool{
		"darwin": {
			"containers":   true,
			"networking":   true,
			"volumes":      true,
			"port_forward": true,
			"nested_virt":  true, // QEMU with Hypervisor.framework
		},
		"linux": {
			"containers":   true,
			"networking":   true,
			"volumes":      true,
			"port_forward": true,
			"nested_virt":  true, // KVM supports nested virtualization
		},
		"windows": {
			"containers":   true,
			"networking":   true,
			"volumes":      true,
			"port_forward": true,
			"nested_virt":  false, // Depends on backend, but conservative default
		},
	}

	expected := expectedCapabilities[runtime.GOOS]
	for capability, expectedValue := range expected {
		if actual, exists := info.Capabilities[capability]; !exists {
			t.Errorf("Missing capability: %s", capability)
		} else if actual != expectedValue {
			t.Errorf("Expected %s=%v, got %v", capability, expectedValue, actual)
		}
	}

	t.Logf("✅ VM capabilities test passed on %s", runtime.GOOS)
}

// BenchmarkVMCreation benchmarks VM creation performance
func BenchmarkVMCreation(b *testing.B) {
	config := DefaultVMConfig("benchmark-vm")
	config.Memory = 512 // Minimal for benchmarking
	config.CPUs = 1

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		provider, err := GetVMProvider(config)
		if err != nil {
			b.Fatalf("Failed to get VM provider: %v", err)
		}
		b.StartTimer()

		if err := provider.Create(config); err != nil {
			b.Fatalf("Failed to create VM: %v", err)
		}

		b.StopTimer()
		provider.Destroy()
		b.StartTimer()
	}
}

// TestVMSSHConnectivity tests SSH connectivity after VM start
func TestVMSSHConnectivity(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping SSH connectivity test in short mode")
	}

	config := DefaultVMConfig("ssh-test-vm")
	config.Memory = 1024
	config.CPUs = 1

	provider, err := GetVMProvider(config)
	if err != nil {
		t.Fatalf("Failed to get VM provider: %v", err)
	}

	if err := provider.Create(config); err != nil {
		t.Fatalf("Failed to create VM: %v", err)
	}

	defer provider.Destroy()

	if err := provider.Start(); err != nil {
		t.Fatalf("Failed to start VM: %v", err)
	}

	defer provider.Stop()

	// Wait for SSH to be ready (with timeout)
	maxWait := 120 * time.Second
	start := time.Now()
	sshReady := false

	for time.Since(start) < maxWait {
		info, err := provider.GetInfo()
		if err == nil && info.Capabilities["ssh_access"] {
			sshReady = true
			break
		}
		time.Sleep(5 * time.Second)
	}

	if !sshReady {
		t.Errorf("SSH not ready after %v on %s", maxWait, runtime.GOOS)
		return
	}

	t.Logf("✅ SSH connectivity test passed on %s", runtime.GOOS)
}

// TestContainerOperations tests basic container operations in VM
func TestContainerOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping container operations test in short mode")
	}

	config := DefaultVMConfig("container-test-vm")
	config.Memory = 2048 // More memory for container operations
	config.CPUs = 2

	provider, err := GetVMProvider(config)
	if err != nil {
		t.Fatalf("Failed to get VM provider: %v", err)
	}

	if err := provider.Create(config); err != nil {
		t.Fatalf("Failed to create VM: %v", err)
	}

	defer provider.Destroy()

	if err := provider.Start(); err != nil {
		t.Fatalf("Failed to start VM: %v", err)
	}

	defer provider.Stop()

	// Wait for VM to be ready
	time.Sleep(60 * time.Second)

	// Test container creation
	containerConfig := &ContainerConfig{
		Image:   "hello-world",
		Name:    "test-container",
		Command: []string{},
	}

	result, err := provider.RunContainer(containerConfig)
	if err != nil {
		t.Fatalf("Failed to run container: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("Container exited with code %d, expected 0", result.ExitCode)
		t.Logf("Container output: %s", result.Output)
		t.Logf("Container error: %s", result.Error)
	}

	// Test container listing
	containers, err := provider.ListContainers()
	if err != nil {
		t.Fatalf("Failed to list containers: %v", err)
	}

	found := false
	for _, container := range containers {
		if container.Name == "test-container" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Test container not found in container list")
	}

	t.Logf("✅ Container operations test passed on %s", runtime.GOOS)
}

// TestPlatformSpecificFeatures tests features specific to each platform
func TestPlatformSpecificFeatures(t *testing.T) {
	config := DefaultVMConfig("platform-test-vm")
	provider, err := GetVMProvider(config)
	if err != nil {
		t.Fatalf("Failed to get VM provider: %v", err)
	}

	switch runtime.GOOS {
	case "darwin":
		// Test macOS-specific features
		info, _ := provider.GetInfo()
		if !strings.Contains(info.Provider, "QEMU") && !strings.Contains(info.Provider, "Virtualization") {
			t.Errorf("Expected macOS provider to use QEMU or Virtualization.framework")
		}

	case "linux":
		// Test Linux-specific features
		info, _ := provider.GetInfo()
		if !strings.Contains(info.Provider, "KVM") {
			t.Errorf("Expected Linux provider to use KVM")
		}

	case "windows":
		// Test Windows-specific features
		info, _ := provider.GetInfo()
		if !strings.Contains(info.Provider, "HYPERV") && !strings.Contains(info.Provider, "WSL2") && !strings.Contains(info.Provider, "VIRTUALBOX") {
			t.Errorf("Expected Windows provider to use Hyper-V, WSL2, or VirtualBox")
		}
	}

	t.Logf("✅ Platform-specific features test passed on %s", runtime.GOOS)
}

// TestVMResourceLimits tests VM resource allocation and limits
func TestVMResourceLimits(t *testing.T) {
	testCases := []struct {
		name   string
		memory int
		cpus   int
		valid  bool
	}{
		{"minimal", 512, 1, true},
		{"normal", 2048, 2, true},
		{"high", 4096, 4, true},
		{"invalid-memory", 0, 1, false},
		{"invalid-cpus", 1024, 0, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := DefaultVMConfig(fmt.Sprintf("resource-test-%s", tc.name))
			config.Memory = tc.memory
			config.CPUs = tc.cpus

			provider, err := GetVMProvider(config)
			if err != nil {
				if tc.valid {
					t.Fatalf("Unexpected error getting provider: %v", err)
				}
				return
			}

			err = provider.Create(config)
			if tc.valid && err != nil {
				t.Errorf("Expected valid config to succeed, got error: %v", err)
			} else if !tc.valid && err == nil {
				t.Error("Expected invalid config to fail, but it succeeded")
			}

			if err == nil {
				provider.Destroy()
			}
		})
	}
}

// Helper function to check if we're in a CI environment
func isCI() bool {
	return os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != ""
}

// TestCIEnvironment tests VM functionality in CI environments
func TestCIEnvironment(t *testing.T) {
	if !isCI() {
		t.Skip("Skipping CI test when not in CI environment")
	}

	// In CI, we test with minimal resources and shorter timeouts
	config := DefaultVMConfig("ci-test-vm")
	config.Memory = 512
	config.CPUs = 1

	provider, err := GetVMProvider(config)
	if err != nil {
		t.Fatalf("Failed to get VM provider in CI: %v", err)
	}

	// Test creation only in CI (starting VMs might not work in all CI environments)
	if err := provider.Create(config); err != nil {
		t.Fatalf("Failed to create VM in CI: %v", err)
	}

	defer provider.Destroy()

	info, err := provider.GetInfo()
	if err != nil {
		t.Fatalf("Failed to get VM info in CI: %v", err)
	}

	t.Logf("✅ CI test passed on %s with provider %s", runtime.GOOS, info.Provider)
}
