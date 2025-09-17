#!/usr/bin/env python3
"""
Simple test script to verify VM API endpoints
"""

import requests
import json

def test_vm_api():
    base_url = "http://127.0.0.1:5555"
    
    # Test VM status endpoint
    try:
        print("Testing VM status endpoint...")
        print(f"Connecting to: {base_url}/api/vm/status")
        response = requests.get(f"{base_url}/api/vm/status", timeout=5)
        print(f"Status Code: {response.status_code}")
        
        if response.status_code == 200:
            try:
                data = response.json()
                print(f"Response: {json.dumps(data, indent=2)}")
            except:
                print(f"Response (text): {response.text}")
        else:
            print(f"Error response: {response.text}")
            
    except requests.exceptions.ConnectionError as e:
        print(f"Connection Error: {e}")
        print("Make sure Flask server is running on port 5555.")
    except requests.exceptions.Timeout:
        print("Error: Request timed out.")
    except Exception as e:
        print(f"Unexpected Error: {e}")

if __name__ == "__main__":
    test_vm_api()