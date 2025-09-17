#!/bin/bash
# Servin Release Creation Script
# Creates GitHub releases with proper artifacts and documentation

set -e

# Configuration
GITHUB_REPO="immyemperor/servin"
RELEASE_BRANCH="main"
BUILD_DIR="dist"
TEMP_DIR=""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

print_header() {
    echo -e "${CYAN}================================================${NC}"
    echo -e "${CYAN}   Servin Release Creator${NC}"
    echo -e "${CYAN}================================================${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_info() {
    echo -e "${BLUE}→ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

show_usage() {
    cat << EOF
Usage: $0 <version> [options]

Arguments:
    version             Version to release (e.g., v1.0.0)

Options:
    --draft             Create as draft release
    --prerelease        Mark as pre-release
    --skip-build        Skip building artifacts
    --skip-tests        Skip running tests
    --notes-file FILE   Use custom release notes file
    --help              Show this help

Examples:
    $0 v1.0.0
    $0 v1.0.0-beta.1 --prerelease
    $0 v1.0.0 --draft --notes-file RELEASE_NOTES.md

EOF
}

# Parse command line arguments
parse_args() {
    if [ $# -eq 0 ]; then
        show_usage
        exit 1
    fi
    
    VERSION="$1"
    shift
    
    DRAFT=false
    PRERELEASE=false
    SKIP_BUILD=false
    SKIP_TESTS=false
    NOTES_FILE=""
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --draft)
                DRAFT=true
                shift
                ;;
            --prerelease)
                PRERELEASE=true
                shift
                ;;
            --skip-build)
                SKIP_BUILD=true
                shift
                ;;
            --skip-tests)
                SKIP_TESTS=true
                shift
                ;;
            --notes-file)
                NOTES_FILE="$2"
                shift 2
                ;;
            --help)
                show_usage
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    # Validate version format
    if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.-]+)?$ ]]; then
        print_error "Invalid version format. Use semantic versioning (e.g., v1.0.0)"
        exit 1
    fi
    
    print_info "Creating release: $VERSION"
}

# Check prerequisites
check_prerequisites() {
    print_info "Checking prerequisites..."
    
    # Check if gh CLI is installed
    if ! command -v gh >/dev/null 2>&1; then
        print_error "GitHub CLI (gh) is required but not installed"
        print_info "Install with: brew install gh (macOS) or https://cli.github.com/"
        exit 1
    fi
    
    # Check if authenticated
    if ! gh auth status >/dev/null 2>&1; then
        print_error "GitHub CLI not authenticated"
        print_info "Run: gh auth login"
        exit 1
    fi
    
    # Check if in git repository
    if ! git rev-parse --git-dir >/dev/null 2>&1; then
        print_error "Not in a git repository"
        exit 1
    fi
    
    # Check if on main branch
    local current_branch=$(git rev-parse --abbrev-ref HEAD)
    if [ "$current_branch" != "$RELEASE_BRANCH" ]; then
        print_warning "Not on $RELEASE_BRANCH branch (currently on $current_branch)"
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
    
    # Check if working directory is clean
    if ! git diff-index --quiet HEAD --; then
        print_error "Working directory has uncommitted changes"
        print_info "Commit or stash changes before creating release"
        exit 1
    fi
    
    print_success "Prerequisites checked"
}

# Run tests
run_tests() {
    if [ "$SKIP_TESTS" = true ]; then
        print_warning "Skipping tests"
        return
    fi
    
    print_info "Running tests..."
    
    # Run test suite
    if ! make test >/dev/null 2>&1; then
        print_error "Tests failed"
        exit 1
    fi
    
    # Run VM integration tests
    if ! make test-vm-integration >/dev/null 2>&1; then
        print_warning "VM integration tests failed (continuing anyway)"
    fi
    
    print_success "Tests passed"
}

