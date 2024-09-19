package utils

// import (
// 	"testing"
// 	"time"

// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/require"
// )

// func TestTokenFunc(t *testing.T) {
// 	testCases := []struct {
// 		name    string
// 		runTest func(maker *JWTMaker, t *testing.T)
// 	}{
// 		{
// 			name: "valid token",
// 			runTest: func(maker *JWTMaker, t *testing.T) {
// 				username := "Emilio Cliff"
// 				duration := time.Minute

// 				token, err := maker.CreateToken(username, duration)
// 				require.NoError(t, err)
// 				require.NotEmpty(t, token)

// 				payload, err := maker.VerifyToken(token)
// 				require.NoError(t, err)
// 				require.NotEmpty(t, payload)

// 				require.NotZero(t, payload.ID)
// 				require.Equal(t, username, payload.Username)
// 			},
// 		},
// 		{
// 			name: "expired token",
// 			runTest: func(maker *JWTMaker, t *testing.T) {
// 				username := "Emilio Cliff"
// 				duration := -time.Minute

// 				token, err := maker.CreateToken(username, duration)
// 				require.NoError(t, err)
// 				require.NotEmpty(t, token)

// 				payload, err := maker.VerifyToken(token)
// 				require.Error(t, err)
// 				// require.EqualError(t, err, ErrTokenExpired.Error())
// 				require.Nil(t, payload)
// 			},
// 		},
// 		{
// 			name: "algorithm none",
// 			runTest: func(maker *JWTMaker, t *testing.T) {
// 				uuid, err := uuid.NewRandom()
// 				require.NoError(t, err)
// 				require.NotEmpty(t, uuid)

// 				claims := Payload{
// 					uuid,
// 					"username",
// 					jwt.RegisteredClaims{
// 						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
// 						IssuedAt:  jwt.NewNumericDate(time.Now()),
// 						Issuer:    "authApp",
// 					},
// 				}

// 				jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
// 				token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
// 				require.NoError(t, err)

// 				payload, err := maker.VerifyToken(token)
// 				require.Error(t, err)
// 				// require.EqualError(t, err, ErrInvalidToken.Error())
// 				require.Nil(t, payload)
// 			},
// 		},
// 		{
// 			name: "wrong issuer",
// 			runTest: func(maker *JWTMaker, t *testing.T) {
// 				uuid, err := uuid.NewRandom()
// 				require.NoError(t, err)
// 				require.NotEmpty(t, uuid)

// 				claims := Payload{
// 					uuid,
// 					"username",
// 					jwt.RegisteredClaims{
// 						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
// 						IssuedAt:  jwt.NewNumericDate(time.Now()),
// 						Issuer:    "wrong-issuer",
// 					},
// 				}

// 				jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
// 				token, err := jwtToken.SignedString(maker.PrivateKey)
// 				require.NoError(t, err)

// 				payload, err := maker.VerifyToken(token)
// 				require.Error(t, err)
// 				require.EqualError(t, err, ErrInvalidIssuer.Error())
// 				require.Nil(t, payload)
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			maker, err := NewJWTMaker("./my_rsa_key", "./my_rsa_key.pub.pem")
// 			require.NoError(t, err)

// 			tc.runTest(maker, t)
// 		})
// 	}
// }
