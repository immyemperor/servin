package cri

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"servin/pkg/image"
	"servin/pkg/logger"
	"servin/pkg/state"
)

// MinimalRuntimeService implements a basic CRI RuntimeService interface
type MinimalRuntimeService struct {
	imageManager *image.Manager
	stateManager *state.StateManager
	logger       *logger.Logger
	criBaseDir   string
}

// NewMinimalRuntimeService creates a new minimal CRI runtime service
func NewMinimalRuntimeService(imageManager *image.Manager, stateManager *state.StateManager, logger *logger.Logger, baseDir string) *MinimalRuntimeService {
	criBaseDir := filepath.Join(baseDir, "cri")
	os.MkdirAll(criBaseDir, 0755)

	return &MinimalRuntimeService{
		imageManager: imageManager,
		stateManager: stateManager,
		logger:       logger,
		criBaseDir:   criBaseDir,
	}
}

// Version returns the runtime name, runtime version, and runtime API version
func (s *MinimalRuntimeService) Version(ctx context.Context, req *VersionRequest) (*VersionResponse, error) {
	s.logger.Info("CRI Version called")

	return &VersionResponse{
		Version:           CRIVersion,
		RuntimeName:       ServinRuntimeName,
		RuntimeVersion:    ServinRuntimeVersion,
		RuntimeApiVersion: CRIVersion,
	}, nil
}

// Status returns the status of the runtime
func (s *MinimalRuntimeService) Status(ctx context.Context, req *StatusRequest) (*StatusResponse, error) {
	s.logger.Info("CRI Status called")

	status := &RuntimeStatus{
		Conditions: []RuntimeCondition{
			{
				Type:    "RuntimeReady",
				Status:  true,
				Reason:  "RuntimeReady",
				Message: "Runtime is ready",
			},
			{
				Type:    "NetworkReady",
				Status:  true,
				Reason:  "NetworkReady",
				Message: "Network is ready",
			},
		},
	}

	response := &StatusResponse{
		Status: status,
	}

	if req.Verbose {
		response.Info = map[string]string{
			"runtime":     ServinRuntimeName,
			"version":     ServinRuntimeVersion,
			"api_version": CRIVersion,
			"base_dir":    s.criBaseDir,
		}
	}

	return response, nil
}

// RunPodSandbox creates and starts a pod-level sandbox
func (s *MinimalRuntimeService) RunPodSandbox(ctx context.Context, req *RunPodSandboxRequest) (*RunPodSandboxResponse, error) {
	s.logger.Info("CRI RunPodSandbox called for pod: %s", req.Config.Metadata.Name)

	// Generate pod sandbox ID
	podID := generatePodSandboxID(req.Config.Metadata)

	// Create pod sandbox directory
	podDir := filepath.Join(s.criBaseDir, "pods", podID)
	if err := os.MkdirAll(podDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create pod directory: %v", err)
	}

	// Create pod sandbox configuration
	podConfig := &PodSandbox{
		ID:          podID,
		Metadata:    req.Config.Metadata,
		State:       PodSandboxStateReady,
		CreatedAt:   time.Now().UnixNano(),
		Annotations: req.Config.Annotations,
		Labels:      req.Config.Labels,
		RuntimeInfo: &PodSandboxRuntimeInfo{
			RuntimeName:       ServinRuntimeName,
			RuntimeVersion:    ServinRuntimeVersion,
			RuntimeApiVersion: CRIVersion,
		},
	}

	// Save pod sandbox state
	if err := s.savePodSandboxState(podID, podConfig); err != nil {
		return nil, fmt.Errorf("failed to save pod sandbox state: %v", err)
	}

	s.logger.Info("Created pod sandbox: %s", podID)
	return &RunPodSandboxResponse{PodSandboxId: podID}, nil
}

