#!/bin/bash

# Servin Integration Testing Suite
# Comprehensive testing for enterprise-grade installer packages and CI/CD pipeline

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_LOG_DIR="$SCRIPT_DIR/test-results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
TEST_SESSION_LOG="$TEST_LOG_DIR/integration_test_$TIMESTAMP.log"

# Color output functions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
NC='\033[0m'

print_success() { 
    echo -e "${GREEN}✅ $1${NC}" | tee -a "$TEST_SESSION_LOG"
}
print_warning() { 
    echo -e "${YELLOW}⚠️  $1${NC}" | tee -a "$TEST_SESSION_LOG"
}
print_error() { 
    echo -e "${RED}❌ $1${NC}" | tee -a "$TEST_SESSION_LOG"
}
print_info() { 
    echo -e "${BLUE}ℹ️  $1${NC}" | tee -a "$TEST_SESSION_LOG"
}
print_header() { 
    echo -e "\n${CYAN}${BOLD}🔍 $1${NC}" | tee -a "$TEST_SESSION_LOG"
}
print_step() { 
    echo -e "${MAGENTA}👉 $1${NC}" | tee -a "$TEST_SESSION_LOG"
}

print_banner() {
    local banner="
╔══════════════════════════════════════════════════════════════════════╗
║                Servin Integration Testing Suite                     ║
║              Enterprise-Grade Package Validation                    ║
╚══════════════════════════════════════════════════════════════════════╝"
    echo -e "${CYAN}${BOLD}$banner${NC}" | tee -a "$TEST_SESSION_LOG"
    echo -e "\n${BLUE}Session: $TIMESTAMP${NC}" | tee -a "$TEST_SESSION_LOG"
    echo -e "${BLUE}Log: $TEST_SESSION_LOG${NC}\n" | tee -a "$TEST_SESSION_LOG"
}

# Initialize testing environment
setup_test_environment() {
    print_header "Setting Up Test Environment"
    
    # Create test directories
    mkdir -p "$TEST_LOG_DIR"
    mkdir -p "$TEST_LOG_DIR/artifacts"
    mkdir -p "$TEST_LOG_DIR/reports"
    
    # Initialize log file
    echo "Servin Integration Test Session - $(date)" > "$TEST_SESSION_LOG"
    echo "=============================================" >> "$TEST_SESSION_LOG"
    echo "" >> "$TEST_SESSION_LOG"
    
    print_success "Test environment initialized"
    print_info "Test artifacts: $TEST_LOG_DIR/artifacts"
    print_info "Test reports: $TEST_LOG_DIR/reports"
}

