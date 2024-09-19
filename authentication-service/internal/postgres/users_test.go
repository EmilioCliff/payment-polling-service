package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/postgres/generated"
	mockdb "github.com/EmilioCliff/payment-polling-app/authentication-service/internal/postgres/mock"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	TestTime = time.Date(2024, time.September, 18, 12, 0, 0, 0, time.UTC)
)

func NewTestUserRepository() *UserRepository {
	store := NewStore(pkg.Config{
		HASH_COST: 8,
		ENCRYPTION_KEY: "12345678901234567890123456789012",
	}, pkg.JWTMaker{})
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

	tests := []struct{
		name string
		user repository.User
		buildStubs func(*mockdb.MockQuerier)
		want *repository.User
		wantErr bool
	}{
		{
			name: "success",
			user: repository.User{
					FullName: "Jane",
					Email: "jane@gmail.com",
					Password: "password",
					PaydUsername: "username",
					PaydUsernameKey: "user_key",
					PaydPasswordKey: "pass_key",
					PaydAccountID: "account_id",
				},
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().CreateUser(gomock.Any(), gomock.AssignableToTypeOf(generated.CreateUserParams{})).
					DoAndReturn(func(_ context.Context, params generated.CreateUserParams) (generated.User, error) {
						require.Equal(t, "Jane", params.FullName)
						require.Equal(t, "jane@gmail.com", params.Email)
						require.Equal(t, "username", params.PaydUsername)
						require.Equal(t, "account_id", params.PaydAccountID)
						
						require.NotEmpty(t, params.Password)
						require.NotEmpty(t, params.PaydUsernameKey)
						require.NotEmpty(t, params.PaydPasswordKey)

						hashPassword, _ := pkg.GenerateHashPassword("password", 10)
						
						return generated.User{
								ID: 32,
								FullName: "Jane",
								Email: "jane@gmail.com",
								Password: hashPassword,
								PaydUsername: "username",
								PaydUsernameKey: "user_key",
								PaydPasswordKey: "pass_key",
								PaydAccountID: "account_id",
								CreatedAt: TestTime,
							}, nil
					}).Times(1)
			},
			want: &repository.User{
				FullName: "Jane",
				Email: "jane@gmail.com",
				CreatedAt: TestTime,
			},
			wantErr: false,
		},
		{
			name: "missing_values",
			user: repository.User{},
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			want: nil,
			wantErr: true,
		},
		{
			name: "db_error",
			user: repository.User{
				FullName: "Jane",
				Email: "jane@gmail.com",
				Password: "password",
				PaydUsername: "username",
				PaydUsernameKey: "user_key",
				PaydPasswordKey: "pass_key",
				PaydAccountID: "account_id",
			},
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().CreateUser(gomock.Any(), gomock.AssignableToTypeOf(generated.CreateUserParams{})).
				DoAndReturn(func(_ context.Context, params generated.CreateUserParams) (*generated.User, error){
					require.Equal(t, "Jane", params.FullName)
					require.Equal(t, "jane@gmail.com", params.Email)
					require.Equal(t, "username", params.PaydUsername)
					require.Equal(t, "account_id", params.PaydAccountID)
					
					require.NotEmpty(t, params.Password)
					require.NotEmpty(t, params.PaydUsernameKey)
					require.NotEmpty(t, params.PaydPasswordKey)
					err := errors.New("db error")
					return nil, err
				}).
					Times(1)
			},
			want: nil,
			wantErr: true,
		},
		{
			name: "unique_violation",
			user: repository.User{
				FullName: "Jane",
				Email: "jane@gmail.com",
				Password: "password",
				PaydUsername: "username",
				PaydUsernameKey: "user_key",
				PaydPasswordKey: "pass_key",
				PaydAccountID: "account_id",
			},
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().CreateUser(gomock.Any(), gomock.AssignableToTypeOf(generated.CreateUserParams{})).
				DoAndReturn(func(_ context.Context, params generated.CreateUserParams) (*generated.User, error){
					require.Equal(t, "Jane", params.FullName)
					require.Equal(t, "jane@gmail.com", params.Email)
					require.Equal(t, "username", params.PaydUsername)
					require.Equal(t, "account_id", params.PaydAccountID)
					
					require.NotEmpty(t, params.Password)
					require.NotEmpty(t, params.PaydUsernameKey)
					require.NotEmpty(t, params.PaydPasswordKey)
					pqErr := &pq.Error{
						Code: "23505", // unique_violation
					}
					return nil, pqErr
				}).
					Times(1)
			},
			want: nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries)
			got, err := s.CreateUser(context.Background(), tc.user)
			if !tc.wantErr {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.False(t, tc.wantErr)
				require.Equal(t, got, tc.want)
			} else {
				require.Error(t, err)
				require.Nil(t, got)
				require.Equal(t, got, tc.want)
			}
		})
	}
}