// StopPodSandbox stops any running process that is part of the sandbox
func (s *MinimalRuntimeService) StopPodSandbox(ctx context.Context, req *StopPodSandboxRequest) (*StopPodSandboxResponse, error) {
	s.logger.Info("CRI StopPodSandbox called for pod: %s", req.PodSandboxId)

	// Update pod sandbox state
	if err := s.updatePodSandboxState(req.PodSandboxId, PodSandboxStateNotReady); err != nil {
		return nil, fmt.Errorf("failed to update pod sandbox state: %v", err)
	}

	return &StopPodSandboxResponse{}, nil
}

// RemovePodSandbox removes the sandbox
func (s *MinimalRuntimeService) RemovePodSandbox(ctx context.Context, req *RemovePodSandboxRequest) (*RemovePodSandboxResponse, error) {
	s.logger.Info("CRI RemovePodSandbox called for pod: %s", req.PodSandboxId)

	// Remove pod sandbox directory
	podDir := filepath.Join(s.criBaseDir, "pods", req.PodSandboxId)
	if err := os.RemoveAll(podDir); err != nil {
		return nil, fmt.Errorf("failed to remove pod directory: %v", err)
	}

	return &RemovePodSandboxResponse{}, nil
}

// PodSandboxStatus returns the status of the PodSandbox
func (s *MinimalRuntimeService) PodSandboxStatus(ctx context.Context, req *PodSandboxStatusRequest) (*PodSandboxStatusResponse, error) {
	s.logger.Info("CRI PodSandboxStatus called for pod: %s", req.PodSandboxId)

	podConfig, err := s.loadPodSandboxState(req.PodSandboxId)
	if err != nil {
		return nil, fmt.Errorf("failed to load pod sandbox state: %v", err)
	}

	status := &PodSandboxStatus{
		Id:          podConfig.ID,
		Metadata:    podConfig.Metadata,
		State:       podConfig.State,
		CreatedAt:   podConfig.CreatedAt,
		Labels:      podConfig.Labels,
		Annotations: podConfig.Annotations,
		RuntimeInfo: podConfig.RuntimeInfo,
		Network: &PodSandboxNetworkStatus{
			Ip: "127.0.0.1", // Default for now
		},
	}

	response := &PodSandboxStatusResponse{
		Status: status,
	}

	if req.Verbose {
		response.Info = map[string]string{
			"podDir":  filepath.Join(s.criBaseDir, "pods", req.PodSandboxId),
			"runtime": ServinRuntimeName,
		}
	}

	return response, nil
}

// ListPodSandbox returns a list of PodSandboxes
func (s *MinimalRuntimeService) ListPodSandbox(ctx context.Context, req *ListPodSandboxRequest) (*ListPodSandboxResponse, error) {
	s.logger.Info("CRI ListPodSandbox called")

	podsDir := filepath.Join(s.criBaseDir, "pods")
	if _, err := os.Stat(podsDir); os.IsNotExist(err) {
		return &ListPodSandboxResponse{Items: []*PodSandbox{}}, nil
	}

	entries, err := os.ReadDir(podsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read pods directory: %v", err)
	}

	var pods []*PodSandbox
	for _, entry := range entries {
		if entry.IsDir() {
			podConfig, err := s.loadPodSandboxState(entry.Name())
			if err != nil {
				s.logger.Info("Failed to load pod %s: %v", entry.Name(), err)
				continue
			}

			// Apply filter if provided
			if req.Filter != nil {
				if !s.matchesPodSandboxFilter(podConfig, req.Filter) {
					continue
				}
			}

			pods = append(pods, podConfig)
		}
	}

	return &ListPodSandboxResponse{Items: pods}, nil
}

// CreateContainer creates a new container in specified PodSandbox
func (s *MinimalRuntimeService) CreateContainer(ctx context.Context, req *CreateContainerRequest) (*CreateContainerResponse, error) {
	s.logger.Info("CRI CreateContainer called for container: %s", req.Config.Metadata.Name)

	// Generate container ID
	containerID := generateContainerID(req.Config.Metadata, req.PodSandboxId)

	// For minimal implementation, just return the ID
	// In a full implementation, this would create the actual container

	s.logger.Info("Created container: %s", containerID)
	return &CreateContainerResponse{ContainerId: containerID}, nil
}

