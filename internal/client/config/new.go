package config

import (
	"sync"

	goCFG "github.com/dsbasko/go-cfg"

	"github.com/dsbasko/pass-keeper/pkg/errors"
)

var (
	once sync.Once
	cfg  config
)

func Init(filePath string) (err error) {
	defer errors.ErrorPtrWithOP(&err, "config.Init")

	once.Do(func() {
		if err = goCFG.ReadFile(filePath, &cfg); err != nil {
			return
		}

		if err = goCFG.ReadEnv(&cfg); err != nil {
			return
		}
	})

	return err
}

func MustInit(envFilePath string) {
	if err := Init(envFilePath); err != nil {
		panic(err)
	}
}

func Get() *config {
	return &cfg
}
