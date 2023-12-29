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

echo "[*] Cloning wireguard-go.."
git clone https://git.zx2c4.com/wireguard-go/ --depth 1

echo "[*] Building wireguard-go.."
cd wireguard-go
make -j${CPU_COUNT}

cd ..

echo "[*] Cloning wireguard-tools.."
git clone https://git.zx2c4.com/wireguard-tools/ --depth 1

echo "[*] Building wireguard-tools.."

cd wireguard-tools/src
make

cd ../..

if [ ! -f "../Bash/bash" ]; then
    echo "[*] Building Bash.."

    if [ -d "bash" ]; then
        rm -rf bash
    fi

    git clone https://git.savannah.gnu.org/git/bash.git --depth 1
    cd bash
    ./configure
    make -j${CPU_COUNT}

    cd ..

    mkdir -p ../Bash
    mv bash/bash ../Bash/bash
fi

cd ..

echo "[*] Moving binaries.."

mkdir -p Wireguard

mv ./.deps/wireguard-go/wireguard-go ./Wireguard/
mv ./.deps/wireguard-tools/src/wg ./Wireguard/
mv ./.deps/wireguard-tools/src/wg-quick/darwin.bash ./Wireguard/wg-quick.bash

echo "[*] Build succeeded!"

