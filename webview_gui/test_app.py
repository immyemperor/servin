#!/usr/bin/env python3
"""
Test script for Servin Desktop GUI
"""

import os
import sys
import subprocess

def test_imports():
    """Test if all required modules can be imported"""
    try:
        import flask
        print("✓ Flask module imported successfully")
    except ImportError as e:
        print(f"✗ Failed to import flask: {e}")
        return False
    
    try:
        import webview
        print("✓ Webview module imported successfully")
    except ImportError as e:
        print(f"✗ Failed to import webview: {e}")
        return False
    
    return True

def test_servin_connection():
    """Test Servin runtime connection"""
    try:
        import platform
        
        # Look for servin binary in the parent directory
        parent_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
        
        # Platform-specific binary detection
        system = platform.system().lower()
        machine = platform.machine().lower()
        
        search_paths = []
        
        if system == "windows":
            search_paths.extend([
                os.path.join(parent_dir, "build", "windows-amd64", "servin.exe"),
                os.path.join(parent_dir, "servin.exe"),
                os.path.join(parent_dir, "servin")
            ])
        elif system == "darwin":
            if machine in ["arm64", "aarch64"]:
                search_paths.extend([
                    os.path.join(parent_dir, "build", "darwin-arm64", "servin"),
                    os.path.join(parent_dir, "build", "darwin-amd64", "servin"),
                ])
            else:
                search_paths.extend([
                    os.path.join(parent_dir, "build", "darwin-amd64", "servin"),
                    os.path.join(parent_dir, "build", "darwin-arm64", "servin"),
                ])
            search_paths.extend([
                os.path.join(parent_dir, "servin"),
                os.path.join(parent_dir, "servin.exe")
            ])
        elif system == "linux":
            if machine in ["aarch64", "arm64"]:
                search_paths.extend([
                    os.path.join(parent_dir, "build", "linux-arm64", "servin"),
                    os.path.join(parent_dir, "build", "linux-amd64", "servin"),
                ])
            else:
                search_paths.extend([
                    os.path.join(parent_dir, "build", "linux-amd64", "servin"),
                    os.path.join(parent_dir, "build", "linux-arm64", "servin"),
                ])
            search_paths.extend([
                os.path.join(parent_dir, "servin"),
                os.path.join(parent_dir, "servin.exe")
            ])
        
        # Find the first existing and executable binary
        servin_path = None
        for path in search_paths:
            if os.path.exists(path) and os.access(path, os.X_OK):
                servin_path = path
                break
        
        if servin_path and os.path.exists(servin_path):
            result = subprocess.run([servin_path, '--help'], 
                                  capture_output=True, text=True, timeout=10)
            if result.returncode == 0:
                print("✓ Servin runtime is accessible")
                return True
        
        print("✗ Servin runtime not found or not accessible")
        return False
    except Exception as e:
        print(f"✗ Servin runtime not accessible: {e}")
        return False

def test_flask_app():
    """Test if Flask app can be imported"""
    try:
        current_dir = os.path.dirname(os.path.abspath(__file__))
        sys.path.insert(0, current_dir)
        from app import app
        print("✓ Flask app imported successfully")
        return True
    except Exception as e:
        print(f"✗ Failed to import Flask app: {e}")
        return False

def test_servin_client():
    """Test if Servin client can be imported"""
    try:
        current_dir = os.path.dirname(os.path.abspath(__file__))
        sys.path.insert(0, current_dir)
        from servin_client import ServinClient
        print("✓ Servin client imported successfully")
        return True
    except Exception as e:
        print(f"✗ Failed to import Servin client: {e}")
        return False

def main():
    print("Servin Desktop GUI - Test Suite")
    print("=" * 40)
    
    all_tests_passed = True
    
    print("\n1. Testing module imports...")
    if not test_imports():
        all_tests_passed = False
    
    print("\n2. Testing Servin client...")
    if not test_servin_client():
        all_tests_passed = False
    
    print("\n3. Testing Flask app...")
    if not test_flask_app():
        all_tests_passed = False
    
    print("\n4. Testing Servin runtime...")
    if not test_servin_connection():
        all_tests_passed = False
    
    print("\n" + "=" * 40)
    if all_tests_passed:
        print("✓ All tests passed! The application should work correctly.")
        print("\nTo run the application:")
        print("python main.py")
        print("\nOr run the demo:")
        print("python demo.py")
    else:
        print("✗ Some tests failed. Please fix the issues before running the application.")
    
    print("\nNote: Make sure the Servin binary is available in the parent directory.")

if __name__ == "__main__":
    main()
