#!/usr/bin/env python3

import platform
import os

# Test platform detection logic from the installer
def detect_platform():
    system = platform.system().lower()
    machine = platform.machine().lower()
    
    if system == "darwin":
        if machine in ["arm64", "aarch64"]:
            return "darwin-arm64"
        elif machine in ["x86_64", "amd64"]:
            return "darwin-amd64"
        else:
            return f"darwin-{machine}"
    else:
        return f"{system}-{machine}"

def test_build_directory_detection():
    script_dir = os.path.dirname(os.path.abspath(__file__))
    project_root = os.path.dirname(os.path.dirname(script_dir))  # Go up two levels from installers/macos/
    
    platform_name = detect_platform()
    build_dir = os.path.join(project_root, "build", platform_name)
    
    print(f"Script directory: {script_dir}")
    print(f"Project root: {project_root}")
    print(f"Detected platform: {platform_name}")
    print(f"Expected build directory: {build_dir}")
    print(f"Build directory exists: {os.path.exists(build_dir)}")
    
    if os.path.exists(build_dir):
        print(f"Contents of build directory:")
        for file in os.listdir(build_dir):
            file_path = os.path.join(build_dir, file)
            if os.path.isfile(file_path):
                print(f"  - {file} (executable: {os.access(file_path, os.X_OK)})")

if __name__ == "__main__":
    test_build_directory_detection()
