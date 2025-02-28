#!/bin/bash
set -e

ACTION=$1
ENVIRONMENT=$2
GIT_REPO=$3
GIT_BRANCH=$4
TEMPLATE_FILE=$5
RESOURCE_GROUP=$6


ENVIRONMENT_FILE="/tmp/${ENVIRONMENT}.env"

build () {
  git clone -c advice.detachedHead=false --single-branch --branch "${GIT_BRANCH}" ${GIT_REPO} remote-templates
  AZURE_RESOURCE_GROUP_NAME=$RESOURCE_GROUP packer build remote-templates/${TEMPLATE_FILE}
}

envup() {
  if [ -f ${ENVIRONMENT_FILE} ]; then
    set -a
    source ${ENVIRONMENT_FILE}
    set +a
  fi
}

envup

case $ACTION in
  build )
    build
    ;;
esac
