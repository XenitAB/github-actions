package terraform.analysis

import rego.v1

input_create_fake_resource := {"resource_changes": [{
	"change": {"actions": ["create"]},
	"type": "fake_resource",
}]}

input_delete_fake_resource := {"resource_changes": [{
	"change": {"actions": ["delete"]},
	"type": "fake_resource",
}]}

input_update_fake_resource := {"resource_changes": [{
	"change": {"actions": ["update"]},
	"type": "fake_resource",
}]}

test_create_fake_resource if {
	authz with input as input_create_fake_resource
}

test_delete_fake_resource if {
	not authz with input as input_delete_fake_resource
}

test_update_fake_resource if {
	authz with input as input_update_fake_resource
}
