locals {

## Generic
  tags = {
    environment = var.environment
    location    = var.location
  }
  ### Location Short Codes
  location_short_codes = {
    "swedencentral" = "swc"
    "northeurope"   = "neu"
    "westeurope"    = "weu"
  }

  ### Location Shortname
  location_shortname = lookup(local.location_short_codes, lower(var.location), "unknown")

## Resource Group
  resource_group_name = "rg-${var.environment}-${local.location_shortname}-tfstate${var.suffix}"
  
## Storage Account
  storage_account_name = "sa${var.environment}${local.location_shortname}tfstate${var.suffix}"

## Key Vault
  keyvault_name = "kv-${var.environment}-${local.location_shortname}-tfstate${var.suffix}"
  backend_kv_key = "sops"

}