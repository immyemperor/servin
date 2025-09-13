package cri

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"servin/pkg/image"
	"servin/pkg/logger"
)

// ServinImageService implements the CRI ImageService interface
type ServinImageService struct {
	imageManager *image.Manager
	logger       *logger.Logger
	criBaseDir   string
}

// NewServinImageService creates a new Servin CRI image service
func NewServinImageService(imageManager *image.Manager, logger *logger.Logger, baseDir string) *ServinImageService {
	criBaseDir := filepath.Join(baseDir, "cri")
	os.MkdirAll(criBaseDir, 0755)

	return &ServinImageService{
		imageManager: imageManager,
		logger:       logger,
		criBaseDir:   criBaseDir,
	}
}

// ListImages lists existing images
func (s *ServinImageService) ListImages(ctx context.Context, req *ListImagesRequest) (*ListImagesResponse, error) {
	s.logger.Info("CRI ListImages called")

	// Get images from Servin image manager
	images, err := s.imageManager.ListImages()
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %v", err)
	}

	var criImages []*Image
	for _, img := range images {
		criImage := s.convertServinImageToCRI(img)

		// Apply filter if provided
		if req.Filter != nil {
			if !s.matchesImageFilter(criImage, req.Filter) {
				continue
			}
		}

		criImages = append(criImages, criImage)
	}

	return &ListImagesResponse{Images: criImages}, nil
}

// ImageStatus returns the status of the image
func (s *ServinImageService) ImageStatus(ctx context.Context, req *ImageStatusRequest) (*ImageStatusResponse, error) {
	s.logger.Info("CRI ImageStatus called for image: %s", req.Image.Image)

	// Get image from Servin image manager
	img, err := s.imageManager.GetImage(req.Image.Image)
	if err != nil {
		return &ImageStatusResponse{Image: nil}, nil // Image not found
	}

	criImage := s.convertServinImageToCRI(img)

	response := &ImageStatusResponse{
		Image: criImage,
	}

	if req.Verbose {
		response.Info = map[string]string{
			"imageId":  img.ID,
			"size":     fmt.Sprintf("%d", img.Size),
			"created":  img.Created.Format(time.RFC3339),
			"repoTags": strings.Join(img.RepoTags, ","),
			"rootfs":   img.RootFSPath,
		}
	}

	return response, nil
}

// PullImage pulls an image with authentication config
func (s *ServinImageService) PullImage(ctx context.Context, req *PullImageRequest) (*PullImageResponse, error) {
	s.logger.Info("CRI PullImage called for image: %s", req.Image.Image)

	// For now, we'll simulate image pulling by checking if the image exists locally
	// In a real implementation, this would pull from a registry
	imageName := req.Image.Image

	// Check if image already exists
	_, err := s.imageManager.GetImage(imageName)
	if err == nil {
		s.logger.Info("Image %s already exists locally", imageName)
		return &PullImageResponse{ImageRef: imageName}, nil
	}

	// Simulate pulling - in reality this would download from a registry
	s.logger.Info("Simulating pull for image: %s", imageName)

	// For demonstration, we'll create a placeholder image entry
	// In a real implementation, this would involve:
	// 1. Authenticating with registry using req.Auth
	// 2. Downloading image layers
	// 3. Extracting and storing the image

	imageRef := fmt.Sprintf("%s@sha256:placeholder", imageName)

	return &PullImageResponse{ImageRef: imageRef}, nil
}

// RemoveImage removes the image
func (s *ServinImageService) RemoveImage(ctx context.Context, req *RemoveImageRequest) (*RemoveImageResponse, error) {
	s.logger.Info("CRI RemoveImage called for image: %s", req.Image.Image)

	err := s.imageManager.RemoveImage(req.Image.Image)
	if err != nil {
		return nil, fmt.Errorf("failed to remove image: %v", err)
	}

	return &RemoveImageResponse{}, nil
}

// ImageFsInfo returns information of the filesystem that is used to store images
func (s *ServinImageService) ImageFsInfo(ctx context.Context, req *ImageFsInfoRequest) (*ImageFsInfoResponse, error) {
	s.logger.Info("CRI ImageFsInfo called")

	// Get filesystem information for the image storage directory
	imageDir := s.imageManager.GetImageDir()

	// Get directory size and usage
	usage, err := s.getDirUsage(imageDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory usage: %v", err)
	}

	filesystemUsage := &FilesystemUsage{
		Timestamp: time.Now().UnixNano(),
		FsId: &FilesystemIdentifier{
			Mountpoint: imageDir,
		},
		UsedBytes: &UInt64Value{
			Value: usage,
		},
		InodesUsed: &UInt64Value{
			Value: 0, // TODO: Calculate actual inode usage
		},
	}

	return &ImageFsInfoResponse{
		ImageFilesystems: []*FilesystemUsage{filesystemUsage},
	}, nil
}

// Helper methods

// convertServinImageToCRI converts a Servin image to CRI image format
func (s *ServinImageService) convertServinImageToCRI(img *image.Image) *Image {
	// Use repository tags from the image
	var repoTags []string
	if len(img.RepoTags) > 0 {
		repoTags = img.RepoTags
	} else {
		// Fallback to ID if no repo tags
		repoTags = []string{img.ID}
	}

	// Generate repository digests (for now, use ID as digest)
	var repoDigests []string
	repoDigests = append(repoDigests, fmt.Sprintf("sha256:%s", img.ID))

	// Use first repo tag as image spec
	imageSpec := ""
	if len(repoTags) > 0 {
		imageSpec = repoTags[0]
	}

	return &Image{
		ID:          img.ID,
		RepoTags:    repoTags,
		RepoDigests: repoDigests,
		Size:        uint64(img.Size),
		Spec: &ImageSpec{
			Image: imageSpec,
		},
		Pinned: false, // Default to not pinned
	}
}

// matchesImageFilter checks if an image matches the given filter
func (s *ServinImageService) matchesImageFilter(image *Image, filter *ImageFilter) bool {
	if filter.Image != nil {
		// Check if any of the repo tags match the filter image
		imageSpec := filter.Image.Image
		found := false
		for _, tag := range image.RepoTags {
			if tag == imageSpec || strings.Contains(tag, imageSpec) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// getDirUsage calculates the total size of a directory
func (s *ServinImageService) getDirUsage(dirPath string) (uint64, error) {
	var size uint64

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += uint64(info.Size())
		}
		return nil
	})

	return size, err
}
