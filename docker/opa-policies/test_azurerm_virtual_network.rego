package terraform.analysis

import rego.v1

input_create_azurerm_virtual_network := {"resource_changes": [{
	"change": {"actions": ["create"]},
	"type": "azurerm_virtual_network",
}]}

input_delete_azurerm_virtual_network := {"resource_changes": [{
	"change": {"actions": ["delete"]},
	"type": "azurerm_virtual_network",
}]}

input_update_azurerm_virtual_network := {"resource_changes": [{
	"change": {"actions": ["update"]},
	"type": "azurerm_virtual_network",
}]}

test_create_azurerm_virtual_network if {
	authz with input as input_create_azurerm_virtual_network
}

test_delete_azurerm_virtual_network if {
	not authz with input as input_delete_azurerm_virtual_network
}

test_update_azurerm_virtual_network if {
	authz with input as input_update_azurerm_virtual_network
}
