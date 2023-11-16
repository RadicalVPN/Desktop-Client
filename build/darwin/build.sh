#!/bin/bash

# exit when any command fails
set -e


SIGN_CERT=""
while getopts ":c:" opt; do
  case $opt in
    c) SIGN_CERT="$OPTARG"
    ;;
  esac
done

if [ -z "${SIGN_CERT}" ]; then
  echo "Usage:"
  echo "    $0 -c <APPLE_DEVELOPER_CERTIFICATE>"
  echo "    Example: $0 -c WXXXXXXXXN"
  exit 1
fi


SCRIPT_DIR="$( pwd )"
DAEMON_PATH="${SCRIPT_DIR}/radicalvpnd"
BUILD_PATH="${SCRIPT_DIR}/build/darwin"

if [ ! -d ${DAEMON_PATH} ]; then
  echo "[!] ERROR: Execute this script from the root directory"
  exit 1
fi

if [ ! -d ${BUILD_PATH} ]; then
  echo "[!] ERROR: Failed to find build directory: ${BUILD_PATH}"
  exit 1
fi

echo "[+] Build RadicalVPN Desktop for Darwin (MacOS).."

echo "[+] Build Daemon.."
cd ${DAEMON_PATH}
go build .

echo "[+] Build Daemon Dependencies.."
cd deps/Darwin/scripts
./build-darwin.sh

echo "[+] Build Electron Frontend.."
cd ../../../../gui
npm i
npm run build


echo "[+] Build Daemon Installer.."
cd ${BUILD_PATH}/daemon-installer
bash build.sh -c ${SIGN_CERT}

echo "[+] Build Daemon Booter.."
cd ${BUILD_PATH}/daemon-boot
bash build.sh -c ${SIGN_CERT}
cd ${SCRIPT_DIR}

echo ======================================================
echo =================== Preparing DMG ====================
echo ======================================================

COMPILEDFOLDER="mac"
ARCH="$( uname -m )"
if [ ${ARCH} = "arm64" ]; then
  COMPILEDFOLDER="mac-arm64"
  echo "[+] Frontend Path ${COMPILEDFOLDER}. Arch: ${ARCH}"
else
  echo "[+] Frontend Path ${COMPILEDFOLDER}. Default Arch (${ARCH})"
fi

FRONTEND_COMPILED_NAME="RadicalVPN.app"
FRONTEND_BINARY_PATH="${SCRIPT_DIR}/gui/dist/${COMPILEDFOLDER}/${FRONTEND_COMPILED_NAME}"

if [ ! -d ${FRONTEND_BINARY_PATH} ]; then
  echo "[!] ERROR: Failed to find frontend binary: ${FRONTEND_BINARY_PATH}"
  exit 1
fi

mkdir -p "${BUILD_PATH}/_image"

echo "[+] Copying frontend binary to build folder.."
cp -a "${FRONTEND_BINARY_PATH}" "${BUILD_PATH}/_image/"

echo "[+] Copying daemon binary to build folder.."
cp -a "${DAEMON_PATH}/radicalvpnd" "${BUILD_PATH}/_image/RadicalVPN.app/Contents/MacOS/RadicalVPN Daemon"
chmod 0700 "${BUILD_PATH}/_image/RadicalVPN.app/Contents/MacOS/RadicalVPN Daemon"

echo "[+] Copying wireguard binaries.."
cp -a "${DAEMON_PATH}/deps/Darwin/Wireguard" "${BUILD_PATH}/_image/RadicalVPN.app/Contents/MacOS/Wireguard/"
chmod -R 0700 "${BUILD_PATH}/_image/RadicalVPN.app/Contents/MacOS/Wireguard"

echo "[+] Copying Bash binaries.."
mkdir -p "${BUILD_PATH}/_image/RadicalVPN.app/Contents/MacOS/Bash"
cp -a "${DAEMON_PATH}/deps/Darwin/Bash/bash" "${BUILD_PATH}/_image/RadicalVPN.app/Contents/MacOS/Bash/bash"
chmod 0700 "${BUILD_PATH}/_image/RadicalVPN.app/Contents/MacOS/Bash/bash"

echo "[+] Copying daemon installer binary.."
cp -a "${BUILD_PATH}/daemon-installer/bin/RadicalVPN-Installer.app" "${BUILD_PATH}/_image/RadicalVPN.app/Contents/MacOS"

