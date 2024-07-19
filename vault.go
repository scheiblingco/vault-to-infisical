package main

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/vault-client-go"
	"path/filepath"
	"time"
)

type VaultConnection struct {
	Config VaultConfig

	client *vault.Client
}

func (vc *VaultConnection) Connect() error {
	vaultOptions := []vault.ClientOption{
		vault.WithAddress(vc.Config.Address),
	}

	if vc.Config.IgnoreTls {
		vaultOptions = append(vaultOptions, vault.WithTLS(
			vault.TLSConfiguration{
				InsecureSkipVerify: true,
			},
		))
	}

	vlt, err := vault.New(vaultOptions...)
	if err != nil {
		return err
	}

	if err := vlt.SetToken(vc.Config.Token); err != nil {
		return err
	}

	vc.client = vlt

	return nil
}

func (vc *VaultConnection) ListSecrets(path string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	secrets, err := vc.client.Secrets.KvV2List(ctx, path, vault.WithMountPath(vc.Config.StoreName))
	if err != nil {
		return nil, err
	}

	return secrets.Data.Keys, nil
}

func (vc *VaultConnection) GetSecret(path, item string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	fullpath := filepath.Join(path, item)

	secret, err := vc.client.Secrets.KvV2Read(ctx, fullpath, vault.WithMountPath(vc.Config.StoreName))
	if err != nil {
		return nil, err
	}

	js, err := json.Marshal(secret.Data.Data)
	return js, err
}
