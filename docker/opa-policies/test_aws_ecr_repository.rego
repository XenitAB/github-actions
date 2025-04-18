package terraform.analysis

import rego.v1

input_create_aws_ecr_repository := {"resource_changes": [{
	"change": {"actions": ["create"]},
	"type": "aws_ecr_repository",
}]}

input_delete_aws_ecr_repository := {"resource_changes": [{
	"change": {"actions": ["delete"]},
	"type": "aws_ecr_repository",
}]}

input_update_aws_ecr_repository := {"resource_changes": [{
	"change": {"actions": ["update"]},
	"type": "aws_ecr_repository",
}]}

test_create_aws_ecr_repository if {
	authz with input as input_create_aws_ecr_repository
}

test_delete_aws_ecr_repository if {
	not authz with input as input_delete_aws_ecr_repository
}

test_update_aws_ecr_repository if {
	authz with input as input_update_aws_ecr_repository
}