echo "[+] Copying daemon booter binary.."
#make sure to use the package name for the daemon booter -> https://developer.apple.com/documentation/servicemanagement/1431078-smjobbless
cp -a "${BUILD_PATH}/daemon-boot/daemon_boot" "${BUILD_PATH}/_image/RadicalVPN.app/Contents/MacOS/RadicalVPN-Installer.app/Contents/Library/LaunchServices/com.radicalvpn.booter.helper"

echo "[+] Copying DMG Background.."
mkdir -p "${BUILD_PATH}/_image/.background"
cp -a "${BUILD_PATH}/assets/bg-dmg.png" "${BUILD_PATH}/_image/.background/bg-dmg.png"


echo ======================================================
echo ================= Signing Binaries ===================
echo ======================================================

bash ${BUILD_PATH}/sign.sh -c ${SIGN_CERT}


echo ======================================================
echo =================== Crafting DMG =====================
echo ======================================================

PATH_COMPILED_FOLDER=${BUILD_PATH}/_compiled
PATH_DMG_FILE="${PATH_COMPILED_FOLDER}/RadicalVPN.dmg"
if [ ${ARCH} != "x86_64" ]; then
  PATH_DMG_FILE="${PATH_COMPILED_FOLDER}/RadicalVPN.dmg"
fi


PATH_TMP_DMG_FILE="${PATH_COMPILED_FOLDER}/radicalvpn.temp.dmg"
APPLICATION_NAME="RadicalVPN.app"

mkdir -p ${PATH_COMPILED_FOLDER} 

rm -f ${PATH_DMG_FILE}
rm -f ${PATH_TMP_DMG_FILE}

IMAGE_TITLE="RadicalVPN"
IMAGE_SIZE=1009600
IMAGE_SOURCE="${BUILD_PATH}/_image"

echo "[+] Creating a new temp DMG ..."
hdiutil create -srcfolder "${IMAGE_SOURCE}" -volname "${IMAGE_TITLE}" -fs HFS+ \
      -fsargs "-c c=64,a=16,e=16" -format UDRW -size ${IMAGE_SIZE}k ${PATH_TMP_DMG_FILE}

echo "[+] Mounting temp DMG ..."
IMAGE_DEVICE=$(hdiutil attach -readwrite -noverify -noautoopen ${PATH_TMP_DMG_FILE} | \
         egrep '^/dev/' | sed 1q | awk '{print $1}')

echo "[i] Mpunted Temp DMG: ${IMAGE_DEVICE}"

sleep 5

echo "[+] Writing content to tmp DMG Image ..."
echo '
   tell application "Finder"
     tell disk "'${IMAGE_TITLE}'"
           open
           set current view of container window to icon view
           set toolbar visible of container window to false
           set statusbar visible of container window to false
           set the bounds of container window to {200, 200, 758, 680}
           set theViewOptions to the icon view options of container window
           set arrangement of theViewOptions to not arranged
           set icon size of theViewOptions to 108
           set background picture of theViewOptions to file ".background:bg-dmg.png"
           make new alias file at container window to POSIX file "/Applications" with properties {name:"Applications"}
           set position of item "'${FRONTEND_COMPILED_NAME}'" of container window to {120, 110}
           set position of item "Applications" of container window to {420, 110}
           set position of item ".background" of container window to {120, 500}
           set position of item ".fseventsd" of container window to {420, 500}
           update without registering applications
           delay 3
           close
     end tell
   end tell
' | osascript

echo "[+] chmod tmp DMG image ..."
chmod -Rf go-w "/Volumes/${IMAGE_TITLE}"
sync
sync

echo "[+] detaching tmp DMG image ${IMAGE_DEVICE} ..."
hdiutil detach ${IMAGE_DEVICE}

echo "[+] coverting tmp DMG Image to Final DMG Image ..."
hdiutil convert ${PATH_TMP_DMG_FILE} -format UDZO -imagekey zlib-level=9 -o "${PATH_DMG_FILE}"

echo "[+] Deleting tmp DMG Image..."
rm -f ${PATH_TMP_DMG_FILE}

open ${PATH_COMPILED_FOLDER}
cd ${SCRIPT_DIR}
