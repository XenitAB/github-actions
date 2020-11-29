#!/bin/bash
set -e

ACTION=$1
DIR=$2
ENVIRONMENT=$3
PREFIX=$4

RG_LOCATION_SHORT="we"
BACKEND_KEY="${ENVIRONMENT}.terraform.tfstate"
BACKEND_RG="rg-${ENVIRONMENT}-${RG_LOCATION_SHORT}-tfstate"
BACKEND_NAME="sa${ENVIRONMENT}${RG_LOCATION_SHORT}${PREFIX}tfstate"
CONTAINER_NAME="tfstate-${DIR}"

plan () {
  az lock create --name DoNotDelete --resource-group ${BACKEND_RG} --lock-type CanNotDelete --resource-type Microsoft.Storage/storageAccounts --resource ${BACKEND_NAME}

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
  terraform plan -input=false -var-file="variables/${ENVIRONMENT}.tfvars" -var-file="variables/common.tfvars" -out=".terraform/plans/${ENVIRONMENT}"
  opa test /opt/opa-policies -v
}

apply () {
  terraform apply ".terraform/plans/dev"
}

destroy () {
  echo "destroy"
}

tfenv install 0.13.5
tfenv use 0.13.5
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
esac
