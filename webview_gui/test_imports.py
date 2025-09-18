#!/usr/bin/env python3
"""
Test script to verify PyInstaller build includes all required modules
"""

def test_imports():
    """Test that all required modules can be imported"""
    test_results = {}
    
    modules_to_test = [
        'app',
        'servin_client', 
        'mock_servin_client',
        'flask',
        'flask_cors',
        'flask_socketio',
        'webview',
        'tkinter',
        'threading',
        'queue',
        '_thread',
    ]
    
    # Async libraries (optional but preferred)
    async_modules = [
        'eventlet',
        'eventlet.wsgi',
        'gevent',
        'socketio',
        'engineio',
    ]
    
    print("[TEST] Testing core module imports...")
    
    for module_name in modules_to_test:
        try:
            __import__(module_name)
            test_results[module_name] = "[OK]"
            print(f"  {module_name}: [OK]")
        except ImportError as e:
            test_results[module_name] = f"[FAILED]: {e}"
            print(f"  {module_name}: [FAILED] {e}")
        except Exception as e:
            test_results[module_name] = f"[ERROR]: {e}"
            print(f"  {module_name}: [ERROR] {e}")
    
    print("\n[TEST] Testing async libraries...")
    
    for module_name in async_modules:
        try:
            __import__(module_name)
            test_results[module_name] = "[OK]"
            print(f"  {module_name}: [OK]")
        except ImportError as e:
            test_results[module_name] = f"[OPTIONAL]: {e}"
            print(f"  {module_name}: [OPTIONAL] {e}")
        except Exception as e:
            test_results[module_name] = f"[ERROR]: {e}"
            print(f"  {module_name}: [ERROR] {e}")
    
    print("\n[SUMMARY] Test Results Summary:")
    failed_count = 0
    for module, result in test_results.items():
        print(f"  {module}: {result}")
        if "FAILED" in result:
            failed_count += 1
    
    if failed_count == 0:
        print(f"\n[SUCCESS] All {len(modules_to_test)} modules imported successfully!")
        return True
    else:
        print(f"\n[FAILURE] {failed_count} modules failed to import")
        return False

if __name__ == "__main__":
    test_imports()