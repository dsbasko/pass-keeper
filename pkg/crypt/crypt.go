// Package crypt предоставляет функции шифрования и дешифрования данных.
package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"

	"github.com/dsbasko/pass-keeper/pkg/errors"
)

// Encrypt - шифрует данные
//
// На вход принимает ключ и данные в формате []byte.
// Отдает зашифрованные данные в формате []byte и ошибку.
//
// Для шифрования используется AES-GCM.
//
//	encryptData, err := crypt.Encrypt(someKey, originalData)
//	if err != nil {
//		panic(err)
//	}
func Encrypt(key, data []byte) (encryptData []byte, err error) {
	defer errors.ErrorPtrWithOP(&err, "crypt.Encrypt")

	aesBlock, err := aes.NewCipher(hash(key))
	if err != nil {
		return []byte{}, errors.ErrorWithOP(err, "aes.NewCipher")
	}

	aesGCM, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return []byte{}, errors.ErrorWithOP(err, "cipher.NewGCM")
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		panic(err)
	}

	cipherText := aesGCM.Seal(nonce, nonce, data, nil)
	return cipherText, nil
}

// Decrypt - дешифрует данные
//
// На вход принимает ключ и зашифрованные данные в формате []byte.
// Отдает расшифрованные данные в формате []byte и ошибку.
//
// Для дешифрования используется AES-GCM.
//
//	decryptData, err := crypt.Decrypt(someKey, encryptedData)
//	if err != nil {
//			panic(err)
//	}
func Decrypt(key, encryptData []byte) (decryptData []byte, err error) {
	defer errors.ErrorPtrWithOP(&err, "crypt.Decrypt")

	aesBlock, err := aes.NewCipher(hash(key))
	if err != nil {
		return []byte{}, errors.ErrorWithOP(err, "aes.NewCipher")
	}

	aesGCM, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return []byte{}, errors.ErrorWithOP(err, "cipher.NewGCM")
	}

	nonceSize := aesGCM.NonceSize()
	nonce, cipherText := encryptData[:nonceSize], encryptData[nonceSize:]

	decryptData, err = aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return []byte{}, errors.ErrorWithOP(err, "aesGCM.Open")
	}

	return decryptData, nil
}

// hash - хэширует данные
//
// На вход принимает данные в формате []byte.
// Отдает хэш в формате []byte.
func hash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}
