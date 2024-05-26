package provider

import (
	"fmt"
	"github.com/dsbasko/pass-keeper/internal/client/models"
	"time"
)

var TestList = []models.SecretForView{{
	ID:             "1",
	VaultID:        "1",
	Name:           "xxxxxx",
	Secret:         []byte("1243345"),
	Comment:        "some comment",
	UnPackedSecret: "543143221",
	CreateAt:       time.Now(),
	UpdateAt:       time.Now(),
	Type:           "FILE",
}, {
	ID:             "2",
	VaultID:        "1",
	Name:           "zzzzzzz",
	Secret:         []byte("12315tr4"),
	Comment:        "some comment",
	UnPackedSecret: "12312312",
	CreateAt:       time.Now(),
	UpdateAt:       time.Now(),
	Type:           "STRING",
}, {
	ID:             "3",
	VaultID:        "1",
	Name:           "wwwwwww",
	Secret:         []byte("f423g4342"),
	Comment:        "some comment",
	UnPackedSecret: "d1q2wt2t2",
	CreateAt:       time.Now(),
	UpdateAt:       time.Now(),
	Type:           "KEYPASS",
}, {
	ID:             "4",
	VaultID:        "2",
	Name:           "rrrrr",
	Secret:         []byte("f23g43gf2"),
	Comment:        "some comment",
	UnPackedSecret: "d1d12d1221",
	CreateAt:       time.Now(),
	UpdateAt:       time.Now(),
	Type:           "STRING",
}, {
	ID:             "5",
	VaultID:        "2",
	Name:           "qqqqqqqq",
	Secret:         []byte("e12cfw"),
	Comment:        "some comment",
	UnPackedSecret: "d12e12121t34",
	CreateAt:       time.Now(),
	UpdateAt:       time.Now(),
	Type:           "CARD",
}, {
	ID:             "6",
	VaultID:        "2",
	Name:           "dwqdwqdwqas",
	Secret:         []byte("qwdqw134213"),
	Comment:        "some comment",
	UnPackedSecret: "d12e12121t34",
	CreateAt:       time.Now(),
	UpdateAt:       time.Now(),
	Type:           "CARD",
}, {
	ID:             "7",
	VaultID:        "3",
	Name:           "13131s",
	Secret:         []byte("jiofj923jv32ji29"),
	Comment:        "some comment",
	UnPackedSecret: "fjif3j2j9fj93ffr3",
	CreateAt:       time.Now(),
	UpdateAt:       time.Now(),
	Type:           "CARD",
}, {
	ID:             "8",
	VaultID:        "3",
	Name:           "89329ujufj84",
	Secret:         []byte("qwdqw134213"),
	Comment:        "some comment",
	UnPackedSecret: "cqwwfq13f1",
	CreateAt:       time.Now(),
	UpdateAt:       time.Now(),
	Type:           "CARD",
}, {
	ID:             "9",
	VaultID:        "3",
	Name:           "liujpmnmnnkl",
	Secret:         []byte("qwdqw134213"),
	Comment:        "some comment",
	UnPackedSecret: "f,k23fk3023",
	CreateAt:       time.Now(),
	UpdateAt:       time.Now(),
	Type:           "CARD",
}}

type MockProvider struct {
}

func (provider *MockProvider) CreateSecret(secret models.SecretForView) error {
	return nil
}

func (provider *MockProvider) GetSecret(id string) (models.SecretForView, error) {
	secrets, err := provider.GetSecrets("1")
	if err != nil {
		return models.SecretForView{}, err
	}
	for _, s := range secrets {
		if s.ID == id {
			return s, nil
		}
	}
	return models.SecretForView{}, models.ErrNotFoundSecret
}

func (provider *MockProvider) GetSecrets(id string) ([]models.SecretForView, error) {
	var list []models.SecretForView
	for _, s := range TestList {
		if s.VaultID == id {
			list = append(list, s)
		}
	}
	return list, nil
}

func (provider *MockProvider) GetVaults() ([]models.VaultForView, error) {
	var vaults []models.VaultForView
	for i := 1; i < 4; i++ {
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
