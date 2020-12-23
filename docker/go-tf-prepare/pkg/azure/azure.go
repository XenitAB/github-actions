package azure

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/arm/keyvault/2019-09-01/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/resources/2020-06-01/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/storage/2019-06-01/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2016-09-01/locks"
	"github.com/go-logr/logr"
	"github.com/jongio/azidext/go/azidext"
)

// CreateResourceGroup creates Azure Resource Group (if it doesn't exist) or returns error
func CreateResourceGroup(ctx context.Context, defaultAzureCredentialOptions azidentity.DefaultAzureCredentialOptions, resourceGroupName, resourceGroupLocation, subscriptionID string) error {
	log := logr.FromContext(ctx)

	cred, err := azidentity.NewDefaultAzureCredential(&defaultAzureCredentialOptions)
	if err != nil {
		log.Error(err, "azidentity.NewDefaultAzureCredential")
		return err
	}

	client := armresources.NewResourceGroupsClient(armcore.NewDefaultConnection(cred, nil), subscriptionID)
	resourceGroupExists, err := client.CheckExistence(ctx, resourceGroupName, &armresources.ResourceGroupsCheckExistenceOptions{})
	if err != nil {
		log.Error(err, "client.CheckExistence")
		return err
	}
	if !resourceGroupExists.Success {
		_, err = client.CreateOrUpdate(ctx, resourceGroupName, armresources.ResourceGroup{
			Location: to.StringPtr(resourceGroupLocation),
		}, nil)
		if err != nil {
			log.Error(err, "client.CreateOrUpdate")
			return err
		}

		log.Info("Azure Resource Group created", "resourceGroupName", resourceGroupName)
		return nil
	}

	log.Info("Azure Resource Group already exists", "resourceGroupName", resourceGroupName)
	return nil
}

// CreateStorageAccount creates Azure Storage Account (if it doesn't exist) or returns error
func CreateStorageAccount(ctx context.Context, defaultAzureCredentialOptions azidentity.DefaultAzureCredentialOptions, resourceGroupName, resourceGroupLocation, storageAccountName, subscriptionID string) error {
	log := logr.FromContext(ctx)

	cred, err := azidentity.NewDefaultAzureCredential(&defaultAzureCredentialOptions)
	if err != nil {
		log.Error(err, "azidentity.NewDefaultAzureCredential")
		return err
	}
	client := armstorage.NewStorageAccountsClient(armcore.NewDefaultConnection(cred, nil), subscriptionID)
	_, err = client.GetProperties(ctx, resourceGroupName, storageAccountName, nil)

	if err == nil {
		log.Info("Azure Storage Account already exists", "storageAccountName", storageAccountName)
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
			log.Error(err, "client.CheckNameAvailability")
			return err
		}

		if !*res.CheckNameAvailabilityResult.NameAvailable {
			log.Error(err, "client.CheckNameAvailability: Azure Storage Account Name not available", "storageAccountName", storageAccountName)
			return err
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
				Kind:     armstorage.KindStorageV2.ToPtr(),
				Location: to.StringPtr(resourceGroupLocation),
				Properties: &armstorage.StorageAccountPropertiesCreateParameters{
					AccessTier:            armstorage.AccessTierHot.ToPtr(),
					AllowBlobPublicAccess: to.BoolPtr(false),
					MinimumTLSVersion:     armstorage.MinimumTLSVersionTLS12.ToPtr(),
				},
			}, nil)

		if err != nil {
			log.Error(err, "client.BeginCreate")
			return err
		}

		_, err = poller.PollUntilDone(ctx, 30*time.Second)
		if err != nil {
			log.Error(err, "poller.PollUntilDone")
			return err
		}

		log.Info("Azure Storage Account created", "storageAccountName", storageAccountName)
		return nil
	}

	log.Error(err, "client.GetProperties")
	return err
}