// StartContainer starts the container
func (s *MinimalRuntimeService) StartContainer(ctx context.Context, req *StartContainerRequest) (*StartContainerResponse, error) {
	s.logger.Info("CRI StartContainer called for container: %s", req.ContainerId)

	// For minimal implementation, just log the action
	// In a full implementation, this would start the actual container

	return &StartContainerResponse{}, nil
}

// StopContainer stops a running container with a grace period
func (s *MinimalRuntimeService) StopContainer(ctx context.Context, req *StopContainerRequest) (*StopContainerResponse, error) {
	s.logger.Info("CRI StopContainer called for container: %s", req.ContainerId)

	// For minimal implementation, just log the action
	// In a full implementation, this would stop the actual container

	return &StopContainerResponse{}, nil
}

// RemoveContainer removes the container
func (s *MinimalRuntimeService) RemoveContainer(ctx context.Context, req *RemoveContainerRequest) (*RemoveContainerResponse, error) {
	s.logger.Info("CRI RemoveContainer called for container: %s", req.ContainerId)

	// For minimal implementation, just log the action
	// In a full implementation, this would remove the actual container

	return &RemoveContainerResponse{}, nil
}

// ListContainers lists all containers by filters
func (s *MinimalRuntimeService) ListContainers(ctx context.Context, req *ListContainersRequest) (*ListContainersResponse, error) {
	s.logger.Info("CRI ListContainers called")

	// For minimal implementation, return empty list
	// In a full implementation, this would list actual containers

	return &ListContainersResponse{Containers: []*Container{}}, nil
}

// ContainerStatus returns status of the container
func (s *MinimalRuntimeService) ContainerStatus(ctx context.Context, req *ContainerStatusRequest) (*ContainerStatusResponse, error) {
	s.logger.Info("CRI ContainerStatus called for container: %s", req.ContainerId)

	// For minimal implementation, return a basic status
	// In a full implementation, this would return actual container status

	status := &ContainerStatus{
		Id: req.ContainerId,
		Metadata: &ContainerMetadata{
			Name:    "test-container",
			Attempt: 0,
		},
		State:     ContainerStateRunning,
		CreatedAt: time.Now().UnixNano(),
		StartedAt: time.Now().UnixNano(),
		Image: &ImageSpec{
			Image: "alpine:latest",
		},
		ImageRef: "alpine:latest",
	}

	response := &ContainerStatusResponse{
		Status: status,
	}

	if req.Verbose {
		response.Info = map[string]string{
			"containerId": req.ContainerId,
			"runtime":     ServinRuntimeName,
		}
	}

	return response, nil
}

// ContainerStats returns stats of the container
func (s *MinimalRuntimeService) ContainerStats(ctx context.Context, req *ContainerStatsRequest) (*ContainerStatsResponse, error) {
	s.logger.Info("CRI ContainerStats called for container: %s", req.ContainerId)

	// Return basic stats
	stats := &ContainerStats{
		Attributes: &ContainerAttributes{
			ID: req.ContainerId,
		},
		Cpu: &CpuUsage{
			Timestamp: time.Now().UnixNano(),
		},
		Memory: &MemoryUsage{
			Timestamp: time.Now().UnixNano(),
		},
	}

	return &ContainerStatsResponse{Stats: stats}, nil
}

// ListContainerStats returns stats of all running containers
func (s *MinimalRuntimeService) ListContainerStats(ctx context.Context, req *ListContainerStatsRequest) (*ListContainerStatsResponse, error) {
	s.logger.Info("CRI ListContainerStats called")

	// Return empty stats list for minimal implementation
	return &ListContainerStatsResponse{Stats: []*ContainerStats{}}, nil
}

// ExecSync runs a command in a container synchronously
func (s *MinimalRuntimeService) ExecSync(ctx context.Context, req *ExecSyncRequest) (*ExecSyncResponse, error) {
	s.logger.Info("CRI ExecSync called for container: %s", req.ContainerId)

	return &ExecSyncResponse{
		Stdout:   []byte("exec not implemented\n"),
		Stderr:   []byte(""),
		ExitCode: 0,
	}, nil
}

