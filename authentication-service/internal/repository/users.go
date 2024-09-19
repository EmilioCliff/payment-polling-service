package repository

import (
	"context"
	"time"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
)

type User struct {
	ID              int64     `json:"id"`
	FullName        string    `json:"full_name"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	PaydUsername    string    `json:"payd_username"`
	PaydAccountID   string    `json:"payd_account_id"`
	PaydUsernameKey string    `json:"payd_username_key"`
	PaydPasswordKey string    `json:"payd_password_key"`
	CreatedAt       time.Time `json:"created_at"`
}

func (u *User) Validate() error {
	if u.FullName == "" {
		return pkg.Errorf(pkg.INVALID_ERROR, "full_name is required")
	}
	if u.Email == "" {
		return pkg.Errorf(pkg.INVALID_ERROR, "email is required")
	}
	if u.Password == "" {
		return pkg.Errorf(pkg.INVALID_ERROR, "password is required")
	}
	if u.PaydUsername == "" {
		return pkg.Errorf(pkg.INVALID_ERROR, "payd_username is required")
	}
	if u.PaydAccountID == "" {
		return pkg.Errorf(pkg.INVALID_ERROR, "payd_account_id is required")
	}
	if u.PaydUsernameKey == "" {
		return pkg.Errorf(pkg.INVALID_ERROR, "payd_username_key is required")
	}
	if u.PaydPasswordKey == "" {
		return pkg.Errorf(pkg.INVALID_ERROR, "payd_password_key is required")
	}

	return nil
}

type UserRepository interface {
	CreateUser(context.Context, User) (*User, error)
	GetUser(context.Context, string) (*User, error)
	GetUserByID(context.Context, int64) (*User, error)
}