// CreateStorageAccountContainer creates Storage Account Container (if it doesn't exist) or returns error
func CreateStorageAccountContainer(ctx context.Context, defaultAzureCredentialOptions azidentity.DefaultAzureCredentialOptions, resourceGroupName, storageAccountName, storageAccountContainer, subscriptionID string) error {
	log := logr.FromContext(ctx)

	cred, err := azidentity.NewDefaultAzureCredential(&defaultAzureCredentialOptions)
	if err != nil {
		log.Error(err, "azidentity.NewDefaultAzureCredential")
		return err
	}
	client := armstorage.NewBlobContainersClient(armcore.NewDefaultConnection(cred, nil), subscriptionID)
	_, err = client.Get(
		ctx,
		resourceGroupName,
		storageAccountName,
		storageAccountContainer, nil)

	if err == nil {
		log.Info("Azure Storage Account Container already exists", "storageAccountContainer", storageAccountContainer)
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
			log.Error(err, "client.Create")
			return err
		}

		log.Info("Azure Storage Account Container created", "storageAccountContainer", storageAccountContainer)
		return nil
	}

	log.Error(err, "armstorage.NewBlobContainersClient")
	return err
}

// CreateKeyVault creates Azure Key Vault (if it doesn't exist) or returns error
func CreateKeyVault(ctx context.Context, defaultAzureCredentialOptions azidentity.DefaultAzureCredentialOptions, resourceGroupName, resourceGroupLocation, keyVaultName, subscriptionID, tenantID string) error {
	log := logr.FromContext(ctx)

	cred, err := azidentity.NewDefaultAzureCredential(&defaultAzureCredentialOptions)
	if err != nil {
		log.Error(err, "azidentity.NewDefaultAzureCredential")
		return err
	}
	client := armkeyvault.NewVaultsClient(armcore.NewDefaultConnection(cred, nil), subscriptionID)

	_, err = client.Get(ctx, resourceGroupName, keyVaultName, nil)
	if err == nil {
		log.Info("Azure KeyVault already exists", "keyVaultName", keyVaultName)
		return nil
	}

	if err != nil && strings.Contains(err.Error(), "ResourceNotFound") {
		keyVaultNameAvailable, err := client.CheckNameAvailability(ctx, armkeyvault.VaultCheckNameAvailabilityParameters{Name: to.StringPtr(keyVaultName), Type: to.StringPtr("Microsoft.KeyVault/vaults")}, nil)
		if err != nil {
			log.Error(err, "client.CheckNameAvailability")
			return err
		}

		if !*keyVaultNameAvailable.CheckNameAvailabilityResult.NameAvailable {
			log.Error(err, "client.CheckNameAvailability: Azure KeyVault Name not available", "keyVaultName", keyVaultName)
			return err
		}

		poll, err := client.BeginCreateOrUpdate(
			ctx,
			resourceGroupName,
			keyVaultName,
			armkeyvault.VaultCreateOrUpdateParameters{
				Location: to.StringPtr(resourceGroupLocation),
				Properties: &armkeyvault.VaultProperties{
					TenantID: to.StringPtr(tenantID),
					SKU: &armkeyvault.SKU{
						Family: armkeyvault.SKUFamilyA.ToPtr(),
						Name:   armkeyvault.SKUNameStandard.ToPtr(),
					},
					AccessPolicies: &[]armkeyvault.AccessPolicyEntry{},
				},
			}, nil)
		if err != nil {
			log.Error(err, "client.BeginCreateOrUpdate")
			return err
		}
		_, err = poll.PollUntilDone(ctx, 5*time.Second)
		if err != nil {
			log.Error(err, "poll.PollUntilDone")
			return err
		}

		log.Info("Azure KeyVault created", "keyVaultName", keyVaultName)
		return nil
	}

	return fmt.Errorf("Failed Azure/CreateKeyVault/client.Get: %v", err)
}

