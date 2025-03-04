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

cd $(mktemp -d)
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64-${VERSION}.zip" -o "awscliv2.zip"
unzip awscliv2.zip
./aws/install
