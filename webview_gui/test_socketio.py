#!/usr/bin/env python3
"""
Test SocketIO async mode configuration for Windows PyInstaller builds
"""

import sys
import os

def test_socketio_async_modes():
    """Test different SocketIO async modes to ensure compatibility"""
    print("Testing SocketIO async mode configurations...")
    
    try:
        from flask import Flask
        from flask_socketio import SocketIO
        print("✓ Flask and SocketIO imported successfully")
    except ImportError as e:
        print(f"✗ Failed to import Flask/SocketIO: {e}")
        return False
    
    # Test different async modes
    test_app = Flask(__name__)
    
    async_modes = ['threading', 'eventlet', 'gevent']
    
    for mode in async_modes:
        try:
            print(f"Testing async_mode='{mode}'...")
            socketio = SocketIO(test_app, async_mode=mode, logger=False, engineio_logger=False)
            print(f"✓ SocketIO with async_mode='{mode}' works")
            return True  # If any mode works, we're good
        except ValueError as e:
            print(f"✗ SocketIO with async_mode='{mode}' failed: {e}")
        except Exception as e:
            print(f"✗ Unexpected error with async_mode='{mode}': {e}")
    
    # Test auto-detection as fallback
    try:
        print("Testing auto-detection mode...")
        socketio = SocketIO(test_app, logger=False, engineio_logger=False)
        print("✓ SocketIO with auto-detection works")
        return True
    except Exception as e:
        print(f"✗ SocketIO auto-detection failed: {e}")
        return False

def test_app_import():
    """Test importing the main app module"""
    print("Testing app.py import...")
    
    try:
        # Add current directory to path
        current_dir = os.path.dirname(os.path.abspath(__file__))
        if current_dir not in sys.path:
            sys.path.insert(0, current_dir)
        
        import app
        print("✓ app.py imported successfully")
        
        # Test if socketio is configured
        if hasattr(app, 'socketio'):
            print("✓ SocketIO instance found in app")
            return True
        else:
            print("✗ SocketIO instance not found in app")
            return False
            
    except Exception as e:
        print(f"✗ Failed to import app.py: {e}")
        return False

def main():
    """Run all tests"""
    print("=" * 50)
    print("Servin WebView GUI - SocketIO Test")
    print("=" * 50)
    
    # Detect PyInstaller environment
    if hasattr(sys, '_MEIPASS'):
        print(f"Running from PyInstaller bundle: {sys._MEIPASS}")
    else:
        print("Running from source")
    
    print(f"Platform: {sys.platform}")
    print(f"Python version: {sys.version}")
    print()
    
    tests_passed = 0
    total_tests = 2
    
    # Test SocketIO async modes
    if test_socketio_async_modes():
        tests_passed += 1
    
    print()
    
    # Test app import
    if test_app_import():
        tests_passed += 1
    
    print()
    print("=" * 50)
    print(f"Tests passed: {tests_passed}/{total_tests}")
    
    if tests_passed == total_tests:
        print("✓ All tests passed! SocketIO configuration is working.")
        return 0
    else:
        print("✗ Some tests failed. Check the configuration.")
        return 1

if __name__ == "__main__":
    sys.exit(main())