package cri

import (
	"context"
)

// Container Runtime Interface (CRI) types and interfaces
// Based on Kubernetes CRI specification v1alpha2/v1beta1

// RuntimeService defines the interface for container runtime operations
type RuntimeService interface {
	// Version returns the runtime name, runtime version, and runtime API version.
	Version(ctx context.Context, req *VersionRequest) (*VersionResponse, error)

	// RunPodSandbox creates and starts a pod-level sandbox.
	RunPodSandbox(ctx context.Context, req *RunPodSandboxRequest) (*RunPodSandboxResponse, error)

	// StopPodSandbox stops any running process that is part of the sandbox and reclaims network resources.
	StopPodSandbox(ctx context.Context, req *StopPodSandboxRequest) (*StopPodSandboxResponse, error)

	// RemovePodSandbox removes the sandbox.
	RemovePodSandbox(ctx context.Context, req *RemovePodSandboxRequest) (*RemovePodSandboxResponse, error)

	// PodSandboxStatus returns the status of the PodSandbox.
	PodSandboxStatus(ctx context.Context, req *PodSandboxStatusRequest) (*PodSandboxStatusResponse, error)

	// ListPodSandbox returns a list of PodSandboxes.
	ListPodSandbox(ctx context.Context, req *ListPodSandboxRequest) (*ListPodSandboxResponse, error)

	// CreateContainer creates a new container in specified PodSandbox.
	CreateContainer(ctx context.Context, req *CreateContainerRequest) (*CreateContainerResponse, error)

	// StartContainer starts the container.
	StartContainer(ctx context.Context, req *StartContainerRequest) (*StartContainerResponse, error)

	// StopContainer stops a running container with a grace period.
	StopContainer(ctx context.Context, req *StopContainerRequest) (*StopContainerResponse, error)

	// RemoveContainer removes the container.
	RemoveContainer(ctx context.Context, req *RemoveContainerRequest) (*RemoveContainerResponse, error)

	// ListContainers lists all containers by filters.
	ListContainers(ctx context.Context, req *ListContainersRequest) (*ListContainersResponse, error)

	// ContainerStatus returns status of the container.
	ContainerStatus(ctx context.Context, req *ContainerStatusRequest) (*ContainerStatusResponse, error)

	// ContainerStats returns stats of the container.
	ContainerStats(ctx context.Context, req *ContainerStatsRequest) (*ContainerStatsResponse, error)

	// ListContainerStats returns stats of all running containers.
	ListContainerStats(ctx context.Context, req *ListContainerStatsRequest) (*ListContainerStatsResponse, error)

	// ExecSync runs a command in a container synchronously.
	ExecSync(ctx context.Context, req *ExecSyncRequest) (*ExecSyncResponse, error)

	// Exec prepares a streaming endpoint to execute a command in the container.
	Exec(ctx context.Context, req *ExecRequest) (*ExecResponse, error)

	// Attach prepares a streaming endpoint to attach to a running container.
	Attach(ctx context.Context, req *AttachRequest) (*AttachResponse, error)

	// PortForward prepares a streaming endpoint to forward ports from a PodSandbox.
	PortForward(ctx context.Context, req *PortForwardRequest) (*PortForwardResponse, error)

	// UpdateRuntimeConfig updates the runtime configuration based on the given request.
	UpdateRuntimeConfig(ctx context.Context, req *UpdateRuntimeConfigRequest) (*UpdateRuntimeConfigResponse, error)

	// Status returns the status of the runtime.
	Status(ctx context.Context, req *StatusRequest) (*StatusResponse, error)
}

// ImageService defines the interface for image operations
type ImageService interface {
	// ListImages lists existing images.
	ListImages(ctx context.Context, req *ListImagesRequest) (*ListImagesResponse, error)

	// ImageStatus returns the status of the image.
	ImageStatus(ctx context.Context, req *ImageStatusRequest) (*ImageStatusResponse, error)

	// PullImage pulls an image with authentication config.
	PullImage(ctx context.Context, req *PullImageRequest) (*PullImageResponse, error)

	// RemoveImage removes the image.
	RemoveImage(ctx context.Context, req *RemoveImageRequest) (*RemoveImageResponse, error)

	// ImageFsInfo returns information of the filesystem that is used to store images.
	ImageFsInfo(ctx context.Context, req *ImageFsInfoRequest) (*ImageFsInfoResponse, error)
}

