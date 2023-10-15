SetCompressor lzma

!include "MUI.nsh"
!include "LogicLib.nsh"
!include "StrFunc.nsh"
!include "x64.nsh"
!include "WinVer.nsh"
!include "winmessages.nsh"

!define NAME "RadicalVPN Client"
!define MUI_WELCOMEPAGE_TITLE "Welcome to the ${NAME} Setup Wizard!"

!define MUI_HEADERIMAGE
!define MUI_HEADERIMAGE_RIGHT
!define MUI_HEADERIMAGE_BITMAP "logo.bmp"

!define MUI_ICON "logo.ico"
!define MUI_UNICON "logo.ico"

#!define MUI_FINISHPAGE_RUN_FUNCTION ExecAppFile

LicenseForceSelection checkbox "I Agree"


!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE License.txt
!insertmacro MUI_PAGE_INSTFILES


!insertmacro MUI_UNPAGE_WELCOME
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_UNPAGE_FINISH

!insertmacro MUI_LANGUAGE "English"

OutFile "${OUT_FILE}"
Name "${NAME}"
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

Section "${NAME}" RadicalVPN
  SetRegView 64
  SetOutPath "$INSTDIR"

  DetailPrint "Installing RadicalVPN Daemon..."
  nsExec::ExecToLog '"$SYSDIR\sc.exe" create "RadicalVPN" binPath= "\"$INSTDIR\RadicalVPN.exe\"" start= auto'

SectionEnd

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