#!/usr/bin/env python3
"""
Test script for enhanced VM status display with color-coded engine status.
This script tests the VM status API response and verifies the new engine status fields.
"""

import requests
import json
import time

def test_vm_status_api():
    """Test the VM status API and display the response"""
    print("Testing Enhanced VM Status Display")
    print("=" * 50)
    
    try:
        # Test VM status endpoint
        response = requests.get('http://127.0.0.1:5555/api/vm/status', timeout=10)
        
        if response.status_code == 200:
            data = response.json()
            print("✅ VM Status API Response:")
            print(json.dumps(data, indent=2))
            
            # Test status interpretation
            print("\n🔍 Status Analysis:")
            if data.get('available', False):
                if data.get('running', False):
                    print("   Engine Status: 🟢 RUNNING")
                    print("   Color: Green (Success)")
                elif data.get('enabled', False):
                    print("   Engine Status: 🔴 STOPPED")
                    print("   Color: Red (Danger)")
                else:
                    print("   Engine Status: ⚪ DISABLED")
                    print("   Color: Gray (Secondary)")
            else:
                print("   Engine Status: ❌ UNAVAILABLE")
                print("   Color: Gray (Secondary)")
                
            print(f"\n📊 Details:")
            print(f"   Provider: {data.get('provider', 'Unknown')}")
            print(f"   Platform: {data.get('platform', 'Unknown')}")
            print(f"   Containers: {data.get('containers', 0)}")
            
            details = data.get('details', {})
            if details:
                print(f"   VM Name: {details.get('name', 'N/A')}")
                print(f"   IP Address: {details.get('ip', 'N/A')}")
                
        else:
            print(f"❌ API Error: {response.status_code}")
            print(f"Response: {response.text}")
            
    except requests.exceptions.RequestException as e:
        print(f"❌ Connection Error: {e}")
        print("Make sure the Flask app is running on http://127.0.0.1:5555")

def test_vm_operations():
    """Test VM operation endpoints"""
    print("\n" + "=" * 50)
    print("VM Operations Test (Visual Only)")
    print("=" * 50)
    
    operations = [
        ("Start VM", "POST", "/api/vm/start"),
        ("Stop VM", "POST", "/api/vm/stop"),
        ("Restart VM", "POST", "/api/vm/restart")
    ]
    
    for name, method, endpoint in operations:
        print(f"\n🔧 {name}:")
        print(f"   Endpoint: {method} http://127.0.0.1:5555{endpoint}")
        print(f"   Expected Transitional State: 'Starting'/'Stopping'/'Restarting'")
        print(f"   Expected Color: Orange (Warning) with pulse animation")

if __name__ == "__main__":
    test_vm_status_api()
    test_vm_operations()
    
    print("\n" + "=" * 50)
    print("🎨 Enhanced VM Status Features:")
    print("=" * 50)
    print("✅ Engine Status Indicator with color coding")
    print("✅ Separate engine status display in details")
    print("✅ Transitional states during operations")
    print("✅ Color-coded status text")
    print("✅ Improved visual feedback")
    print("\n🌈 Color Scheme:")
    print("   🟢 Running: Green with glow")
    print("   🔴 Stopped: Red with glow") 
    print("   🟠 Starting/Stopping/Restarting: Orange with pulse")
    print("   ⚪ Disabled/Unavailable: Gray")