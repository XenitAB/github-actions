package terraform.analysis

import rego.v1

input_create_azurerm_resource_group := {"resource_changes": [{
	"change": {"actions": ["create"]},
	"type": "azurerm_resource_group",
}]}

input_delete_azurerm_resource_group := {"resource_changes": [{
	"change": {"actions": ["delete"]},
	"type": "azurerm_resource_group",
}]}

input_update_azurerm_resource_group := {"resource_changes": [{
	"change": {"actions": ["update"]},
	"type": "azurerm_resource_group",
}]}

test_create_azurerm_resource_group if {
	authz with input as input_create_azurerm_resource_group
}

test_delete_azurerm_resource_group if {
	not authz with input as input_delete_azurerm_resource_group
}

test_update_azurerm_resource_group if {
	authz with input as input_update_azurerm_resource_group
}
