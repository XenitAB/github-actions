#!/bin/bash
set -e

while [ $# -gt 0 ]; do
  case "$1" in
    --latest-terraform-version=*)
      LATEST_TERRAFORM_VERSION="${1#*=}"
      ;;
    --tfenv-version=*)
      TFENV_VERSION="${1#*=}"
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

git clone -b ${TFENV_VERSION} https://github.com/tfutils/tfenv.git /opt/tfenv
ln -s /opt/tfenv/bin/* /usr/local/bin

tfenv list-remote | grep -v "-" | grep ${LATEST_TERRAFORM_VERSION} -A 4 | xargs -t -I % tfenv install %
tfenv use ${LATEST_TERRAFORM_VERSION}

chown -R ${USER}:${GROUP} /opt/tfenv
