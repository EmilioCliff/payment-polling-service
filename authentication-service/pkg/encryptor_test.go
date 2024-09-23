package pkg

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func TestEncryptor(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		encryptWord func(t *testing.T, word string, key string) string
		decryptWord func(t *testing.T, ciphertext string, word string, key string)
	}{
		{
			name: "success",
			key:  "abcdefghijklmnopqrstuvwxyzabcdef",
			encryptWord: func(t *testing.T, plainText string, key string) string {
				ciphertext, err := Encrypt(plainText, []byte(key))
				require.NoError(t, err)
				require.NotEmpty(t, ciphertext)

				return ciphertext
			},
			decryptWord: func(t *testing.T, cipherText string, word string, key string) {
				decrypted, err := Decrypt(cipherText, []byte(key))
				require.NoError(t, err)
				require.NotEmpty(t, decrypted)
				require.Equal(t, word, decrypted)
			},
		},
		{
			name: "invalid key",
			key:  "INVALID_KEY",
			encryptWord: func(t *testing.T, word string, key string) string {
				ciphertext, err := Encrypt(word, []byte(key))
				require.Error(t, err)
				require.Empty(t, ciphertext)

				return ciphertext
			},
			decryptWord: func(t *testing.T, cipherText string, _ string, key string) {
				decrypted, err := Decrypt(cipherText, []byte(key))
				require.Error(t, err)
				require.Empty(t, decrypted)
			},
		},
		{
			name: "empty ciphertext",
			key:  "abcdefghijklmnopqrstuvwxyzabcdef",
			encryptWord: func(_ *testing.T, _ string, _ string) string {
				return ""
			},
			decryptWord: func(t *testing.T, cipherText string, _ string, key string) {
				decrypted, err := Decrypt(cipherText, []byte(key))
				require.Error(t, err)
				require.Empty(t, decrypted)
			},
		},
		{
			name: "invalid ciphertext",
			key:  "abcdefghijklmnopqrstuvwxyzabcdef",
			encryptWord: func(_ *testing.T, _ string, _ string) string {
				return "INVALID_CIPHERTEXT"
			},
			decryptWord: func(t *testing.T, cipherText string, _ string, key string) {
				decrypted, err := Decrypt(cipherText, []byte(key))
				require.Error(t, err)
				require.Empty(t, decrypted)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			word := gofakeit.Word()
			ciphertext := tc.encryptWord(t, word, tc.key)
			tc.decryptWord(t, ciphertext, word, tc.key)
		})
	}
}
