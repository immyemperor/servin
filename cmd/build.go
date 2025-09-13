package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"servin/pkg/errors"
	"servin/pkg/image"
	"servin/pkg/logger"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build [OPTIONS] PATH",
	Short: "Build an image from a Buildfile",
	Long: `Build a container image from a Buildfile (similar to Dockerfile).
The build context is the PATH where the Buildfile is located.

Example Buildfile:
  FROM alpine:latest
  RUN apk add --no-cache curl
  COPY . /app
  WORKDIR /app
  CMD ["./app"]

Examples:
  servin build .
  servin build -t myapp:v1.0 .
  servin build -f MyBuildfile .`,
	Args: cobra.ExactArgs(1),
	RunE: runBuild,
}

var (
	buildTag     string
	buildFile    string
	buildNoCache bool
	buildQuiet   bool
	buildArgs    []string
	buildLabels  []string
)

func init() {
	rootCmd.AddCommand(buildCmd)

	// Add flags for build options
	buildCmd.Flags().StringVarP(&buildTag, "tag", "t", "", "Name and optionally a tag in the 'name:tag' format")
	buildCmd.Flags().StringVarP(&buildFile, "file", "f", "Buildfile", "Name of the Buildfile (default 'Buildfile')")
	buildCmd.Flags().BoolVar(&buildNoCache, "no-cache", false, "Do not use cache when building the image")
	buildCmd.Flags().BoolVarP(&buildQuiet, "quiet", "q", false, "Suppress the build output and print image ID on success")
	buildCmd.Flags().StringArrayVar(&buildArgs, "build-arg", []string{}, "Set build-time variables")
	buildCmd.Flags().StringArrayVar(&buildLabels, "label", []string{}, "Set metadata for an image")
}

func runBuild(cmd *cobra.Command, args []string) error {
	buildContext := args[0]

	logger.Debug("Building image from context: %s", buildContext)

	// Resolve build context path
	buildContextPath, err := filepath.Abs(buildContext)
	if err != nil {
		logger.Error("Failed to resolve build context path: %v", err)
		return errors.NewValidationError("build", fmt.Sprintf("invalid build context path: %v", err))
	}

	// Check if build context exists
	if _, err := os.Stat(buildContextPath); os.IsNotExist(err) {
		logger.Error("Build context does not exist: %s", buildContextPath)
		return errors.NewNotFoundError("build", fmt.Sprintf("build context '%s' not found", buildContextPath))
	}

	// Resolve Buildfile path
	buildfilePath := filepath.Join(buildContextPath, buildFile)
	if _, err := os.Stat(buildfilePath); os.IsNotExist(err) {
		logger.Error("Buildfile does not exist: %s", buildfilePath)
		return errors.NewNotFoundError("build", fmt.Sprintf("Buildfile '%s' not found", buildfilePath))
	}

	logger.Debug("Using Buildfile: %s", buildfilePath)

	// Parse build arguments
	buildArgMap := make(map[string]string)
	for _, arg := range buildArgs {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			buildArgMap[parts[0]] = parts[1]
		} else {
			buildArgMap[parts[0]] = ""
		}
	}

	// Parse labels
	labelMap := make(map[string]string)
	for _, label := range buildLabels {
		parts := strings.SplitN(label, "=", 2)
		if len(parts) == 2 {
			labelMap[parts[0]] = parts[1]
		}
	}

	// Create build configuration
	buildConfig := &BuildConfig{
		ContextPath: buildContextPath,
		Buildfile:   buildfilePath,
		Tag:         buildTag,
		NoCache:     buildNoCache,
		Quiet:       buildQuiet,
		BuildArgs:   buildArgMap,
		Labels:      labelMap,
	}

	// Execute the build
	builder := NewImageBuilder()
	imageID, err := builder.Build(buildConfig)
	if err != nil {
		logger.Error("Build failed: %v", err)
		return errors.NewImageError("build", fmt.Sprintf("image build failed: %v", err))
	}

	if buildQuiet {
		fmt.Println(imageID)
	} else {
		fmt.Printf("Successfully built image: %s\n", imageID)
		if buildTag != "" {
			fmt.Printf("Successfully tagged: %s\n", buildTag)
		}
	}

	return nil
}

