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

  ## Backend Key Name
  backend_kv_key = "sops"

}