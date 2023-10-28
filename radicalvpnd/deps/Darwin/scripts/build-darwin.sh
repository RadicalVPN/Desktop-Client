#!/bin/sh

# exit on any non 0 return code
set -e

scriptdir=$(dirname "$0")/../.deps

mkdir -p $scriptdir
cd $scriptdir

echo "### Start compiling Wireguard  ###"

echo "[*] Build cleanup.."
rm -rf wireguard-go
rm -rf wireguard-tools

echo "[*] Cloning wireguard-go.."
git clone https://git.zx2c4.com/wireguard-go/

echo "[*] Building wireguard-go.."
cd wireguard-go
make

cd ..

echo "[*] Cloning wireguard-tools.."
git clone https://git.zx2c4.com/wireguard-tools/

echo "[*] Building wireguard-tools.."

cd wireguard-tools/src
make

echo "[*] Moving binaries.."

cd ../../..
mkdir -p Wireguard

mv ./.deps/wireguard-go/wireguard-go ./Wireguard/
mv ./.deps/wireguard-tools/src/wg ./Wireguard/
mv ./.deps/wireguard-tools/src/wg-quick/darwin.bash ./Wireguard/wg-quick.bash



echo "[*] Build succeeded!"

