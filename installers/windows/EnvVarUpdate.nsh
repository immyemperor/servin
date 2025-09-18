/*
EnvVarUpdate.nsh
Environment Variable Update Functions for NSIS
Provides functions to add, remove, and update PATH and other environment variables
*/

!verbose push
!verbose 3
!ifndef _EnvVarUpdate_nsh
!define _EnvVarUpdate_nsh
!verbose pop

!include "LogicLib.nsh"

!define HWND_BROADCAST      0xffff
!define WM_WININICHANGE     0x001A
!define SMTO_ABORTIFHUNG    0x0002

!macro _EnvVarUpdateConstructor ResultVar EnvVarName Action Regloc PathString
  Push "${PathString}"
  Push "${Regloc}"
  Push "${Action}"
  Push "${EnvVarName}"
  Call EnvVarUpdate
  Pop "${ResultVar}"
!macroend

!define EnvVarUpdate '!insertmacro "_EnvVarUpdateConstructor"'

!macro _unEnvVarUpdateConstructor ResultVar EnvVarName Action Regloc PathString
  Push "${PathString}"
  Push "${Regloc}"
  Push "${Action}"
  Push "${EnvVarName}"
  Call un.EnvVarUpdate
  Pop "${ResultVar}"
!macroend

!define un.EnvVarUpdate '!insertmacro "_unEnvVarUpdateConstructor"'

; ---- Fix for Windows Vista/7 UAC Policy (Modified) ----
!macro _EnvVarUpdate_Modify RegLoc EnvVarName PathString Regloc Action
  DetailPrint "Updating environment variable: ${EnvVarName}"
  
  ReadRegStr $R5 ${RegLoc} "Environment" "${EnvVarName}"
  
  StrCpy $R2 $R5 1 -1 ; copy last char
  ${If} "$R2" == ";"
    StrCpy $R5 $R5 -1 ; remove last char
  ${EndIf}
  
  ${If} "${Action}" == "A"
    ; Add to path
    ${If} $R5 == ""
      StrCpy $R5 "${PathString}"
    ${Else}
      StrCpy $R5 "$R5;${PathString}"
    ${EndIf}
  ${ElseIf} "${Action}" == "P"
    ; Prepend to path
    ${If} $R5 == ""
      StrCpy $R5 "${PathString}"
    ${Else}
      StrCpy $R5 "${PathString};$R5"
    ${EndIf}
  ${ElseIf} "${Action}" == "R"
    ; Remove from path
    StrCpy $R1 $R5
    ${Do}
      StrCpy $R3 $R1 "" "${PathString}"
      ${If} $R3 != $R1
        StrCpy $R1 $R3
      ${Else}
        ${Break}
      ${EndIf}
    ${Loop}
    StrCpy $R5 $R1
    
    ; Clean up multiple semicolons
    ${Do}
      StrCpy $R3 $R5
      ${StrRep} $R5 $R5 ";;" ";"
      ${If} $R5 == $R3
        ${Break}
      ${EndIf}
    ${Loop}
    
    ; Remove leading/trailing semicolons
    StrCpy $R2 $R5 1
    ${If} "$R2" == ";"
      StrCpy $R5 $R5 "" 1
    ${EndIf}
    StrCpy $R2 $R5 1 -1
    ${If} "$R2" == ";"
      StrCpy $R5 $R5 -1
    ${EndIf}
  ${EndIf}
  
  WriteRegExpandStr ${RegLoc} "Environment" "${EnvVarName}" $R5
  
  ; Notify all windows of environment block change
  SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000
!macroend

Function EnvVarUpdate
  Push $R0
  Push $R1
  Push $R2
  Push $R3
  Push $R4
  Push $R5
  
  Exch 5
  Exch $R4 ; EnvVarName
  Exch 4
  Exch $R3 ; Action
  Exch 3
  Exch $R2 ; Regloc
  Exch 2
  Exch $R1 ; PathString
  Exch
  Exch $R0 ; ResultVar
  
  ${If} $R2 == "HKLM"
    !insertmacro _EnvVarUpdate_Modify HKLM $R4 $R1 HKLM $R3
  ${ElseIf} $R2 == "HKCU"
    !insertmacro _EnvVarUpdate_Modify HKCU $R4 $R1 HKCU $R3
  ${EndIf}
  
  StrCpy $R0 "0"
  
  Pop $R5
  Pop $R4
  Pop $R3
  Pop $R2
  Pop $R1
  Exch $R0
FunctionEnd

Function un.EnvVarUpdate
  Push $R0
  Push $R1
  Push $R2
  Push $R3
  Push $R4
  Push $R5
  
  Exch 5
  Exch $R4 ; EnvVarName
  Exch 4
  Exch $R3 ; Action
  Exch 3
  Exch $R2 ; Regloc
  Exch 2
  Exch $R1 ; PathString
  Exch
  Exch $R0 ; ResultVar
  
  ${If} $R2 == "HKLM"
    !insertmacro _EnvVarUpdate_Modify HKLM $R4 $R1 HKLM $R3
  ${ElseIf} $R2 == "HKCU"
    !insertmacro _EnvVarUpdate_Modify HKCU $R4 $R1 HKCU $R3
  ${EndIf}
  
  StrCpy $R0 "0"
  
  Pop $R5
  Pop $R4
  Pop $R3
  Pop $R2
  Pop $R1
  Exch $R0
FunctionEnd

!verbose push
!verbose 3
!endif ; _EnvVarUpdate_nsh
!verbose pop