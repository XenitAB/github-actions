data "azurerm_client_config" "current" {}

terraform {
  required_version = ">= 1.9.0, < 2.0.0"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.10"  # Updated to a newer 4.x version
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.5"  # Ensuring latest within 3.x range
    }
    modtm = {
      source  = "Azure/modtm"
      version = "0.3.2"
    }
  }
}

provider "azurerm" {
  subscription_id = var.subscription_id
  features {}
}

# Importing the Azure naming module to ensure resources have unique CAF compliant names.
module "naming" {
  source  = "Azure/naming/azurerm"
  version = "0.4.2"
  suffix = [var.environment, local.location_shortname, "tfstate${var.suffix}"]
}

module "avm-res-resources-resourcegroup" {
  source    = "Azure/avm-res-resources-resourcegroup/azurerm"
  version   = "0.2.1"
  name      = module.naming.resource_group.name
  location  = var.location
}

# Create Storage Account
module "avm-res-storage-storageaccount" {
  source  = "Azure/avm-res-storage-storageaccount/azurerm"
  version = "0.5.0"
  name                = module.naming.storage_account.name
  location            = var.location
  resource_group_name = module.avm-res-resources-resourcegroup.name
  tags                = local.tags
  lock                = var.resource_locks == true ? { kind = "CanNotDelete" } : null
  shared_access_key_enabled = true
}

# Create Key Vault
module "avm-res-keyvault-vault" {
  source              = "Azure/avm-res-keyvault-vault/azurerm"
  version             = "0.10.0"
  name                = module.naming.key_vault.name
  resource_group_name = module.avm-res-resources-resourcegroup.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  location            = var.location
  tags                = local.tags
  lock                = var.resource_locks == true ? { kind = "CanNotDelete" } : null
  network_acls = {
    bypass = "AzureServices"
    default_action = "Deny"
    ip_rules = ["98.128.229.108"]
  }
  keys = {
    key1 = {
      name     = "sops"
      key_type = "RSA"
      key_size = 2048
      key_opts = ["decrypt", "encrypt"]
    }
  }
}

resource "azuread_group" "tfstate-group" {
  display_name = "tfstatetest"
  members = [ data.azurerm_client_config.current.object_id ]
  security_enabled = true
}

module "avm-res-authorization-roleassignment" {
  source  = "Azure/avm-res-authorization-roleassignment/azurerm"
  version = "0.2.0"

  role_assignments_azure_resource_manager = {
    keyVault = {
      principal_id = azuread_group.tfstate-group.id
      principal_type = "Group"
      role_definition_name = "Key Vault Administrator"
      scope = module.avm-res-resources-resourcegroup.resource_id
    },
    storageAccount = {
      principal_id = azuread_group.tfstate-group.id
      principal_type = "Group"
      role_definition_name = "Storage Blob Data Contributor"
      scope = module.avm-res-resources-resourcegroup.resource_id
    }
  }
}