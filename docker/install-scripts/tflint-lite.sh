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
wget -nv https://github.com/terraform-linters/tflint/releases/download/${VERSION}/tflint_linux_${ARCHITECTURE}.zip

unzip tflint_linux_${ARCHITECTURE}.zip
rm tflint_linux_${ARCHITECTURE}.zip
mv tflint /usr/local/bin/tflint
mkdir -p /work/.tflint.d