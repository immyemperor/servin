@echo off
REM Servin Container Runtime - Windows Installer
REM This script installs Servin as a Windows Service

echo ================================================
echo   Servin Container Runtime - Windows Installer
echo ================================================
echo.

REM Check if running as administrator
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo ERROR: This installer must be run as Administrator
    echo Please right-click and select "Run as Administrator"
    pause
    exit /b 1
)

REM Set installation directory
set INSTALL_DIR=C:\Program Files\Servin
set DATA_DIR=C:\ProgramData\Servin
set LOG_DIR=%DATA_DIR%\logs
set CONFIG_DIR=%DATA_DIR%\config
set VOLUMES_DIR=%DATA_DIR%\volumes
set IMAGES_DIR=%DATA_DIR%\images

echo Installing Servin to: %INSTALL_DIR%
echo Data directory: %DATA_DIR%

REM Create directories
echo Creating directories...
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"
if not exist "%DATA_DIR%" mkdir "%DATA_DIR%"
if not exist "%LOG_DIR%" mkdir "%LOG_DIR%"
if not exist "%CONFIG_DIR%" mkdir "%CONFIG_DIR%"
if not exist "%VOLUMES_DIR%" mkdir "%VOLUMES_DIR%"
if not exist "%IMAGES_DIR%" mkdir "%IMAGES_DIR%"

REM Copy executables
echo Copying executables...
copy /Y "servin.exe" "%INSTALL_DIR%\servin.exe" >nul
copy /Y "servin-gui.exe" "%INSTALL_DIR%\servin-gui.exe" >nul

if not exist "%INSTALL_DIR%\servin.exe" (
    echo ERROR: servin.exe not found in current directory
    echo Please run this installer from the directory containing the executables
    pause
    exit /b 1
)

REM Set permissions (allow full control for SYSTEM and Administrators)
echo Setting permissions...
icacls "%INSTALL_DIR%" /grant "SYSTEM:(OI)(CI)F" /grant "Administrators:(OI)(CI)F" >nul
icacls "%DATA_DIR%" /grant "SYSTEM:(OI)(CI)F" /grant "Administrators:(OI)(CI)F" >nul

REM Add to PATH
echo Adding to PATH...
setx PATH "%PATH%;%INSTALL_DIR%" /M >nul

REM Create default configuration
echo Creating default configuration...
(
echo # Servin Configuration File
echo # Data directory
echo data_dir=%DATA_DIR%
echo.
echo # Log settings
echo log_level=info
echo log_file=%LOG_DIR%\servin.log
echo.
echo # Runtime settings
echo runtime=native
echo.
echo # Network settings
echo bridge_name=servin0
echo.
echo # CRI settings
echo cri_port=10250
echo cri_enabled=false
) > "%CONFIG_DIR%\servin.conf"

REM Create Windows Service
echo Creating Windows Service...

REM Create service wrapper script
(
echo @echo off
echo cd /d "%INSTALL_DIR%"
echo "%INSTALL_DIR%\servin.exe" daemon --config "%CONFIG_DIR%\servin.conf"
) > "%INSTALL_DIR%\servin-service.bat"

REM Install service using sc command
sc create "ServinRuntime" ^
    binPath= "\"%INSTALL_DIR%\servin-service.bat\"" ^
    DisplayName= "Servin Container Runtime" ^
    Description= "Servin Container Runtime Service providing Docker-compatible container management" ^
    start= auto ^
    depend= "Tcpip"

if %errorLevel% neq 0 (
    echo WARNING: Failed to create Windows Service
    echo You can manually start Servin using: %INSTALL_DIR%\servin.exe
) else (
    echo Windows Service created successfully
)

REM Create uninstaller
echo Creating uninstaller...
(
echo @echo off
echo echo Uninstalling Servin Container Runtime...
echo.
echo REM Stop and remove service
echo sc stop "ServinRuntime" ^>nul 2^>^&1
echo sc delete "ServinRuntime" ^>nul 2^>^&1
echo.
echo REM Remove from PATH
echo for /f "tokens=2*" %%%%a in ^('reg query "HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment" /v PATH'^) do set "PATH=%%%%b"
echo set "PATH=%%PATH:;%INSTALL_DIR%=%%"
echo setx PATH "%%PATH%%" /M ^>nul
echo.
echo REM Remove directories
echo rmdir /s /q "%INSTALL_DIR%" ^>nul 2^>^&1
echo rmdir /s /q "%DATA_DIR%" ^>nul 2^>^&1
echo.
echo echo Servin has been uninstalled.
echo pause
) > "%INSTALL_DIR%\uninstall.bat"

REM Create desktop shortcuts
echo Creating shortcuts...
powershell "$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut('%PUBLIC%\Desktop\Servin GUI.lnk'); $Shortcut.TargetPath = '%INSTALL_DIR%\servin-gui.exe'; $Shortcut.Save()"

REM Create Start Menu entries
if not exist "%ProgramData%\Microsoft\Windows\Start Menu\Programs\Servin" mkdir "%ProgramData%\Microsoft\Windows\Start Menu\Programs\Servin"
powershell "$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut('%ProgramData%\Microsoft\Windows\Start Menu\Programs\Servin\Servin GUI.lnk'); $Shortcut.TargetPath = '%INSTALL_DIR%\servin-gui.exe'; $Shortcut.Save()"
powershell "$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut('%ProgramData%\Microsoft\Windows\Start Menu\Programs\Servin\Uninstall Servin.lnk'); $Shortcut.TargetPath = '%INSTALL_DIR%\uninstall.bat'; $Shortcut.Save()"

echo.
echo ================================================
echo   Installation completed successfully!
echo ================================================
echo.
echo Installation directory: %INSTALL_DIR%
echo Data directory: %DATA_DIR%
echo Configuration: %CONFIG_DIR%\servin.conf
echo.
echo Next steps:
echo 1. Start the service: sc start "ServinRuntime"
echo 2. Or run manually: %INSTALL_DIR%\servin.exe
echo 3. Use GUI: %INSTALL_DIR%\servin-gui.exe
echo.
echo The service is configured to start automatically on boot.
echo Check logs at: %LOG_DIR%\servin.log
echo.
pause