# Build artifacts
build_artifacts() {
    if [ "$SKIP_BUILD" = true ]; then
        print_warning "Skipping build"
        return
    fi
    
    print_info "Building release artifacts..."
    
    # Clean previous builds
    rm -rf "$BUILD_DIR"
    
    # Build distribution
    if ! ./build-vm-distribution.sh --all; then
        print_error "Build failed"
        exit 1
    fi
    
    # Verify artifacts
    if [ ! -d "$BUILD_DIR" ]; then
        print_error "Build directory not found"
        exit 1
    fi
    
    # Generate checksums
    print_info "Generating checksums..."
    find "$BUILD_DIR" -type f \( -name "*.tar.gz" -o -name "*.zip" -o -name "*.deb" -o -name "*.rpm" \) -exec sha256sum {} \; > "$BUILD_DIR/checksums.txt"
    
    print_success "Artifacts built"
}

# Create git tag
create_tag() {
    print_info "Creating git tag..."
    
    # Check if tag already exists
    if git tag --list | grep -q "^$VERSION$"; then
        print_warning "Tag $VERSION already exists"
        read -p "Delete existing tag and continue? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            git tag -d "$VERSION"
            git push origin --delete "$VERSION" 2>/dev/null || true
        else
            exit 1
        fi
    fi
    
    # Create annotated tag
    git tag -a "$VERSION" -m "Release $VERSION"
    
    # Push tag
    git push origin "$VERSION"
    
    print_success "Tag created and pushed"
}

