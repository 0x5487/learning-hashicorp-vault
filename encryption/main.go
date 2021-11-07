package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
	"gocloud.dev/secrets/hashivault"
)

func main() {
	ctx := context.Background()

	// Get a client to use with the Vault API.
	client, err := hashivault.Dial(ctx, &hashivault.Config{
		Token: "s.oR23Vk0shqtcIfgQ0RnNRwz0",
		APIConfig: api.Config{
			Address: "http://127.0.0.1:8200",
		},
	})
	if err != nil {
		panic(err)
	}

	keyID := "node1"
	svc := NewVaultService(client)
	chiperBytes, err := svc.Encrypt(ctx, "hello", keyID)
	if err != nil {
		panic(err)
	}

	plainText, err := svc.Decrypt(ctx, chiperBytes, keyID)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(plainText))
}

type VaultService struct {
	client *api.Client
}

func NewVaultService(client *api.Client) *VaultService {
	return &VaultService{
		client: client,
	}
}

func (v *VaultService) Encrypt(ctx context.Context, content string, keyID string) ([]byte, error) {
	keeper := hashivault.OpenKeeper(v.client, keyID, nil)
	defer keeper.Close()

	return keeper.Encrypt(ctx, []byte(content))
}

func (v *VaultService) Decrypt(ctx context.Context, cipherText []byte, keyID string) ([]byte, error) {
	keeper := hashivault.OpenKeeper(v.client, keyID, nil)
	defer keeper.Close()

	return keeper.Decrypt(ctx, cipherText)
}
