package pkg

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	password := gofakeit.Password(true, true, true, true, false, 8)

	hashPassword, err := GenerateHashPassword(password, 10)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword)

	err = ComparePasswordAndHash(hashPassword, password)
	require.NoError(t, err)
}
