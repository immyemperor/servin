; Servin Container Runtime - Windows Installer Wizard
; Built with NSIS (Nullsoft Scriptable Install System)

;--------------------------------
; Includes
!include "MUI2.nsh"
!include "LogicLib.nsh"
!include "WinVer.nsh"
!include "x64.nsh"

;--------------------------------
; General Settings
Name "Servin Container Runtime"
OutFile "ServinSetup-1.0.0.exe"
Unicode True

; Default installation directory
InstallDir "$PROGRAMFILES\Servin"
InstallDirRegKey HKLM "Software\Servin" "InstallPath"

; Request admin privileges
RequestExecutionLevel admin

; Compression
SetCompressor /SOLID lzma

;--------------------------------
; Variables
Var StartMenuFolder
Var ServiceRunning

;--------------------------------
; Version Information
VIProductVersion "1.0.0.0"
VIAddVersionKey "ProductName" "Servin Container Runtime"
VIAddVersionKey "Comments" "Docker-compatible container runtime with GUI"
VIAddVersionKey "CompanyName" "Servin Project"
VIAddVersionKey "LegalCopyright" "© 2025 Servin Project"
VIAddVersionKey "FileDescription" "Servin Container Runtime Installer"
VIAddVersionKey "FileVersion" "1.0.0.0"
VIAddVersionKey "ProductVersion" "1.0.0"
VIAddVersionKey "InternalName" "ServinSetup"

;--------------------------------
; Modern UI Configuration
!define MUI_ABORTWARNING
!define MUI_ICON "servin.ico"
!define MUI_UNICON "servin.ico"

; Header image
; !define MUI_HEADERIMAGE
; !define MUI_HEADERIMAGE_BITMAP "header.bmp"
; !define MUI_HEADERIMAGE_RIGHT

; Welcome page
!define MUI_WELCOMEPAGE_TITLE "Welcome to Servin Container Runtime Setup"
!define MUI_WELCOMEPAGE_TEXT "This wizard will guide you through the installation of Servin Container Runtime.$\r$\n$\r$\nServin is a lightweight, Docker-compatible container runtime with a modern GUI interface.$\r$\n$\r$\nClick Next to continue."

; License page
!define MUI_LICENSEPAGE_TEXT_TOP "Please review the license terms before installing Servin Container Runtime."
!define MUI_LICENSEPAGE_TEXT_BOTTOM "If you accept the terms of the agreement, click I Agree to continue. You must accept the agreement to install Servin Container Runtime."

; Components page
!define MUI_COMPONENTSPAGE_SMALLDESC

; Directory page
!define MUI_DIRECTORYPAGE_TEXT_TOP "Setup will install Servin Container Runtime in the following folder. To install in a different folder, click Browse and select another folder. Click Next to continue."

; Start Menu page
!define MUI_STARTMENUPAGE_DEFAULTFOLDER "Servin Container Runtime"
!define MUI_STARTMENUPAGE_REGISTRY_ROOT "HKLM"
!define MUI_STARTMENUPAGE_REGISTRY_KEY "Software\Servin"
!define MUI_STARTMENUPAGE_REGISTRY_VALUENAME "Start Menu Folder"

; Finish page
!define MUI_FINISHPAGE_TITLE "Servin Container Runtime Installation Complete"
!define MUI_FINISHPAGE_TEXT "Servin Container Runtime has been successfully installed on your computer.$\r$\n$\r$\nThe Servin service will start automatically. You can launch the GUI from the Start Menu or desktop shortcut."
!define MUI_FINISHPAGE_RUN "$INSTDIR\servin-gui.exe"
!define MUI_FINISHPAGE_RUN_TEXT "Launch Servin GUI"
!define MUI_FINISHPAGE_SHOWREADME "$INSTDIR\README.txt"
!define MUI_FINISHPAGE_SHOWREADME_TEXT "Show README"

;--------------------------------
; Pages
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "LICENSE.txt"
!insertmacro MUI_PAGE_COMPONENTS
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_STARTMENU Application $StartMenuFolder
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_WELCOME
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_UNPAGE_FINISH

;--------------------------------
; Languages
!insertmacro MUI_LANGUAGE "English"

;--------------------------------
; Functions
Function .onInit
  ; Check Windows version
  ${IfNot} ${AtLeastWin10}
    MessageBox MB_OK|MB_ICONSTOP "Servin Container Runtime requires Windows 10 or later."
    Abort
  ${EndIf}
  
  ; Check if 64-bit
  ${IfNot} ${RunningX64}
    MessageBox MB_OK|MB_ICONSTOP "Servin Container Runtime requires a 64-bit version of Windows."
    Abort
  ${EndIf}
  
  ; Check admin rights
  UserInfo::GetAccountType
  pop $0
  ${If} $0 != "admin"
    MessageBox MB_OK|MB_ICONSTOP "Administrator privileges are required to install Servin Container Runtime."
    Abort
  ${EndIf}
FunctionEnd

Function .onInstSuccess
  ; Start the service
  ExecWait 'sc start "ServinRuntime"'
