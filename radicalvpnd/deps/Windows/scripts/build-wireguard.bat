@ECHO OFF

setlocal
set SCRIPTDIR=%~dp0
set WG_VERSION=v0.5.3

echo ### Start compiling Wireguard  ###

if exist "%SCRIPTDIR%..\WireGuard\x86_64" (
  del /f /q /s "%SCRIPTDIR%..\WireGuard\x86_64\*" >nul 2>&1 || exit /b 1
)

if not exist "%SCRIPTDIR%..\.deps\wireguard-windows\.deps\prepared" (
  if not exist "%SCRIPTDIR%..\.deps" (
    mkdir "%SCRIPTDIR%..\.deps" || exit /b 1
    cd "%SCRIPTDIR%..\.deps" 	|| exit /b 1
  )

  if exist "%SCRIPTDIR%..\.deps\wireguard-windows" (
    rd /s /q "%SCRIPTDIR%..\.deps\wireguard-windows" || exit /b 1
    sleep 2
  )

  cd "%SCRIPTDIR%..\.deps"

  echo [*] Cloning wireguard-windows..
  git clone git://git.zx2c4.com/wireguard-windows || exit /b 1
  cd wireguard-windows || exit /b 1

  echo [*] Checking out version [%WG_VERSION%]...
  git checkout %WG_VERSION% >nul 2>&1 || exit /b 1
    echo [*] Building wireguard-windows from NEW source..
) else (
  echo [*] Building wireguard-windows from OLD source..
  cd "%SCRIPTDIR%..\.deps\wireguard-windows" 	|| exit /b 1
)

call build.bat
if not %errorlevel%==0 (
    echo [!] ERROR: Failed to compile wireguard-windows
    exit /b 1
)

echo [*] WireGuard compiled! :=)

if not exist "%SCRIPTDIR%..\WireGuard" 			mkdir "%SCRIPTDIR%..\WireGuard" 		|| exit /b 1
if not exist "%SCRIPTDIR%..\WireGuard\x86_64" 	mkdir "%SCRIPTDIR%..\WireGuard\x86_64"	|| exit /b 1

copy /y "%SCRIPTDIR%..\.deps/wireguard-windows\amd64\wg.exe" 		"%SCRIPTDIR%..\WireGuard\x86_64\wg.exe" 		>nul 2>&1 || exit /b 1
copy /y "%SCRIPTDIR%..\.deps\wireguard-windows\amd64\wireguard.exe" "%SCRIPTDIR%..\WireGuard\x86_64\wireguard.exe" 	>nul 2>&1 || exit /b 1
