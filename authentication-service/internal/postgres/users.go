package postgres

import (
	"context"
	"database/sql"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/postgres/generated"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
	"github.com/lib/pq"
)

var _ repository.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *Store
	queries generated.Querier
}

func NewUserService(db *Store) *UserRepository {
	queries := generated.New(db.conn)

	return &UserRepository{
		db: db,
		queries: queries,
	}
}

func (s *UserRepository) CreateUser(ctx context.Context, u repository.User) (*repository.User, error) {
	err := u.Validate()
	if err != nil {
		return nil, err
	}

	hashPassword, err := pkg.GenerateHashPassword(u.Password, s.db.config.HASH_COST)
	if err != nil {
        return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error hashing password: %s", err)
	}

	encryptPassApi, err := pkg.Encrypt(u.PaydPasswordKey, []byte(s.db.config.ENCRYPTION_KEY))
	if err != nil {
        return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error encrypting password api key: %s", err)
	}

	encryptUserApi, err := pkg.Encrypt(u.PaydUsernameKey, []byte(s.db.config.ENCRYPTION_KEY))
	if err != nil {
        return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error encrypting user api key: %s", err)
	}

	user, err := s.queries.CreateUser(ctx, generated.CreateUserParams{
		FullName:        u.FullName,
		Email:           u.Email,
		Password:        hashPassword,
		PaydUsername:    u.PaydUsername,
		PaydUsernameKey: encryptUserApi,
		PaydPasswordKey: encryptPassApi,
		PaydAccountID: u.PaydAccountID,
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

	return &repository.User{
		FullName: user.FullName,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *UserRepository) GetUser(ctx context.Context, email string) (*repository.User, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "user not found: %s", err)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error geting user: %s", err)
	}

	return &repository.User{
		ID: user.ID,
		FullName: user.FullName,
		Email: user.Email,
		Password: user.Password,
		PaydUsername: user.PaydUsername,
		PaydAccountID: user.PaydAccountID,
		PaydUsernameKey: user.PaydUsernameKey,
		PaydPasswordKey: user.PaydPasswordKey,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *UserRepository) GetUserByID(ctx context.Context, id int64) (*repository.User, error) {
	user, err := s.queries.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "user not found: %s", err)
		}
		return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "error geting user: %s", err)
	}

	return &repository.User{
		ID: user.ID,
		FullName: user.FullName,
		Email: user.Email,
		Password: user.Password,
		PaydUsername: user.PaydUsername,
		PaydAccountID: user.PaydAccountID,
		PaydUsernameKey: user.PaydUsernameKey,
		PaydPasswordKey: user.PaydPasswordKey,
		CreatedAt: user.CreatedAt,
	}, nil
}