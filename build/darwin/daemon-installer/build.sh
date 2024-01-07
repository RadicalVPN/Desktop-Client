#!/bin/bash

SIGN_CERT=""
VERSION="0.0.0"
while getopts ":c:v:" opt; do
  case $opt in
    c) SIGN_CERT="$OPTARG"
    ;;
	v) VERSION="$OPTARG"
	;;
  esac
done


mkdir -p bin

plutil -replace SMPrivilegedExecutables -xml \
        "<dict> \
      		<key>com.radicalvpn.booter.helper</key> \
      		<string>identifier com.radicalvpn.booter.helper and certificate leaf[subject.OU] = ${SIGN_CERT}</string> \
      	</dict>" "RadicalVPN-Installer.plist"

plutil -replace CFBundleShortVersionString -xml "<string>${VERSION}</string>" "RadicalVPN-Installer.plist"
plutil -replace CFBundleVersion -xml "<string>${VERSION}</string>" "RadicalVPN-Installer.plist"

cc -m64 -framework Foundation \
		-mmacosx-version-min=10.6 \
		-D IS_INSTALLER=0 \
		-framework ServiceManagement \
		-framework Security \
		-Xlinker -sectcreate -Xlinker __TEXT -Xlinker __info_plist -Xlinker "RadicalVPN-Installer.plist" \
		daemon-installer.c \
		-o "bin/RadicalVPN-Installer"

mkdir -p "bin/RadicalVPN-Installer.app/Contents/Library/LaunchServices"
mkdir -p "bin/RadicalVPN-Installer.app/Contents/MacOS" 
cp "bin/RadicalVPN-Installer" "bin/RadicalVPN-Installer.app/Contents/MacOS" 
cp "RadicalVPN-Installer.plist" "bin/RadicalVPN-Installer.app/Contents/Info.plist" 