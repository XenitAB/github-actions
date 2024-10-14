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
wget -nv https://github.com/getsops/sops/releases/download/${VERSION}/sops-${VERSION}.linux.${ARCHITECTURE}

chmod +x sops-${VERSION}.linux.${ARCHITECTURE}
mv sops-${VERSION}.linux.${ARCHITECTURE} /usr/local/bin/sops
