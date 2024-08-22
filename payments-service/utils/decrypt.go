package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func Decrypt(secret string, key []byte) (string, error) {
	ciphertextBytes, err := base64.URLEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(secret) < aes.BlockSize {
		return "", err
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}
