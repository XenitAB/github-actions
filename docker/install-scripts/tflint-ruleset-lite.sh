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
   *)
      echo "Error: Invalid argument."
      exit 1
  esac
  shift
done

ARCHITECTURE=$(dpkg --print-architecture)
wget -nv https://github.com/terraform-linters/tflint-ruleset-${RULESET}/releases/download/${VERSION}/tflint-ruleset-${RULESET}_linux_${ARCHITECTURE}.zip

unzip tflint-ruleset-${RULESET}_linux_${ARCHITECTURE}.zip
rm tflint-ruleset-${RULESET}_linux_${ARCHITECTURE}.zip
mkdir -p /work/.tflint.d/plugins/
mv tflint-ruleset-${RULESET} /work/.tflint.d/plugins/tflint-ruleset-${RULESET}
