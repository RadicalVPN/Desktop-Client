SCRIPT_DIR="$( pwd )"
SIGN_CERT=""
# reading version info from arguments
while getopts ":c:" opt; do
  case $opt in
    c) SIGN_CERT="$OPTARG"
    ;;
  esac
done

echo "[i] Signing with cert: '${SIGN_CERT}'"

cd "${SCRIPT_DIR}/build/darwin"

ThirdPartyBinaries=(
"_image/RadicalVPN.app/Contents/MacOS/Wireguard/wg"
"_image/RadicalVPN.app/Contents/MacOS/Wireguard/wireguard-go"
"_image/RadicalVPN.app/Contents/MacOS/RadicalVPN-Installer.app/Contents/Library/LaunchServices/com.radicalvpn.booter.helper"
)

Binaries=(
"_image/RadicalVPN.app/Contents/MacOS/RadicalVPN"
"_image/RadicalVPN.app/Contents/MacOS/RadicalVPN Daemon"
"_image/RadicalVPN.app/Contents/MacOS/RadicalVPN-Installer.app/Contents/MacOS/RadicalVPN-Installer"
"_image/RadicalVPN.app/Contents/MacOS/RadicalVPN-Installer.app"
"_image/RadicalVPN.app"
)

echo "[+] Signing third-party binaries..."
for f in "${ThirdPartyBinaries[@]}";
do
  echo "    signing: [" $f "]";
  codesign --verbose=4 --force --sign "${SIGN_CERT}" --options runtime "$f"
done


echo "[+] Signing compiled binaries..."
for f in "${Binaries[@]}";
do
  echo "    signing: [" $f "]";
  codesign --verbose=4 --force --deep --sign "${SIGN_CERT}" --entitlements security.plist --options runtime "$f"
done

cd "${SCRIPT_DIR}"