// CreateKeyVaultAccessPolicy creates Azure Key Vault Access Policy (if it doesn't exist) or returns error
func CreateKeyVaultAccessPolicy(ctx context.Context, defaultAzureCredentialOptions azidentity.DefaultAzureCredentialOptions, resourceGroupName, resourceGroupLocation, keyVaultName, subscriptionID, tenantID, servicePrincipalObjectID string) error {
	log := logr.FromContext(ctx)

	var currentUserObjectID string
	if servicePrincipalObjectID == "" {
		var err error
		currentUserObjectID, err = getCurrentUserObjectID(ctx, defaultAzureCredentialOptions, tenantID)
		if err != nil {
			log.Error(err, "getCurrentUserObjectID")
			return err
		}
	}
	if servicePrincipalObjectID != "" {
		currentUserObjectID = servicePrincipalObjectID
	}

	cred, err := azidentity.NewDefaultAzureCredential(&defaultAzureCredentialOptions)
	if err != nil {
		log.Error(err, "azidentity.NewDefaultAzureCredential")
		return err
	}

	client := armkeyvault.NewVaultsClient(armcore.NewDefaultConnection(cred, nil), subscriptionID)

	keyPermissions := armkeyvault.Permissions{
		Keys: &[]armkeyvault.KeyPermissions{
			armkeyvault.KeyPermissionsUpdate,
			armkeyvault.KeyPermissionsCreate,
			armkeyvault.KeyPermissionsGet,
			armkeyvault.KeyPermissionsList,
			armkeyvault.KeyPermissionsEncrypt,
			armkeyvault.KeyPermissionsDecrypt,
		},
	}

	accessPolicies := []armkeyvault.AccessPolicyEntry{
		{
			TenantID:    &tenantID,
			ObjectID:    &currentUserObjectID,
			Permissions: &keyPermissions,
		},
	}

	properties := armkeyvault.VaultAccessPolicyProperties{AccessPolicies: &accessPolicies}
	parameters := armkeyvault.VaultAccessPolicyParameters{Properties: &properties}
	options := armkeyvault.VaultsUpdateAccessPolicyOptions{}

	kv, err := client.Get(ctx, resourceGroupName, keyVaultName, nil)
	if err != nil {
		log.Error(err, "client.Get")
		return err
	}

	// Loop through all access policies
	for _, accessPolicy := range *kv.Vault.Properties.AccessPolicies {
		// Check if the current object id for the access policy is the same as the current user object id
		if *accessPolicy.ObjectID == currentUserObjectID {
			// Check if the Key Permissions in the access policy are the same as the required Key Permissions
			if keyPermissionsEqual(*accessPolicy.Permissions.Keys, *keyPermissions.Keys) {
				// If the correct Key Permissions already exists, return early
				log.Info("Azure KeyVault Access Policy already correct", "currentUserObjectID", currentUserObjectID)
				return nil
			}
		}
	}

	_, err = client.UpdateAccessPolicy(ctx, resourceGroupName, keyVaultName, armkeyvault.AccessPolicyUpdateKindAdd, parameters, &options)
	if err != nil {
		log.Error(err, "client.UpdateAccessPolicy")
		return err
	}

	log.Info("Azure KeyVault Access Policy created or updated", "currentUserObjectID", currentUserObjectID)

	return nil
}

// CreateKeyVaultKey creates Azure Key Vault Key (if it doesn't exist) or returns error
func CreateKeyVaultKey(ctx context.Context, defaultAzureCredentialOptions azidentity.DefaultAzureCredentialOptions, resourceGroupName, keyVaultName, keyName, subscriptionID string) error {
	log := logr.FromContext(ctx)

	cred, err := azidentity.NewDefaultAzureCredential(&defaultAzureCredentialOptions)
	if err != nil {
		log.Error(err, "azidentity.NewDefaultAzureCredential")
		return err
	}

	client := armkeyvault.NewKeysClient(armcore.NewDefaultConnection(cred, nil), subscriptionID)

	_, err = client.Get(ctx, resourceGroupName, keyVaultName, keyName, nil)
	if err == nil {
		log.Info("Azure KeyVault Key already exists", "keyName", keyName)
		return nil
	}

	_, err = client.CreateIfNotExist(
		ctx,
		resourceGroupName,
		keyVaultName,
		keyName,
		armkeyvault.KeyCreateParameters{
			Properties: &armkeyvault.KeyProperties{
				Attributes: &armkeyvault.Attributes{
					Enabled: to.BoolPtr(true),
				},
				KeySize: to.Int32Ptr(2048),
				KeyOps: &[]armkeyvault.JSONWebKeyOperation{
					armkeyvault.JSONWebKeyOperationEncrypt,
					armkeyvault.JSONWebKeyOperationDecrypt,
				},
				Kty: armkeyvault.JSONWebKeyTypeRsa.ToPtr(),
			}}, nil)
	if err != nil {
		log.Error(err, "armkeyvault.NewKeysClient")
		return err
	}

	log.Info("Azure KeyVault Key created", "keyName", keyName)
	return nil
}

