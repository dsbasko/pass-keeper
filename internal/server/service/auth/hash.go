package auth

import (
	"golang.org/x/crypto/bcrypt"

	errWrapper "github.com/dsbasko/pass-keeper/pkg/err-wrapper"
)

func hashPassword(password []byte) (_ []byte, err error) {
	defer errWrapper.PtrWithOP(&err, "auth.hashPassword")

	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, errWrapper.WithOP(err, "bcrypt.GenerateFromPassword")
	}

	return hash, nil
}

func compareHashAndPassword(hash, password []byte) (_ bool, err error) {
	defer errWrapper.PtrWithOP(&err, "auth.compareHashAndPassword")

	if err = bcrypt.CompareHashAndPassword(hash, password); err != nil {
		return false, errWrapper.WithOP(err, "bcrypt.CompareHashAndPassword")
	}

	return true, nil
}
