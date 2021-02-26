#!/bin/bash
set -e

while [ $# -gt 0 ]; do
  case "$1" in
    --ruleset=*)
      RULESET="${1#*=}"
      ;;
    --version=*)
      VERSION="${1#*=}"
      ;;
    --sha=*)
      SHA="${1#*=}"
      ;;
    --user=*)
      USER="${1#*=}"
      ;;
    --group=*)
      GROUP="${1#*=}"
      ;;
    *)
      echo "Error: Invalid argument."
      exit 1
  esac
  shift
done

wget https://github.com/terraform-linters/tflint-ruleset-${RULESET}/releases/download/${VERSION}/tflint-ruleset-${RULESET}_linux_amd64.zip

DOWNLOAD_SHA=$(openssl sha1 -sha256 tflint-ruleset-${RULESET}_linux_amd64.zip | awk '{print $2}')
if [[ "${SHA}" != "${DOWNLOAD_SHA}" ]]; then
    echo "Downloaded checksum (${DOWNLOAD_SHA}) does not match expected value: ${SHA}"
    exit 1
fi

unzip tflint-ruleset-${RULESET}_linux_amd64.zip
rm tflint-ruleset-${RULESET}_linux_amd64.zip
mkdir -p /home/${USER}/.tflint.d/plugins/
mv tflint-ruleset-${RULESET} /home/${USER}/.tflint.d/plugins/tflint-ruleset-${RULESET}
chown -R ${USER}:${GROUP} /home/${USER}/.tflint.d