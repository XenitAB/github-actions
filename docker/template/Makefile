.ONESHELL:
SHELL:=/bin/bash

SUFFIX="tfstate[random 4 digits]"
IMAGE="ghcr.io/xenitab/github-actions/tools:[tag]"
AWS_ENABLED:=false

# User mapping compat when using podman as alias for docker (with Debian update-alternatives or Red Hat alternatives)
export PODMAN_USERNS=keep-id

OPA_BLAST_RADIUS := $(if $(OPA_BLAST_RADIUS),$(OPA_BLAST_RADIUS),50)
RG_LOCATION_SHORT:=we
RG_LOCATION_LONG:=westeurope
AZURE_CONFIG_DIR := $(if $(AZURE_CONFIG_DIR),$(AZURE_CONFIG_DIR),"$${HOME}/.azure")
TTY_OPTIONS=$(shell [ -t 0 ] && echo '-it')
TEMP_ENV_FILE:=$(shell mktemp)
TEMP_SHELL_RC_FILE:=$(shell mktemp)

ifndef ENV
$(error Need to set ENV)
endif
ifndef DIR
$(error Need to set DIR)
endif

AZURE_DIR_MOUNT:=-v $(AZURE_CONFIG_DIR):/work/.azure
DOCKER_ENTRYPOINT:=/opt/terraform.sh
DOCKER_OPTS:=--user $(shell id -u) $(TTY_OPTIONS) --rm --env-file $(TEMP_ENV_FILE)
DOCKER_MOUNTS:=-v "$${PWD}/$(DIR)":"/tmp/$(DIR)" -v "$${PWD}/global.tfvars":"/tmp/global.tfvars"
DOCKER_RUN:=docker run $(DOCKER_OPTS) --entrypoint $(DOCKER_ENTRYPOINT) $(AZURE_DIR_MOUNT) $(DOCKER_MOUNTS) $(IMAGE)
DOCKER_SHELL:=docker run $(DOCKER_OPTS) --entrypoint /bin/bash $(AZURE_DIR_MOUNT) $(DOCKER_MOUNTS) -v "$(TEMP_SHELL_RC_FILE)":"/work/.bashrc" $(IMAGE)
CLEANUP_COMMAND:=$(MAKE) --no-print-directory teardown TEMP_ENV_FILE=$(TEMP_ENV_FILE) TEMP_SHELL_RC_FILE=$(TEMP_SHELL_RC_FILE)

.PHONY: setup
.SILENT: setup
setup:
	set -e

	mkdir -p $(AZURE_CONFIG_DIR)
	export AZURE_CONFIG_DIR="$(AZURE_CONFIG_DIR)"

	if [ -n "$${servicePrincipalId}" ]; then
		echo ARM_CLIENT_ID=$${servicePrincipalId} >> $(TEMP_ENV_FILE)
		echo ARM_CLIENT_SECRET=$${servicePrincipalKey} >> $(TEMP_ENV_FILE)
		echo ARM_TENANT_ID=$${tenantId} >> $(TEMP_ENV_FILE)
	fi

	echo ARM_SUBSCRIPTION_ID=$$(az account show -o tsv --query 'id') >> $(TEMP_ENV_FILE)

	if [ "$(AWS_ENABLED)" == "true" ]; then
		if [ -z "$${AWS_ACCESS_KEY_ID}" ]; then
			AWS_ROLE_ARN=$$(aws configure get role_arn)
			aws sts assume-role --role-arn $${AWS_ROLE_ARN} --role-session-name awscli --output text --query 'Credentials' | \
				{
					read -r AWS_ACCESS_KEY_ID TIMESTAMP AWS_SECRET_ACCESS_KEY AWS_SESSION_TOKEN
					cat >> $(TEMP_ENV_FILE) <<EOF
					AWS_ACCESS_KEY_ID=$${AWS_ACCESS_KEY_ID}
					AWS_SECRET_ACCESS_KEY=$${AWS_SECRET_ACCESS_KEY}
					AWS_SESSION_TOKEN=$${AWS_SESSION_TOKEN}
					AWS_DEFAULT_REGION=$$(aws configure get region)
					EOF
				}
		else
		    if [ -n "$${AWS_SESSION_TOKEN}" ]; then
				echo "AWS_SESSION_TOKEN=$${AWS_SESSION_TOKEN}" >> $(TEMP_ENV_FILE)
				echo "AWS_SECURITY_TOKEN=$${AWS_SECURITY_TOKEN}" >> $(TEMP_ENV_FILE)
			fi
			echo AWS_ACCESS_KEY_ID=$${AWS_ACCESS_KEY_ID} >> $(TEMP_ENV_FILE)
			echo AWS_SECRET_ACCESS_KEY=$${AWS_SECRET_ACCESS_KEY} >> $(TEMP_ENV_FILE)
			echo AWS_DEFAULT_REGION=$${AWS_DEFAULT_REGION} >> $(TEMP_ENV_FILE)
		fi

		echo AWS_DEFAULT_OUTPUT="json" >> $(TEMP_ENV_FILE)
		echo AWS_PAGER= >> $(TEMP_ENV_FILE)
	fi

	echo RG_LOCATION_SHORT=$(RG_LOCATION_SHORT) >> $(TEMP_ENV_FILE)
	echo RG_LOCATION_LONG=$(RG_LOCATION_LONG) >> $(TEMP_ENV_FILE)

.PHONY: teardown
.SILENT: teardown
teardown:
	-rm -f $(TEMP_ENV_FILE)

.PHONY: prepare
prepare: setup
	trap '$(CLEANUP_COMMAND)' EXIT
	$(DOCKER_RUN) prepare $(DIR) $(ENV) $(SUFFIX)

.PHONY: plan
plan: setup
	trap '$(CLEANUP_COMMAND)' EXIT
	$(DOCKER_RUN) plan $(DIR) $(ENV) $(SUFFIX) $(OPA_BLAST_RADIUS)

.PHONY: apply
apply: setup
	trap '$(CLEANUP_COMMAND)' EXIT
	$(DOCKER_RUN) apply $(DIR) $(ENV) $(SUFFIX)

.PHONY: destroy
destroy: setup
	trap '$(CLEANUP_COMMAND)' EXIT
	$(DOCKER_RUN) destroy $(DIR) $(ENV) $(SUFFIX)

.PHONY: state-remove
state-remove: setup
	trap '$(CLEANUP_COMMAND)' EXIT
	$(DOCKER_RUN) state-remove $(DIR) $(ENV) $(SUFFIX)

.PHONY: validate
validate: setup
	trap '$(CLEANUP_COMMAND)' EXIT
	$(DOCKER_RUN) validate $(DIR) $(ENV) $(SUFFIX)

.PHONY: shell
shell: setup
	trap '$(CLEANUP_COMMAND)' EXIT
	cat << EOF > $(TEMP_SHELL_RC_FILE)
	cd /tmp/$(DIR)
	/opt/terraform.sh init $(DIR) $(ENV) $(SUFFIX)
	EOF
	$(DOCKER_SHELL)
