#!/bin/bash

VERSION="0.0.1"
SIGN_CERT=""
while getopts ":c:v:" opt; do
  case $opt in
    c) SIGN_CERT="$OPTARG"
    ;;
    v) VERSION="$OPTARG"
    ;;
  esac
done

PLIST_INFO="Daemon-Boot-Info.plist"
PLIST_LAUNCHD="Daemon-Boot-Launchd.plist"

echo "SIGN_CERT: ${SIGN_CERT}"

plutil -replace SMAuthorizedClients -xml \
        "<array> \
          <string>identifier com.radicalvpn.installer and certificate leaf[subject.OU] = ${SIGN_CERT}</string>\
        </array>" "${PLIST_INFO}"

plutil -replace CFBundleShortVersionString -xml "<string>${VERSION}</string>" "${PLIST_INFO}"
plutil -replace CFBundleVersion -xml "<string>${VERSION}</string>" "${PLIST_INFO}"

go build -ldflags='-extldflags "-sectcreate __TEXT __info_plist '${PLIST_INFO}' -sectcreate __TEXT __launchd_plist '${PLIST_LAUNCHD}'"' .

