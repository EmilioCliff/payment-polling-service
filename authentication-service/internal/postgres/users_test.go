package postgres

import (
	"context"
	"database/sql"
	"errors"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/postgres/generated"
	mockdb "github.com/EmilioCliff/payment-polling-app/authentication-service/internal/postgres/mock"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var TestTime = time.Date(2024, time.September, 18, 12, 0, 0, 0, time.UTC)

func NewTestUserRepository() *UserRepository {
	store := NewStore(pkg.Config{
		HASH_COST:      8,
		ENCRYPTION_KEY: "12345678901234567890123456789012",
	})
	store.conn = nil
	u := NewUserService(store)

	return u
}

func TestUserRepository_CreateUser(t *testing.T) {
	s := NewTestUserRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)

	s.queries = mockQueries

	tests := []struct {
		name       string
		user       repository.User
		buildStubs func(*mockdb.MockQuerier, repository.User)
		wantErr    bool
	}{
		{
			name: "success",
			user: randomUser(),
			buildStubs: func(mockQueries *mockdb.MockQuerier, req repository.User) {
				mockQueries.EXPECT().CreateUser(gomock.Any(), gomock.AssignableToTypeOf(generated.CreateUserParams{})).
					DoAndReturn(func(_ context.Context, params generated.CreateUserParams) (generated.User, error) {
						require.Equal(t, req.FullName, params.FullName)
						require.Equal(t, req.Email, params.Email)
						require.Equal(t, req.PaydUsername, params.PaydUsername)
						require.Equal(t, req.PaydAccountID, params.PaydAccountID)

						require.NotEmpty(t, params.Password)
						require.NotEmpty(t, params.PaydUsernameKey)
						require.NotEmpty(t, params.PaydPasswordKey)

						hashPassword, _ := pkg.GenerateHashPassword("password", 10)

						return generated.User{
							ID:              32,
							FullName:        req.FullName,
							Email:           req.Email,
							Password:        hashPassword,
							PaydUsername:    req.PaydUsername,
							PaydUsernameKey: req.PaydUsernameKey,
							PaydPasswordKey: req.PaydPasswordKey,
							PaydAccountID:   req.PaydAccountID,
							CreatedAt:       TestTime,
						}, nil
					}).Times(1)
			},
			wantErr: false,
		},
		{
			name: "missing_values",
			user: repository.User{},
			buildStubs: func(mockQueries *mockdb.MockQuerier, _ repository.User) {
				mockQueries.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			wantErr: true,
		},
		{
			name: "db_error",
			user: randomUser(),
			buildStubs: func(mockQueries *mockdb.MockQuerier, req repository.User) {
				mockQueries.EXPECT().CreateUser(gomock.Any(), gomock.AssignableToTypeOf(generated.CreateUserParams{})).
					DoAndReturn(func(_ context.Context, params generated.CreateUserParams) (generated.User, error) {
						require.Equal(t, req.FullName, params.FullName)
						require.Equal(t, req.Email, params.Email)
						require.Equal(t, req.PaydUsername, params.PaydUsername)
						require.Equal(t, req.PaydAccountID, params.PaydAccountID)

						require.NotEmpty(t, params.Password)
						require.NotEmpty(t, params.PaydUsernameKey)
						require.NotEmpty(t, params.PaydPasswordKey)

						err := errors.New("db error")

						return generated.User{}, err
					}).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "unique_violation",
			user: randomUser(),
			buildStubs: func(mockQueries *mockdb.MockQuerier, req repository.User) {
				mockQueries.EXPECT().CreateUser(gomock.Any(), gomock.AssignableToTypeOf(generated.CreateUserParams{})).
					DoAndReturn(func(_ context.Context, params generated.CreateUserParams) (generated.User, error) {
						require.Equal(t, req.FullName, params.FullName)
						require.Equal(t, req.Email, params.Email)
						require.Equal(t, req.PaydUsername, params.PaydUsername)
						require.Equal(t, req.PaydAccountID, params.PaydAccountID)

						require.NotEmpty(t, params.Password)
						require.NotEmpty(t, params.PaydUsernameKey)
						require.NotEmpty(t, params.PaydPasswordKey)

						pqErr := &pq.Error{
							Code: "23505", // unique_violation
						}

						return generated.User{}, pqErr
					}).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries, tc.user)

			got, err := s.CreateUser(context.Background(), tc.user)
			if !tc.wantErr {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.False(t, tc.wantErr)
				require.Equal(t, got.FullName, tc.user.FullName)
				require.Equal(t, got.Email, tc.user.Email)
				require.Equal(t, got.CreatedAt, tc.user.CreatedAt)
			} else {
				require.Error(t, err)
				require.Equal(t, got, (*repository.User)(nil))
			}
		})
	}
}

func TestUserRepository_GetUser(t *testing.T) {
	s := NewTestUserRepository()

	ctrl := gomock.NewController(t)

	mockUserRepo := mockdb.NewMockQuerier(ctrl)

	s.queries = mockUserRepo

	tests := []struct {
		user       repository.User
		name       string
		buildStubs func(*mockdb.MockQuerier, repository.User)
		wantErr    bool
	}{
		{
			name: "success",
			user: randomUser(),
			buildStubs: func(mockQueries *mockdb.MockQuerier, req repository.User) {
				mockQueries.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(req.Email)).
					Return(repositoryUserToGenerated(req), nil).Times(1)
			},
			wantErr: false,
		},
		{
			name: "no user",
			user: randomUser(),
			buildStubs: func(mockQueries *mockdb.MockQuerier, req repository.User) {
				mockQueries.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(req.Email)).
					Return(generated.User{}, sql.ErrNoRows).Times(1)
			},
			wantErr: true,
		},
		{
			name: "db error",
			user: randomUser(),
			buildStubs: func(mockQueries *mockdb.MockQuerier, req repository.User) {
				mockQueries.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(req.Email)).
					Return(generated.User{}, errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockUserRepo, tc.user)

			got, err := s.GetUser(context.Background(), tc.user.Email)
			if !tc.wantErr {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.False(t, tc.wantErr)
				require.Equal(t, got, &tc.user)
			} else {
				require.Error(t, err)
				require.Nil(t, got)
				require.Equal(t, got, (*repository.User)(nil))
			}
		})
	}
}

