package azure

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/arm/resources/2020-06-01/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/storage/2019-06-01/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

// CreateResourceGroup creates Azure Resource Group (if it doesn't exist) or returns error
func CreateResourceGroup(ctx context.Context, resourceGroupName, resourceGroupLocation, subscriptionID string) error {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return fmt.Errorf("Failed Azure/CreateResourceGroup/azidentity.NewDefaultAzureCredential: %v", err)
	}

	client := armresources.NewResourceGroupsClient(armcore.NewDefaultConnection(cred, nil), subscriptionID)
	resourceGroupExists, err := client.CheckExistence(ctx, resourceGroupName, &armresources.ResourceGroupsCheckExistenceOptions{})
	if err != nil {
		return fmt.Errorf("Failed Azure/CreateResourceGroup/client.CheckExistence: %v", err)
	}
	if !resourceGroupExists.Success {
		_, err = client.CreateOrUpdate(ctx, resourceGroupName, armresources.ResourceGroup{
			Location: to.StringPtr(resourceGroupLocation),
		}, nil)
		if err != nil {
			return fmt.Errorf("Failed Azure/CreateResourceGroup/client.CreateOrUpdate: %v", err)
		}
		fmt.Printf("INFO: Azure Resource Group (%s) created.\n", resourceGroupName)
		return nil
	}

	fmt.Printf("INFO: Azure Resource Group (%s) already exists.\n", resourceGroupName)
	return nil
}

// CreateStorageAccount creates Azure Storage Account (if it doesn't exist) or returns error
func CreateStorageAccount(ctx context.Context, resourceGroupName, resourceGroupLocation, storageAccountName, subscriptionID string) error {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return fmt.Errorf("Failed Azure/CreateStorageAccount/azidentity.NewDefaultAzureCredential: %v", err)
	}
	client := armstorage.NewStorageAccountsClient(armcore.NewDefaultConnection(cred, nil), subscriptionID)
	_, err = client.GetProperties(ctx, resourceGroupName, storageAccountName, nil)

	if err == nil {
		fmt.Printf("INFO: Azure Storage Account (%s) already exists.\n", storageAccountName)
		return nil
	}

	if err != nil && strings.Contains(err.Error(), "ResourceNotFound") {
		res, err := client.CheckNameAvailability(
			ctx,
			armstorage.StorageAccountCheckNameAvailabilityParameters{
				Name: to.StringPtr(storageAccountName),
				Type: to.StringPtr("Microsoft.Storage/storageAccounts"),
			},
			nil)

		if err != nil {
			return fmt.Errorf("Failed Azure/CreateStorageAccount/client.CheckNameAvailability: %v", err)
		}

		if !*res.CheckNameAvailabilityResult.NameAvailable {
			return fmt.Errorf("Failed Azure/CreateStorageAccount/client.CheckNameAvailability: Storage Account Name (%s) not available", storageAccountName)
		}

		poller, err := client.BeginCreate(
			ctx,
			resourceGroupName,
			storageAccountName,
			armstorage.StorageAccountCreateParameters{
				SKU: &armstorage.SKU{
					Name: armstorage.SKUNameStandardGrs.ToPtr(),
					Tier: armstorage.SKUTierStandard.ToPtr(),
				},
				Kind:     armstorage.KindBlobStorage.ToPtr(),
				Location: to.StringPtr(resourceGroupLocation),
				Properties: &armstorage.StorageAccountPropertiesCreateParameters{
					AccessTier: armstorage.AccessTierCool.ToPtr(),
				},
			}, nil)

		if err != nil {
			log.Fatalf("failed to obtain a response: %v", err)
		}

		_, err = poller.PollUntilDone(context.Background(), 30*time.Second)
		if err != nil {
			return fmt.Errorf("Failed Azure/CreateStorageAccount/poller.PollUntilDone: %v", err)
		}

		fmt.Printf("INFO: Azure Storage Account (%s) created.\n", storageAccountName)
		return nil
	}

	return fmt.Errorf("Failed Azure/CreateStorageAccount/client.GetProperties: %v", err)
}

// CreateStorageAccountContainer creates Storage Account Container (if it doesn't exist) or returns error
func CreateStorageAccountContainer(ctx context.Context, resourceGroupName, storageAccountName, storageAccountContainer, subscriptionID string) error {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return fmt.Errorf("Failed Azure/CreateStorageAccountContainer/azidentity.NewDefaultAzureCredential: %v", err)
	}
	client := armstorage.NewBlobContainersClient(armcore.NewDefaultConnection(cred, nil), subscriptionID)
	_, err = client.Get(
		context.Background(),
		resourceGroupName,
		storageAccountName,
		storageAccountContainer, nil)

	if err == nil {
		fmt.Printf("INFO: Azure Storage Account Container (%s) already exists.\n", storageAccountContainer)
		return nil
	}

	if err != nil && strings.Contains(err.Error(), "The specified container does not exist") {
		_, err := client.Create(
			ctx,
			resourceGroupName,
			storageAccountName,
			storageAccountContainer,
			armstorage.BlobContainer{}, nil)

		if err != nil {
			return fmt.Errorf("Failed Azure/CreateStorageAccountContainer/client.Create: %v", err)
		}

		fmt.Printf("INFO: Azure Storage Account Container (%s) created.\n", storageAccountContainer)
		return nil
	}

	return fmt.Errorf("Failed Azure/CreateStorageAccountContainer/armstorage.NewBlobContainersClient: %v", err)
}
