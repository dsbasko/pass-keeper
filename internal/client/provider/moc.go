package provider

import (
	"fmt"
	"github.com/dsbasko/pass-keeper/internal/client/models"
	"time"
)

type MockProvider struct {
}

func (provider *MockProvider) GetSecrets(id string) ([]models.SecretForView, error) {
	TestList := []models.SecretForView{{
		ID:             "1",
		VaultID:        id,
		Name:           "xxxxxx",
		Secret:         []byte("1243345"),
		Comment:        "some comment",
		UnPackedSecret: "543143221",
		CreateAt:       time.Now(),
		UpdateAt:       time.Now(),
		Type:           "FILE",
	}, {
		ID:             "2",
		VaultID:        id,
		Name:           "zzzzzzz",
		Secret:         []byte("12315tr4"),
		Comment:        "some comment",
		UnPackedSecret: "12312312",
		CreateAt:       time.Now(),
		UpdateAt:       time.Now(),
		Type:           "STRING",
	}, {
		ID:             "3",
		VaultID:        id,
		Name:           "wwwwwww",
		Secret:         []byte("f423g4342"),
		Comment:        "some comment",
		UnPackedSecret: "d1q2wt2t2",
		CreateAt:       time.Now(),
		UpdateAt:       time.Now(),
		Type:           "KEYPASS",
	}, {
		ID:             "4",
		VaultID:        id,
		Name:           "rrrrr",
		Secret:         []byte("f23g43gf2"),
		Comment:        "some comment",
		UnPackedSecret: "d1d12d1221",
		CreateAt:       time.Now(),
		UpdateAt:       time.Now(),
		Type:           "STRING",
	}, {
		ID:             "5",
		VaultID:        id,
		Name:           "qqqqqqqq",
		Secret:         []byte("e12cfw"),
		Comment:        "some comment",
		UnPackedSecret: "d12e12121t34",
		CreateAt:       time.Now(),
		UpdateAt:       time.Now(),
		Type:           "CARD",
	}}
	return TestList, nil
}

func (provider *MockProvider) GetVaults() ([]models.VaultForView, error) {
	var vaults []models.VaultForView
	for i := 0; i < 10; i++ {
		vaults = append(vaults, models.VaultForView{
			ID:       fmt.Sprintf("%d", i),
			Name:     fmt.Sprintf("someVault%d", i),
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
			Comment:  "i'm comment, mlem",
		})
	}
	return vaults, nil
}
