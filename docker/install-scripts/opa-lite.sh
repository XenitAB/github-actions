#!/bin/bash
set -e

while [ $# -gt 0 ]; do
  case "$1" in
    --version=*)
      VERSION="${1#*=}"
      ;;
    *)
      echo "Error: Invalid argument."
      exit 1
  esac
  shift
done

ARCHITECTURE=$(dpkg --print-architecture)

wget https://github.com/open-policy-agent/opa/releases/download/${VERSION}/opa_linux_${ARCHITECTURE}_static

chmod +x opa_linux_${ARCHITECTURE}_static
mv opa_linux_${ARCHITECTURE}_static /usr/local/bin/opa
