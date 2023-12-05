package azure

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armlocks"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/go-logr/logr"

	adapter "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

// CreateResourceGroup creates Azure Resource Group (if it doesn't exist) or returns error
func CreateResourceGroup(ctx context.Context, cred azcore.TokenCredential, config azureConfig) error {
	resourceGroupName := config.ResourceGroupName
	resourceGroupLocation := config.ResourceGroupLocation
	subscriptionID := config.SubscriptionID

	log, err := logr.FromContext(ctx)
	if err != nil {
		return err
	}

	client, err := armresources.NewResourceGroupsClient(subscriptionID, cred, &arm.ClientOptions{})
	if err != nil {
		log.Error(err, "armresources.NewResourceGroupsClient")
		return err
	}
	resourceGroupExists, err := client.CheckExistence(ctx, resourceGroupName, &armresources.ResourceGroupsClientCheckExistenceOptions{})
	if err != nil {
		log.Error(err, "client.CheckExistence")
		return err
	}
	if !resourceGroupExists.Success {
		_, err = client.CreateOrUpdate(ctx, resourceGroupName, armresources.ResourceGroup{
			Location: to.Ptr(resourceGroupLocation),
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
func CreateStorageAccount(ctx context.Context, cred azcore.TokenCredential, config azureConfig) error {
	resourceGroupName := config.ResourceGroupName
	resourceGroupLocation := config.ResourceGroupLocation
	storageAccountName := config.StorageAccountName
	subscriptionID := config.SubscriptionID
	log, err := logr.FromContext(ctx)
	if err != nil {
		return err
	}

	client, err := armstorage.NewAccountsClient(subscriptionID, cred, &arm.ClientOptions{})
	if err != nil {
		log.Error(err, "armstorage.NewAccountsClient")
		return err
	}
	_, err = client.GetProperties(ctx, resourceGroupName, storageAccountName, nil)

	if err == nil {
		log.Info("Azure Storage Account already exists", "storageAccountName", storageAccountName)
		return nil
	}

	if err != nil && strings.Contains(err.Error(), "ResourceNotFound") {
		err := registerResourceProviderIfNeeded(ctx, cred, config, "Microsoft.Storage")
		if err != nil {
			return err
		}

		res, err := client.CheckNameAvailability(
			ctx,
			armstorage.AccountCheckNameAvailabilityParameters{
				Name: to.Ptr(storageAccountName),
				Type: to.Ptr("Microsoft.Storage/storageAccounts"),
			},
			nil)

		if err != nil {
			log.Error(err, "client.CheckNameAvailability")
			return err
		}

		if !*res.CheckNameAvailabilityResult.NameAvailable {
			err := fmt.Errorf("Azure Storage Account Name '%s' not available", storageAccountName)
			log.Error(err, "azure.CreateStorageAccount")
			return err
		}

		poller, err := client.BeginCreate(
			ctx,
			resourceGroupName,
			storageAccountName,
			armstorage.AccountCreateParameters{
				SKU: &armstorage.SKU{
					Name: to.Ptr(armstorage.SKUNameStandardGRS),
					Tier: to.Ptr(armstorage.SKUTierStandard),
				},
				Kind:     to.Ptr(armstorage.KindStorageV2),
				Location: to.Ptr(resourceGroupLocation),
				Properties: &armstorage.AccountPropertiesCreateParameters{
					AccessTier:            to.Ptr(armstorage.AccessTierHot),
					AllowBlobPublicAccess: to.Ptr(false),
					MinimumTLSVersion:     to.Ptr(armstorage.MinimumTLSVersionTLS12),
				},
			}, nil)

		if err != nil {
			log.Error(err, "client.BeginCreate")
			return err
		}

		_, err = poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
			Frequency: 30 * time.Second,
		})
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

func registerResourceProviderIfNeeded(ctx context.Context, cred azcore.TokenCredential, config azureConfig, resourceProviderNamespace string) error {
	subscriptionID := config.SubscriptionID
	log, err := logr.FromContext(ctx)
	if err != nil {
		return err
	}

	client, err := armresources.NewProvidersClient(subscriptionID, cred, &arm.ClientOptions{})
	if err != nil {
		log.Error(err, "armresources.NewProvidersClient")
		return err
	}

	res, err := client.Get(ctx, resourceProviderNamespace, &armresources.ProvidersClientGetOptions{})
	if err != nil {
		log.Error(err, "client.Get")
		return err
	}

	if *res.RegistrationState == "Registered" {
		log.Info("Azure Resource Provider already registered", "resourceProviderNamespace", resourceProviderNamespace, "registrationState", *res.RegistrationState)
		return nil
	}

	regRes, err := client.Register(ctx, resourceProviderNamespace, &armresources.ProvidersClientRegisterOptions{})
	if err != nil {
		log.Error(err, "client.Register")
		return err
	}

	log.Info("Registering Azure Resource Provider", "resourceProviderNamespace", resourceProviderNamespace, "registrationState", *regRes.RegistrationState)

	if *regRes.RegistrationState == "Registered" {
		log.Info("Azure Resource Provider registered", "resourceProviderNamespace", resourceProviderNamespace, "registrationState", *regRes.RegistrationState)
		return nil
	}

	currentRegistrationState := ""
	for i := 1; i < 10; i++ {
		res, err := client.Get(ctx, resourceProviderNamespace, &armresources.ProvidersClientGetOptions{})
		if err != nil {
			log.Error(err, "client.Get")
			return err
		}

		currentRegistrationState = *res.RegistrationState
		if currentRegistrationState == "Registered" {
			log.Info("Azure Resource Provider registered", "resourceProviderNamespace", resourceProviderNamespace, "registrationState", currentRegistrationState)
			return nil
		}

		if currentRegistrationState != "Registering" {
			err := fmt.Errorf("unknown registration state")
			log.Error(err, "Azure Resource Provider in unknown registration state", "resourceProviderNamespace", resourceProviderNamespace, "registrationState", currentRegistrationState)
			return err
		}

		log.Info("Registering Azure Resource Provider", "resourceProviderNamespace", resourceProviderNamespace, "registrationState", currentRegistrationState, "retryCounter", i)

		time.Sleep(time.Duration(i*5) * time.Second)
	}

	err = fmt.Errorf("registration not completed within the specified time period")
	log.Error(err, "Unable to register Azure Resource Provider", resourceProviderNamespace, "registrationState", currentRegistrationState)

	return err
}

// CreateStorageAccountContainer creates Storage Account Container (if it doesn't exist) or returns error
func CreateStorageAccountContainer(ctx context.Context, cred azcore.TokenCredential, config azureConfig) error {
	resourceGroupName := config.ResourceGroupName
	storageAccountName := config.StorageAccountName
	storageAccountContainer := config.StorageAccountContainer
	subscriptionID := config.SubscriptionID
	log, err := logr.FromContext(ctx)
	if err != nil {
		return err
	}

	client, err := armstorage.NewBlobContainersClient(subscriptionID, cred, &arm.ClientOptions{})
	if err != nil {
		log.Error(err, "armstorage.NewBlobContainersClient")
		return err
	}
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
func CreateKeyVault(ctx context.Context, cred azcore.TokenCredential, config azureConfig) error {
	resourceGroupName := config.ResourceGroupName
	resourceGroupLocation := config.ResourceGroupLocation
	keyVaultName := config.KeyVaultName
	subscriptionID := config.SubscriptionID
	tenantID := config.TenantID
	log, err := logr.FromContext(ctx)
	if err != nil {
		return err
	}

	client, err := armkeyvault.NewVaultsClient(subscriptionID, cred, &arm.ClientOptions{})
	if err != nil {
		return err
	}

	_, err = client.Get(ctx, resourceGroupName, keyVaultName, nil)
	if err == nil {
		log.Info("Azure KeyVault already exists", "keyVaultName", keyVaultName)
		return nil
	}

	if err != nil && strings.Contains(err.Error(), "ResourceNotFound") {
		keyVaultNameAvailable, err := client.CheckNameAvailability(ctx, armkeyvault.VaultCheckNameAvailabilityParameters{Name: to.Ptr(keyVaultName), Type: to.Ptr("Microsoft.KeyVault/vaults")}, nil)
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
				Location: to.Ptr(resourceGroupLocation),
				Properties: &armkeyvault.VaultProperties{
					TenantID: to.Ptr(tenantID),
					SKU: &armkeyvault.SKU{
						Family: to.Ptr(armkeyvault.SKUFamilyA),
						Name:   to.Ptr(armkeyvault.SKUNameStandard),
					},
					AccessPolicies: []*armkeyvault.AccessPolicyEntry{},
				},
			}, nil)
		if err != nil {
			log.Error(err, "client.BeginCreateOrUpdate")
			return err
		}
		_, err = poll.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
			Frequency: 5 * time.Second,
		})
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
func CreateKeyVaultAccessPolicy(ctx context.Context, cred azcore.TokenCredential, config azureConfig) error {
	resourceGroupName := config.ResourceGroupName
	keyVaultName := config.KeyVaultName
	subscriptionID := config.SubscriptionID
	tenantID := config.TenantID
	servicePrincipalObjectID := config.ServicePrincipalObjectID
	log, err := logr.FromContext(ctx)
	if err != nil {
		return err
	}

	var currentUserObjectID string
	if servicePrincipalObjectID == "" {
		var err error
		currentUserObjectID, err = getCurrentUserObjectID(ctx, cred, tenantID)
		if err != nil {
			log.Error(err, "getCurrentUserObjectID")
			return err
		}
	}
	if servicePrincipalObjectID != "" {
		currentUserObjectID = servicePrincipalObjectID
	}

	client, err := armkeyvault.NewVaultsClient(subscriptionID, cred, &arm.ClientOptions{})
	if err != nil {
		return err
	}

	keyPermissions := armkeyvault.Permissions{
		Keys: []*armkeyvault.KeyPermissions{
			to.Ptr(armkeyvault.KeyPermissionsUpdate),
			to.Ptr(armkeyvault.KeyPermissionsCreate),
			to.Ptr(armkeyvault.KeyPermissionsGet),
			to.Ptr(armkeyvault.KeyPermissionsList),
			to.Ptr(armkeyvault.KeyPermissionsEncrypt),
			to.Ptr(armkeyvault.KeyPermissionsDecrypt),
		},
	}

	accessPolicies := []*armkeyvault.AccessPolicyEntry{
		{
			TenantID:    &tenantID,
			ObjectID:    &currentUserObjectID,
			Permissions: &keyPermissions,
		},
	}

	properties := armkeyvault.VaultAccessPolicyProperties{AccessPolicies: accessPolicies}
	parameters := armkeyvault.VaultAccessPolicyParameters{Properties: &properties}
	options := armkeyvault.VaultsClientUpdateAccessPolicyOptions{}

	kv, err := client.Get(ctx, resourceGroupName, keyVaultName, nil)
	if err != nil {
		log.Error(err, "client.Get")
		return err
	}

	// Loop through all access policies
	for _, accessPolicy := range kv.Vault.Properties.AccessPolicies {
		// Check if the current object id for the access policy is the same as the current user object id
		if *accessPolicy.ObjectID == currentUserObjectID {
			// Check if the Key Permissions in the access policy are the same as the required Key Permissions
			if keyPermissionsEqual(accessPolicy.Permissions.Keys, keyPermissions.Keys) {
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
func CreateKeyVaultKey(ctx context.Context, cred azcore.TokenCredential, config azureConfig) error {
	resourceGroupName := config.ResourceGroupName
	keyVaultName := config.KeyVaultName
	keyName := config.KeyVaultKeyName
	subscriptionID := config.SubscriptionID
	log, err := logr.FromContext(ctx)
	if err != nil {
		return err
	}

	client, err := armkeyvault.NewKeysClient(subscriptionID, cred, &arm.ClientOptions{})
	if err != nil {
		log.Error(err, "armkeyvault.NewKeysClient")
		return err
	}

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
				Attributes: &armkeyvault.KeyAttributes{
					Enabled: to.Ptr(true),
				},
				KeySize: to.Ptr[int32](2048),
				KeyOps: []*armkeyvault.JSONWebKeyOperation{
					to.Ptr(armkeyvault.JSONWebKeyOperationEncrypt),
					to.Ptr(armkeyvault.JSONWebKeyOperationDecrypt),
				},
				Kty: to.Ptr(armkeyvault.JSONWebKeyTypeRSA),
			}}, nil)
	if err != nil {
		log.Error(err, "armkeyvault.NewKeysClient")
		return err
	}

	log.Info("Azure KeyVault Key created", "keyName", keyName)
	return nil
}

// fixApiVersionTransporter is needed for NewManagementLocksClient since it uses an api-version that doesn't seem to work
type fixApiVersionTransporter struct {
	apiVersion string
	httpClient *http.Client
}

func (t *fixApiVersionTransporter) Do(req *http.Request) (*http.Response, error) {
	reqQP := req.URL.Query()
	reqQP.Set("api-version", t.apiVersion)
	req.URL.RawQuery = reqQP.Encode()

	return t.httpClient.Do(req)
}

// CreateResourceLock creates Azure Resource Lock (if it doesn't exist) or return error
func CreateResourceLock(ctx context.Context, cred azcore.TokenCredential, config azureConfig, resourceProviderNamespace, parentResourcePath, resourceType, resourceName, lockName string) error {
	resourceGroupName := config.ResourceGroupName
	subscriptionID := config.SubscriptionID
	log, err := logr.FromContext(ctx)
	if err != nil {
		return err
	}

	client, err := armlocks.NewManagementLocksClient(subscriptionID, cred, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: &fixApiVersionTransporter{
				apiVersion: "2016-09-01",
				httpClient: http.DefaultClient,
			},
		},
	})
	if err != nil {
		log.Error(err, "armlocks.NewManagementLocksClient")
		return err
	}

	_, err = client.GetAtResourceLevel(ctx, resourceGroupName, resourceProviderNamespace, parentResourcePath, resourceType, resourceName, lockName, &armlocks.ManagementLocksClientGetAtResourceLevelOptions{})
	if err == nil {
		log.Info("Azure Resource Lock already exists", "resourceGroupName", resourceGroupName, "resourceProviderNamespace", resourceProviderNamespace, "resourceType", resourceType, "resourceName", resourceName)
		return nil
	}

	if err != nil && strings.Contains(err.Error(), "LockNotFound") {
		_, err = client.CreateOrUpdateAtResourceLevel(ctx, resourceGroupName, resourceProviderNamespace, parentResourcePath, resourceType, resourceName, lockName, armlocks.ManagementLockObject{Properties: &armlocks.ManagementLockProperties{Level: to.Ptr(armlocks.LockLevelCanNotDelete), Notes: to.Ptr("CanNotDelete")}}, &armlocks.ManagementLocksClientCreateOrUpdateAtResourceLevelOptions{})
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

func getCurrentUserObjectID(ctx context.Context, cred azcore.TokenCredential, tenantID string) (string, error) {
	log, err := logr.FromContext(ctx)
	if err != nil {
		return "", err
	}

	auth, err := adapter.NewAzureIdentityAuthenticationProviderWithScopes(cred, []string{"https://graph.microsoft.com"})
	if err != nil {
		log.Error(err, "adapter.NewAzureIdentityAuthenticationProviderWithScopes")
		return "", err
	}

	adapter, err := msgraphsdk.NewGraphRequestAdapter(auth)
	if err != nil {
		log.Error(err, "msgraphsdk.NewGraphRequestAdapter")
		return "", err
	}

	client := msgraphsdk.NewGraphServiceClient(adapter)

	me, err := client.Me().Get(ctx, &users.UserItemRequestBuilderGetRequestConfiguration{})
	if err != nil {
		log.Error(err, "client.Me().Get()")
		return "", err
	}
	me.GetId()

	id := me.GetId()

	return *id, nil
}

func keyPermissionsEqual(a, b []*armkeyvault.KeyPermissions) bool {
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
			if strings.ToLower(string(*i)) == strings.ToLower(string(*j)) {
				continue OUTER
			}
		}
		return false
	}

	return true
}
