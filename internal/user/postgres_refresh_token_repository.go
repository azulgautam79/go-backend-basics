package user

import (
	"context"
	"database/sql"
)

type PostgresRefreshTokenRepository struct {
	db *sql.DB
}

func NewPostgresRefreshTokenRepository(
	db *sql.DB,
) *PostgresRefreshTokenRepository {
	return &PostgresRefreshTokenRepository{
		db: db,
	}
}

// ! Create refresh token
func (r *PostgresRefreshTokenRepository) Create(
	ctx context.Context,
	token RefreshToken,
) (RefreshToken, error) {

	query := `
		INSERT INTO refresh_tokens (
			user_id,
			token_hash,
			expires_at
		)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
	).Scan(
		&token.ID,
		&token.CreatedAt,
	)

	if err != nil {
		return RefreshToken{}, err
	}

	return token, nil
}

// ! Get by Token Hash
func (r *PostgresRefreshTokenRepository) GetByTokenHash(
	ctx context.Context,
	tokenHash string,
) (RefreshToken, error) {

	query := `
		SELECT 
			id,
			user_id,
			token_hash,
			expires_at,
			revoked_at,
			created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	var token RefreshToken

	err := r.db.QueryRowContext(
		ctx,
		query,
		tokenHash,
	).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.RevokedAt,
		&token.CreatedAt,
	)

	if err != nil {
		return RefreshToken{}, err
	}

	return token, nil
}

// ! Revoke Refresh Token
func (r *PostgresRefreshTokenRepository) Revoke(
	ctx context.Context,
	id int64,
) error {

	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE id = $1
			AND revoked_at IS NULL
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		id,
	)
	return err
}
