package config

import (
	"sync"

	goCFG "github.com/dsbasko/go-cfg"

	errWrapper "github.com/dsbasko/pass-keeper/pkg/err-wrapper"
)

var (
	once sync.Once
	cfg  config
)

func Init(filePath string) (err error) {
	defer errWrapper.PtrWithOP(&err, "config.Init")

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
