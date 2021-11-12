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

#pip --no-cache-dir install azure-cli==${VERSION}
curl -sL https://aka.ms/InstallAzureCLIDeb | bash
az extension add --name azure-devops
az extension add --name managementpartner
