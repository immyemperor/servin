#!/usr/bin/env python3
"""
Quick verification script for servin-gui executable
This can be run to test the built executable
"""

import sys
import os
import subprocess
import time

def test_executable(executable_path):
    """Test the executable by running it briefly"""
    print(f"🧪 Testing executable: {executable_path}")
    
    if not os.path.exists(executable_path):
        print(f"❌ Executable not found: {executable_path}")
        return False
    
    if not os.access(executable_path, os.X_OK):
        print(f"❌ Executable is not executable: {executable_path}")
        return False
    
    print(f"✅ Executable exists and is executable")
    print(f"📊 Size: {os.path.getsize(executable_path) / (1024*1024):.1f} MB")
    
    # Try to run the executable for a short time
    try:
        print("🚀 Starting executable test (5 second timeout)...")
        process = subprocess.Popen(
            [executable_path],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )
        
        # Wait a few seconds for it to start
        try:
            stdout, stderr = process.communicate(timeout=5)
            print(f"🔍 Process exited with code: {process.returncode}")
            if stdout:
                print(f"📤 stdout: {stdout[:200]}...")
            if stderr:
                print(f"📤 stderr: {stderr[:200]}...")
        except subprocess.TimeoutExpired:
            print("⏰ Executable started successfully (timeout as expected)")
            process.terminate()
            process.wait()
            return True
            
    except Exception as e:
        print(f"❌ Failed to run executable: {e}")
        return False
    
    return True

if __name__ == "__main__":
    if len(sys.argv) > 1:
        executable_path = sys.argv[1]
    else:
        # Default paths for different platforms
        if sys.platform.startswith('win'):
            executable_path = 'dist/servin-gui.exe'
        else:
            executable_path = 'dist/servin-gui'
    
    test_executable(executable_path)