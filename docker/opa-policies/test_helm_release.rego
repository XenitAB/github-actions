package terraform.analysis

input_create_helm_release = {
  "resource_changes": [
    {
      "change": {
        "actions": [
          "create"
        ]
      },
      "type": "helm_release"
    }
  ]
}

input_delete_helm_release = {
  "resource_changes": [
    {
      "change": {
        "actions": [
          "delete"
        ]
      },
      "type": "helm_release"
    }
  ]
}

input_update_helm_release = {
  "resource_changes": [
    {
      "change": {
        "actions": [
          "update"
        ]
      },
      "type": "helm_release"
    }
  ]
}

test_create_helm_release {
    authz with input as input_create_helm_release
}

test_delete_helm_release {
    not authz with input as input_delete_helm_release
}

test_update_helm_release {
    authz with input as input_update_helm_release
}