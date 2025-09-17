# -*- mode: python ; coding: utf-8 -*-
import sys
import os

block_cipher = None

# Platform-specific executable name
exe_name = 'servin-gui'
if sys.platform.startswith('win'):
    exe_name = 'servin-gui.exe'

# Get the current working directory to ensure all paths are correct
current_dir = os.getcwd()
print(f"[SPEC] Building from directory: {current_dir}")

# Ensure all Python files are explicitly included as data files
additional_datas = [
    ('templates', 'templates'),
    ('static', 'static'),
    ('app.py', '.'),
    ('servin_client.py', '.'),
    ('mock_servin_client.py', '.'),
]

a = Analysis(
    ['main.py'],
    pathex=[current_dir],  # Explicitly add current directory to path
    binaries=[],
    datas=additional_datas,
    hiddenimports=[
        'app',
        'servin_client', 
        'mock_servin_client',
        'flask',
        'flask_cors',
        'webview',
        'tkinter',
        'tkinter.ttk',
        'tkinter.messagebox',
        'threading',
        'webbrowser',
        'subprocess',
        'json',
        'urllib.parse',
        'urllib.request',
    ],
    hookspath=[current_dir],  # Use current directory for hooks
    hooksconfig={},
    runtime_hooks=[],
    excludes=[
        'matplotlib',
        'numpy', 
        'pandas',
        'scipy',
        'pytest',
        'setuptools',
        'pip',
        'wheel',
    ],
    win_no_prefer_redirects=False,
    win_private_assemblies=False,
    cipher=block_cipher,
    noarchive=False,
)

pyz = PYZ(a.pure, a.zipped_data, cipher=block_cipher)

exe = EXE(
    pyz,
    a.scripts,
    a.binaries,
    a.zipfiles,
    a.datas,
    [],
    name=exe_name,
    debug=False,  # Set to True for debugging, False for release
    bootloader_ignore_signals=False,
    strip=False,
    upx=True,
    upx_exclude=[],
    runtime_tmpdir=None,
    console=False,  # Set to True for debugging, False for release
    disable_windowed_traceback=False,
    argv_emulation=False,
    target_arch=None,
    codesign_identity=None,
    entitlements_file=None,
    windowed=True,
)