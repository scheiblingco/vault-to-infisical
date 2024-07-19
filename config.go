package main

type InfisicalConfig struct {
	Url          string `json:"url" env:"INFISICAL_URL"`
	ProjectId    string `json:"projectId" env:"INFISICAL_PROJECTID"`
	ClientId     string `json:"clientId" env:"INFISICAL_CLIENTID"`
	ClientSecret string `json:"clientSecret" env:"INFISICAL_CLIENTSECRET"`
	Environment  string `json:"environment" env:"INFISICAL_ENVIRONMENT"`
}

type VaultAuthMethod string

const VaultAuthToken VaultAuthMethod = "token"

type VaultConfig struct {
	Address    string          `json:"addr" env:"VAULT_ADDR"`
	AuthMethod VaultAuthMethod `json:"authMethod" env:"VAULT_AUTHMETHOD"`
	Token      string          `json:"token" env:"VAULT_TOKEN"`
	StoreName  string          `json:"storeName" env:"VAULT_STORENAME"`
	IgnoreTls  bool            `json:"ignoreTls" env:"VAULT_IGNORETLS"`
}

type TransferSettings struct {
	Overwrite bool `json:"overwrite" env:"SETTINGS_OVERWRITE"`
}

type Config struct {
	Infisical InfisicalConfig  `json:"infisical"`
	Vault     VaultConfig      `json:"vault"`
	Settings  TransferSettings `json:"settings"`
}
