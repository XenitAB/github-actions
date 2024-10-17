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
wget -nv https://github.com/cli/cli/releases/download/v${VERSION}/gh_${VERSION}_linux_${ARCHITECTURE}.tar.gz

tar xzvf gh_${VERSION}_linux_${ARCHITECTURE}.tar.gz
mv gh_${VERSION}_linux_${ARCHITECTURE}/bin/gh /usr/local/bin/gh
