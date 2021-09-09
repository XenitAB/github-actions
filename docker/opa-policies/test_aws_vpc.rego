package terraform.analysis

input_create_aws_vpc = {
  "resource_changes": [
    {
      "change": {
        "actions": [
          "create"
        ]
      },
      "type": "aws_vpc"
    }
  ]
}

input_delete_aws_vpc = {
  "resource_changes": [
    {
      "change": {
        "actions": [
          "delete"
        ]
      },
      "type": "aws_vpc"
    }
  ]
}

input_update_aws_vpc = {
  "resource_changes": [
    {
      "change": {
        "actions": [
          "update"
        ]
      },
      "type": "aws_vpc"
    }
  ]
}

test_create_aws_vpc {
    authz with input as input_create_aws_vpc
}

test_delete_aws_vpc {
    not authz with input as input_delete_aws_vpc
}

test_update_aws_vpc {
    authz with input as input_update_aws_vpc
}