package main

import (
	"context"
	"fmt"
	"os"

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
	flag.Parse()

	ctx := context.Background()
	err := azure.CreateResourceGroup(ctx, *resourceGroupName, *resourceGroupLocation, *subscriptionID)
	if err != nil {
		fmt.Printf("ERROR: %v \n", err)
		os.Exit(1)
	}

	err = azure.CreateStorageAccount(ctx, *resourceGroupName, *resourceGroupLocation, *storageAccountName, *subscriptionID)
	if err != nil {
		fmt.Printf("ERROR: %v \n", err)
		os.Exit(1)
	}

	err = azure.CreateStorageAccountContainer(ctx, *resourceGroupName, *storageAccountName, *storageAccountContainer, *subscriptionID)
	if err != nil {
		fmt.Printf("ERROR: %v \n", err)
		os.Exit(1)
	}

	err = azure.CreateKeyVault(ctx, *resourceGroupName, *resourceGroupLocation, *keyVaultName, *subscriptionID, *tenantID)
	if err != nil {
		fmt.Printf("ERROR: %v \n", err)
		os.Exit(1)
	}

	os.Exit(0)

}
