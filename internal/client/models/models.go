package models

import (
	"errors"
	"time"
)

type Providerer interface {
	GetVaults() ([]VaultForView, error)
	GetSecrets(vaultID string) ([]SecretForView, error)
	GetSecret(secretID string) (SecretForView, error)
	//GetVault(vaultID string) (models.VaultForView, error)
}

type VaultForView struct {
	ID       string
	Name     string
	CreateAt time.Time
	UpdateAt time.Time
	Comment  string
}

type SecretForView struct {
	ID             string
	Name           string
	Secret         []byte
	UnPackedSecret string
	Comment        string
	CreateAt       time.Time
	UpdateAt       time.Time
	VaultID        string
	Type           string
}

var (
	ErrNotFoundSecret = errors.New("not found secret")
)