# Generate release notes
generate_release_notes() {
    print_info "Generating release notes..."
    
    TEMP_DIR=$(mktemp -d)
    local notes_file="$TEMP_DIR/release_notes.md"
    
    if [ -n "$NOTES_FILE" ] && [ -f "$NOTES_FILE" ]; then
        cp "$NOTES_FILE" "$notes_file"
        print_success "Using custom release notes from $NOTES_FILE"
        return
    fi
    
    # Get previous tag
    local prev_tag=$(git describe --tags --abbrev=0 HEAD~1 2>/dev/null || echo "")
    
    # Generate notes
    cat > "$notes_file" << EOF
# Servin $VERSION

## What's Changed

EOF
    
    # Add commits since last tag
    if [ -n "$prev_tag" ]; then
        echo "### Changes since $prev_tag" >> "$notes_file"
        echo "" >> "$notes_file"
        git log --pretty=format:"- %s (%an)" "$prev_tag..HEAD" >> "$notes_file"
    else
        echo "### Changes in this release" >> "$notes_file"
        echo "" >> "$notes_file"
        git log --pretty=format:"- %s (%an)" --max-count=20 >> "$notes_file"
    fi
    
    cat >> "$notes_file" << EOF

## Installation

### Universal Installer
\`\`\`bash
curl -fsSL https://get.servin.dev | sh
\`\`\`

### Package Managers
\`\`\`bash
# Homebrew (macOS/Linux)
brew install servin

# Winget (Windows)
winget install servin

# APT (Debian/Ubuntu)
curl -fsSL https://packages.servin.dev/gpg | sudo apt-key add -
echo "deb https://packages.servin.dev/apt stable main" | sudo tee /etc/apt/sources.list.d/servin.list
sudo apt update && sudo apt install servin
\`\`\`

## Checksums

EOF
    
    # Add checksums if available
    if [ -f "$BUILD_DIR/checksums.txt" ]; then
        echo "\`\`\`" >> "$notes_file"
        cat "$BUILD_DIR/checksums.txt" >> "$notes_file"
        echo "\`\`\`" >> "$notes_file"
    fi
    
    echo "" >> "$notes_file"
    echo "**Full Changelog**: https://github.com/$GITHUB_REPO/compare/$prev_tag...$VERSION" >> "$notes_file"
    
    print_success "Release notes generated"
}

# Create GitHub release
create_github_release() {
    print_info "Creating GitHub release..."
    
    local release_args=""
    
    if [ "$DRAFT" = true ]; then
        release_args="$release_args --draft"
    fi
    
    if [ "$PRERELEASE" = true ]; then
        release_args="$release_args --prerelease"
    fi
    
    # Create release
    gh release create "$VERSION" \
        --repo "$GITHUB_REPO" \
        --title "Servin $VERSION" \
        --notes-file "$TEMP_DIR/release_notes.md" \
        $release_args
    
    print_success "GitHub release created"
}

# Upload artifacts
upload_artifacts() {
    if [ "$SKIP_BUILD" = true ] || [ ! -d "$BUILD_DIR" ]; then
        print_warning "No artifacts to upload"
        return
    fi
    
    print_info "Uploading artifacts..."
    
    # Find and upload artifacts
    local artifacts=(
        "$BUILD_DIR/packages/"*.tar.gz
        "$BUILD_DIR/packages/"*.zip
        "$BUILD_DIR/packages/"*.deb
        "$BUILD_DIR/packages/"*.rpm
        "$BUILD_DIR/installers/"*
        "$BUILD_DIR/checksums.txt"
    )
    
    for artifact in "${artifacts[@]}"; do
        if [ -f "$artifact" ]; then
            gh release upload "$VERSION" "$artifact" --repo "$GITHUB_REPO"
            print_success "Uploaded $(basename "$artifact")"
        fi
    done
    
    print_success "All artifacts uploaded"
}

# Update package managers
update_package_managers() {
    print_info "Updating package managers..."
    
    # Update Homebrew formula (if script exists)
    if [ -f "scripts/update-homebrew.sh" ]; then
        print_info "Updating Homebrew formula..."
        ./scripts/update-homebrew.sh "$VERSION" || print_warning "Homebrew update failed"
    fi
    
    # Update other package managers
    if [ -f "scripts/update-packages.sh" ]; then
        print_info "Updating other package managers..."
        ./scripts/update-packages.sh "$VERSION" || print_warning "Package manager updates failed"
    fi
    
    print_info "Package manager updates initiated"
}

# Print next steps
print_next_steps() {
    echo ""
    echo -e "${CYAN}================================================${NC}"
    echo -e "${CYAN}   Release Created Successfully!${NC}"
    echo -e "${CYAN}================================================${NC}"
    echo ""
    echo -e "${GREEN}Release Details:${NC}"
    echo "• Version: $VERSION"
    echo "• GitHub: https://github.com/$GITHUB_REPO/releases/tag/$VERSION"
    if [ "$DRAFT" = true ]; then
        echo -e "• Status: ${YELLOW}Draft${NC} (remember to publish)"
    elif [ "$PRERELEASE" = true ]; then
        echo -e "• Status: ${YELLOW}Pre-release${NC}"
    else
        echo -e "• Status: ${GREEN}Published${NC}"
    fi
    echo ""
    
    echo -e "${GREEN}Next Steps:${NC}"
    echo "1. Review the release on GitHub"
    echo "2. Test the installation process:"
    echo -e "   ${YELLOW}curl -fsSL https://get.servin.dev | sh${NC}"
    echo "3. Monitor download statistics"
    echo "4. Update documentation if needed"
    echo "5. Announce the release"
    echo ""
    
    if [ "$DRAFT" = true ]; then
        echo -e "${YELLOW}Remember to publish the draft release when ready!${NC}"
        echo ""
    fi
}

# Cleanup
cleanup() {
    if [ -n "$TEMP_DIR" ] && [ -d "$TEMP_DIR" ]; then
        rm -rf "$TEMP_DIR"
    fi
}

# Main function
main() {
    # Setup trap for cleanup
    trap cleanup EXIT
    
    print_header
    parse_args "$@"
    check_prerequisites
    run_tests
    build_artifacts
    create_tag
    generate_release_notes
    create_github_release
    upload_artifacts
    update_package_managers
    print_next_steps
}

# Run main function with all arguments
main "$@"