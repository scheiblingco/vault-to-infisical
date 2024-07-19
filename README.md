# Vault to Infisical
This was written to aid in the transfer of data from Hashicorp Vault (or OpenBAO) to Infisical Secrets Manager.

## Usage
You can configure either via environment variables or configuration file. 
The following options are currently available, we are planning to add more authentication methods for the Vault part once we get some free time.

### Configuration Notes
- Environment variables will take precedence over what's configured in the configuration file

### ENV Configuration
You can configure all variables below by specifying the path in ENV notation, e.g.:

```bash
# This represents "infisical": { "url": "https://app.infisical.com" }
INFISICAL_URL=https://app.infisical.com

# This represents "vault": { "storeName": "kv2" }
VAULT_STORENAME=kv2
```

### Json Configuration
```json
{
  // Connection settings for infisical
  "infisical": {
    
    // The URL, generally app.infisical.com if you're not hosting on-prem
    "url": "https://app.infisical.com",
    
    // The project ID of the target project
    "projectId": "",
    
    // The client ID and client Secret for the Universal Auth
    "clientId": "",
    "clientSecret": "",
    
    // The target environment for the secrets from vault
    "environment": "prod"
  },
  
  // The connection to hashicorp vault
  "vault": {
    
    // The URL to your vault instance
    "url": "",
    
    // The authentication method
    // Currently, only token is supported
    "authMethod": "token",
    
    // The root token
    "token": "",
    
    // The KV2 store name
    // Currently only the KV2 store is supported
    // Default is kv2
    // If you have multiple set up or changed the name you may need to change this
    "storeName": "kv2",
    
    // Ignore TLS Certificate Errors
    "ignoreTls": false
  },
  "settings": {
    // Overwrite existing keys in Infisical
    "overwrite": true
  }
}
```