FunctionEnd

Function un.onInit
  ; Check admin rights for uninstall
  UserInfo::GetAccountType
  pop $0
  ${If} $0 != "admin"
    MessageBox MB_OK|MB_ICONSTOP "Administrator privileges are required to uninstall Servin Container Runtime."
    Abort
  ${EndIf}
FunctionEnd

;--------------------------------
; Installation Sections

; Core Runtime (Required)
Section "Core Runtime" SecCore
  SectionIn RO  ; Read-only (required)
  
  SetOutPath "$INSTDIR"
  
  ; Stop service if running
  ExecWait 'sc stop "ServinRuntime"'
  
  ; Install main executable
  File "servin.exe"
  
  ; Install TUI (Terminal User Interface)
  File "servin-desktop.exe"
  
  ; Create data directories
  CreateDirectory "$APPDATA\Servin"
  CreateDirectory "$APPDATA\Servin\config"
  CreateDirectory "$APPDATA\Servin\data"
  CreateDirectory "$APPDATA\Servin\logs"
  CreateDirectory "$APPDATA\Servin\volumes"
  CreateDirectory "$APPDATA\Servin\images"
  
  ; Install configuration file
  SetOutPath "$APPDATA\Servin\config"
  File "servin.conf"
  
  ; Install documentation
  SetOutPath "$INSTDIR"
  File "README.txt"
  File "LICENSE.txt"
  
  ; Create uninstaller
  WriteUninstaller "$INSTDIR\Uninstall.exe"
  
  ; Registry entries
  WriteRegStr HKLM "Software\Servin" "InstallPath" "$INSTDIR"
  WriteRegStr HKLM "Software\Servin" "Version" "1.0.0"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin" "DisplayName" "Servin Container Runtime"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin" "UninstallString" "$INSTDIR\Uninstall.exe"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin" "InstallLocation" "$INSTDIR"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin" "DisplayIcon" "$INSTDIR\servin.exe"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin" "Publisher" "Servin Project"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin" "DisplayVersion" "1.0.0"
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin" "NoModify" 1
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin" "NoRepair" 1
  
  ; Add to PATH
  Push "$INSTDIR"
  Call AddToPath
SectionEnd

; GUI Application
Section "Desktop GUI" SecGUI
  SetOutPath "$INSTDIR"
  File "servin-gui.exe"
  
  ; Create desktop shortcut
  CreateShortcut "$DESKTOP\Servin GUI.lnk" "$INSTDIR\servin-gui.exe" "" "$INSTDIR\servin-gui.exe" 0
SectionEnd

; Windows Service
Section "Windows Service" SecService
  ; Install service using Windows sc command
  ExecWait 'sc create "ServinRuntime" binpath= "$INSTDIR\servin.exe daemon" displayname= "Servin Container Runtime" start= auto'
  
  ; Set service description
  ExecWait 'sc description "ServinRuntime" "Servin Container Runtime Service providing Docker-compatible container management"'
  
  ; Set service to restart on failure
  ExecWait 'sc failure "ServinRuntime" reset= 86400 actions= restart/60000/restart/60000/restart/60000'
SectionEnd

; Start Menu Shortcuts
Section "Start Menu Shortcuts" SecStartMenu
  !insertmacro MUI_STARTMENU_WRITE_BEGIN Application
  
  CreateDirectory "$SMPROGRAMS\$StartMenuFolder"
  CreateShortcut "$SMPROGRAMS\$StartMenuFolder\Servin GUI.lnk" "$INSTDIR\servin-gui.exe"
  CreateShortcut "$SMPROGRAMS\$StartMenuFolder\Servin Desktop (TUI).lnk" "$INSTDIR\servin-desktop.exe"
  CreateShortcut "$SMPROGRAMS\$StartMenuFolder\Servin Command Prompt.lnk" "cmd.exe" "/k cd /d $\"$INSTDIR$\""
  CreateShortcut "$SMPROGRAMS\$StartMenuFolder\Uninstall Servin.lnk" "$INSTDIR\Uninstall.exe"
  CreateShortcut "$SMPROGRAMS\$StartMenuFolder\README.lnk" "$INSTDIR\README.txt"
  
  !insertmacro MUI_STARTMENU_WRITE_END
SectionEnd

;--------------------------------
; Section Descriptions
!insertmacro MUI_FUNCTION_DESCRIPTION_BEGIN
  !insertmacro MUI_DESCRIPTION_TEXT ${SecCore} "Core Servin container runtime (required)"
  !insertmacro MUI_DESCRIPTION_TEXT ${SecGUI} "Desktop GUI application for managing containers"
  !insertmacro MUI_DESCRIPTION_TEXT ${SecService} "Install as Windows Service for automatic startup"
  !insertmacro MUI_DESCRIPTION_TEXT ${SecStartMenu} "Create Start Menu shortcuts"
!insertmacro MUI_FUNCTION_DESCRIPTION_END

