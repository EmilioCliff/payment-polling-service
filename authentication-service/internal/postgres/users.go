package postgres

import (
	"context"
	"time"
    "database/sql"

    "github.com/lib/pq"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/postgres/generated"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
)

type RegisterUserRequest struct {
    FullName       string `json:"full_name" binding:"required"`
	PaydUsername   string `json:"payd_username" binding:"required"`
	Email          string `json:"email" binding:"required"`
	Password       string `json:"password" binding:"required"`
	UsernameApiKey string `json:"payd_username_api_key" binding:"required"`
	PasswordApiKey string `json:"payd_password_api_key" binding:"required"`
}

type RegisterUserResponse struct {
    FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Store) RegisterUser(ctx context.Context, req RegisterUserRequest) (*RegisterUserResponse, *pkg.Error) {
    hashPassword, err := pkg.GenerateHashPassword(req.Password, s.config.HASH_COST)
	if err != nil {
        return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error hashing password: %s", err)
	}

	encryptPassApi, err := pkg.Encrypt(req.PasswordApiKey, []byte(s.config.ENCRYPTION_KEY))
	if err != nil {
        return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error encrypting password api key: %s", err)
	}

	encryptUserApi, err := pkg.Encrypt(req.UsernameApiKey, []byte(s.config.ENCRYPTION_KEY))
	if err != nil {
        return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error encrypting user api key: %s", err)
	}

	if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
        return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error creating user: %s", err)
	}

	user, err := s.Queries.CreateUser(ctx, generated.CreateUserParams{
		FullName:        req.FullName,
		Email:           req.Email,
		Password:        hashPassword,
		PaydUsername:    req.PaydUsername,
		PaydUsernameKey: encryptUserApi,
		PaydPasswordKey: encryptPassApi,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
                return nil, pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "user already exists: %s", err)
			}
		}

        return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error creating user: %s", err)
	}

	return &RegisterUserResponse{
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

type LoginUserRequest struct {
    Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponse struct {
    AccessToken  string    `json:"access_token"`
	ExpirationAt time.Time `json:"expiration_at"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
}

func (s *Store) LoginUser(ctx context.Context, req LoginUserRequest) (*LoginUserResponse, *pkg.Error) {
    user, err := s.Queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
            pkg.Errorf(pkg.INTERNAL_ERROR, "error creating access token: %s", err)
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "user not found: %s", err)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting user by email: %s", err)
	}

	err = pkg.ComparePasswordAndHash(user.Password, req.Password)
	if err != nil {
		return nil, pkg.Errorf(pkg.AUTHENTICATION_ERROR, "invalid credentials: %s", err)
	}

	accessToken, err := s.maker.CreateToken(user.Email, s.config.TOKEN_DURATION)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error creating access token: %s", err)
	}

	return &LoginUserResponse{
		AccessToken:  accessToken,
		ExpirationAt: time.Now().Add(s.config.TOKEN_DURATION),
		FullName:     user.FullName,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt,
	}, nil
}