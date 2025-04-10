package terraform.analysis

import rego.v1

input_create_azurerm_kubernetes_cluster := {"resource_changes": [{
	"change": {"actions": ["create"]},
	"type": "azurerm_kubernetes_cluster",
}]}

input_delete_azurerm_kubernetes_cluster := {"resource_changes": [{
	"change": {"actions": ["delete"]},
	"type": "azurerm_kubernetes_cluster",
}]}

input_update_azurerm_kubernetes_cluster := {"resource_changes": [{
	"change": {"actions": ["update"]},
	"type": "azurerm_kubernetes_cluster",
}]}

test_create_azurerm_kubernetes_cluster if {
	authz with input as input_create_azurerm_kubernetes_cluster
}

test_delete_azurerm_kubernetes_cluster if {
	not authz with input as input_delete_azurerm_kubernetes_cluster
}

test_update_azurerm_kubernetes_cluster if {
	authz with input as input_update_azurerm_kubernetes_cluster
}