;--------------------------------
; Uninstaller Section
Section "Uninstall"
  ; Stop and remove service
  ExecWait 'sc stop "ServinRuntime"'
  ExecWait 'sc delete "ServinRuntime"'
  
  ; Remove files
  Delete "$INSTDIR\servin.exe"
  Delete "$INSTDIR\servin-gui.exe"
  Delete "$INSTDIR\README.txt"
  Delete "$INSTDIR\LICENSE.txt"
  Delete "$INSTDIR\Uninstall.exe"
  
  ; Remove directories
  RMDir "$INSTDIR"
  
  ; Remove data directories (ask user)
  MessageBox MB_YESNO|MB_ICONQUESTION "Do you want to remove all container data, images, and volumes? This cannot be undone." IDNO skip_data_removal
  RMDir /r "$APPDATA\Servin"
  skip_data_removal:
  
  ; Remove shortcuts
  Delete "$DESKTOP\Servin GUI.lnk"
  
  ; Remove Start Menu shortcuts
  !insertmacro MUI_STARTMENU_GETFOLDER Application $StartMenuFolder
  Delete "$SMPROGRAMS\$StartMenuFolder\Servin GUI.lnk"
  Delete "$SMPROGRAMS\$StartMenuFolder\Servin Command Prompt.lnk"
  Delete "$SMPROGRAMS\$StartMenuFolder\Uninstall Servin.lnk"
  Delete "$SMPROGRAMS\$StartMenuFolder\README.lnk"
  RMDir "$SMPROGRAMS\$StartMenuFolder"
  
  ; Remove from PATH
  Push "$INSTDIR"
  Call un.RemoveFromPath
  
  ; Remove registry entries
  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin"
  DeleteRegKey HKLM "Software\Servin"
SectionEnd

;--------------------------------
; PATH Functions
!define Environ 'HKLM "SYSTEM\CurrentControlSet\Control\Session Manager\Environment"'

Function AddToPath
  Exch $0
  Push $1
  Push $2
  
  Call IsNT
  Pop $1
  ${If} $1 == 1
    ReadRegStr $1 ${Environ} "PATH"
    ${If} $1 != ""
      StrCpy $0 "$1;$0"
    ${EndIf}
    WriteRegExpandStr ${Environ} "PATH" $0
    SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000
  ${EndIf}
  
  Pop $2
  Pop $1
  Pop $0
FunctionEnd

Function un.RemoveFromPath
  Exch $0
  Push $1
  Push $2
  
  Call un.IsNT
  Pop $1
  ${If} $1 == 1
    ReadRegStr $1 ${Environ} "PATH"
    Push $1
    Push $0
    Call un.StrStr
    Pop $2
    ${If} $2 == ""
      Goto done
    ${EndIf}
    
    StrLen $2 $0
    StrCpy $1 $1 $2
    StrCpy $0 $0 "" $2
    StrCpy $0 "$1$0"
    WriteRegExpandStr ${Environ} "PATH" $0
    SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000
  ${EndIf}
  
  done:
  Pop $2
  Pop $1
  Pop $0
FunctionEnd

Function IsNT
  Push $0
  ReadRegStr $0 HKLM "SOFTWARE\Microsoft\Windows NT\CurrentVersion" CurrentVersion
  ${If} $0 == ""
    Push 0
  ${Else}
    Push 1
  ${EndIf}
  Exch
  Pop $0
FunctionEnd

Function un.IsNT
  Push $0
  ReadRegStr $0 HKLM "SOFTWARE\Microsoft\Windows NT\CurrentVersion" CurrentVersion
  ${If} $0 == ""
    Push 0
  ${Else}
    Push 1
  ${EndIf}
  Exch
  Pop $0
FunctionEnd

Function StrStr
  Exch $R1
  Exch
  Exch $R2
  Push $R3
  Push $R4
  Push $R5
  StrLen $R3 $R1
  StrCpy $R4 0
  
  loop:
    StrCpy $R5 $R2 $R3 $R4
    ${If} $R5 == $R1
      StrCpy $R0 $R2 $R4
      Goto done
    ${EndIf}
    ${If} $R3 >= $R4
      StrCpy $R0 ""
      Goto done
    ${EndIf}
    IntOp $R4 $R4 + 1
    Goto loop
  
  done:
  Pop $R5
  Pop $R4
  Pop $R3
  Pop $R2
  Exch $R0
  Pop $R1
FunctionEnd

Function un.StrStr
  Exch $R1
  Exch
  Exch $R2
  Push $R3
  Push $R4
  Push $R5
  StrLen $R3 $R1
  StrCpy $R4 0
  
  loop:
    StrCpy $R5 $R2 $R3 $R4
    ${If} $R5 == $R1
      StrCpy $R0 $R2 $R4
      Goto done
    ${EndIf}
    ${If} $R3 >= $R4
      StrCpy $R0 ""
      Goto done
    ${EndIf}
    IntOp $R4 $R4 + 1
    Goto loop
  
  done:
  Pop $R5
  Pop $R4
  Pop $R3
  Pop $R2
  Exch $R0
  Pop $R1
FunctionEnd
