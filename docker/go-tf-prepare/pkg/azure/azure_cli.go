package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/go-playground/validator/v10"
	"github.com/urfave/cli/v2"
)

type azureConfig struct {
	ServicePrincipalObjectID      string `validate:"omitempty,uuid"`
	SubscriptionID                string `validate:"uuid"`
	TenantID                      string `validate:"uuid"`
	ResourceGroupName             string `validate:"resourcegroup,min=1,max=90"`
	ResourceGroupLocation         string `validate:"alphanum,lowercase"`
	StorageAccountName            string `validate:"alphanum,lowercase,min=3,max=24"`
	StorageAccountContainer       string `validate:"storageaccountcontainer,min=3,max=24"`
	KeyVaultName                  string `validate:"keyvault,min=3,max=24"`
	KeyVaultKeyName               string `validate:"keyvaultkey,min=1,max=127"`
	ResourceLocks                 bool
	DefaultAzureCredentialOptions azidentity.DefaultAzureCredentialOptions
}

func (config azureConfig) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("resourcegroup", validateResourceGroupName)
	validate.RegisterValidation("storageaccountcontainer", validateStorageAccountContainerName)
	validate.RegisterValidation("keyvault", validateKeyVaultName)
	validate.RegisterValidation("keyvaultkey", validateKeyVaultKeyName)
	err := validate.Struct(config)
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
		&cli.BoolFlag{
			Name:    "exclude-cli-credential",
			Usage:   "Should Azure CLI authentication be excluded from authentication chain?",
			Value:   false,
			EnvVars: []string{"AZURE_EXCLUDE_CLI_CREDENTIAL"},
		},
		&cli.BoolFlag{
			Name:    "exclude-environment-credential",
			Usage:   "Should Azure Environment authentication be excluded from authentication chain?",
			Value:   true,
			EnvVars: []string{"AZURE_EXCLUDE_ENVIRONMENT_CREDENTIAL"},
		},
		&cli.BoolFlag{
			Name:    "exclude-msi-credential",
			Usage:   "Should Azure MSI authentication be excluded from authentication chain?",
			Value:   true,
			EnvVars: []string{"AZURE_EXCLUDE_MSI_CREDENTIAL"},
		},
	}
	return flags
}

// Action executes the Azure action
func Action(ctx context.Context, cli *cli.Context) error {
	config := azureConfig{
		ServicePrincipalObjectID: cli.String("service-principal-object-id"),
		SubscriptionID:           cli.String("subscription-id"),
		TenantID:                 cli.String("tenant-id"),
		ResourceGroupName:        cli.String("resource-group-name"),
		ResourceGroupLocation:    cli.String("resource-group-location"),
		StorageAccountName:       cli.String("storage-account-name"),
		StorageAccountContainer:  cli.String("storage-account-container"),
		KeyVaultName:             cli.String("keyvault-name"),
		KeyVaultKeyName:          cli.String("keyvault-key-name"),
		ResourceLocks:            cli.Bool("resource-locks"),
		DefaultAzureCredentialOptions: azidentity.DefaultAzureCredentialOptions{
			ExcludeAzureCLICredential:    cli.Bool("exclude-cli-credential"),
			ExcludeEnvironmentCredential: cli.Bool("exclude-environment-credential"),
			ExcludeMSICredential:         cli.Bool("exclude-msi-credential"),
		},
	}

	err := config.Validate()
	if err != nil {
		return err
	}

	err = CreateResourceGroup(ctx, config)
	if err != nil {
		return err
	}

	err = CreateStorageAccount(ctx, config)
	if err != nil {
		return err
	}

	if config.ResourceLocks {
		err = CreateResourceLock(ctx, config, "Microsoft.Storage", "", "storageAccounts", config.StorageAccountName, "DoNotDelete")
		if err != nil {
			return err
		}
	}

	err = CreateStorageAccountContainer(ctx, config)
	if err != nil {
		return err
	}

	err = CreateKeyVault(ctx, config)
	if err != nil {
		return err
	}

	if config.ResourceLocks {
		err = CreateResourceLock(ctx, config, "Microsoft.KeyVault", "", "vaults", config.KeyVaultName, "DoNotDelete")
		if err != nil {
			return err
		}
	}

	err = CreateKeyVaultAccessPolicy(ctx, config)
	if err != nil {
		return err
	}

	err = CreateKeyVaultKey(ctx, config)
	if err != nil {
		return err
	}

	return nil
}
