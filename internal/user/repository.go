package user

import "context"

type UserRepository interface {
	Create(ctx context.Context, user User) (User, error)
	GetByID(ctx context.Context, id int64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, token RefreshToken) (RefreshToken, error)

	GetByTokenHash(
		ctx context.Context,
		tokenHash string,
	) (RefreshToken, error)

	Revoke(ctx context.Context, id int64) error
}