// BuildConfig represents build configuration
type BuildConfig struct {
	ContextPath string
	Buildfile   string
	Tag         string
	NoCache     bool
	Quiet       bool
	BuildArgs   map[string]string
	Labels      map[string]string
}

// BuildStep represents a single step in the Buildfile
type BuildStep struct {
	Instruction string
	Arguments   []string
	RawLine     string
}

// ImageBuilder handles the image building process
type ImageBuilder struct {
	imgManager *image.Manager
}

// NewImageBuilder creates a new image builder
func NewImageBuilder() *ImageBuilder {
	return &ImageBuilder{
		imgManager: image.NewManager(),
	}
}

// Build executes the image build process
func (b *ImageBuilder) Build(config *BuildConfig) (string, error) {
	logger.Info("Starting image build")
	logger.Debug("Build context: %s", config.ContextPath)
	logger.Debug("Buildfile: %s", config.Buildfile)

	// Parse the Buildfile
	steps, err := b.parseBuildfile(config.Buildfile, config.BuildArgs)
	if err != nil {
		return "", fmt.Errorf("failed to parse Buildfile: %v", err)
	}

	logger.Debug("Parsed %d build steps", len(steps))

	// Create a new image
	img := &image.Image{
		ID:         generateImageID(),
		Created:    time.Now(),
		Size:       0,
		Layers:     []string{},
		RootFSType: "layer",
		Config: image.ImageConfig{
			Env:          []string{},
			Cmd:          []string{},
			Entrypoint:   []string{},
			WorkingDir:   "/",
			User:         "root",
			ExposedPorts: make(map[string]struct{}),
			Labels:       config.Labels,
		},
		Metadata: make(map[string]string),
	}

	// Add build metadata
	img.Metadata["build.context"] = config.ContextPath
	img.Metadata["build.buildfile"] = config.Buildfile
	img.Metadata["build.timestamp"] = time.Now().Format(time.RFC3339)

	// Process each step
	var fromProcessed bool
	for i, step := range steps {
		if !config.Quiet {
			fmt.Printf("Step %d/%d : %s\n", i+1, len(steps), step.RawLine)
		}

		logger.Debug("Executing step %d: %s %v", i+1, step.Instruction, step.Arguments)

		switch strings.ToUpper(step.Instruction) {
		case "FROM":
			_, err = b.processFrom(step, img)
			fromProcessed = true
		case "RUN":
			err = b.processRun(step, img, config.ContextPath)
		case "COPY":
			err = b.processCopy(step, img, config.ContextPath)
		case "ADD":
			err = b.processAdd(step, img, config.ContextPath)
		case "WORKDIR":
			err = b.processWorkdir(step, img)
		case "ENV":
			err = b.processEnv(step, img)
		case "EXPOSE":
			err = b.processExpose(step, img)
		case "CMD":
			err = b.processCmd(step, img)
		case "ENTRYPOINT":
			err = b.processEntrypoint(step, img)
		case "LABEL":
			err = b.processLabel(step, img)
		case "USER":
			err = b.processUser(step, img)
		case "VOLUME":
			err = b.processVolume(step, img)
		default:
			logger.Warn("Unknown instruction: %s", step.Instruction)
			if !config.Quiet {
				fmt.Printf("Warning: Unknown instruction '%s' - skipping\n", step.Instruction)
			}
		}

		if err != nil {
			return "", fmt.Errorf("step %d failed: %v", i+1, err)
		}
	}

	// If no FROM instruction was processed, create a minimal image
	if !fromProcessed {
		if !config.Quiet {
			fmt.Println("Warning: No FROM instruction found, creating minimal image")
		}
		img.Layers = []string{"scratch"}
	}

	// Set image tag if specified
	if config.Tag != "" {
		img.RepoTags = []string{config.Tag}
	} else {
		img.RepoTags = []string{"<none>:<none>"}
	}

	// Save the image
	err = b.imgManager.SaveImage(img)
	if err != nil {
		return "", fmt.Errorf("failed to save image: %v", err)
	}

	logger.Info("Image build completed successfully: %s", img.ID)
	return img.ID, nil
}