// Pod and Container States
type PodSandboxState int32

const (
	PodSandboxStateReady    PodSandboxState = 0
	PodSandboxStateNotReady PodSandboxState = 1
)

type ContainerState int32

const (
	ContainerStateCreated ContainerState = 0
	ContainerStateRunning ContainerState = 1
	ContainerStateExited  ContainerState = 2
	ContainerStateUnknown ContainerState = 3
)

// Runtime and Image information
type RuntimeStatus struct {
	Conditions []RuntimeCondition `json:"conditions,omitempty"`
}

type RuntimeCondition struct {
	Type    string `json:"type,omitempty"`
	Status  bool   `json:"status,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// PodSandbox represents a pod-level sandbox
type PodSandbox struct {
	ID          string                 `json:"id,omitempty"`
	Metadata    *PodSandboxMetadata    `json:"metadata,omitempty"`
	State       PodSandboxState        `json:"state,omitempty"`
	CreatedAt   int64                  `json:"created_at,omitempty"`
	Annotations map[string]string      `json:"annotations,omitempty"`
	Labels      map[string]string      `json:"labels,omitempty"`
	RuntimeInfo *PodSandboxRuntimeInfo `json:"runtime_info,omitempty"`
}

type PodSandboxMetadata struct {
	Name      string `json:"name,omitempty"`
	UID       string `json:"uid,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Attempt   uint32 `json:"attempt,omitempty"`
}

type PodSandboxRuntimeInfo struct {
	RuntimeName       string `json:"runtime_name,omitempty"`
	RuntimeVersion    string `json:"runtime_version,omitempty"`
	RuntimeApiVersion string `json:"runtime_api_version,omitempty"`
}

// Container represents a container within a pod
type Container struct {
	ID           string             `json:"id,omitempty"`
	PodSandboxID string             `json:"pod_sandbox_id,omitempty"`
	Metadata     *ContainerMetadata `json:"metadata,omitempty"`
	Image        *ImageSpec         `json:"image,omitempty"`
	ImageRef     string             `json:"image_ref,omitempty"`
	State        ContainerState     `json:"state,omitempty"`
	CreatedAt    int64              `json:"created_at,omitempty"`
	Labels       map[string]string  `json:"labels,omitempty"`
	Annotations  map[string]string  `json:"annotations,omitempty"`
}

type ContainerMetadata struct {
	Name    string `json:"name,omitempty"`
	Attempt uint32 `json:"attempt,omitempty"`
}

type ImageSpec struct {
	Image string `json:"image,omitempty"`
}

// Network and Volume configurations
type PodSandboxConfig struct {
	Metadata     *PodSandboxMetadata    `json:"metadata,omitempty"`
	Hostname     string                 `json:"hostname,omitempty"`
	LogDirectory string                 `json:"log_directory,omitempty"`
	DnsConfig    *DNSConfig             `json:"dns_config,omitempty"`
	PortMappings []*PortMapping         `json:"port_mappings,omitempty"`
	Labels       map[string]string      `json:"labels,omitempty"`
	Annotations  map[string]string      `json:"annotations,omitempty"`
	Linux        *LinuxPodSandboxConfig `json:"linux,omitempty"`
}

type DNSConfig struct {
	Servers  []string `json:"servers,omitempty"`
	Searches []string `json:"searches,omitempty"`
	Options  []string `json:"options,omitempty"`
}

type PortMapping struct {
	Protocol      Protocol `json:"protocol,omitempty"`
	ContainerPort int32    `json:"container_port,omitempty"`
	HostPort      int32    `json:"host_port,omitempty"`
	HostIP        string   `json:"host_ip,omitempty"`
}

type Protocol int32

const (
	ProtocolTCP  Protocol = 0
	ProtocolUDP  Protocol = 1
	ProtocolSCTP Protocol = 2
)

type LinuxPodSandboxConfig struct {
	CgroupParent    string                     `json:"cgroup_parent,omitempty"`
	SecurityContext *PodSandboxSecurityContext `json:"security_context,omitempty"`
	Sysctls         map[string]string          `json:"sysctls,omitempty"`
}

