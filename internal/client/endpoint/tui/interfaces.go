package tui

import "github.com/dsbasko/pass-keeper/internal/client/models"

type Providerer interface {
	GetVaults() ([]models.VaultForView, error)
	GetSecrets(vaultID string) ([]models.SecretForView, error)
	//GetVault(vaultID string) (models.VaultForView, error)
}
