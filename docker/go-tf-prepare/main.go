package main

import (
	"context"
	stdlog "log"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	flag "github.com/spf13/pflag"
	"github.com/xenitab/github-actions/docker/go-tf-prepare/pkg/azure"
)

func main() {
	subscriptionID := flag.String("subscription-id", "", "Azure Subscription ID")
	tenantID := flag.String("tenant-id", "", "Azure Tenant ID")
	resourceGroupName := flag.String("resource-group-name", "", "Azure Resource Group Name")
	resourceGroupLocation := flag.String("resource-group-location", "", "Azure Resource Group Location")
	storageAccountName := flag.String("storage-account-name", "", "Azure Storage Account Name")
	storageAccountContainer := flag.String("storage-account-container", "", "Azure Storage Container")
	keyVaultName := flag.String("keyvault-name", "", "Azure KeyVault Name")
	keyVaultKeyName := flag.String("keyvault-key-name", "", "Azure KeyVault Key Name")
	flag.Parse()

	stdr.SetVerbosity(1)
	log := stdr.New(stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile))
	log = log.WithName("go-tf-preparer")

	ctx := logr.NewContext(context.Background(), log)

	err := azure.CreateResourceGroup(ctx, *resourceGroupName, *resourceGroupLocation, *subscriptionID)
	if err != nil {
		os.Exit(1)
	}

	err = azure.CreateStorageAccount(ctx, *resourceGroupName, *resourceGroupLocation, *storageAccountName, *subscriptionID)
	if err != nil {
		os.Exit(1)
	}

	err = azure.CreateResourceLock(ctx, *resourceGroupName, "Microsoft.Storage", "", "storageAccounts", *storageAccountName, "DoNotDelete", *subscriptionID)
	if err != nil {
		os.Exit(1)
	}

	err = azure.CreateStorageAccountContainer(ctx, *resourceGroupName, *storageAccountName, *storageAccountContainer, *subscriptionID)
	if err != nil {
		os.Exit(1)
	}

	err = azure.CreateKeyVault(ctx, *resourceGroupName, *resourceGroupLocation, *keyVaultName, *subscriptionID, *tenantID)
	if err != nil {
		os.Exit(1)
	}

	err = azure.CreateResourceLock(ctx, *resourceGroupName, "Microsoft.KeyVault", "", "vaults", *keyVaultName, "DoNotDelete", *subscriptionID)
	if err != nil {
		os.Exit(1)
	}

	err = azure.CreateKeyVaultAccessPolicy(ctx, *resourceGroupName, *resourceGroupLocation, *keyVaultName, *subscriptionID, *tenantID)
	if err != nil {
		os.Exit(1)
	}

	err = azure.CreateKeyVaultKey(ctx, *resourceGroupName, *keyVaultName, *keyVaultKeyName, *subscriptionID)
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
