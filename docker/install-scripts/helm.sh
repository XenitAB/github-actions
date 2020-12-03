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

wget https://get.helm.sh/helm-${VERSION}-linux-amd64.tar.gz

DOWNLOAD_SHA=$(openssl sha1 -sha256 helm-${VERSION}-linux-amd64.tar.gz | awk '{print $2}')
if [[ "${SHA}" != "${DOWNLOAD_SHA}" ]]; then
    echo "Downloaded checksum (${DOWNLOAD_SHA}) does not match expected value: ${SHA}"
    exit 1
fi

tar xzvf helm-${VERSION}-linux-amd64.tar.gz
chmod +x linux-amd64/helm
mv linux-amd64/helm /usr/local/bin/helm
