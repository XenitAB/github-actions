package terraform.analysis

import rego.v1

input_create_azurerm_storage_account := {"resource_changes": [{
	"change": {"actions": ["create"]},
	"type": "azurerm_storage_account",
}]}

input_delete_azurerm_storage_account := {"resource_changes": [{
	"change": {"actions": ["delete"]},
	"type": "azurerm_storage_account",
}]}

input_update_azurerm_storage_account := {"resource_changes": [{
	"change": {"actions": ["update"]},
	"type": "azurerm_storage_account",
}]}

test_create_azurerm_storage_account if {
	authz with input as input_create_azurerm_storage_account
}

test_delete_azurerm_storage_account if {
	not authz with input as input_delete_azurerm_storage_account
}

test_update_azurerm_storage_account if {
	authz with input as input_update_azurerm_storage_account
}