func TestUserRepository_GetUser(t *testing.T) {
	s := NewTestUserRepository()

	ctrl := gomock.NewController(t)

	mockUserRepo := mockdb.NewMockQuerier(ctrl)

	s.queries = mockUserRepo

	tests := []struct{
		name string
		email string
		buildStubs func(*mockdb.MockQuerier)
		want *repository.User
		wantErr bool
	}{
		{
			name: "success",
			email: "jane@gmail",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq("jane@gmail")).
					Return(generated.User{
						ID: 32,
						FullName: "Jane",
						Email: "jane@gmail",
						Password: "password",
						PaydUsername: "username",
						PaydUsernameKey: "user_key",
						PaydPasswordKey: "pass_key",
						PaydAccountID: "account_id",
						CreatedAt: TestTime,
					}, nil).Times(1)
			},
			want: &repository.User{
				ID: 32,
				FullName: "Jane",
				Email: "jane@gmail",
				Password: "password",
				PaydUsername: "username",
				PaydUsernameKey: "user_key",
				PaydPasswordKey: "pass_key",
				PaydAccountID: "account_id",
				CreatedAt: TestTime,
			},
			wantErr: false,
		},
		{
			name: "no user",
			email: "jane@gmail",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq("jane@gmail")).
					Return(generated.User{}, sql.ErrNoRows).Times(1)
			},
			want: nil,
			wantErr: true,
		},
		{
			name: "db error",
			email: "jane@gmail",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq("jane@gmail")).
					Return(generated.User{}, errors.New("db error")).Times(1)
			},
			want: nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockUserRepo)
			got, err := s.GetUser(context.Background(), tc.email)
			if !tc.wantErr {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.False(t, tc.wantErr)
				require.Equal(t, got, tc.want)
			} else {
				require.Error(t, err)
				require.Nil(t, got)
				require.Equal(t, got, tc.want)
			}
		})
	}
}

func TestUserRepository_GetUserByID(t *testing.T) {
	s := NewTestUserRepository()

	ctrl := gomock.NewController(t)

	mockUserRepo := mockdb.NewMockQuerier(ctrl)

	s.queries = mockUserRepo

	tests := []struct{
		name string
		id int64
		buildStubs func(*mockdb.MockQuerier)
		want *repository.User
		wantErr bool
	}{
		{
			name: "success",
			id: 32,
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().GetUser(gomock.Any(), gomock.Eq(int64(32))).
					Return(generated.User{
						ID: 32,
						FullName: "Jane",
						Email: "jane@gmail",
						Password: "password",
						PaydUsername: "username",
						PaydUsernameKey: "user_key",
						PaydPasswordKey: "pass_key",
						PaydAccountID: "account_id",
						CreatedAt: TestTime,
					}, nil).Times(1)
			},
			want: &repository.User{
				ID: 32,
				FullName: "Jane",
				Email: "jane@gmail",
				Password: "password",
				PaydUsername: "username",
				PaydUsernameKey: "user_key",
				PaydPasswordKey: "pass_key",
				PaydAccountID: "account_id",
				CreatedAt: TestTime,
			},
			wantErr: false,
		},
		{
			name: "no user",
			id: 33,
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().GetUser(gomock.Any(), gomock.Eq(int64(33))).
					Return(generated.User{}, sql.ErrNoRows).Times(1)
			},
			want: nil,
			wantErr: true,
		},
		{
			name: "db error",
			id: 32,
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().GetUser(gomock.Any(), gomock.Eq(int64(32))).
					Return(generated.User{}, errors.New("db error")).Times(1)
			},
			want: nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockUserRepo)
			got, err := s.GetUserByID(context.Background(), tc.id)
			if !tc.wantErr {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.False(t, tc.wantErr)
				require.Equal(t, got, tc.want)
			} else {
				require.Error(t, err)
				require.Nil(t, got)
				require.Equal(t, got, tc.want)
			}
		})
	}
}