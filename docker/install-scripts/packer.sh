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

wget https://releases.hashicorp.com/packer/${VERSION}/packer_${VERSION}_linux_amd64.zip
DOWNLOAD_SHA=$(openssl sha1 -sha256 packer_${VERSION}_linux_amd64.zip | awk '{print $2}')
if [[ "${SHA}" != "${DOWNLOAD_SHA}" ]]; then
    echo "Downloaded checksum (${DOWNLOAD_SHA}) does not match expected value: ${SHA}"
    exit 1
fi

unzip packer_${VERSION}_linux_amd64.zip
rm packer_${VERSION}_linux_amd64.zip
chmod +x packer
mv packer /usr/local/bin/packer
