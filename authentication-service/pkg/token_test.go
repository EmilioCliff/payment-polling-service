package pkg

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTokenFunc(t *testing.T) {
	testCases := []struct {
		name    string
		runTest func(maker *JWTMaker, t *testing.T)
	}{
		{
			name: "valid token",
			runTest: func(maker *JWTMaker, t *testing.T) {
				username := "Emilio Cliff"
				userID := int64(1)
				duration := time.Minute

				token, err := maker.CreateToken(username, userID, duration)
				require.NoError(t, err)
				require.NotEmpty(t, token)

				payload, err := maker.VerifyToken(token)
				require.NoError(t, err)
				require.NotEmpty(t, payload)

				require.NotZero(t, payload.ID)
				require.Equal(t, username, payload.Username)
				require.Equal(t, userID, payload.UserID)
			},
		},
		{
			name: "expired token",
			runTest: func(maker *JWTMaker, t *testing.T) {
				username := "Emilio Cliff"
				duration := -time.Minute

				token, err := maker.CreateToken(username, 1, duration)
				require.NoError(t, err)
				require.NotEmpty(t, token)

				payload, err := maker.VerifyToken(token)
				require.Error(t, err)
				require.Nil(t, payload)
			},
		},
		{
			name: "algorithm none",
			runTest: func(maker *JWTMaker, t *testing.T) {
				uuid, err := uuid.NewRandom()
				require.NoError(t, err)
				require.NotEmpty(t, uuid)

				claims := Payload{
					uuid,
					"username",
					1,
					jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
						Issuer:    "authApp",
					},
				}

				jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
				token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
				require.NoError(t, err)

				payload, err := maker.VerifyToken(token)
				require.Error(t, err)
				require.Nil(t, payload)
			},
		},
		{
			name: "wrong issuer",
			runTest: func(maker *JWTMaker, t *testing.T) {
				uuid, err := uuid.NewRandom()
				require.NoError(t, err)
				require.NotEmpty(t, uuid)

				claims := Payload{
					uuid,
					"username",
					1,
					jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
						Issuer:    "wrong-issuer",
					},
				}

				jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
				token, err := jwtToken.SignedString(maker.PrivateKey)
				require.NoError(t, err)

				payload, err := maker.VerifyToken(token)
				require.Error(t, err)
				require.EqualError(t, err, ErrInvalidIssuer.Error())
				require.Nil(t, payload)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
			require.NoError(t, err)

			publicKey := &privateKey.PublicKey

			maker := &JWTMaker{
				PrivateKey: privateKey,
				PublicKey:  publicKey,
			}

			tc.runTest(maker, t)
		})
	}
}
