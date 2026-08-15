package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrEmailAlreadyUsed = errors.New("email already exists")
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

// ! Create New User
func (r *PostgresUserRepository) Create(ctx context.Context, user User) (User, error) {

	query := `
		INSERT INTO users (
			name,
			email,
			password,
			role
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" {
				return User{}, ErrEmailAlreadyUsed
			}
		}

		return User{}, err
	}

	return user, nil
}

// ! Get User by email
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (User, error) {

	query := `
		SELECT 
			id, 
			name,
			email,
			password,
			role,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`

	var user User

	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

// ! Get User by ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id int64) (User, error) {

	query := `
		SELECT
			id, 
			name,
			email,
			password,
			role,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	var user User

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}
