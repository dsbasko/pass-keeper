package provider

import "github.com/dsbasko/pass-keeper/internal/client/models"

const (
	Mock   = "mock"
	SQLite = "sqlite"
)

func Init(provType string) models.Providerer {
	// Возвращать интерфейсы плохо, но я не придумал иного решения
	switch provType {
	default:
		return &MockProvider{}
	}
}
