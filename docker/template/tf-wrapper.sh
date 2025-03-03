#!/bin/bash

set -e

# Inputs
IMAGE_TAG="latest"
DIR=""

# Variables
SUBSCRIPTION_NAME=$(az account show --query 'name' -o tsv)
IMAGE="test_arm64:$IMAGE_TAG"

# Validate inputs and check Docker before proceeding
if [ -z "$1" ] || ! [[ "$1" =~ ^(setup|teardown|prepare|plan|apply|destroy|state-remove|validate|shell)$ ]]; then
    echo "Usage: $0 {setup|teardown|prepare|plan|apply|destroy|state-remove|validate|shell} --DIR <directory> (--imageTag <tag> (Optional))"
    exit 1
fi

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
        return 0
    else
        echo "Docker image '$image' not found locally. Attempting to pull..."
        if ! docker pull "$image" &> /dev/null; then
            echo "Error: Docker container '$image' is not reachable. Please check the image name and your network connection."
            exit 1
        fi
    fi
}

# Function to confirm action with the user
function confirm_action {
    local action="$1"
    if [ -z "$SUBSCRIPTION_NAME" ]; then
        echo "Error: Unable to retrieve Azure subscription name. Please ensure you are logged into Azure CLI."
        exit 1
    fi
    echo "You are about to run '$action' on subscription '$SUBSCRIPTION_NAME'. Are you sure? (y/n)"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        echo "Action canceled."
        exit 0
    fi
}

while [[ "$#" -gt 0 ]]; do
    case "$1" in
        --imageTag)
            IMAGE_TAG="$2"
            shift 2
            ;;
        --DIR)
            DIR="$2"
            shift 2
            ;;
        *)
            COMMAND="$1"
            shift
            ;;
    esac
done

# Check if Docker is installed
check_docker_installed

# Check if Docker container is reachable
check_docker_container_reachable "$IMAGE"

# Proceed with the rest of the script
export PODMAN_USERS=keep-id

OPA_BLAST_RADIUS=${OPA_BLAST_RADIUS:-50}
AZURE_CONFIG_DIR=${AZURE_CONFIG_DIR:-"$HOME/.azure"}
TTY_OPTIONS=$( [ -t 0 ] && echo "-it" )


DOCKER_MOUNTS="-v $(pwd)/${DIR}:/tmp/${DIR} -v $(pwd)/global.tfvars:/tmp/global.tfvars"

AZURE_DIR_MOUNT="-v ${AZURE_CONFIG_DIR}:/work/.azure"
DOCKER_ENTRYPOINT="/opt/terraform.sh"
DOCKER_OPTS="--user $(id -u) ${TTY_OPTIONS} --rm"
DOCKER_RUN="docker run ${DOCKER_OPTS} --entrypoint ${DOCKER_ENTRYPOINT} ${AZURE_DIR_MOUNT} ${DOCKER_MOUNTS} ${IMAGE}"
DOCKER_SHELL="docker run ${DOCKER_OPTS} --entrypoint /bin/bash ${AZURE_DIR_MOUNT} ${DOCKER_MOUNTS} ${IMAGE}"

function teardown {
    confirm_action teardown
    echo "Teardown executed."
}

function prepare {
    confirm_action prepare
    eval "$DOCKER_RUN prepare"
}

function plan {
    if [ -z "$DIR" ]; then
        echo "Error: Missing required argument --DIR for plan command"
        exit 1
    fi
    confirm_action plan
    setup
    eval "$DOCKER_RUN plan $DIR $OPA_BLAST_RADIUS"
}

function apply {
    if [ -z "$DIR" ]; then
        echo "Error: Missing required argument --DIR for apply command"
        exit 1
    fi
    confirm_action apply
    setup
    eval "$DOCKER_RUN apply $DIR"
}

function destroy {
    if [ -z "$DIR" ]; then
        echo "Error: Missing required argument --DIR for destroy command"
        exit 1
    fi
    confirm_action destroy
    setup
    eval "$DOCKER_RUN destroy $DIR"
}

function state_remove {
    if [ -z "$DIR" ]; then
        echo "Error: Missing required argument --DIR for state-remove command"
        exit 1
    fi
    confirm_action state-remove
    setup
    eval "$DOCKER_RUN state-remove $DIR"
}

function validate {
    if [ -z "$DIR" ]; then
        echo "Error: Missing required argument --DIR for validate command"
        exit 1
    fi
    setup
    eval "$DOCKER_RUN validate $DIR"
}

function shell {
    if [ -z "$DIR" ]; then
        echo "Error: Missing required argument --DIR for shell command"
        exit 1
    fi
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