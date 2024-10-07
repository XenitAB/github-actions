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

echo "Adding keys to the keyring..."
mkdir -p /etc/apt/keyrings
curl -sLS https://packages.microsoft.com/keys/microsoft.asc |
  gpg --dearmor | tee /etc/apt/keyrings/microsoft.gpg > /dev/null
echo "Keys added to the keyring, setting permissions..."
chmod go+r /etc/apt/keyrings/microsoft.gpg

#AZ_DIST=$(lsb_release -cs)
AZ_DIST=$(grep -ioP '^VERSION_CODENAME=\K.+' /etc/os-release)
ARCHITECTURE=$(dpkg --print-architecture)
echo "Adding sources to the sources list, DIST=${AZ_DIST} and ARCH=..."

echo "Types: deb
URIs: https://packages.microsoft.com/repos/azure-cli/
Suites: ${AZ_DIST}
Components: main
Architectures: ${ARCHITECTURE}
Signed-by: /etc/apt/keyrings/microsoft.gpg" | tee /etc/apt/sources.list.d/azure-cli.sources

echo "Sources added to the sources list, updating apt and installing AZ CLI..."
apt-get update
apt-get install -y azure-cli=${VERSION}-1~${AZ_DIST}

echo "AZ CLI installed..."
az version

echo "Adding DEVOPS extension..."
az extension add --yes --allow-preview false --upgrade --name azure-devops
echo "Adding MANAGEMENTPARTNER extension..."
az extension add --yes --allow-preview false --upgrade --name managementpartner
echo "AZ CLI installation complete."