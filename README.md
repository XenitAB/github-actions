# github-actions

This repo builds a Docker image which collects tooling used to setup environments. It serves as documentation of what tools and what versions are in use. Currently it contains primarily:

- Terraform
- Packer

## Using

An example Makefile that makes use of this container might look like this:
```
SHELL:=/bin/bash

SUFFIX="tfstate<4 random digits>"
IMAGE="ghcr.io/xenitab/github-actions/tools:latest"
OPA_BLAST_RADIUS := $(if $(OPA_BLAST_RADIUS), $(OPA_BLAST_RADIUS), 50)
AZURE_CONFIG_DIR := $(if $(AZURE_CONFIG_DIR), $(AZURE_CONFIG_DIR), "$${HOME}/.azure")
TTY_OPTIONS=$(shell [ -t 0 ] && echo '-it')
ifndef ENV
$(error Need to set ENV)
endif
ifndef DIR
$(error Need to set DIR)
endif

prepare:
        docker run --user $(shell id -u) $(TTY_OPTIONS) --entrypoint "/opt/terraform.sh" -v $${PWD}/$(DIR)/.terraform/$(ENV).env:/tmp/$(ENV).env -v $(AZURE_CONFIG_DIR):/work/.azure -v $${PWD}/$(DIR):/tmp/$(DIR) -v $${PWD}/global.tfvars:/tmp/global.tfvars $(IMAGE) prepare $(DIR) $(ENV) $(SUFFIX)

plan:
        docker run --user $(shell id -u) $(TTY_OPTIONS) --entrypoint "/opt/terraform.sh" -v $${PWD}/$(DIR)/.terraform/$(ENV).env:/tmp/$(ENV).env -v $(AZURE_CONFIG_DIR):/work/.azure -v $${PWD}/$(DIR):/tmp/$(DIR) -v $${PWD}/global.tfvars:/tmp/global.tfvars $(IMAGE) plan $(DIR) $(ENV) $(SUFFIX) $(OPA_BLAST_RADIUS)

apply:
        docker run --user $(shell id -u) $(TTY_OPTIONS) --entrypoint "/opt/terraform.sh" -v $${PWD}/$(DIR)/.terraform/$(ENV).env:/tmp/$(ENV).env -v $(AZURE_CONFIG_DIR):/work/.azure -v $${PWD}/$(DIR):/tmp/$(DIR) -v $${PWD}/global.tfvars:/tmp/global.tfvars $(IMAGE) apply $(DIR) $(ENV) $(SUFFIX)

destroy:
        docker run --user $(shell id -u) $(TTY_OPTIONS) --entrypoint "/opt/terraform.sh" -v $${PWD}/$(DIR)/.terraform/$(ENV).env:/tmp/$(ENV).env -v $(AZURE_CONFIG_DIR):/work/.azure -v $${PWD}/$(DIR):/tmp/$(DIR) -v $${PWD}/global.tfvars:/tmp/global.tfvars $(IMAGE) destroy $(DIR) $(ENV) $(SUFFIX)

state-remove:
        docker run --user $(shell id -u) $(TTY_OPTIONS) --entrypoint "/opt/terraform.sh" -v $${PWD}/$(DIR)/.terraform/$(ENV).env:/tmp/$(ENV).env -v $(AZURE_CONFIG_DIR):/work/.azure -v $${PWD}/$(DIR):/tmp/$(DIR) -v $${PWD}/global.tfvars:/tmp/global.tfvars $(IMAGE) state-remove $(DIR) $(ENV) $(SUFFIX)

validate:
        docker run --user $(shell id -u) $(TTY_OPTIONS) --entrypoint "/opt/terraform.sh" -v $${PWD}/$(DIR)/.terraform/$(ENV).env:/tmp/$(ENV).env -v $(AZURE_CONFIG_DIR):/work/.azure -v $${PWD}/$(DIR):/tmp/$(DIR) -v $${PWD}/global.tfvars:/tmp/global.tfvars $(IMAGE) validate $(DIR) $(ENV) $(SUFFIX)
```

## Building

Build the tooling image:s
```shell
docker build -t dev docker/
```

## Releasing

In order to push a new image to the container registry, you create a new release from the GitHub UI or via the API. The publish release event will trigger a GitHub pipeline that deploys the container from the release tag.

If you need to push a custom image to the registry, you need to go to your GitHub [personal access tokens](https://github.com/settings/tokens) page and create an access token. That token is your password when logging in:
```
docker login ghcr.io --username <GITHUB_USERNAME>
docker build -t ghcr.io/xenitab/github-actions/tools:<TAG> ./docker
docker push ghcr.io/xenitab/github-actions/tools:<TAG>
```