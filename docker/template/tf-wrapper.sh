#!/bin/bash

# TODO: rename terraform.sh to tf-prepare.sh

set -e

# Function to check if Docker is installed
function check_docker_installed {
    if ! command -v docker &> /dev/null; then
        echo "Error: Docker is not installed. Please install Docker and try again."
        exit 1
    fi
}

# Function to check if Docker container is reachable
function check_docker_container_reachable {
    local image="$1"
    if docker image inspect "$image" &> /dev/null; then
        echo "Docker image '$image' found locally."
        exit 1
    fi
    # if ! docker pull "$image" &> /dev/null; then
    #     echo "Error: Docker container '$image' is not reachable. Please check the image name and your network connection."
    #     exit 1
    # fi
}

# Function to confirm action with the user
function confirm_action {
    local action="$1"
    local target="$2"
    local subscription_name
    subscription_name=$(az account show --query 'name' -o tsv)
    if [ -z "$subscription_name" ]; then
        echo "Error: Unable to retrieve Azure subscription name. Please ensure you are logged into Azure CLI."
        exit 1
    fi
    echo "You are about to run '$action' on subscription '$subscription_name'. Are you sure? (y/n)"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        echo "Action canceled."
        exit 0
    fi
}

# Validate inputs and check Docker before proceeding
if [ -z "$1" ]; then
    echo "Usage: $0 {setup|teardown|prepare|plan|apply|destroy|state-remove|validate|shell} --imageTag <tag> --stateSuffix <suffix>"
    exit 1
fi

# Inputs
IMAGE_TAG="latest"
ARCH=$(uname -m) # "" # amd64, arm64, arm #ARCH=$(uname -m)
STATE_SUFFIX=""

while [[ "$#" -gt 0 ]]; do
    case "$1" in
        --imageTag)
            IMAGE_TAG="$2"
            shift 2
            ;;
        --stateSuffix)
            STATE_SUFFIX="$2"
            shift 2
            ;;
        *)
            COMMAND="$1"
            shift
            ;;
    esac
done

if [ -z "$IMAGE_TAG" ] || [ -z "$STATE_SUFFIX" ]; then
    echo "Error: Missing required arguments --imageTag or --stateSuffix"
    exit 1
fi

# Check if Docker is installed
check_docker_installed

# Define Docker image name
#IMAGE="ghcr.io/xenitab/github-actions/tf-wrapper:$IMAGE_TAG"
IMAGE="test:$IMAGE_TAG"
# Check if Docker container is reachable
check_docker_container_reachable "$IMAGE"

# Proceed with the rest of the script
SUFFIX="tfstate$STATE_SUFFIX"
export PODMAN_USERNS=keep-id

OPA_BLAST_RADIUS=${OPA_BLAST_RADIUS:-50}
RG_LOCATION_SHORT=we
RG_LOCATION_LONG=westeurope
AZURE_CONFIG_DIR=${AZURE_CONFIG_DIR:-"$HOME/.azure"}
TTY_OPTIONS=$( [ -t 0 ] && echo "-it" )

if [ -z "$ENV" ]; then
    echo "Need to set ENV"
    exit 1
fi
if [ -z "$DIR" ]; then
    echo "Need to set DIR"
    exit 1
fi

AZURE_DIR_MOUNT="-v ${AZURE_CONFIG_DIR}:/work/.azure"
DOCKER_ENTRYPOINT="/opt/terraform.sh"
DOCKER_OPTS="--user $(id -u) ${TTY_OPTIONS} --rm"
DOCKER_MOUNTS="-v $(pwd)/${DIR}:/tmp/${DIR} -v $(pwd)/global.tfvars:/tmp/global.tfvars"
DOCKER_RUN="docker run ${DOCKER_OPTS} --entrypoint ${DOCKER_ENTRYPOINT} ${AZURE_DIR_MOUNT} ${DOCKER_MOUNTS} ${IMAGE}"
DOCKER_SHELL="docker run ${DOCKER_OPTS} --entrypoint /bin/bash ${AZURE_DIR_MOUNT} ${DOCKER_MOUNTS} ${IMAGE}"

function setup {
    mkdir -p "$AZURE_CONFIG_DIR"
    export AZURE_CONFIG_DIR="$AZURE_CONFIG_DIR"
    
    if [ -n "$servicePrincipalId" ]; then
        echo "ARM_CLIENT_ID=$servicePrincipalId"
        echo "ARM_CLIENT_SECRET=$servicePrincipalKey"
        echo "ARM_TENANT_ID=$tenantId"
    fi

    echo "ARM_SUBSCRIPTION_ID=$(az account show -o tsv --query 'id')"
    echo "RG_LOCATION_SHORT=$RG_LOCATION_SHORT"
    echo "RG_LOCATION_LONG=$RG_LOCATION_LONG"
}

function teardown {
    echo "Teardown executed."
}

function prepare {
    setup
    eval "$DOCKER_RUN prepare $DIR $ENV $SUFFIX"
}

function plan {
    setup
    eval "$DOCKER_RUN plan $DIR $ENV $SUFFIX $OPA_BLAST_RADIUS"
}

function apply {
    confirm_action "apply" "$DIR in environment $ENV"
    setup
    eval "$DOCKER_RUN apply $DIR $ENV $SUFFIX"
}

function destroy {
    confirm_action "destroy" "$DIR in environment $ENV"
    setup
    eval "$DOCKER_RUN destroy $DIR $ENV $SUFFIX"
}

function state_remove {
    confirm_action "state-remove" "$DIR in environment $ENV"
    setup
    eval "$DOCKER_RUN state-remove $DIR $ENV $SUFFIX"
}

function validate {
    setup
    eval "$DOCKER_RUN validate $DIR $ENV $SUFFIX"
}

function shell {
    setup
    eval "$DOCKER_SHELL"
}

case "$COMMAND" in
    setup) setup ;;
    teardown) teardown ;;
    prepare) prepare ;;
    plan) plan ;;
    apply) apply ;;
    destroy) destroy ;;
    state-remove) state_remove ;;
    validate) validate ;;
    shell) shell ;;
    *) echo "Invalid command: $COMMAND"; exit 1 ;;
esac