func TestUserRepository_GetUserByID(t *testing.T) {
	s := NewTestUserRepository()

	ctrl := gomock.NewController(t)

	mockUserRepo := mockdb.NewMockQuerier(ctrl)

	s.queries = mockUserRepo

	tests := []struct {
		name       string
		user       repository.User
		buildStubs func(*mockdb.MockQuerier, repository.User)
		wantErr    bool
	}{
		{
			name: "success",
			user: randomUser(),
			buildStubs: func(mockQueries *mockdb.MockQuerier, req repository.User) {
				mockQueries.EXPECT().GetUser(gomock.Any(), gomock.Eq(req.ID)).
					Return(repositoryUserToGenerated(req), nil).Times(1)
			},
			wantErr: false,
		},
		{
			name: "no user",
			user: randomUser(),
			buildStubs: func(mockQueries *mockdb.MockQuerier, req repository.User) {
				mockQueries.EXPECT().GetUser(gomock.Any(), gomock.Eq(req.ID)).
					Return(generated.User{}, sql.ErrNoRows).Times(1)
			},
			wantErr: true,
		},
		{
			name: "db error",
			user: randomUser(),
			buildStubs: func(mockQueries *mockdb.MockQuerier, req repository.User) {
				mockQueries.EXPECT().GetUser(gomock.Any(), gomock.Eq(req.ID)).
					Return(generated.User{}, errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockUserRepo, tc.user)

			got, err := s.GetUserByID(context.Background(), tc.user.ID)
			if !tc.wantErr {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.False(t, tc.wantErr)
				require.Equal(t, got, &tc.user)
			} else {
				require.Error(t, err)
				require.Nil(t, got)
				require.Equal(t, got, (*repository.User)(nil))
			}
		})
	}
}

func randomUser() repository.User {
	return repository.User{
		ID:              int64(rand.IntN(100)),
		FullName:        gofakeit.Name(),
		Email:           gofakeit.Email(),
		Password:        gofakeit.Password(true, true, true, true, true, 7),
		PaydUsername:    gofakeit.Username(),
		PaydAccountID:   gofakeit.UUID(),
		PaydUsernameKey: gofakeit.UUID(),
		PaydPasswordKey: gofakeit.UUID(),
		CreatedAt:       TestTime,
	}
}

func repositoryUserToGenerated(user repository.User) generated.User {
	return generated.User{
		ID:              user.ID,
		FullName:        user.FullName,
		Email:           user.Email,
		Password:        user.Password,
		PaydUsername:    user.PaydUsername,
		PaydAccountID:   user.PaydAccountID,
		PaydUsernameKey: user.PaydUsernameKey,
		PaydPasswordKey: user.PaydPasswordKey,
		CreatedAt:       user.CreatedAt,
	}
}
