package terraform.analysis

import rego.v1

import input as tfplan

########################
# Parameters for Policy
########################

# More info here: https://www.openpolicyagent.org/docs/latest/terraform/

# Create test using:
# cat [...].tfplan.json | jq "{resource_changes: [{change: {actions: .resource_changes[_].change.actions}, type: .resource_changes[_].type}]}" > test.json

# acceptable score for automated authorization
blast_radius := data.blast_radius

# weights assigned for each operation on each resource-type
weights := {
	"kubernetes_namespace": {"delete": 100, "create": 1, "modify": 1},
	"kubernetes_service_account": {"delete": 100, "create": 1, "modify": 1},
	"azuread_group": {"delete": 100, "create": 1, "modify": 1},
	"azurerm_container_registry": {"delete": 100, "create": 1, "modify": 1},
	"azurerm_kubernetes_cluster": {"delete": 100, "create": 1, "modify": 1},
	"azurerm_resource_group": {"delete": 200, "create": 1, "modify": 1},
	"azurerm_storage_account": {"delete": 100, "create": 1, "modify": 1},
	"azurerm_virtual_network": {"delete": 100, "create": 1, "modify": 1},
	"azurerm_virtual_machine": {"delete": 100, "create": 1, "modify": 1},
	"azuread_application_password": {"delete": 100, "create": 1, "modify": 100},
	"azurerm_user_assigned_identity": {"delete": 100, "create": 1, "modify": 100},
	"helm_release": {"delete": 100, "create": 1, "modify": 1},
	"aws_ecr_repository": {"delete": 100, "create": 1, "modify": 1},
	"aws_eks_cluster": {"delete": 100, "create": 1, "modify": 1},
	"aws_s3_bucket": {"delete": 100, "create": 1, "modify": 1},
	"aws_vpc": {"delete": 100, "create": 1, "modify": 1},
}

resource_types := {r | weights[r]}

other_resource_types contains type if {
	type := tfplan.resource_changes[_].type
	not resource_types[type]
}

#########
# Policy
#########

# Authorization holds if score for the plan is acceptable and no changes are made to IAM
default authz := false

authz if {
	score < blast_radius
	# not touches_iam
}

# Compute the score for a Terraform plan as the weighted sum of deletions, creations, modifications
score := s if {
	all := [x |
		some resource_type
		crud := weights[resource_type]
		del := crud.delete * num_deletes[resource_type]
		new := crud.create * num_creates[resource_type]
		mod := crud.modify * num_modifies[resource_type]
		x := (del + new) + mod
	]
	others := [x |
		some resource_type
		crud := {"delete": 100, "create": 1, "modify": 1}
		del := crud.delete * other_num_deletes[resource_type]
		new := crud.create * other_num_creates[resource_type]
		mod := crud.modify * other_num_modifies[resource_type]
		x := (del + new) + mod
	]
	s := sum(all) + sum(others)
}

# Whether there is any change to IAM
# touches_iam {
#     all := resources["aws_iam"]
#     count(all) > 0
# }

####################
# Terraform Library
####################

# list of all resources of a given type
resources[resource_type] := all if {
	some resource_type
	resource_types[resource_type]
	all := [name |
		name := tfplan.resource_changes[_]
		name.type == resource_type
	]
}

other_resources[resource_type] := all if {
	some resource_type
	other_resource_types[resource_type]
	all := [name |
		name := tfplan.resource_changes[_]
		name.type == resource_type
	]
}

# number of creations of resources of a given type
num_creates[resource_type] := num if {
	some resource_type
	resource_types[resource_type]
	all := resources[resource_type]
	creates := [res | res := all[_]; res.change.actions[_] == "create"]
	num := count(creates)
}

other_num_creates[resource_type] := num if {
	some resource_type
	other_resource_types[resource_type]
	all := other_resources[resource_type]
	creates := [res | res := all[_]; res.change.actions[_] == "create"]
	num := count(creates)
}

# number of deletions of resources of a given type
num_deletes[resource_type] := num if {
	some resource_type
	resource_types[resource_type]
	all := resources[resource_type]
	deletions := [res | res := all[_]; res.change.actions[_] == "delete"]
	num := count(deletions)
}

other_num_deletes[resource_type] := num if {
	some resource_type
	other_resource_types[resource_type]
	all := other_resources[resource_type]
	deletions := [res | res := all[_]; res.change.actions[_] == "delete"]
	num := count(deletions)
}

# number of modifications to resources of a given type
num_modifies[resource_type] := num if {
	some resource_type
	resource_types[resource_type]
	all := resources[resource_type]
	modifies := [res | res := all[_]; res.change.actions[_] == "update"]
	num := count(modifies)
}

other_num_modifies[resource_type] := num if {
	some resource_type
	other_resource_types[resource_type]
	all := other_resources[resource_type]
	modifies := [res | res := all[_]; res.change.actions[_] == "update"]
	num := count(modifies)
}
