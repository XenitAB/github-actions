#!/bin/bash
set -e

ACTION=$1
DIR=$2
ENVIRONMENT=$3
SUFFIX=$4
OPA_BLAST_RADIUS=$5

RG_LOCATION_SHORT="we"
RG_LOCATION_LONG="westeurope"
BACKEND_KEY="${ENVIRONMENT}.terraform.tfstate"
BACKEND_RG="rg-${ENVIRONMENT}-${RG_LOCATION_SHORT}-${SUFFIX}"
BACKEND_KV="kv-${ENVIRONMENT}-${RG_LOCATION_SHORT}-${SUFFIX}"
BACKEND_KV_KEY="sops"
BACKEND_NAME="sa${ENVIRONMENT}${RG_LOCATION_SHORT}${SUFFIX}"
CONTAINER_NAME="tfstate-${DIR}"
ENVIRONMENT_FILE="/tmp/${ENVIRONMENT}.env"

if [ -z "${OPA_BLAST_RADIUS}" ]; then
  OPA_BLAST_RADIUS=50
fi

set_azure_keyvault_permissions() {
  echo "Assigning permissions to Azure KeyVault ${BACKEND_KV}"
  AZ_ACCOUNT_TYPE="$(az account show --query user.type --output tsv)"
  if [[ "${AZ_ACCOUNT_TYPE}" = "user" ]]; then
    AZ_USER_OBJECT_ID="$(az ad signed-in-user show --query objectId --output tsv)"
    az keyvault set-policy --name ${BACKEND_KV} --resource-group ${BACKEND_RG} --object-id ${AZ_USER_OBJECT_ID} --key-permissions create get list encrypt decrypt 1>/dev/null
  elif [[ "${AZ_ACCOUNT_TYPE}" = "servicePrincipal" ]]; then
    AZ_SPN="$(az account show --query user.name --output tsv)"
    az keyvault set-policy --name ${BACKEND_KV} --resource-group ${BACKEND_RG} --spn ${AZ_SPN} --key-permissions create get list encrypt decrypt 1>/dev/null
  fi
}

prepare () {
  if [ $(az group exists --name ${BACKEND_RG}) = false ]; then
    echo "INFO: Creating resource group ${BACKEND_RG} in location ${RG_LOCATION_LONG}"
    az group create --name ${BACKEND_RG} --location ${RG_LOCATION_LONG}
  fi

  if ! $(az storage account show --resource-group ${BACKEND_RG} --name ${BACKEND_NAME} --output none); then
    echo "Creating Azure Storage Account ${BACKEND_NAME} in location ${RG_LOCATION_LONG} / resource group ${BACKEND_RG}"
    az storage account create --resource-group ${BACKEND_RG} --name ${BACKEND_NAME} 1>/dev/null
  fi

  if ! $(az storage container show --account-name ${BACKEND_NAME} --name ${CONTAINER_NAME} --output none); then
    echo "Creating Azure Storage Container ${CONTAINER_NAME} in Storage Account ${BACKEND_NAME}"
    az storage container create --account-name ${BACKEND_NAME} --name ${CONTAINER_NAME} 1>/dev/null
  fi

  if ! $(az keyvault show --name ${BACKEND_KV} --output none); then
    echo "Creating Azure KeyVault ${BACKEND_KV} in location ${RG_LOCATION_LONG} / resource group ${BACKEND_RG}"
    az keyvault create --name ${BACKEND_KV} --resource-group ${BACKEND_RG} --location ${RG_LOCATION_LONG} 1>/dev/null
  fi

  set +e
  KEYVAULT_KEY_TEST="$(az keyvault key show --vault-name ${BACKEND_KV} --name ${BACKEND_KV_KEY} --output none 2>&1)"
  if echo ${KEYVAULT_KEY_TEST} | grep KeyNotFound; then
    echo "Creating Azure KeyVault key in ${BACKEND_KV}"
    az keyvault key create --name ${BACKEND_KV_KEY} --vault-name ${BACKEND_KV} --protection software --ops encrypt decrypt 1>/dev/null
    set_azure_keyvault_permissions
  elif echo ${KEYVAULT_KEY_TEST} | grep "does not have keys get permission on key vault"; then
    set_azure_keyvault_permissions
  fi
  set -e

  az lock create --name DoNotDelete --resource-group ${BACKEND_RG} --lock-type CanNotDelete --resource-type Microsoft.Storage/storageAccounts --resource ${BACKEND_NAME} 1>/dev/null
  az lock create --name DoNotDelete --resource-group ${BACKEND_RG} --lock-type CanNotDelete --resource-type Microsoft.KeyVault/vaults --resource ${BACKEND_KV} 1>/dev/null
}

