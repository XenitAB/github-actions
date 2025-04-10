package terraform.analysis

input_create_aws_s3_bucket = {
  "resource_changes": [
    {
      "change": {
        "actions": [
          "create"
        ]
      },
      "type": "aws_s3_bucket"
    }
  ]
}

input_delete_aws_s3_bucket = {
  "resource_changes": [
    {
      "change": {
        "actions": [
          "delete"
        ]
      },
      "type": "aws_s3_bucket"
    }
  ]
}

input_update_aws_s3_bucket = {
  "resource_changes": [
    {
      "change": {
        "actions": [
          "update"
        ]
      },
      "type": "aws_s3_bucket"
    }
  ]
}

test_create_aws_s3_bucket {
    authz with input as input_create_aws_s3_bucket
}

test_delete_aws_s3_bucket {
    not authz with input as input_delete_aws_s3_bucket
}

test_update_aws_s3_bucket {
    authz with input as input_update_aws_s3_bucket
}