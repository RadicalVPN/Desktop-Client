@echo off

SET NSIS="C:\Program Files (x86)\NSIS\makensis.exe"
set VERSION=%1

set SCRIPTDIR=%~dp0
SET OUT_DIR=%SCRIPTDIR%bin
set TMP_DIR=%OUT_DIR%\temp

:: this is a hack to get the commit hash on windows
for /f %%i in ('git rev-parse HEAD') do set COMMIT_HASH=%%i

if not exist %NSIS% (
    echo [!] NSIS not found [%NSIS%]
	echo [!] Install NSIS [https://nsis.sourceforge.io/] and try again.
	goto :error
)

call :build_wireguard
call :build_daemon
call :build_gui

call :copy_files
call :build_nsis_installer

goto :success


:copy_files
    echo [*] Copying files...

    ::create only if exists
    IF not exist %TMP_DIR% (
        echo [*] Creating temp directory [%TMP_DIR%]
        mkdir %TMP_DIR%
    ) else (
        echo [*] Temp directory already exists [%TMP_DIR%]
        rmdir /s /q %TMP_DIR%
    )

    xcopy /e /i "gui\dist\win-unpacked" "%TMP_DIR%\gui"
    xcopy /e /i "radicalvpnd\deps\Windows\WireGuard\x86_64" "%TMP_DIR%\wireguard"

    copy /b "radicalvpnd\radicalvpnd.exe" "%TMP_DIR%\radicalvpnd.exe"
    copy /b "build\windows\logo.ico" "%TMP_DIR%\logo.ico"

    IF not ERRORLEVEL 0 (
        echo [!] Failed to copy files.
        exit /b 1
    )

    goto :eof

:build_nsis_installer
    echo [*] Building NSIS installer...

    mkdir %OUT_DIR%

    set FILE_OUT="%OUT_DIR%/RadicalVPN-Setup.exe"
    %NSIS% /DPRODUCT_VERSION="0.0.1" /DOUT_FILE=%FILE_OUT% /DSOURCE_DIR=%TMP_DIR% "build/windows/RadicalVPN Client.nsi"
    IF not ERRORLEVEL 0 (
		echo [!] Failed to build NSIS installer.
		exit /b 1
	)

    goto :eof

:build_wireguard
    call "radicalvpnd\deps\Windows\scripts\build-wireguard.bat"
    echo %SCRIPTDIR%

    goto :eof

:build_daemon
    call cd radicalvpnd
    call go build -ldflags "-X 'radicalvpnd/version.version=%VERSION%' -X 'radicalvpnd/version.commitHash=%COMMIT_HASH%'"  .
    call cd ..

    goto :eof

:build_gui
    call cd gui
    call npm install
    call npm run build
    call cd ..

    goto :eof


:success
	echo [*] Installer created successfully.
	exit /b 0


:error
	echo [!] Installation failed: #%errorlevel%.
	exit /b %errorlevel%