# github-actions

This repo builds a Docker image which collects tooling used to setup environments. It serves as documentation of what tools and what versions are in use. Currently it contains primarily:

- Terraform
- Packer

## Using

An example Makefile that makes use of this container might look like this:

```Makefile
.ONESHELL:
SHELL:=/bin/bash

SUFFIX="tfstate[random 4 digits]"
IMAGE="ghcr.io/xenitab/github-actions/tools:[tag]"
AWS_ENABLED:=false

OPA_BLAST_RADIUS := $(if $(OPA_BLAST_RADIUS),$(OPA_BLAST_RADIUS),50)
AZURE_CONFIG_DIR := $(if $(AZURE_CONFIG_DIR),$(AZURE_CONFIG_DIR),"$${HOME}/.azure")
TTY_OPTIONS=$(shell [ -t 0 ] && echo '-it')
TEMP_ENV_FILE:=$(shell mktemp)

ifndef ENV
$(error Need to set ENV)
endif
ifndef DIR
$(error Need to set DIR)
endif

AZURE_DIR_MOUNT:=$(shell echo "-v $(AZURE_CONFIG_DIR):/work/.azure")
COMMAND:=$(shell echo "docker run --user $(shell id -u) $(TTY_OPTIONS) --entrypoint "/opt/terraform.sh" --env-file $(TEMP_ENV_FILE) $(AZURE_DIR_MOUNT) -v $${PWD}/$(DIR):/tmp/$(DIR) -v $${PWD}/global.tfvars:/tmp/global.tfvars $(IMAGE)")
CLEANUP_COMMAND:=$(shell echo "$(MAKE) --no-print-directory teardown TEMP_ENV_FILE=$(TEMP_ENV_FILE)")

.PHONY: setup
.SILENT: setup
setup:
	set -e

	mkdir -p $(AZURE_CONFIG_DIR)
	export AZURE_CONFIG_DIR="$(AZURE_CONFIG_DIR)"

	if [ "$${CI}" == "true" ]; then
		echo ARM_CLIENT_ID=$${servicePrincipalId} >> $(TEMP_ENV_FILE)
		echo ARM_CLIENT_SECRET=$${servicePrincipalKey} >> $(TEMP_ENV_FILE)
		echo ARM_TENANT_ID=$${tenantId} >> $(TEMP_ENV_FILE)
	fi

	echo ARM_SUBSCRIPTION_ID=$$(az account show -o tsv --query 'id') >> $(TEMP_ENV_FILE)

	if [ "$(AWS_ENABLED)" == "true" ]; then
		if [ "$${CI}" == "true" ]; then
			echo AWS_ACCESS_KEY_ID=$${AWS_ACCESS_KEY_ID} >> $(TEMP_ENV_FILE)
			echo AWS_SECRET_ACCESS_KEY=$${AWS_SECRET_ACCESS_KEY} >> $(TEMP_ENV_FILE)
			echo AWS_DEFAULT_REGION=$${AWS_DEFAULT_REGION} >> $(TEMP_ENV_FILE)
		else
			AWS_ROLE_ARN=$$(aws configure get role_arn)
			AWS_VARS=$$(aws sts assume-role --role-arn $${AWS_ROLE_ARN} --role-session-name awscli --output json | grep -E "AccessKeyId|SecretAccessKey|SessionToken" | sed "s/[ |,|\"]//g")
			echo AWS_ACCESS_KEY_ID=$$(echo "$${AWS_VARS}" | grep AccessKeyId | sed "s/AccessKeyId://g") >> $(TEMP_ENV_FILE)
			echo AWS_SECRET_ACCESS_KEY=$$(echo "$${AWS_VARS}" | grep SecretAccessKey | sed "s/SecretAccessKey://g") >> $(TEMP_ENV_FILE)
			echo AWS_SESSION_TOKEN=$$(echo "$${AWS_VARS}" | grep SessionToken | sed "s/SessionToken://g") >> $(TEMP_ENV_FILE)
			echo AWS_DEFAULT_REGION=$$(aws configure get region) >> $(TEMP_ENV_FILE)
		fi

		echo AWS_DEFAULT_OUTPUT="json" >> $(TEMP_ENV_FILE)
		echo AWS_PAGER= >> $(TEMP_ENV_FILE)
	fi

.PHONY: teardown
.SILENT: teardown
teardown:
	-rm -f $(TEMP_ENV_FILE)

.PHONY: prepare
prepare: setup
	trap '$(CLEANUP_COMMAND)' EXIT
	$(COMMAND) prepare $(DIR) $(ENV) $(SUFFIX)

.PHONY: plan
plan: setup
	trap '$(CLEANUP_COMMAND)' EXIT
	$(COMMAND) plan $(DIR) $(ENV) $(SUFFIX) $(OPA_BLAST_RADIUS)

.PHONY: apply
apply: setup
	trap '$(CLEANUP_COMMAND)' EXIT
	$(COMMAND) apply $(DIR) $(ENV) $(SUFFIX)

.PHONY: destroy
destroy: setup
	trap '$(CLEANUP_COMMAND)' EXIT
	$(COMMAND) destroy $(DIR) $(ENV) $(SUFFIX)

.PHONY: state-remove
state-remove: setup
	trap '$(CLEANUP_COMMAND)' EXIT
	$(COMMAND) state-remove $(DIR) $(ENV) $(SUFFIX)

.PHONY: validate
validate: setup
	trap '$(CLEANUP_COMMAND)' EXIT
	$(COMMAND) validate $(DIR) $(ENV) $(SUFFIX)
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
