#!/bin/bash

# # Load configuration from config.yaml
# CONFIG_FILE="config.yaml"
# if [[ ! -f "$CONFIG_FILE" ]]; then
#   echo "Error: Configuration file $CONFIG_FILE not found."
#   exit 1
# fi

# # Parse YAML file (requires `yq` - install via `pip install yq` or `brew install yq`)
# SUBSCRIPTION_ID=$(yq e '.subscription_id' $CONFIG_FILE)
# TENANT_ID=$(yq e '.tenant_id' $CONFIG_FILE)
# RESOURCE_GROUP_NAME=$(yq e '.resource_group_name' $CONFIG_FILE)
# RESOURCE_GROUP_LOCATION=$(yq e '.resource_group_location' $CONFIG_FILE)
# STORAGE_ACCOUNT_NAME=$(yq e '.storage_account_name' $CONFIG_FILE)
# STORAGE_ACCOUNT_CONTAINER=$(yq e '.storage_account_container' $CONFIG_FILE)
# KEY_VAULT_NAME=$(yq e '.key_vault_name' $CONFIG_FILE)
# KEY_VAULT_KEY_NAME=$(yq e '.key_vault_key_name' $CONFIG_FILE)
# RESOURCE_LOCKS=$(yq e '.resource_locks' $CONFIG_FILE)
# SERVICE_PRINCIPAL_OBJECT_ID=$(yq e '.service_principal_object_id' $CONFIG_FILE)

# Login to Azure (if not already logged in)
az account show > /dev/null 2>&1
if [[ $? -ne 0 ]]; then
  echo "Logging in to Azure..."
  az login --tenant $TENANT_ID
fi

# Set the active subscription
az account set --subscription $SUBSCRIPTION_ID

# Validate Resource Group Name
if ! [[ $RESOURCE_GROUP_NAME =~ ^[a-zA-Z0-9_\(\)-\.]+$ ]]; then
  echo "Error: Invalid Resource Group Name."
  exit 1
fi

# Create Resource Group
echo "Creating Resource Group: $RESOURCE_GROUP_NAME"
az group create \
  --name $RESOURCE_GROUP_NAME \
  --location $RESOURCE_GROUP_LOCATION

# Register Resource Provider (if needed)
echo "Registering Resource Provider: Microsoft.Storage"
az provider register --namespace Microsoft.Storage --wait

# Create Storage Account
echo "Creating Storage Account: $STORAGE_ACCOUNT_NAME"
az storage account create \
  --name $STORAGE_ACCOUNT_NAME \
  --resource-group $RESOURCE_GROUP_NAME \
  --location $RESOURCE_GROUP_LOCATION \
  --sku Standard_GRS \
  --kind StorageV2 \
  --access-tier Hot \
  --https-only true \
  --min-tls-version TLS1_2

# Create Storage Account Container
echo "Creating Storage Account Container: $STORAGE_ACCOUNT_CONTAINER"
az storage container create \
  --name $STORAGE_ACCOUNT_CONTAINER \
  --account-name $STORAGE_ACCOUNT_NAME

# Create Key Vault
echo "Creating Key Vault: $KEY_VAULT_NAME"
az keyvault create \
  --name $KEY_VAULT_NAME \
  --resource-group $RESOURCE_GROUP_NAME \
  --location $RESOURCE_GROUP_LOCATION \
  --sku standard

# Add Key Vault Access Policy
if [[ -z "$SERVICE_PRINCIPAL_OBJECT_ID" ]]; then
  echo "Retrieving current user's Object ID..."
  CURRENT_USER_OBJECT_ID=$(az ad signed-in-user show --query id -o tsv)
  if [[ -z "$CURRENT_USER_OBJECT_ID" ]]; then
    echo "Error: Unable to retrieve current user's Object ID."
    exit 1
  fi
  SERVICE_PRINCIPAL_OBJECT_ID=$CURRENT_USER_OBJECT_ID
fi

echo "Adding Key Vault Access Policy for Object ID: $SERVICE_PRINCIPAL_OBJECT_ID"
az keyvault set-policy \
  --name $KEY_VAULT_NAME \
  --resource-group $RESOURCE_GROUP_NAME \
  --object-id $SERVICE_PRINCIPAL_OBJECT_ID \
  --key-permissions get list create update encrypt decrypt

# Create Key Vault Key
echo "Creating Key Vault Key: $KEY_VAULT_KEY_NAME"
az keyvault key create \
  --name $KEY_VAULT_KEY_NAME \
  --vault-name $KEY_VAULT_NAME \
  --kty RSA \
  --size 2048 \
  --ops encrypt decrypt

# Add Resource Locks (if enabled)
if [[ "$RESOURCE_LOCKS" == "true" ]]; then
  echo "Adding Resource Lock for Storage Account: $STORAGE_ACCOUNT_NAME"
  az lock create \
    --name DoNotDelete \
    --resource-group $RESOURCE_GROUP_NAME \
    --resource-type Microsoft.Storage/storageAccounts \
    --resource-name $STORAGE_ACCOUNT_NAME \
    --lock-type CanNotDelete

  echo "Adding Resource Lock for Key Vault: $KEY_VAULT_NAME"
  az lock create \
    --name DoNotDelete \
    --resource-group $RESOURCE_GROUP_NAME \
    --resource-type Microsoft.KeyVault/vaults \
    --resource-name $KEY_VAULT_NAME \
    --lock-type CanNotDelete
fi

echo "tf-prepare.sh ran successfully!"