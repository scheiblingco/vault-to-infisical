package main

import (
	"fmt"
	infisical "github.com/infisical/go-sdk"
	"strings"
)

type InfisicalConnection struct {
	Config   InfisicalConfig
	Settings TransferSettings

	client infisical.InfisicalClientInterface
}

func (ic *InfisicalConnection) Connect() error {
	ic.client = infisical.NewInfisicalClient(infisical.Config{
		SiteUrl: ic.Config.Url,
	})

	_, err := ic.client.Auth().UniversalAuthLogin(ic.Config.ClientId, ic.Config.ClientSecret)
	if err != nil {
		return err
	}

	return nil
}

func (ic *InfisicalConnection) CreateFolder(path, name string) error {
	_, err := ic.client.Folders().Create(
		infisical.CreateFolderOptions{
			Environment: ic.Config.Environment,
			ProjectID:   ic.Config.ProjectId,
			Path:        path,
			Name:        name,
		},
	)

	return err
}

func (ic *InfisicalConnection) EnsureFolderExists(path string) error {
	pathParts := strings.Split(path, "/")
	checkPathParts := []string{}

	for i := 0; i < len(pathParts); i++ {
		if pathParts[i] == "" {
			continue
		}

		checkPath := "/" + strings.Join(checkPathParts, "/")
		fold, err := ic.client.Folders().List(infisical.ListFoldersOptions{
			Environment: ic.Config.Environment,
			ProjectID:   ic.Config.ProjectId,
			Path:        checkPath,
		})
		if err != nil {
			return err
		}

		var exists bool = false

		for _, f := range fold {
			if f.Name == pathParts[i] {
				exists = true
				checkPathParts = append(checkPathParts, pathParts[i])
				break
			}
		}

		if exists {
			continue
		}

		err = ic.CreateFolder(checkPath, pathParts[i])
		if err != nil {
			return err
		}

		checkPathParts = append(checkPathParts, pathParts[i])
	}

	return nil
}

func (ic *InfisicalConnection) IngestSecret(path, item string, data []byte) error {
	_, err := ic.client.Secrets().Retrieve(infisical.RetrieveSecretOptions{
		SecretKey:   item,
		ProjectID:   ic.Config.ProjectId,
		Environment: ic.Config.Environment,
		SecretPath:  path,
	})

	if err == nil {
		if !ic.Settings.Overwrite {
			fmt.Sprintf("secret %s at path %s already exists", item, path)
			return nil
		}

		_, err = ic.client.Secrets().Update(infisical.UpdateSecretOptions{
			SecretKey:                item,
			ProjectID:                ic.Config.ProjectId,
			Environment:              ic.Config.Environment,
			SecretPath:               "/" + path,
			NewSecretValue:           string(data),
			NewSkipMultilineEncoding: false,
		})

	}

	err = ic.EnsureFolderExists(path)
	if err != nil {
		return err
	}

	_, err = ic.client.Secrets().Create(infisical.CreateSecretOptions{
		SecretKey:             item,
		ProjectID:             ic.Config.ProjectId,
		Environment:           ic.Config.Environment,
		SecretPath:            path + "/",
		SecretComment:         "Imported from HCP Vault",
		SkipMultiLineEncoding: false,
		SecretValue:           string(data),
	})

	return err
}