type PodSandboxSecurityContext struct {
	NamespaceOptions   *NamespaceOption `json:"namespace_options,omitempty"`
	SELinuxOptions     *SELinuxOption   `json:"selinux_options,omitempty"`
	RunAsUser          *Int64Value      `json:"run_as_user,omitempty"`
	RunAsGroup         *Int64Value      `json:"run_as_group,omitempty"`
	ReadonlyRootfs     bool             `json:"readonly_rootfs,omitempty"`
	SupplementalGroups []int64          `json:"supplemental_groups,omitempty"`
	Privileged         bool             `json:"privileged,omitempty"`
	SeccompProfilePath string           `json:"seccomp_profile_path,omitempty"`
}

type NamespaceOption struct {
	Network  NamespaceMode `json:"network,omitempty"`
	Pid      NamespaceMode `json:"pid,omitempty"`
	Ipc      NamespaceMode `json:"ipc,omitempty"`
	TargetId string        `json:"target_id,omitempty"`
}

type NamespaceMode int32

const (
	NamespaceModePod       NamespaceMode = 0
	NamespaceModeContainer NamespaceMode = 1
	NamespaceModeNode      NamespaceMode = 2
	NamespaceModeTarget    NamespaceMode = 3
)

type SELinuxOption struct {
	User  string `json:"user,omitempty"`
	Role  string `json:"role,omitempty"`
	Type  string `json:"type,omitempty"`
	Level string `json:"level,omitempty"`
}

type Int64Value struct {
	Value int64 `json:"value,omitempty"`
}

// Container configuration
type ContainerConfig struct {
	Metadata    *ContainerMetadata    `json:"metadata,omitempty"`
	Image       *ImageSpec            `json:"image,omitempty"`
	Command     []string              `json:"command,omitempty"`
	Args        []string              `json:"args,omitempty"`
	WorkingDir  string                `json:"working_dir,omitempty"`
	Envs        []*KeyValue           `json:"envs,omitempty"`
	Mounts      []*Mount              `json:"mounts,omitempty"`
	Devices     []*Device             `json:"devices,omitempty"`
	Labels      map[string]string     `json:"labels,omitempty"`
	Annotations map[string]string     `json:"annotations,omitempty"`
	LogPath     string                `json:"log_path,omitempty"`
	Stdin       bool                  `json:"stdin,omitempty"`
	StdinOnce   bool                  `json:"stdin_once,omitempty"`
	Tty         bool                  `json:"tty,omitempty"`
	Linux       *LinuxContainerConfig `json:"linux,omitempty"`
}

type KeyValue struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type Mount struct {
	ContainerPath  string           `json:"container_path,omitempty"`
	HostPath       string           `json:"host_path,omitempty"`
	Readonly       bool             `json:"readonly,omitempty"`
	SelinuxRelabel bool             `json:"selinux_relabel,omitempty"`
	Propagation    MountPropagation `json:"propagation,omitempty"`
}

type MountPropagation int32

const (
	MountPropagationPrivate         MountPropagation = 0
	MountPropagationHostToContainer MountPropagation = 1
	MountPropagationBidirectional   MountPropagation = 2
)

type Device struct {
	ContainerPath string `json:"container_path,omitempty"`
	HostPath      string `json:"host_path,omitempty"`
	Permissions   string `json:"permissions,omitempty"`
}

type LinuxContainerConfig struct {
	Resources       *LinuxContainerResources       `json:"resources,omitempty"`
	SecurityContext *LinuxContainerSecurityContext `json:"security_context,omitempty"`
}

type LinuxContainerResources struct {
	CpuPeriod          int64             `json:"cpu_period,omitempty"`
	CpuQuota           int64             `json:"cpu_quota,omitempty"`
	CpuShares          int64             `json:"cpu_shares,omitempty"`
	MemoryLimitInBytes int64             `json:"memory_limit_in_bytes,omitempty"`
	OomScoreAdj        int64             `json:"oom_score_adj,omitempty"`
	CpusetCpus         string            `json:"cpuset_cpus,omitempty"`
	CpusetMems         string            `json:"cpuset_mems,omitempty"`
	HugepageLimits     []*HugepageLimit  `json:"hugepage_limits,omitempty"`
	Unified            map[string]string `json:"unified,omitempty"`
}

type HugepageLimit struct {
	PageSize string `json:"page_size,omitempty"`
	Limit    uint64 `json:"limit,omitempty"`
}

