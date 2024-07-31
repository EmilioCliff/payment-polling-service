package utils

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken  = errors.New("token is invalid")
	ErrInvalidIssuer = errors.New("invalid issuer")
	ErrTokenExpired  = fmt.Errorf("token is expired")
)

type Payload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"fullname"`
	jwt.RegisteredClaims
}

type JWTMaker struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

// NewJWTMaker creates a new JWTMaker instance.
//
// It reads the private key from the file "./my_rsa_key" and the public key from the file "./my_rsa_key.pub.pem".
// If any error occurs during the reading or parsing of the keys, an error is returned.
//
// Returns a pointer to a JWTMaker instance and an error.
func NewJWTMaker(privateKeyPath string, publicKeyPath string) (*JWTMaker, error) {
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %s", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %s", err)
	}

	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %s", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %s", err)
	}

	maker := &JWTMaker{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	return maker, nil
}

// CreateToken generates a JSON Web Token (JWT) for the given username and duration.
//
// Parameters:
// - username: the username for which the token is being generated.
// - duration: the duration for which the token is valid.
//
// Returns:
// - string: the generated JWT.
// - error: an error if there was a problem generating the JWT.
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("error generating token uuid")
	}

	claims := Payload{
		uuid,
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "authApp",
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token, err := jwtToken.SignedString(maker.PrivateKey)
	return token, err
}

// VerifyToken verifies the validity of a JSON Web Token (JWT) and returns its payload if valid.
//
// Parameters:
// - token: the JWT to be verified.
//
// Returns:
// - *Payload: the payload of the valid JWT.
// - error: an error if the JWT is invalid or expired.
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, ErrInvalidToken
		}

		return maker.PublicKey, nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	if payload.RegisteredClaims.Issuer != "authApp" {
		return nil, ErrInvalidIssuer
	}

	if payload.RegisteredClaims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	return payload, nil
}