// CreateResourceLock creates Azure Resource Lock (if it doesn't exist) or return error
func CreateResourceLock(ctx context.Context, defaultAzureCredentialOptions azidentity.DefaultAzureCredentialOptions, resourceGroupName, resourceProviderNamespace, parentResourcePath, resourceType, resourceName, lockName, subscriptionID string) error {
	log := logr.FromContext(ctx)

	client := locks.NewManagementLocksClient(subscriptionID)

	tokenRequestOptions := azcore.TokenRequestOptions{Scopes: []string{"https://management.azure.com/.default"}}
	authenticationPolicy := azcore.AuthenticationPolicyOptions{Options: tokenRequestOptions}
	credentialOptions := azidext.DefaultAzureCredentialOptions{DefaultCredential: &defaultAzureCredentialOptions, AuthenticationPolicy: &authenticationPolicy}
	authorizer, err := azidext.NewDefaultAzureCredentialAdapter(&credentialOptions)
	if err != nil {
		log.Error(err, "azidext.NewDefaultAzureCredentialAdapter")
		return err
	}

	client.Authorizer = authorizer

	_, err = client.GetAtResourceLevel(ctx, resourceGroupName, resourceProviderNamespace, parentResourcePath, resourceType, resourceName, lockName)
	if err == nil {
		log.Info("Azure Resource Lock already exists", "resourceGroupName", resourceGroupName, "resourceProviderNamespace", resourceProviderNamespace, "resourceType", resourceType, "resourceName", resourceName)
		return nil
	}

	if err != nil && strings.Contains(err.Error(), "LockNotFound") {
		_, err = client.CreateOrUpdateAtResourceLevel(ctx, resourceGroupName, resourceProviderNamespace, parentResourcePath, resourceType, resourceName, lockName, locks.ManagementLockObject{ManagementLockProperties: &locks.ManagementLockProperties{Level: "CanNotDelete", Notes: to.StringPtr("CanNotDelete")}})
		if err != nil {
			log.Error(err, "client.CreateOrUpdateAtResourceLevel")
			return err
		}

		log.Info("Azure Resource Lock created", "resourceGroupName", resourceGroupName, "resourceProviderNamespace", resourceProviderNamespace, "resourceType", resourceType, "resourceName", resourceName)
		return nil
	}

	log.Error(err, "client.GetAtResourceLevel")
	return err
}

func getCurrentUserObjectID(ctx context.Context, defaultAzureCredentialOptions azidentity.DefaultAzureCredentialOptions, tenantID string) (string, error) {
	log := logr.FromContext(ctx)

	client := graphrbac.NewSignedInUserClient(tenantID)

	tokenRequestOptions := azcore.TokenRequestOptions{Scopes: []string{"https://graph.windows.net/.default"}}
	authenticationPolicy := azcore.AuthenticationPolicyOptions{Options: tokenRequestOptions}
	credentialOptions := azidext.DefaultAzureCredentialOptions{DefaultCredential: &defaultAzureCredentialOptions, AuthenticationPolicy: &authenticationPolicy}
	authorizer, err := azidext.NewDefaultAzureCredentialAdapter(&credentialOptions)
	if err != nil {
		log.Error(err, "azidext.NewDefaultAzureCredentialAdapter")
		return "", err
	}

	client.Authorizer = authorizer

	currentUser, err := client.Get(ctx)
	if err != nil {
		log.Error(err, "client.Get")
		return "", err
	}

	return *currentUser.ObjectID, nil
}

func keyPermissionsEqual(a, b []armkeyvault.KeyPermissions) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

OUTER:
	for _, i := range a {
		for _, j := range b {
			// a may have the first letters uppercase while b always have them lowercase
			if strings.ToLower(string(i)) == strings.ToLower(string(j)) {
				continue OUTER
			}
		}
		return false
	}

	return true
}