// Exec prepares a streaming endpoint to execute a command in the container
func (s *MinimalRuntimeService) Exec(ctx context.Context, req *ExecRequest) (*ExecResponse, error) {
	s.logger.Info("CRI Exec called for container: %s", req.ContainerId)

	return &ExecResponse{
		Url: "http://localhost:8080/exec",
	}, nil
}

// Attach prepares a streaming endpoint to attach to a running container
func (s *MinimalRuntimeService) Attach(ctx context.Context, req *AttachRequest) (*AttachResponse, error) {
	s.logger.Info("CRI Attach called for container: %s", req.ContainerId)

	return &AttachResponse{
		Url: "http://localhost:8080/attach",
	}, nil
}

// PortForward prepares a streaming endpoint to forward ports from a PodSandbox
func (s *MinimalRuntimeService) PortForward(ctx context.Context, req *PortForwardRequest) (*PortForwardResponse, error) {
	s.logger.Info("CRI PortForward called for pod: %s", req.PodSandboxId)

	return &PortForwardResponse{
		Url: "http://localhost:8080/portforward",
	}, nil
}

// UpdateRuntimeConfig updates the runtime configuration based on the given request
func (s *MinimalRuntimeService) UpdateRuntimeConfig(ctx context.Context, req *UpdateRuntimeConfigRequest) (*UpdateRuntimeConfigResponse, error) {
	s.logger.Info("CRI UpdateRuntimeConfig called")

	return &UpdateRuntimeConfigResponse{}, nil
}

// Helper methods for minimal runtime service

// generatePodSandboxID generates a unique ID for a pod sandbox
func generatePodSandboxID(metadata *PodSandboxMetadata) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s-%s-%s-%d",
		metadata.Name, metadata.Namespace, metadata.UID, metadata.Attempt)))
	return fmt.Sprintf("pod-%x", hash[:8])
}

// generateContainerID generates a unique ID for a container
func generateContainerID(metadata *ContainerMetadata, podSandboxId string) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s-%s-%d-%d",
		metadata.Name, podSandboxId, metadata.Attempt, time.Now().UnixNano())))
	return fmt.Sprintf("ctr-%x", hash[:8])
}

// savePodSandboxState saves pod sandbox state to disk
func (s *MinimalRuntimeService) savePodSandboxState(podID string, podConfig *PodSandbox) error {
	podDir := filepath.Join(s.criBaseDir, "pods", podID)
	stateFile := filepath.Join(podDir, "state.json")

	data, err := json.MarshalIndent(podConfig, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(stateFile, data, 0644)
}

// loadPodSandboxState loads pod sandbox state from disk
func (s *MinimalRuntimeService) loadPodSandboxState(podID string) (*PodSandbox, error) {
	stateFile := filepath.Join(s.criBaseDir, "pods", podID, "state.json")

	data, err := os.ReadFile(stateFile)
	if err != nil {
		return nil, err
	}

	var podConfig PodSandbox
	err = json.Unmarshal(data, &podConfig)
	if err != nil {
		return nil, err
	}

	return &podConfig, nil
}

// updatePodSandboxState updates the state of a pod sandbox
func (s *MinimalRuntimeService) updatePodSandboxState(podID string, state PodSandboxState) error {
	podConfig, err := s.loadPodSandboxState(podID)
	if err != nil {
		return err
	}

	podConfig.State = state
	return s.savePodSandboxState(podID, podConfig)
}

// matchesPodSandboxFilter checks if a pod sandbox matches the given filter
func (s *MinimalRuntimeService) matchesPodSandboxFilter(pod *PodSandbox, filter *PodSandboxFilter) bool {
	if filter.Id != "" && filter.Id != pod.ID {
		return false
	}

	if filter.State != nil && filter.State.State != pod.State {
		return false
	}

	if filter.LabelSelector != nil {
		for key, value := range filter.LabelSelector {
			if pod.Labels[key] != value {
				return false
			}
		}
	}

	return true
}
