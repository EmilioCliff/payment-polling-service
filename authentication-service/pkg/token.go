package pkg

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
	Username string    `json:"username"`
	UserID   int64     `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTMaker struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

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

func (maker *JWTMaker) CreateToken(username string, userID int64, duration time.Duration) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("error generating token uuid")
	}

	claims := Payload{
		id,
		username,
		userID,
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

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
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
