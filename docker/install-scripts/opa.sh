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
      printf "***************************\n"
      printf "* Error: Invalid argument.*\n"
      printf "***************************\n"
      exit 1
  esac
  shift
done

wget https://github.com/open-policy-agent/opa/releases/download/${VERSION}/opa_linux_amd64

DOWNLOAD_SHA=$(openssl sha1 -sha256 opa_linux_amd64 | awk '{print $2}')
if [[ "${SHA}" != "${DOWNLOAD_SHA}" ]]; then
    echo "Downloaded checksum (${DOWNLOAD_SHA}) does not match expected value: ${SHA}"
    exit 1
fi

chmod +x opa_linux_amd64
mv opa_linux_amd64 /usr/local/bin/opa
