#!/usr/bin/env python3
"""
Test script to check VM status and trigger VM section load
"""

import requests
import json

def test_vm_gui():
    base_url = "http://127.0.0.1:5555"
    
    # Test VM status endpoint
    try:
        print("Testing VM status endpoint...")
        response = requests.get(f"{base_url}/api/vm/status", timeout=5)
        print(f"Status Code: {response.status_code}")
        
        if response.status_code == 200:
            data = response.json()
            print(f"Response: {json.dumps(data, indent=2)}")
            
            # Check if VM is available and enabled
            if data.get('available', False):
                print("✅ VM is available")
                if data.get('enabled', False):
                    print("✅ VM mode is enabled")
                    if data.get('running', False):
                        print("✅ VM is running")
                    else:
                        print("🔶 VM is stopped - Start button should be enabled")
                else:
                    print("🔶 VM mode is disabled - Enable button should be enabled")
            else:
                print("❌ VM is not available")
        else:
            print(f"Error response: {response.text}")
            
    except requests.exceptions.ConnectionError as e:
        print(f"Connection Error: {e}")
    except Exception as e:
        print(f"Unexpected Error: {e}")

if __name__ == "__main__":
    test_vm_gui()