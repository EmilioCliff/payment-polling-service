package mock

import (
	"context"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/repository"
)

var _ repository.UserRepository = (*MockUsersRepositry)(nil)

type MockUsersRepositry struct {
	CreateUserFunc func(repository.User) (*repository.User, error)
	GetUserFunc func(string) (*repository.User, error)
	GetUserByIDFunc func(int64) (*repository.User, error)
}

func (u *MockUsersRepositry) GetUser(ctx context.Context, email string) (*repository.User, error) {
	return u.GetUserFunc(email)
}

func (u *MockUsersRepositry) CreateUser(ctx context.Context, user repository.User) (*repository.User, error) {
	return u.CreateUserFunc(user)
}

func (u *MockUsersRepositry)  GetUserByID(ctx context.Context, id int64) (*repository.User, error) {
	return u.GetUserByIDFunc(id)
}