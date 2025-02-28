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

wget https://github.com/jqlang/jq/releases/download/jq-${VERSION}/jq-linux64

DOWNLOAD_SHA=$(openssl sha1 -sha256 jq-linux64 | awk '{print $2}')
if [[ "${SHA}" != "${DOWNLOAD_SHA}" ]]; then
    echo "Downloaded checksum (${DOWNLOAD_SHA}) does not match expected value: ${SHA}"
    exit 1
fi

chmod +x jq-linux64
mv jq-linux64 /usr/local/bin/jq
