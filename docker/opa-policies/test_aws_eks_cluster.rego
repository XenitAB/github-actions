package terraform.analysis

import rego.v1

input_create_aws_eks_cluster := {"resource_changes": [{
	"change": {"actions": ["create"]},
	"type": "aws_eks_cluster",
}]}

input_delete_aws_eks_cluster := {"resource_changes": [{
	"change": {"actions": ["delete"]},
	"type": "aws_eks_cluster",
}]}

input_update_aws_eks_cluster := {"resource_changes": [{
	"change": {"actions": ["update"]},
	"type": "aws_eks_cluster",
}]}

test_create_aws_eks_cluster if {
	authz with input as input_create_aws_eks_cluster
}

test_delete_aws_eks_cluster if {
	not authz with input as input_delete_aws_eks_cluster
}

test_update_aws_eks_cluster if {
	authz with input as input_update_aws_eks_cluster
}
