SetCompressor lzma

!include "MUI.nsh"
!include "LogicLib.nsh"
!include "StrFunc.nsh"
!include "x64.nsh"
!include "WinVer.nsh"
!include "winmessages.nsh"

!define NAME "RadicalVPN"
!define MUI_WELCOMEPAGE_TITLE "Welcome to the ${NAME} Setup Wizard!"

!define MUI_HEADERIMAGE
!define MUI_HEADERIMAGE_RIGHT
!define MUI_HEADERIMAGE_BITMAP "logo.bmp"

!define MUI_ICON "logo.ico"
!define MUI_UNICON "logo.ico"

!define APP_RUN_PATH "$INSTDIR\gui\RadicalVPN.exe"

#!define MUI_FINISHPAGE_RUN_FUNCTION ExecAppFile

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE License.txt
!insertmacro MUI_PAGE_INSTFILES


!insertmacro MUI_UNPAGE_WELCOME
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_UNPAGE_FINISH


Name "${NAME}"
OutFile "${OUT_FILE}"
InstallDir "$PROGRAMFILES64\${NAME}"

RequestExecutionLevel admin

; variables
Var HEADLINE_FONT


!macro COMMON_INIT
  ; install for  'all users'
  SetShellVarContext all

  SetRegView 64
  DetailPrint "Running on architecture: x86_64"
!macroend

Function .onInit
  !insertmacro COMMON_INIT

  CreateFont $HEADLINE_FONT "$(^Font)" "12" "600"

  Call IsOSSupported

  ClearErrors
FunctionEnd

Function un.onInit
  !insertmacro COMMON_INIT
FunctionEnd

Function IsOSSupported
    ${If} ${AtLeastWin7}
        goto archcheck
    ${EndIf}
    MessageBox MB_ICONSTOP|MB_OK "Unsupported Windows Version.$\nRadicalVPN can only be installed on Windows 7 and above."
    Quit
archcheck:
    ${If} ${RunningX64}
        goto end
    ${EndIf}
    MessageBox MB_ICONSTOP|MB_OK "Unsupported architecture.$\nRadicalVPN can only be installed on 64-bit Windows."
    Quit
end:
FunctionEnd

!define MUI_FINISHPAGE_NOAUTOCLOSE
!define MUI_FINISHPAGE_RUN "$INSTDIR\RadicalVPN.exe"
!define MUI_FINISHPAGE_RUN_TEXT "Run ${NAME} now"
!define MUI_FINISHPAGE_RUN_FUNCTION exec_app

; Checkbox on finish page: create shortcut on desktop
; using unused 'readme' check box for this
!define MUI_FINISHPAGE_SHOWREADME ""
!define MUI_FINISHPAGE_SHOWREADME_NOTCHECKED
!define MUI_FINISHPAGE_SHOWREADME_TEXT "Create a desktop shortcut"
!define MUI_FINISHPAGE_SHOWREADME_FUNCTION finish_page_creation
Function finish_page_creation
CreateShortcut "$DESKTOP\RadicalVPN.lnk" "${APP_RUN_PATH}"
FunctionEnd

LicenseForceSelection checkbox "I Agree"

!define MUI_STARTMENUPAGE_REGISTRY_ROOT "HKLM"
!define MUI_STARTMENUPAGE_REGISTRY_KEY "Software\${NAME}"
!define MUI_STARTMENUPAGE_REGISTRY_VALUENAME "Start Menu Folder"

!define MUI_ABORTWARNING

!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_WELCOME
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_UNPAGE_FINISH

!insertmacro MUI_LANGUAGE "English"

Function exec_app
    Exec "${APP_RUN_PATH}"
    Sleep 500

    StrCpy $R1 0
    ${While} $R1 < 50
        IntOp $R1 $R1 + 1
        System::Call user32::GetForegroundWindow()i.r0

        ${If} $0 != $hwndparent
            Return
        ${EndIf}

        Sleep 100
    ${EndWhile}

FunctionEnd


Section "${NAME}" RadicalVPN
  SetRegView 64
  SetOutPath "$INSTDIR"

  DetailPrint "Installing RadicalVPN..."
  File /r "${SOURCE_DIR}\*.*"

  CreateDirectory "$SMPROGRAMS\$StartMenuFolder"
  CreateShortCut "$SMPROGRAMS\$StartMenuFolder\${NAME}.lnk" "$INSTDIR\gui\RadicalVPN.exe"

  nsExec::ExecToLog '"$SYSDIR\sc.exe" stop "RadicalVPN Daemon"'

  DetailPrint "Installing RadicalVPN Daemon Service..."
  nsExec::ExecToLog '"$SYSDIR\sc.exe" create "RadicalVPN Daemon" binPath= "\"$INSTDIR\radicalvpnd.exe\"" start= auto'
  nsExec::ExecToLog '"$SYSDIR\sc.exe" sdset "RadicalVPN Daemon" "D:(A;;CCLCSWRPWPDTLOCRRC;;;SY)(A;;CCDCLCSWRPWPDTLOCRSDRCWDWO;;;BA)(A;;CCLCSWLOCRRC;;;IU)(A;;CCLCSWLOCRRC;;;SU)(A;;RPWPDTLO;;;S-1-1-0)"'

  DetailPrint "Starting RadicalVPN Daemon Service..."
  nsExec::ExecToLog '"$SYSDIR\sc.exe" start "RadicalVPN Daemon"'

SectionEnd