plan () {
  rm -f .terraform/plans/${ENVIRONMENT}
  terraform init -input=false -backend-config="key=${BACKEND_KEY}" -backend-config="resource_group_name=${BACKEND_RG}" -backend-config="storage_account_name=${BACKEND_NAME}" -backend-config="container_name=${CONTAINER_NAME}" -backend-config="snapshot=true"

  set +e
  terraform workspace select ${ENVIRONMENT} 2> /dev/null
  if [ $? -ne 0 ]; then
    terraform workspace new ${ENVIRONMENT}
    terraform workspace select ${ENVIRONMENT}
  fi
  set -e

  mkdir -p .terraform/plans
  terraform plan -input=false -var-file="variables/${ENVIRONMENT}.tfvars" -var-file="variables/common.tfvars" -var-file="../global.tfvars" -out=".terraform/plans/${ENVIRONMENT}"
  terraform show -json .terraform/plans/${ENVIRONMENT} > .terraform/plans/${ENVIRONMENT}.json
  cat /opt/opa-policies/data.json | jq ".blast_radius = ${OPA_BLAST_RADIUS}" > /tmp/opa-data.json
  opa test /opt/opa-policies -v
  OPA_AUTHZ=$(opa eval --format pretty --data /tmp/opa-data.json --data /opt/opa-policies/terraform.rego --input .terraform/plans/${ENVIRONMENT}.json "data.terraform.analysis.authz")
  OPA_SCORE=$(opa eval --format pretty --data /tmp/opa-data.json --data /opt/opa-policies/terraform.rego --input .terraform/plans/${ENVIRONMENT}.json "data.terraform.analysis.score")
  if [[ "${OPA_AUTHZ}" == "true" ]]; then
    echo "INFO: OPA Authorization: true (score: ${OPA_SCORE} / blast_radius: ${OPA_BLAST_RADIUS})"
    rm -rf .terraform/plans/${ENVIRONMENT}.json
  else
    echo "ERROR: OPA Authorization: false (score: ${OPA_SCORE} / blast_radius: ${OPA_BLAST_RADIUS})"
    rm -rf .terraform/plans/${ENVIRONMENT}.json
    rm -rf .terraform/plans/${ENVIRONMENT}
    exit 1
  fi
  SOPS_KEY_ID="$(az keyvault key show --name ${BACKEND_KV_KEY} --vault-name ${BACKEND_KV} --query key.kid --output tsv)"
  sops --encrypt --azure-kv ${SOPS_KEY_ID} .terraform/plans/${ENVIRONMENT} > .terraform/plans/${ENVIRONMENT}.enc
  rm -rf .terraform/plans/${ENVIRONMENT}
}

apply () {
  terraform init -input=false -backend-config="key=${BACKEND_KEY}" -backend-config="resource_group_name=${BACKEND_RG}" -backend-config="storage_account_name=${BACKEND_NAME}" -backend-config="container_name=${CONTAINER_NAME}" -backend-config="snapshot=true"
  terraform workspace select ${ENVIRONMENT}
  SOPS_KEY_ID="$(az keyvault key show --name ${BACKEND_KV_KEY} --vault-name ${BACKEND_KV} --query key.kid --output tsv)"
  sops --decrypt --azure-kv ${SOPS_KEY_ID} .terraform/plans/${ENVIRONMENT}.enc > .terraform/plans/${ENVIRONMENT}
  rm -rf .terraform/plans/${ENVIRONMENT}.enc
  set +e
  terraform apply ".terraform/plans/${ENVIRONMENT}"
  EXIT_CODE=$?
  set -e
  rm -rf .terraform/plans/${ENVIRONMENT}
  exit $EXIT_CODE
}

