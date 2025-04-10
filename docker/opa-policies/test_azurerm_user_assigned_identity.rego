package terraform.analysis

input_create_azurerm_user_assigned_identity = {
  "resource_changes": [
    {
      "change": {
        "actions": [
          "create"
        ]
      },
      "type": "azurerm_user_assigned_identity"
    }
  ]
}

input_delete_azurerm_user_assigned_identity = {
  "resource_changes": [
    {
      "change": {
        "actions": [
          "delete"
        ]
      },
      "type": "azurerm_user_assigned_identity"
    }
  ]
}

input_update_azurerm_user_assigned_identity = {
  "resource_changes": [
    {
      "change": {
        "actions": [
          "update"
        ]
      },
      "type": "azurerm_user_assigned_identity"
    }
  ]
}

test_create_azurerm_user_assigned_identity {
    authz with input as input_create_azurerm_user_assigned_identity
}

test_delete_azurerm_user_assigned_identity {
    not authz with input as input_delete_azurerm_user_assigned_identity
}

test_update_azurerm_user_assigned_identity {
    not authz with input as input_update_azurerm_user_assigned_identity
}