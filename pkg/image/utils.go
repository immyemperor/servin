package image

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// generateImageID creates a unique ID for an image
func generateImageID(name, tag string) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%d", name, tag, os.Getpid())))
	return fmt.Sprintf("%x", hash)[:16]
}

// extractTarball extracts a tarball to the specified directory
func extractTarball(tarballPath, destDir string) error {
	file, err := os.Open(tarballPath)
	if err != nil {
		return fmt.Errorf("failed to open tarball: %v", err)
	}
	defer file.Close()

	var reader io.Reader = file

	// Check if it's gzipped
	if strings.HasSuffix(tarballPath, ".gz") || strings.HasSuffix(tarballPath, ".tgz") {
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %v", err)
		}
		defer gzReader.Close()
		reader = gzReader
	}

	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %v", err)
		}

		targetPath := filepath.Join(destDir, header.Name)

		// Security check: prevent path traversal
		if !strings.HasPrefix(targetPath, destDir) {
			return fmt.Errorf("invalid file path in tarball: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", targetPath, err)
			}

		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %v", targetPath, err)
			}

			file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %v", targetPath, err)
			}

			if _, err := io.Copy(file, tarReader); err != nil {
				file.Close()
				return fmt.Errorf("failed to extract file %s: %v", targetPath, err)
			}
			file.Close()

		case tar.TypeSymlink:
			if err := os.Symlink(header.Linkname, targetPath); err != nil {
				// Don't fail on symlink errors, just warn
				fmt.Printf("Warning: failed to create symlink %s -> %s: %v\n", targetPath, header.Linkname, err)
			}

		default:
			fmt.Printf("Warning: unsupported file type %c for %s\n", header.Typeflag, header.Name)
		}
	}

	return nil
}

// createTarball creates a tarball from a directory
func createTarball(sourceDir, tarballPath string) error {
	file, err := os.Create(tarballPath)
	if err != nil {
		return fmt.Errorf("failed to create tarball file: %v", err)
	}
	defer file.Close()

	var writer io.Writer = file

	// Create gzip writer if .gz extension
	if strings.HasSuffix(tarballPath, ".gz") || strings.HasSuffix(tarballPath, ".tgz") {
		gzWriter := gzip.NewWriter(file)
		defer gzWriter.Close()
		writer = gzWriter
	}

	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return fmt.Errorf("failed to create tar header for %s: %v", path, err)
		}

		// Update header name to be relative to source directory
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %v", path, err)
		}
		header.Name = relPath

		// Write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write tar header for %s: %v", path, err)
		}

		// Write file content if it's a regular file
		if info.Mode().IsRegular() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %v", path, err)
			}
			defer file.Close()

			if _, err := io.Copy(tarWriter, file); err != nil {
				return fmt.Errorf("failed to write file %s to tar: %v", path, err)
			}
		}

		return nil
	})
}

// downloadImage downloads an image from a registry (placeholder)
func downloadImage(imageRef string) error {
	// This is a placeholder for future registry support
	// For now, we'll return an error suggesting to use tarballs
	return fmt.Errorf("registry pulling not yet implemented. Please use 'servin image import' with a tarball")
}

// parseImageReference parses an image reference into name and tag
func parseImageReference(ref string) (string, string) {
	parts := strings.Split(ref, ":")
	if len(parts) == 1 {
		return parts[0], "latest"
	}
	return strings.Join(parts[:len(parts)-1], ":"), parts[len(parts)-1]
}

// formatSize formats bytes into human-readable size
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
