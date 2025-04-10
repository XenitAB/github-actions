package terraform.analysis

import rego.v1

input_create_azurerm_container_registry := {"resource_changes": [{
	"change": {"actions": ["create"]},
	"type": "azurerm_container_registry",
}]}

input_delete_azurerm_container_registry := {"resource_changes": [{
	"change": {"actions": ["delete"]},
	"type": "azurerm_container_registry",
}]}

input_update_azurerm_container_registry := {"resource_changes": [{
	"change": {"actions": ["update"]},
	"type": "azurerm_container_registry",
}]}

test_create_azurerm_container_registry if {
	authz with input as input_create_azurerm_container_registry
}

test_delete_azurerm_container_registry if {
	not authz with input as input_delete_azurerm_container_registry
}

test_update_azurerm_container_registry if {
	authz with input as input_update_azurerm_container_registry
}
