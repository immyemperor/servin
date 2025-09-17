# PyInstaller hook file for Servin GUI
# This ensures all local modules are properly included

from PyInstaller.utils.hooks import collect_all

# Collect all data files and submodules for the app
datas, binaries, hiddenimports = collect_all('app', include_py_files=True)

# Add specific modules that might be missed
hiddenimports += [
    'app',
    'servin_client',
    'mock_servin_client',
]

# Ensure templates and static files are included
datas += [
    ('templates/*', 'templates'),
    ('static/*', 'static'),
]