destroy () {
  terraform init -input=false -backend-config="key=${BACKEND_KEY}" -backend-config="resource_group_name=${BACKEND_RG}" -backend-config="storage_account_name=${BACKEND_NAME}" -backend-config="container_name=${CONTAINER_NAME}" -backend-config="snapshot=true"
  terraform workspace select ${ENVIRONMENT}
  
  echo "-------"
  echo "You are about to run terraform destroy on ${DIR} in ${ENVIRONMENT}"
  echo "-------"
  
  echo -n "Please confirm by writing \"${DIR}/${ENVIRONMENT}\": "
  read VERIFICATION_INPUT

  if [[ "${VERIFICATION_INPUT}" == "${DIR}/${ENVIRONMENT}" ]]; then
    terraform destroy -var-file="variables/${ENVIRONMENT}.tfvars" -var-file="variables/common.tfvars" -var-file="../global.tfvars"
  else
    echo "Wrong input detected (${VERIFICATION_INPUT}). Exiting..."
    exit 1
  fi
}

state_remove () {
  terraform init -input=false -backend-config="key=${BACKEND_KEY}" -backend-config="resource_group_name=${BACKEND_RG}" -backend-config="storage_account_name=${BACKEND_NAME}" -backend-config="container_name=${CONTAINER_NAME}" -backend-config="snapshot=true"
  terraform workspace select ${ENVIRONMENT}
  TF_STATE_OBJECTS=$(terraform state list)
  
  echo "-------"
  echo "You are about to run terraform state rm on ${DIR} in ${ENVIRONMENT}"
  echo "-------"
  
  echo -n "Please confirm by writing \"${DIR}/${ENVIRONMENT}\": "
  read VERIFICATION_INPUT

  if [[ "${VERIFICATION_INPUT}" == "${DIR}/${ENVIRONMENT}" ]]; then
    echo -n "Please enter what to grep regex arguments (default: grep -E \".*\"): "
    read GREP_ARGUMENT
    GREP_ARGUMENT=${GREP_ARGUMENT:-.*}
    TF_STATE_TO_REMOVE=$(echo "${TF_STATE_OBJECTS}" | grep -E "${GREP_ARGUMENT}")
    TF_STATE_TO_REMOVE_COUNT=$(echo "${TF_STATE_TO_REMOVE}" | wc -l)
    
    echo "You are about to remove the following objects from the terraform state: "
    echo ""
    echo "-------"
    echo "${TF_STATE_TO_REMOVE}"
    echo "-------"
    echo ""

    echo -n "Please confirm the number of objects that will be removed (${TF_STATE_TO_REMOVE_COUNT}): "
    read VERIFICATION_INPUT_COUNT
    if [[ ${VERIFICATION_INPUT_COUNT} -eq ${TF_STATE_TO_REMOVE_COUNT} ]]; then
      for TF_STATE_OBJECT in ${TF_STATE_TO_REMOVE}; do
        terraform state rm ${TF_STATE_OBJECT}
      done
    else
      echo "Wrong input detected (${VERIFICATION_INPUT_COUNT}). Exiting..."
      exit 1
    fi
  else
    echo "Wrong input detected (${VERIFICATION_INPUT}). Exiting..."
    exit 1
  fi
}



envup() {
  if [ -f ${ENVIRONMENT_FILE} ]; then
    set -a
    source ${ENVIRONMENT_FILE}
    set +a
  fi
}

envup
cd /tmp/$DIR

case $ACTION in
  plan )
    plan
    ;;

  apply )
    apply
    ;;

  destroy )
    destroy
    ;;

  prepare )
    prepare
    ;;

  state-remove )
    state_remove
    ;;
esac
