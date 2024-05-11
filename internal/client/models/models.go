package models

import "time"

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
