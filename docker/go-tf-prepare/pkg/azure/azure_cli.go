package azure

import (
	"context"

	"github.com/urfave/cli/v2"
)

// Action executes the Azure action
func Action(ctx context.Context, cli *cli.Context) error {
	servicePrincipalObjectID := cli.String("service-principal-object-id")
	subscriptionID := cli.String("subscription-id")
	tenantID := cli.String("tenant-id")
	resourceGroupName := cli.String("resource-group-name")
	resourceGroupLocation := cli.String("resource-group-location")
	storageAccountName := cli.String("storage-account-name")
	storageAccountContainer := cli.String("storage-account-container")
	keyVaultName := cli.String("keyvault-name")
	keyVaultKeyName := cli.String("keyvault-key-name")
	resourceLocks := cli.Bool("resource-locks")

	err := CreateResourceGroup(ctx, resourceGroupName, resourceGroupLocation, subscriptionID)
	if err != nil {
		return err
	}

	err = CreateStorageAccount(ctx, resourceGroupName, resourceGroupLocation, storageAccountName, subscriptionID)
	if err != nil {
		return err
	}

	if resourceLocks {
		err = CreateResourceLock(ctx, resourceGroupName, "Microsoft.Storage", "", "storageAccounts", storageAccountName, "DoNotDelete", subscriptionID)
		if err != nil {
			return err
		}
	}

	err = CreateStorageAccountContainer(ctx, resourceGroupName, storageAccountName, storageAccountContainer, subscriptionID)
	if err != nil {
		return err
	}

	err = CreateKeyVault(ctx, resourceGroupName, resourceGroupLocation, keyVaultName, subscriptionID, tenantID)
	if err != nil {
		return err
	}

	if resourceLocks {
		err = CreateResourceLock(ctx, resourceGroupName, "Microsoft.KeyVault", "", "vaults", keyVaultName, "DoNotDelete", subscriptionID)
		if err != nil {
			return err
		}
	}

	err = CreateKeyVaultAccessPolicy(ctx, resourceGroupName, resourceGroupLocation, keyVaultName, subscriptionID, tenantID, servicePrincipalObjectID)
	if err != nil {
		return err
	}

	err = CreateKeyVaultKey(ctx, resourceGroupName, keyVaultName, keyVaultKeyName, subscriptionID)
	if err != nil {
		return err
	}

	return nil
}

// Flags returns the cli flags for Azure
func Flags() []cli.Flag {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     "service-principal-object-id",
			Usage:    "Service Principal Object ID",
			Required: false,
			EnvVars:  []string{"AZURE_SERVICE_PRINCIPAL_OBJECT_ID"},
		},
		&cli.StringFlag{
			Name:     "subscription-id",
			Usage:    "Azure Subscription ID",
			Required: true,
			EnvVars:  []string{"AZURE_SUBSCRIPTION_ID"},
		},
		&cli.StringFlag{
			Name:     "tenant-id",
			Usage:    "Azure Tenant ID",
			Required: true,
			EnvVars:  []string{"AZURE_TENANT_ID"},
		},
		&cli.StringFlag{
			Name:     "resource-group-name",
			Usage:    "Azure Resource Group Name",
			Required: true,
			EnvVars:  []string{"AZURE_RESOURCE_GROUP_NAME"},
		},
		&cli.StringFlag{
			Name:     "resource-group-location",
			Usage:    "Azure Resource Group Location",
			Required: true,
			EnvVars:  []string{"AZURE_RESOURCE_GROUP_LOCATION"},
		},
		&cli.StringFlag{
			Name:     "storage-account-name",
			Usage:    "Azure Storage Account Name",
			Required: true,
			EnvVars:  []string{"AZURE_STORAGE_ACCOUNT_NAME"},
		},
		&cli.StringFlag{
			Name:     "storage-account-container",
			Usage:    "Azure Storage Account Container",
			Required: true,
			EnvVars:  []string{"AZURE_STORAGE_ACCOUNT_CONTAINER"},
		},
		&cli.StringFlag{
			Name:     "keyvault-name",
			Usage:    "Azure KeyVault Name",
			Required: true,
			EnvVars:  []string{"AZURE_KEYVAULT_NAME"},
		},
		&cli.StringFlag{
			Name:     "keyvault-key-name",
			Usage:    "Azure KeyVault Key Name",
			Required: true,
			EnvVars:  []string{"AZURE_KEYVAULT_KEY_NAME"},
		},
		&cli.BoolFlag{
			Name:    "resource-locks",
			Usage:   "Should Azure Resource Locks be used?",
			Value:   true,
			EnvVars: []string{"AZURE_RESOURCE_LOCKS"},
		},
	}
	return flags
}
