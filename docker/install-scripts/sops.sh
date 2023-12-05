#!/bin/bash
set -e

while [ $# -gt 0 ]; do
  case "$1" in
    --version=*)
      VERSION="${1#*=}"
      ;;
    --sha=*)
      SHA="${1#*=}"
      ;;
    *)
      echo "Error: Invalid argument."
      exit 1
  esac
  shift
done

wget https://github.com/getsops/sops/releases/download/${VERSION}/sops-${VERSION}.linux.amd64

DOWNLOAD_SHA=$(openssl sha1 -sha256 sops-${VERSION}.linux.amd64 | awk '{print $2}')
if [[ "${SHA}" != "${DOWNLOAD_SHA}" ]]; then
    echo "Downloaded checksum (${DOWNLOAD_SHA}) does not match expected value: ${SHA}"
    exit 1
fi

chmod +x sops-${VERSION}.linux.amd64
mv sops-${VERSION}.linux.amd64 /usr/local/bin/sops
