package main

import (
	"context"
	stdlog "log"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	"github.com/urfave/cli/v2"
	"github.com/xenitab/github-actions/docker/go-tf-prepare/pkg/azure"
)

func main() {
	stdr.SetVerbosity(1)
	log := stdr.New(stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile))
	log = log.WithName("go-tf-preparer")

	ctx := logr.NewContext(context.Background(), log)

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "azure",
				Usage: "Terraform prepare for Azure",
				Flags: []cli.Flag{
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
				},
				Action: func(cli *cli.Context) error {
					err := azureAction(ctx, cli)
					if err != nil {
						return err
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err, "CLI execution failed")
		os.Exit(1)
	}

	os.Exit(0)
}

func azureAction(ctx context.Context, cli *cli.Context) error {
	subscriptionID := cli.String("subscription-id")
	tenantID := cli.String("tenant-id")
	resourceGroupName := cli.String("resource-group-name")
	resourceGroupLocation := cli.String("resource-group-location")
	storageAccountName := cli.String("storage-account-name")
	storageAccountContainer := cli.String("storage-account-container")
	keyVaultName := cli.String("keyvault-name")
	keyVaultKeyName := cli.String("keyvault-key-name")

	err := azure.CreateResourceGroup(ctx, resourceGroupName, resourceGroupLocation, subscriptionID)
	if err != nil {
		return err
	}

	err = azure.CreateStorageAccount(ctx, resourceGroupName, resourceGroupLocation, storageAccountName, subscriptionID)
	if err != nil {
		return err
	}

	err = azure.CreateResourceLock(ctx, resourceGroupName, "Microsoft.Storage", "", "storageAccounts", storageAccountName, "DoNotDelete", subscriptionID)
	if err != nil {
		return err
	}

	err = azure.CreateStorageAccountContainer(ctx, resourceGroupName, storageAccountName, storageAccountContainer, subscriptionID)
	if err != nil {
		return err
	}

	err = azure.CreateKeyVault(ctx, resourceGroupName, resourceGroupLocation, keyVaultName, subscriptionID, tenantID)
	if err != nil {
		return err
	}

	err = azure.CreateResourceLock(ctx, resourceGroupName, "Microsoft.KeyVault", "", "vaults", keyVaultName, "DoNotDelete", subscriptionID)
	if err != nil {
		return err
	}

	err = azure.CreateKeyVaultAccessPolicy(ctx, resourceGroupName, resourceGroupLocation, keyVaultName, subscriptionID, tenantID)
	if err != nil {
		return err
	}

	err = azure.CreateKeyVaultKey(ctx, resourceGroupName, keyVaultName, keyVaultKeyName, subscriptionID)
	if err != nil {
		return err
	}

	return nil
}
