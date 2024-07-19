package main

import (
	"encoding/json"
	"fmt"
	"github.com/infisical/go-sdk/packages/errors"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CopyRecursive(path string, vac VaultConnection, ifc InfisicalConnection) error {
	fmtpath := strings.Replace(path, ".", "_", -1)

	list, err := vac.ListSecrets(path)
	if err != nil {
		panic(err)
	}

	for _, item := range list {
		if strings.HasSuffix(item, "/") {
			err := CopyRecursive(filepath.Join(path, item), vac, ifc)
			if err != nil {
				return err
			}

			continue
		}

		secret, err := vac.GetSecret(path, item)
		if err != nil {
			fmt.Printf("Error finding secret %s in path %s: %s", item, path, err.Error())
			continue
		}

		err = ifc.IngestSecret("/"+fmtpath, item, secret)
		if err != nil {
			if ex, ok := err.(*errors.APIError); ok {
				if strings.Contains(*ex.ErrorMessage, "Rate limit") {
					time.Sleep(60 * time.Second)
					err = ifc.IngestSecret("/"+fmtpath, item, secret)
					if err != nil {
						return err
					}
					continue
				}
			}
			return err
		}

		// Insert a sleep here, because by default Infisical has a rate limit of 3ps/60pm
		time.Sleep(1 * time.Second)
	}

	return nil
}

func main() {
	godotenv.Load()

	var config Config

	content, err := os.ReadFile("./config.json")
	if err == nil {
		err := json.Unmarshal(content, &config)
		if err != nil {
			panic(err)
		}
	}

	vac := VaultConnection{Config: config.Vault}
	err = vac.Connect()
	if err != nil {
		panic(err)
	}

	ifc := InfisicalConnection{Config: config.Infisical, Settings: config.Settings}
	err = ifc.Connect()
	if err != nil {
		panic(err)
	}

	err = CopyRecursive("", vac, ifc)
	if err != nil {
		panic(err)
	}
}
