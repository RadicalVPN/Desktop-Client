#!/bin/sh

# exit on any non 0 return code
set -e

scriptdir=$(dirname "$0")/../.deps
CPU_COUNT=$(sysctl -n hw.ncpu)

mkdir -p $scriptdir
cd $scriptdir

echo "### Start compiling Wireguard  ###"

echo "[*] Build cleanup.."
rm -rf wireguard-go
rm -rf wireguard-tools
rm -rf bash

echo "[*] Cloning wireguard-go.."
git clone https://git.zx2c4.com/wireguard-go/

echo "[*] Building wireguard-go.."
cd wireguard-go
make -j${CPU_COUNT}

cd ..

echo "[*] Cloning wireguard-tools.."
git clone https://git.zx2c4.com/wireguard-tools/

echo "[*] Building wireguard-tools.."

cd wireguard-tools/src
make

cd ../..

git clone https://git.savannah.gnu.org/git/bash.git --depth 1
cd bash
./configure
make -j${CPU_COUNT}

cd ../..

echo "[*] Moving binaries.."

mkdir -p Wireguard

mv ./.deps/wireguard-go/wireguard-go ./Wireguard/
mv ./.deps/wireguard-tools/src/wg ./Wireguard/
mv ./.deps/wireguard-tools/src/wg-quick/darwin.bash ./Wireguard/wg-quick.bash

mkdir -p Bash
mv ./.deps/bash/bash ./Bash

echo "[*] Build succeeded!"