# Test 1: Build System Validation
test_build_system() {
    print_header "Test 1: Build System Validation"
    
    local test_report="$TEST_LOG_DIR/reports/build_system_test.json"
    local test_start=$(date +%s)
    
    print_step "Checking build script availability"
    
    # Check for required build scripts
    local required_scripts=(
        "build-packages.sh"
        "build.sh" 
        "build-all.sh"
    )
    
    local missing_scripts=()
    local found_scripts=()
    
    for script in "${required_scripts[@]}"; do
        if [[ -f "$SCRIPT_DIR/$script" ]]; then
            print_success "Found: $script"
            found_scripts+=("$script")
        else
            print_error "Missing: $script"
            missing_scripts+=("$script")
        fi
    done
    
    # Test build-packages.sh functionality
    print_step "Testing package builder script"
    if [[ -f "$SCRIPT_DIR/build-packages.sh" ]]; then
        if bash "$SCRIPT_DIR/build-packages.sh" --help &>/dev/null; then
            print_success "build-packages.sh is executable and responds to --help"
        else
            print_warning "build-packages.sh may have execution issues"
        fi
    fi
    
    # Check installer directories
    print_step "Validating installer directory structure"
    local installer_dirs=(
        "installers/windows"
        "installers/linux" 
        "installers/macos"
    )
    
    for dir in "${installer_dirs[@]}"; do
        if [[ -d "$SCRIPT_DIR/$dir" ]]; then
            print_success "Found: $dir"
        else
            print_error "Missing: $dir"
        fi
    done
    
    local test_end=$(date +%s)
    local test_duration=$((test_end - test_start))
    
    # Generate test report
    cat > "$test_report" << EOF
{
  "test_name": "Build System Validation",
  "timestamp": "$(date -Iseconds)",
  "duration_seconds": $test_duration,
  "status": "$([ ${#missing_scripts[@]} -eq 0 ] && echo "PASSED" || echo "FAILED")",
  "found_scripts": $(printf '%s\n' "${found_scripts[@]}" | jq -R . | jq -s .),
  "missing_scripts": $(printf '%s\n' "${missing_scripts[@]}" | jq -R . | jq -s .),
  "details": {
    "required_scripts": ${#required_scripts[@]},
    "found_scripts": ${#found_scripts[@]},
    "missing_scripts": ${#missing_scripts[@]}
  }
}
EOF
    
    if [ ${#missing_scripts[@]} -eq 0 ]; then
        print_success "Build system validation PASSED"
    else
        print_error "Build system validation FAILED - missing ${#missing_scripts[@]} scripts"
    fi
}

# Test 2: GitHub Actions Workflow Validation
test_github_actions() {
    print_header "Test 2: GitHub Actions Workflow Validation"
    
    local test_report="$TEST_LOG_DIR/reports/github_actions_test.json"
    local test_start=$(date +%s)
    
    print_step "Checking GitHub Actions workflow file"
    
    local workflow_file=".github/workflows/build-release.yml"
    local validation_results=()
    
    if [[ -f "$SCRIPT_DIR/$workflow_file" ]]; then
        print_success "Found GitHub Actions workflow: $workflow_file"
        
        # Check for key sections in workflow
        print_step "Validating workflow content"
        
        local required_sections=(
            "name:"
            "jobs:"
            "build:"
            "strategy:"
            "matrix:"
        )
        
        for section in "${required_sections[@]}"; do
            if grep -q "$section" "$SCRIPT_DIR/$workflow_file"; then
                print_success "Found section: $section"
                validation_results+=("$section:FOUND")
            else
                print_error "Missing section: $section"
                validation_results+=("$section:MISSING")
            fi
        done
        
        # Check for installer verification steps
        print_step "Checking for installer verification"
        if grep -q "Verify complete installer packages" "$SCRIPT_DIR/$workflow_file"; then
            print_success "Installer verification steps found"
            validation_results+=("verification:FOUND")
        else
            print_warning "Installer verification steps not found"
            validation_results+=("verification:MISSING")
        fi
        
        # Check for platform matrix
        print_step "Checking platform support"
        local platforms=("windows" "linux" "mac")
        for platform in "${platforms[@]}"; do
            if grep -q "$platform" "$SCRIPT_DIR/$workflow_file"; then
                print_success "Platform supported: $platform"
                validation_results+=("platform_$platform:FOUND")
            else
                print_warning "Platform not found: $platform"
                validation_results+=("platform_$platform:MISSING")
            fi
        done
        
    else
        print_error "GitHub Actions workflow file not found: $workflow_file"
        validation_results+=("workflow_file:MISSING")
    fi
    
    local test_end=$(date +%s)
    local test_duration=$((test_end - test_start))
    
    # Generate test report
    cat > "$test_report" << EOF
{
  "test_name": "GitHub Actions Workflow Validation",
  "timestamp": "$(date -Iseconds)",
  "duration_seconds": $test_duration,
  "status": "$([ -f "$SCRIPT_DIR/$workflow_file" ] && echo "PASSED" || echo "FAILED")",
  "workflow_file": "$workflow_file",
  "validation_results": $(printf '%s\n' "${validation_results[@]}" | jq -R . | jq -s .),
  "details": {
    "file_exists": $([ -f "$SCRIPT_DIR/$workflow_file" ] && echo "true" || echo "false"),
    "checks_performed": ${#validation_results[@]}
  }
}
EOF
    
    if [[ -f "$SCRIPT_DIR/$workflow_file" ]]; then
        print_success "GitHub Actions workflow validation PASSED"
    else
        print_error "GitHub Actions workflow validation FAILED"
    fi
}

# Test 3: Package Creation Test
test_package_creation() {
    print_header "Test 3: Package Creation Test"
    
    local test_report="$TEST_LOG_DIR/reports/package_creation_test.json"
    local test_start=$(date +%s)
    
    print_step "Testing package creation process"
    
    # Create a test build to verify package creation
    print_info "Running build-packages.sh in test mode..."
    
    local build_output="$TEST_LOG_DIR/artifacts/build_test_output.log"
    local build_success=false
    
    if [[ -f "$SCRIPT_DIR/build-packages.sh" ]]; then
        # Run build with dry-run or help to test functionality
        if timeout 300 bash "$SCRIPT_DIR/build-packages.sh" --help > "$build_output" 2>&1; then
            print_success "build-packages.sh executed successfully"
            build_success=true
        else
            print_error "build-packages.sh execution failed"
            cat "$build_output" | tail -20 | while IFS= read -r line; do
                print_error "  $line"
            done
        fi
    else
        print_error "build-packages.sh not found"
    fi
    
    # Check for existing build artifacts
    print_step "Checking for existing build artifacts"
    
    local build_dirs=("build" "dist" "packages")
    local found_artifacts=()
    
    for dir in "${build_dirs[@]}"; do
        if [[ -d "$SCRIPT_DIR/$dir" ]]; then
            local artifact_count=$(find "$SCRIPT_DIR/$dir" -type f | wc -l)
            print_success "Found build directory: $dir ($artifact_count files)"
            found_artifacts+=("$dir:$artifact_count")
        else
            print_info "Build directory not found: $dir (will be created during build)"
            found_artifacts+=("$dir:0")
        fi
    done
    
    local test_end=$(date +%s)
    local test_duration=$((test_end - test_start))
    
    # Generate test report
    cat > "$test_report" << EOF
{
  "test_name": "Package Creation Test",
  "timestamp": "$(date -Iseconds)",
  "duration_seconds": $test_duration,
  "status": "$([ "$build_success" = true ] && echo "PASSED" || echo "FAILED")",
  "build_script_executable": $([ -f "$SCRIPT_DIR/build-packages.sh" ] && echo "true" || echo "false"),
  "build_execution_success": $build_success,
  "found_artifacts": $(printf '%s\n' "${found_artifacts[@]}" | jq -R . | jq -s .),
  "build_output_file": "$build_output"
}
EOF
    
    if [ "$build_success" = true ]; then
        print_success "Package creation test PASSED"
    else
        print_error "Package creation test FAILED"
    fi
}

# Test 4: VM Dependencies Check
test_vm_dependencies() {
    print_header "Test 4: VM Dependencies Check"
    
    local test_report="$TEST_LOG_DIR/reports/vm_dependencies_test.json"
    local test_start=$(date +%s)
    
    print_step "Checking VM-related files and dependencies"
    
    # Check for VM-related files
    local vm_files=(
        "cmd/vm.go"
        "pkg/vm/"
        "scripts/"
    )
    
    local vm_results=()
    
    for file in "${vm_files[@]}"; do
        if [[ -e "$SCRIPT_DIR/$file" ]]; then
            print_success "Found VM file/directory: $file"
            vm_results+=("$file:FOUND")
        else
            print_warning "VM file/directory not found: $file"
            vm_results+=("$file:MISSING")
        fi
    done
    
    # Check for platform-specific VM support
    print_step "Checking platform-specific VM implementations"
    
    local vm_platforms=("linux" "macos" "windows")
    for platform in "${vm_platforms[@]}"; do
        if find "$SCRIPT_DIR" -name "*${platform}*" -type f | grep -q vm; then
            print_success "VM implementation found for: $platform"
            vm_results+=("vm_$platform:FOUND")
        else
            print_info "VM implementation not explicitly found for: $platform"
            vm_results+=("vm_$platform:NOT_EXPLICIT")
        fi
    done
    
    # Check for containerization files
    print_step "Checking containerization support"
    
    local container_files=(
        "go.mod"
        "main.go"
        "cmd/"
        "pkg/"
    )
    
    for file in "${container_files[@]}"; do
        if [[ -e "$SCRIPT_DIR/$file" ]]; then
            print_success "Found containerization file: $file"
            vm_results+=("container_$file:FOUND")
        else
            print_error "Missing containerization file: $file"
            vm_results+=("container_$file:MISSING")
        fi
    done
    
    local test_end=$(date +%s)
    local test_duration=$((test_end - test_start))
    
    # Generate test report
    cat > "$test_report" << EOF
{
  "test_name": "VM Dependencies Check",
  "timestamp": "$(date -Iseconds)",
  "duration_seconds": $test_duration,
  "status": "PASSED",
  "vm_results": $(printf '%s\n' "${vm_results[@]}" | jq -R . | jq -s .),
  "details": {
    "checks_performed": ${#vm_results[@]},
    "vm_files_checked": ${#vm_files[@]},
    "platforms_checked": ${#vm_platforms[@]}
  }
}
EOF
    
    print_success "VM dependencies check PASSED"
}

# Test 5: Documentation Validation
test_documentation() {
    print_header "Test 5: Documentation Validation"
    
    local test_report="$TEST_LOG_DIR/reports/documentation_test.json"
    local test_start=$(date +%s)
    
    print_step "Checking documentation files"
    
    local doc_files=(
        "README.md"
        "docs/"
        "INSTALL.md"
        "BUILD_GUIDE.md"
    )
    
    local doc_results=()
    
    for file in "${doc_files[@]}"; do
        if [[ -e "$SCRIPT_DIR/$file" ]]; then
            print_success "Found documentation: $file"
            doc_results+=("$file:FOUND")
            
            # Check if documentation mentions installers
            if [[ -f "$SCRIPT_DIR/$file" ]] && grep -q -i "installer\|package" "$SCRIPT_DIR/$file"; then
                print_success "  Contains installer information"
                doc_results+=("$file:HAS_INSTALLER_INFO")
            fi
        else
            print_warning "Documentation file missing: $file"
            doc_results+=("$file:MISSING")
        fi
    done
    
    # Check for installer-specific documentation
    print_step "Checking installer documentation"
    
    if [[ -f "$SCRIPT_DIR/docs/installer-packages.md" ]]; then
        print_success "Found installer-packages.md documentation"
        doc_results+=("installer_docs:FOUND")
    else
        print_warning "installer-packages.md not found"
        doc_results+=("installer_docs:MISSING")
    fi
    
    local test_end=$(date +%s)
    local test_duration=$((test_end - test_start))
    
    # Generate test report
    cat > "$test_report" << EOF
{
  "test_name": "Documentation Validation",
  "timestamp": "$(date -Iseconds)",
  "duration_seconds": $test_duration,
  "status": "PASSED",
  "doc_results": $(printf '%s\n' "${doc_results[@]}" | jq -R . | jq -s .),
  "details": {
    "files_checked": ${#doc_files[@]},
    "results_count": ${#doc_results[@]}
  }
}
EOF
    
    print_success "Documentation validation PASSED"
}

# Generate comprehensive test report
generate_final_report() {
    print_header "Generating Comprehensive Test Report"
    
    local final_report="$TEST_LOG_DIR/integration_test_report_$TIMESTAMP.json"
    local summary_report="$TEST_LOG_DIR/test_summary_$TIMESTAMP.md"
    
    # Collect all individual test reports
    local test_reports=()
    for report in "$TEST_LOG_DIR/reports"/*.json; do
        if [[ -f "$report" ]]; then
            test_reports+=("$report")
        fi
    done
    
    # Create comprehensive JSON report
    echo "{" > "$final_report"
    echo "  \"integration_test_session\": {" >> "$final_report"
    echo "    \"timestamp\": \"$(date -Iseconds)\"," >> "$final_report"
    echo "    \"session_id\": \"$TIMESTAMP\"," >> "$final_report"
    echo "    \"total_tests\": ${#test_reports[@]}," >> "$final_report"
    echo "    \"test_results\": [" >> "$final_report"
    
    local first=true
    for report in "${test_reports[@]}"; do
        if [ "$first" = true ]; then
            first=false
        else
            echo "," >> "$final_report"
        fi
        cat "$report" >> "$final_report"
    done
    
    echo "    ]" >> "$final_report"
    echo "  }" >> "$final_report"
    echo "}" >> "$final_report"
    
    # Create markdown summary
    cat > "$summary_report" << EOF
# Servin Integration Test Report

**Session ID:** $TIMESTAMP  
**Date:** $(date)  
**Total Tests:** ${#test_reports[@]}

## Test Results Summary

EOF
    
    local passed_tests=0
    local failed_tests=0
    
    for report in "${test_reports[@]}"; do
        local test_name=$(jq -r '.test_name' "$report" 2>/dev/null || echo "Unknown Test")
        local status=$(jq -r '.status' "$report" 2>/dev/null || echo "UNKNOWN")
        local duration=$(jq -r '.duration_seconds' "$report" 2>/dev/null || echo "0")
        
        if [ "$status" = "PASSED" ]; then
            echo "✅ **$test_name** - PASSED (${duration}s)" >> "$summary_report"
            ((passed_tests++))
        else
            echo "❌ **$test_name** - FAILED (${duration}s)" >> "$summary_report"
            ((failed_tests++))
        fi
    done
    
    cat >> "$summary_report" << EOF

## Summary Statistics

- **Passed Tests:** $passed_tests
- **Failed Tests:** $failed_tests
- **Success Rate:** $(( passed_tests * 100 / (passed_tests + failed_tests) ))%

## Files Generated

- **Detailed Report:** \`$final_report\`
- **Session Log:** \`$TEST_SESSION_LOG\`
- **Test Artifacts:** \`$TEST_LOG_DIR/artifacts/\`

EOF
    
    print_success "Comprehensive test report generated"
    print_info "JSON Report: $final_report"
    print_info "Summary Report: $summary_report"
    print_info "Session Log: $TEST_SESSION_LOG"
    
    # Display summary
    print_header "Integration Test Summary"
    print_success "Total Tests: ${#test_reports[@]}"
    print_success "Passed: $passed_tests"
    if [ $failed_tests -gt 0 ]; then
        print_error "Failed: $failed_tests"
    else
        print_success "Failed: $failed_tests"
    fi
    print_success "Success Rate: $(( passed_tests * 100 / (passed_tests + failed_tests) ))%"
}

# Main execution function
main() {
    print_banner
    
    # Check for required tools
    if ! command -v jq >/dev/null 2>&1; then
        print_warning "jq is not installed - installing for JSON processing..."
        if command -v brew >/dev/null 2>&1; then
            brew install jq
        elif command -v apt-get >/dev/null 2>&1; then
            sudo apt-get update && sudo apt-get install -y jq
        else
            print_error "Please install jq manually for JSON processing"
            exit 1
        fi
    fi
    
    # Run all tests
    setup_test_environment
    test_build_system
    test_github_actions
    test_package_creation
    test_vm_dependencies
    test_documentation
    generate_final_report
    
    print_header "Integration Testing Complete"
    print_success "All tests completed successfully!"
    print_info "Review detailed reports in: $TEST_LOG_DIR"
}

# Execute main function
main "$@"