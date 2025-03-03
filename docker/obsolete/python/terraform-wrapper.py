import yaml
from azure.identity import DefaultAzureCredential, EnvironmentCredential, ManagedIdentityCredential, AzureCliCredential
from azure.mgmt.resource import ResourceManagementClient
from azure.mgmt.storage import StorageManagementClient
from azure.mgmt.keyvault import KeyVaultManagementClient
from azure.mgmt.keyvault.models import VaultCreateOrUpdateParameters, Sku, VaultProperties, AccessPolicyEntry, Permissions, KeyPermissions
from azure.mgmt.storage.models import StorageAccountCreateParameters, Sku as StorageSku, Kind, AccessTier, MinimumTLSVersion
from azure.mgmt.resource.locks import ManagementLockClient
from azure.mgmt.resource.locks.models import ManagementLockObject, LockLevel

class AzureConfig:
    def __init__(self, **kwargs):
        self.service_principal_object_id = kwargs.get('service_principal_object_id')
        self.subscription_id = kwargs.get('subscription_id')
        self.tenant_id = kwargs.get('tenant_id')
        self.resource_group_name = kwargs.get('resource_group_name')
        self.resource_group_location = kwargs.get('resource_group_location')
        self.storage_account_name = kwargs.get('storage_account_name')
        self.storage_account_container = kwargs.get('storage_account_container')
        self.key_vault_name = kwargs.get('key_vault_name')
        self.key_vault_key_name = kwargs.get('key_vault_key_name')
        self.resource_locks = kwargs.get('resource_locks', True)
        self.exclude_cli_credential = kwargs.get('exclude_cli_credential', False)
        self.exclude_environment_credential = kwargs.get('exclude_environment_credential', True)
        self.exclude_msi_credential = kwargs.get('exclude_msi_credential', True)

def load_config(file_path='config.yaml'):
    """Load configuration from a YAML file."""
    with open(file_path, 'r') as file:
        return yaml.safe_load(file)

def get_credentials(config):
    if not config.exclude_environment_credential:
        try:
            return EnvironmentCredential()
        except Exception as e:
            print(f"EnvironmentCredential failed: {e}")
    if not config.exclude_msi_credential:
        try:
            return ManagedIdentityCredential()
        except Exception as e:
            print(f"ManagedIdentityCredential failed: {e}")
    if not config.exclude_cli_credential:
        try:
            return AzureCliCredential()
        except Exception as e:
            print(f"AzureCliCredential failed: {e}")
    raise Exception("No credentials found")

def create_resource_group(client, config):
    if not client.resource_groups.check_existence(config.resource_group_name):
        client.resource_groups.create_or_update(
            config.resource_group_name,
            {'location': config.resource_group_location}
        )
        print(f"Resource Group {config.resource_group_name} created.")
    else:
        print(f"Resource Group {config.resource_group_name} already exists.")

def create_storage_account(client, config):
    if not client.storage_accounts.check_name_availability(config.storage_account_name).name_available:
        print(f"Storage Account {config.storage_account_name} name not available.")
        return
    poller = client.storage_accounts.begin_create(
        config.resource_group_name,
        config.storage_account_name,
        StorageAccountCreateParameters(
            sku=StorageSku(name='Standard_GRS'),
            kind=Kind.storage_v2,
            location=config.resource_group_location,
            access_tier=AccessTier.hot,
            enable_https_traffic_only=True,
            minimum_tls_version=MinimumTLSVersion.tls1_2
        )
    )
    poller.result()
    print(f"Storage Account {config.storage_account_name} created.")

def create_key_vault(client, config):
    poller = client.vaults.begin_create_or_update(
        config.resource_group_name,
        config.key_vault_name,
        VaultCreateOrUpdateParameters(
            location=config.resource_group_location,
            properties=VaultProperties(
                tenant_id=config.tenant_id,
                sku=Sku(name='standard', family='A'),
                access_policies=[]
            )
        )
    )
    poller.result()
    print(f"Key Vault {config.key_vault_name} created.")

def main():
    # Load configuration from config.yaml
    config_data = load_config()
    config = AzureConfig(**config_data)

    # Get Azure credentials
    credential = get_credentials(config)

    # Initialize clients
    resource_client = ResourceManagementClient(credential, config.subscription_id)
    storage_client = StorageManagementClient(credential, config.subscription_id)
    keyvault_client = KeyVaultManagementClient(credential, config.subscription_id)

    # Create resources
    create_resource_group(resource_client, config)
    create_storage_account(storage_client, config)
    create_key_vault(keyvault_client, config)

if __name__ == "__main__":
    main()