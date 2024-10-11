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
wget -nv https://github.com/aquasecurity/tfsec/releases/download/${VERSION}/tfsec-linux-${ARCHITECTURE}

chmod +x tfsec-linux-${ARCHITECTURE}
mv tfsec-linux-${ARCHITECTURE} /usr/local/bin/tfsec