// parseBuildfile parses the Buildfile and returns build steps
func (b *ImageBuilder) parseBuildfile(buildfilePath string, buildArgs map[string]string) ([]BuildStep, error) {
	file, err := os.Open(buildfilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var steps []BuildStep
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Substitute build arguments
		for arg, value := range buildArgs {
			placeholder := fmt.Sprintf("$%s", arg)
			line = strings.ReplaceAll(line, placeholder, value)
			placeholder = fmt.Sprintf("${%s}", arg)
			line = strings.ReplaceAll(line, placeholder, value)
		}

		// Parse instruction and arguments
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		instruction := strings.ToUpper(parts[0])
		arguments := parts[1:]

		steps = append(steps, BuildStep{
			Instruction: instruction,
			Arguments:   arguments,
			RawLine:     line,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading Buildfile at line %d: %v", lineNum, err)
	}

	return steps, nil
}

// generateImageID generates a unique image ID
func generateImageID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// processFrom handles FROM instruction
func (b *ImageBuilder) processFrom(step BuildStep, img *image.Image) (*image.Image, error) {
	if len(step.Arguments) == 0 {
		return nil, fmt.Errorf("FROM instruction requires an argument")
	}

	baseImageName := step.Arguments[0]
	logger.Debug("Using base image: %s", baseImageName)

	// Handle special case for scratch
	if baseImageName == "scratch" {
		img.Layers = []string{"scratch"}
		return nil, nil
	}

	// Try to find the base image
	baseImage, err := b.imgManager.GetImage(baseImageName)
	if err != nil {
		return nil, fmt.Errorf("base image '%s' not found: %v", baseImageName, err)
	}

	// Copy configuration from base image
	img.Config = baseImage.Config
	img.Layers = append(img.Layers, baseImage.Layers...)
	img.RootFSType = baseImage.RootFSType

	return baseImage, nil
}

// processRun handles RUN instruction
func (b *ImageBuilder) processRun(step BuildStep, img *image.Image, contextPath string) error {
	if len(step.Arguments) == 0 {
		return fmt.Errorf("RUN instruction requires an argument")
	}

	command := strings.Join(step.Arguments, " ")
	logger.Debug("RUN: %s", command)

	// For now, we'll simulate the RUN instruction by adding it as metadata
	// In a full implementation, this would execute the command in a container
	layerID := fmt.Sprintf("run-%d", time.Now().UnixNano())
	img.Layers = append(img.Layers, layerID)
	img.Metadata[fmt.Sprintf("layer.%s.command", layerID)] = command
	img.Metadata[fmt.Sprintf("layer.%s.type", layerID)] = "run"

	return nil
}

// processCopy handles COPY instruction
func (b *ImageBuilder) processCopy(step BuildStep, img *image.Image, contextPath string) error {
	if len(step.Arguments) < 2 {
		return fmt.Errorf("COPY instruction requires at least 2 arguments")
	}

	sources := step.Arguments[:len(step.Arguments)-1]
	dest := step.Arguments[len(step.Arguments)-1]

	logger.Debug("COPY: %v -> %s", sources, dest)

	// Validate source files exist in build context
	for _, src := range sources {
		srcPath := filepath.Join(contextPath, src)
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			return fmt.Errorf("source file '%s' not found in build context", src)
		}
	}

	// Add copy operation as a layer
	layerID := fmt.Sprintf("copy-%d", time.Now().UnixNano())
	img.Layers = append(img.Layers, layerID)
	img.Metadata[fmt.Sprintf("layer.%s.sources", layerID)] = strings.Join(sources, ",")
	img.Metadata[fmt.Sprintf("layer.%s.dest", layerID)] = dest
	img.Metadata[fmt.Sprintf("layer.%s.type", layerID)] = "copy"

	return nil
}

// processAdd handles ADD instruction (similar to COPY but with URL support)
func (b *ImageBuilder) processAdd(step BuildStep, img *image.Image, contextPath string) error {
	if len(step.Arguments) < 2 {
		return fmt.Errorf("ADD instruction requires at least 2 arguments")
	}

	sources := step.Arguments[:len(step.Arguments)-1]
	dest := step.Arguments[len(step.Arguments)-1]

	logger.Debug("ADD: %v -> %s", sources, dest)

	// For now, treat ADD the same as COPY
	// In a full implementation, ADD would support URLs and automatic extraction
	return b.processCopy(step, img, contextPath)
}

// processWorkdir handles WORKDIR instruction
func (b *ImageBuilder) processWorkdir(step BuildStep, img *image.Image) error {
	if len(step.Arguments) == 0 {
		return fmt.Errorf("WORKDIR instruction requires an argument")
	}

	workdir := step.Arguments[0]
	img.Config.WorkingDir = workdir
	logger.Debug("WORKDIR: %s", workdir)

	return nil
}

// processEnv handles ENV instruction
func (b *ImageBuilder) processEnv(step BuildStep, img *image.Image) error {
	if len(step.Arguments) == 0 {
		return fmt.Errorf("ENV instruction requires an argument")
	}

	// Handle both "ENV key value" and "ENV key=value" formats
	if strings.Contains(step.Arguments[0], "=") {
		// ENV key=value format
		for _, arg := range step.Arguments {
			if strings.Contains(arg, "=") {
				parts := strings.SplitN(arg, "=", 2)
				if len(parts) == 2 {
					envVar := fmt.Sprintf("%s=%s", parts[0], parts[1])
					img.Config.Env = append(img.Config.Env, envVar)
					logger.Debug("ENV: %s", envVar)
				}
			}
		}
	} else {
		// ENV key value format
		if len(step.Arguments) < 2 {
			return fmt.Errorf("ENV instruction requires at least 2 arguments")
		}
		key := step.Arguments[0]
		value := strings.Join(step.Arguments[1:], " ")
		envVar := fmt.Sprintf("%s=%s", key, value)
		img.Config.Env = append(img.Config.Env, envVar)
		logger.Debug("ENV: %s", envVar)
	}

	return nil
}

// processExpose handles EXPOSE instruction
func (b *ImageBuilder) processExpose(step BuildStep, img *image.Image) error {
	if len(step.Arguments) == 0 {
		return fmt.Errorf("EXPOSE instruction requires an argument")
	}

	for _, port := range step.Arguments {
		img.Config.ExposedPorts[port] = struct{}{}
		logger.Debug("EXPOSE: %s", port)
	}

	return nil
}

// processCmd handles CMD instruction
func (b *ImageBuilder) processCmd(step BuildStep, img *image.Image) error {
	if len(step.Arguments) == 0 {
		return fmt.Errorf("CMD instruction requires an argument")
	}

	img.Config.Cmd = step.Arguments
	logger.Debug("CMD: %v", step.Arguments)

	return nil
}

// processEntrypoint handles ENTRYPOINT instruction
func (b *ImageBuilder) processEntrypoint(step BuildStep, img *image.Image) error {
	if len(step.Arguments) == 0 {
		return fmt.Errorf("ENTRYPOINT instruction requires an argument")
	}

	img.Config.Entrypoint = step.Arguments
	logger.Debug("ENTRYPOINT: %v", step.Arguments)

	return nil
}

// processLabel handles LABEL instruction
func (b *ImageBuilder) processLabel(step BuildStep, img *image.Image) error {
	if len(step.Arguments) == 0 {
		return fmt.Errorf("LABEL instruction requires an argument")
	}

	for _, arg := range step.Arguments {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			img.Config.Labels[parts[0]] = parts[1]
			logger.Debug("LABEL: %s=%s", parts[0], parts[1])
		}
	}

	return nil
}

// processUser handles USER instruction
func (b *ImageBuilder) processUser(step BuildStep, img *image.Image) error {
	if len(step.Arguments) == 0 {
		return fmt.Errorf("USER instruction requires an argument")
	}

	user := step.Arguments[0]
	img.Config.User = user
	logger.Debug("USER: %s", user)

	return nil
}

// processVolume handles VOLUME instruction
func (b *ImageBuilder) processVolume(step BuildStep, img *image.Image) error {
	if len(step.Arguments) == 0 {
		return fmt.Errorf("VOLUME instruction requires an argument")
	}

	// Add volume information to metadata
	volumes := strings.Join(step.Arguments, ",")
	img.Metadata["volumes"] = volumes
	logger.Debug("VOLUME: %s", volumes)

	return nil
}
