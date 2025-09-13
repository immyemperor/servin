package cri

// CRI Request and Response message types
// These correspond to the protobuf messages defined in the CRI specification

// Version requests and responses
type VersionRequest struct {
	Version string `json:"version,omitempty"`
}

type VersionResponse struct {
	Version           string `json:"version,omitempty"`
	RuntimeName       string `json:"runtime_name,omitempty"`
	RuntimeVersion    string `json:"runtime_version,omitempty"`
	RuntimeApiVersion string `json:"runtime_api_version,omitempty"`
}

// PodSandbox requests and responses
type RunPodSandboxRequest struct {
	Config         *PodSandboxConfig `json:"config,omitempty"`
	RuntimeHandler string            `json:"runtime_handler,omitempty"`
}

type RunPodSandboxResponse struct {
	PodSandboxId string `json:"pod_sandbox_id,omitempty"`
}

type StopPodSandboxRequest struct {
	PodSandboxId string `json:"pod_sandbox_id,omitempty"`
}

type StopPodSandboxResponse struct{}

type RemovePodSandboxRequest struct {
	PodSandboxId string `json:"pod_sandbox_id,omitempty"`
}

type RemovePodSandboxResponse struct{}

type PodSandboxStatusRequest struct {
	PodSandboxId string `json:"pod_sandbox_id,omitempty"`
	Verbose      bool   `json:"verbose,omitempty"`
}

type PodSandboxStatusResponse struct {
	Status *PodSandboxStatus `json:"status,omitempty"`
	Info   map[string]string `json:"info,omitempty"`
}

type PodSandboxStatus struct {
	Id          string                   `json:"id,omitempty"`
	Metadata    *PodSandboxMetadata      `json:"metadata,omitempty"`
	State       PodSandboxState          `json:"state,omitempty"`
	CreatedAt   int64                    `json:"created_at,omitempty"`
	Network     *PodSandboxNetworkStatus `json:"network,omitempty"`
	Linux       *LinuxPodSandboxStatus   `json:"linux,omitempty"`
	Labels      map[string]string        `json:"labels,omitempty"`
	Annotations map[string]string        `json:"annotations,omitempty"`
	RuntimeInfo *PodSandboxRuntimeInfo   `json:"runtime_info,omitempty"`
}

type PodSandboxNetworkStatus struct {
	Ip            string   `json:"ip,omitempty"`
	AdditionalIps []*PodIP `json:"additional_ips,omitempty"`
}

type PodIP struct {
	Ip string `json:"ip,omitempty"`
}

type LinuxPodSandboxStatus struct {
	Namespaces *Namespace `json:"namespaces,omitempty"`
}

type Namespace struct {
	Options *NamespaceOption `json:"options,omitempty"`
}

type PodSandboxFilter struct {
	Id            string                `json:"id,omitempty"`
	State         *PodSandboxStateValue `json:"state,omitempty"`
	LabelSelector map[string]string     `json:"label_selector,omitempty"`
}

type PodSandboxStateValue struct {
	State PodSandboxState `json:"state,omitempty"`
}

type ListPodSandboxRequest struct {
	Filter *PodSandboxFilter `json:"filter,omitempty"`
}

type ListPodSandboxResponse struct {
	Items []*PodSandbox `json:"items,omitempty"`
}

// Container requests and responses
type CreateContainerRequest struct {
	PodSandboxId  string            `json:"pod_sandbox_id,omitempty"`
	Config        *ContainerConfig  `json:"config,omitempty"`
	SandboxConfig *PodSandboxConfig `json:"sandbox_config,omitempty"`
}

type CreateContainerResponse struct {
	ContainerId string `json:"container_id,omitempty"`
}

type StartContainerRequest struct {
	ContainerId string `json:"container_id,omitempty"`
}

type StartContainerResponse struct{}

type StopContainerRequest struct {
	ContainerId string `json:"container_id,omitempty"`
	Timeout     int64  `json:"timeout,omitempty"`
}

type StopContainerResponse struct{}

type RemoveContainerRequest struct {
	ContainerId string `json:"container_id,omitempty"`
}

type RemoveContainerResponse struct{}

type ContainerStatusRequest struct {
	ContainerId string `json:"container_id,omitempty"`
	Verbose     bool   `json:"verbose,omitempty"`
}

type ContainerStatusResponse struct {
	Status *ContainerStatus  `json:"status,omitempty"`
	Info   map[string]string `json:"info,omitempty"`
}

type ContainerStatus struct {
	Id          string              `json:"id,omitempty"`
	Metadata    *ContainerMetadata  `json:"metadata,omitempty"`
	State       ContainerState      `json:"state,omitempty"`
	CreatedAt   int64               `json:"created_at,omitempty"`
	StartedAt   int64               `json:"started_at,omitempty"`
	FinishedAt  int64               `json:"finished_at,omitempty"`
	ExitCode    int32               `json:"exit_code,omitempty"`
	Image       *ImageSpec          `json:"image,omitempty"`
	ImageRef    string              `json:"image_ref,omitempty"`
	Reason      string              `json:"reason,omitempty"`
	Message     string              `json:"message,omitempty"`
	Labels      map[string]string   `json:"labels,omitempty"`
	Annotations map[string]string   `json:"annotations,omitempty"`
	Mounts      []*Mount            `json:"mounts,omitempty"`
	LogPath     string              `json:"log_path,omitempty"`
	Resources   *ContainerResources `json:"resources,omitempty"`
}