type LinuxContainerSecurityContext struct {
	Capabilities       *Capability      `json:"capabilities,omitempty"`
	Privileged         bool             `json:"privileged,omitempty"`
	NamespaceOptions   *NamespaceOption `json:"namespace_options,omitempty"`
	SELinuxOptions     *SELinuxOption   `json:"selinux_options,omitempty"`
	RunAsUser          *Int64Value      `json:"run_as_user,omitempty"`
	RunAsGroup         *Int64Value      `json:"run_as_group,omitempty"`
	RunAsUsername      string           `json:"run_as_username,omitempty"`
	ReadonlyRootfs     bool             `json:"readonly_rootfs,omitempty"`
	SupplementalGroups []int64          `json:"supplemental_groups,omitempty"`
	ApparmorProfile    string           `json:"apparmor_profile,omitempty"`
	SeccompProfilePath string           `json:"seccomp_profile_path,omitempty"`
	NoNewPrivs         bool             `json:"no_new_privs,omitempty"`
	MaskedPaths        []string         `json:"masked_paths,omitempty"`
	ReadonlyPaths      []string         `json:"readonly_paths,omitempty"`
}

type Capability struct {
	AddCapabilities  []string `json:"add_capabilities,omitempty"`
	DropCapabilities []string `json:"drop_capabilities,omitempty"`
}

// Statistics and monitoring
type ContainerStats struct {
	Attributes    *ContainerAttributes `json:"attributes,omitempty"`
	Cpu           *CpuUsage            `json:"cpu,omitempty"`
	Memory        *MemoryUsage         `json:"memory,omitempty"`
	WritableLayer *FilesystemUsage     `json:"writable_layer,omitempty"`
}

type ContainerAttributes struct {
	ID          string             `json:"id,omitempty"`
	Metadata    *ContainerMetadata `json:"metadata,omitempty"`
	Labels      map[string]string  `json:"labels,omitempty"`
	Annotations map[string]string  `json:"annotations,omitempty"`
}

type CpuUsage struct {
	Timestamp            int64        `json:"timestamp,omitempty"`
	UsageCoreNanoSeconds *UInt64Value `json:"usage_core_nano_seconds,omitempty"`
	UsageNanoCores       *UInt64Value `json:"usage_nano_cores,omitempty"`
}

type UInt64Value struct {
	Value uint64 `json:"value,omitempty"`
}

type MemoryUsage struct {
	Timestamp       int64        `json:"timestamp,omitempty"`
	WorkingSetBytes *UInt64Value `json:"working_set_bytes,omitempty"`
	AvailableBytes  *UInt64Value `json:"available_bytes,omitempty"`
	UsageBytes      *UInt64Value `json:"usage_bytes,omitempty"`
	RssBytes        *UInt64Value `json:"rss_bytes,omitempty"`
	PageFaults      *UInt64Value `json:"page_faults,omitempty"`
	MajorPageFaults *UInt64Value `json:"major_page_faults,omitempty"`
}

type FilesystemUsage struct {
	Timestamp  int64                 `json:"timestamp,omitempty"`
	FsId       *FilesystemIdentifier `json:"fs_id,omitempty"`
	UsedBytes  *UInt64Value          `json:"used_bytes,omitempty"`
	InodesUsed *UInt64Value          `json:"inodes_used,omitempty"`
}

type FilesystemIdentifier struct {
	Mountpoint string `json:"mountpoint,omitempty"`
}

// Image types
type Image struct {
	ID          string      `json:"id,omitempty"`
	RepoTags    []string    `json:"repo_tags,omitempty"`
	RepoDigests []string    `json:"repo_digests,omitempty"`
	Size        uint64      `json:"size_,omitempty"`
	UID         *Int64Value `json:"uid,omitempty"`
	Username    string      `json:"username,omitempty"`
	Spec        *ImageSpec  `json:"spec,omitempty"`
	Pinned      bool        `json:"pinned,omitempty"`
}

type AuthConfig struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Auth          string `json:"auth,omitempty"`
	ServerAddress string `json:"server_address,omitempty"`
	IdentityToken string `json:"identity_token,omitempty"`
	RegistryToken string `json:"registry_token,omitempty"`
}

// Server configuration
type CRIServer struct {
	RuntimeService RuntimeService
	ImageService   ImageService
	SocketPath     string
}

// Version information
const (
	ServinRuntimeName    = "servin"
	ServinRuntimeVersion = "0.1.0"
	CRIVersion           = "v1alpha2"
)
