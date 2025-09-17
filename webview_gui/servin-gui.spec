# -*- mode: python ; coding: utf-8 -*-
import sys

block_cipher = None

# Platform-specific executable name
exe_name = 'servin-gui'
if sys.platform.startswith('win'):
    exe_name = 'servin-gui.exe'

a = Analysis(
    ['main.py'],
    pathex=[],
    binaries=[],
    datas=[
        ('templates', 'templates'),
        ('static', 'static'),
    ],
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
    hookspath=[],
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
    debug=False,
    bootloader_ignore_signals=False,
    strip=False,
    upx=True,
    upx_exclude=[],
    runtime_tmpdir=None,
    console=False,
    disable_windowed_traceback=False,
    argv_emulation=False,
    target_arch=None,
    codesign_identity=None,
    entitlements_file=None,
    windowed=True,
)