type ContainerResources struct {
	Linux *LinuxContainerResources `json:"linux,omitempty"`
}

type ContainerFilter struct {
	Id            string               `json:"id,omitempty"`
	State         *ContainerStateValue `json:"state,omitempty"`
	PodSandboxId  string               `json:"pod_sandbox_id,omitempty"`
	LabelSelector map[string]string    `json:"label_selector,omitempty"`
}

type ContainerStateValue struct {
	State ContainerState `json:"state,omitempty"`
}

type ListContainersRequest struct {
	Filter *ContainerFilter `json:"filter,omitempty"`
}

type ListContainersResponse struct {
	Containers []*Container `json:"containers,omitempty"`
}

// Container stats requests and responses
type ContainerStatsRequest struct {
	ContainerId string `json:"container_id,omitempty"`
}

type ContainerStatsResponse struct {
	Stats *ContainerStats `json:"stats,omitempty"`
}

type ContainerStatsFilter struct {
	Id            string            `json:"id,omitempty"`
	PodSandboxId  string            `json:"pod_sandbox_id,omitempty"`
	LabelSelector map[string]string `json:"label_selector,omitempty"`
}

type ListContainerStatsRequest struct {
	Filter *ContainerStatsFilter `json:"filter,omitempty"`
}

type ListContainerStatsResponse struct {
	Stats []*ContainerStats `json:"stats,omitempty"`
}

// Exec and Attach requests and responses
type ExecSyncRequest struct {
	ContainerId string   `json:"container_id,omitempty"`
	Cmd         []string `json:"cmd,omitempty"`
	Timeout     int64    `json:"timeout,omitempty"`
}

type ExecSyncResponse struct {
	Stdout   []byte `json:"stdout,omitempty"`
	Stderr   []byte `json:"stderr,omitempty"`
	ExitCode int32  `json:"exit_code,omitempty"`
}

type ExecRequest struct {
	ContainerId string   `json:"container_id,omitempty"`
	Cmd         []string `json:"cmd,omitempty"`
	Tty         bool     `json:"tty,omitempty"`
	Stdin       bool     `json:"stdin,omitempty"`
	Stdout      bool     `json:"stdout,omitempty"`
	Stderr      bool     `json:"stderr,omitempty"`
}

type ExecResponse struct {
	Url string `json:"url,omitempty"`
}

type AttachRequest struct {
	ContainerId string `json:"container_id,omitempty"`
	Stdin       bool   `json:"stdin,omitempty"`
	Tty         bool   `json:"tty,omitempty"`
	Stdout      bool   `json:"stdout,omitempty"`
	Stderr      bool   `json:"stderr,omitempty"`
}

type AttachResponse struct {
	Url string `json:"url,omitempty"`
}

type PortForwardRequest struct {
	PodSandboxId string  `json:"pod_sandbox_id,omitempty"`
	Port         []int32 `json:"port,omitempty"`
}

type PortForwardResponse struct {
	Url string `json:"url,omitempty"`
}

// Runtime config and status
type UpdateRuntimeConfigRequest struct {
	RuntimeConfig *RuntimeConfig `json:"runtime_config,omitempty"`
}

type UpdateRuntimeConfigResponse struct{}

type RuntimeConfig struct {
	NetworkConfig *NetworkConfig `json:"network_config,omitempty"`
}

type NetworkConfig struct {
	PodCidr string `json:"pod_cidr,omitempty"`
}

type StatusRequest struct {
	Verbose bool `json:"verbose,omitempty"`
}

type StatusResponse struct {
	Status *RuntimeStatus    `json:"status,omitempty"`
	Info   map[string]string `json:"info,omitempty"`
}

// Image service requests and responses
type ImageFilter struct {
	Image *ImageSpec `json:"image,omitempty"`
}

type ListImagesRequest struct {
	Filter *ImageFilter `json:"filter,omitempty"`
}

type ListImagesResponse struct {
	Images []*Image `json:"images,omitempty"`
}

type ImageStatusRequest struct {
	Image   *ImageSpec `json:"image,omitempty"`
	Verbose bool       `json:"verbose,omitempty"`
}

type ImageStatusResponse struct {
	Image *Image            `json:"image,omitempty"`
	Info  map[string]string `json:"info,omitempty"`
}

type PullImageRequest struct {
	Image         *ImageSpec        `json:"image,omitempty"`
	Auth          *AuthConfig       `json:"auth,omitempty"`
	SandboxConfig *PodSandboxConfig `json:"sandbox_config,omitempty"`
}

type PullImageResponse struct {
	ImageRef string `json:"image_ref,omitempty"`
}

type RemoveImageRequest struct {
	Image *ImageSpec `json:"image,omitempty"`
}

type RemoveImageResponse struct{}

type ImageFsInfoRequest struct{}

type ImageFsInfoResponse struct {
	ImageFilesystems []*FilesystemUsage `json:"image_filesystems,omitempty"`
}
