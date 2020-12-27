package azure

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

type azureConfig struct {
	servicePrincipalObjectID      string
	subscriptionID                string
	tenantID                      string
	resourceGroupName             string
	resourceGroupLocation         string
	storageAccountName            string
	storageAccountContainer       string
	keyVaultName                  string
	keyVaultKeyName               string
	resourceLocks                 bool
	defaultAzureCredentialOptions azidentity.DefaultAzureCredentialOptions
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
		servicePrincipalObjectID: getUUID(ctx, cli, "service-principal-object-id"),
		subscriptionID:           getUUID(ctx, cli, "subscription-id"),
		tenantID:                 getUUID(ctx, cli, "tenant-id"),
		resourceGroupName:        cli.String("resource-group-name"),
		resourceGroupLocation:    cli.String("resource-group-location"),
		storageAccountName:       cli.String("storage-account-name"),
		storageAccountContainer:  cli.String("storage-account-container"),
		keyVaultName:             cli.String("keyvault-name"),
		keyVaultKeyName:          cli.String("keyvault-key-name"),
		resourceLocks:            cli.Bool("resource-locks"),
		defaultAzureCredentialOptions: azidentity.DefaultAzureCredentialOptions{
			ExcludeAzureCLICredential:    cli.Bool("exclude-cli-credential"),
			ExcludeEnvironmentCredential: cli.Bool("exclude-environment-credential"),
			ExcludeMSICredential:         cli.Bool("exclude-msi-credential"),
		},
	}

	err := CreateResourceGroup(ctx, config)
	if err != nil {
		return err
	}

	err = CreateStorageAccount(ctx, config)
	if err != nil {
		return err
	}

	if config.resourceLocks {
		err = CreateResourceLock(ctx, config, "Microsoft.Storage", "", "storageAccounts", config.storageAccountName, "DoNotDelete")
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

	if config.resourceLocks {
		err = CreateResourceLock(ctx, config, "Microsoft.KeyVault", "", "vaults", config.keyVaultName, "DoNotDelete")
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

func getUUID(ctx context.Context, cli *cli.Context, flagName string) string {
	log := logr.FromContext(ctx)
	id := cli.String(flagName)

	if id != "" {
		_, err := uuid.Parse(id)
		if err != nil {
			log.Error(err, "Unable to parse UUID", "uuid", id, "flagName", flagName)
			os.Exit(1)
		}
	}
	return id
}
