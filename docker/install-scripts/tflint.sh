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
    --user=*)
      USER="${1#*=}"
      ;;
    *)
      echo "Error: Invalid argument."
      exit 1
  esac
  shift
done

wget https://github.com/terraform-linters/tflint/releases/download/${VERSION}/tflint_linux_amd64.zip

DOWNLOAD_SHA=$(openssl sha1 -sha256 tflint_linux_amd64.zip | awk '{print $2}')
if [[ "${SHA}" != "${DOWNLOAD_SHA}" ]]; then
    echo "Downloaded checksum (${DOWNLOAD_SHA}) does not match expected value: ${SHA}"
    exit 1
fi

unzip tflint_linux_amd64.zip
rm tflint_linux_amd64.zip
mv tflint /usr/local/bin/tflint
mkdir -p /home/${USER}/.tflint.d