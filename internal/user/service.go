package user

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrRefreshTokenInvalid = errors.New("invalid refresh token")
)

type AuthService struct {
	users         UserRepository
	refreshTokens RefreshTokenRepository
	jwtService    *JWTService
}

func NewAuthService(
	users UserRepository,
	refreshTokens RefreshTokenRepository,
	jwtService *JWTService,
) *AuthService {
	return &AuthService{
		users:         users,
		refreshTokens: refreshTokens,
		jwtService:    jwtService,
	}
}

// ! Register
func (s *AuthService) Register(
	ctx context.Context,
	name string,
	email string,
	password string,
) (User, error) {
	passwordHash, err := HashPassword(password)

	if err != nil {
		return User{}, err
	}

	user := User{
		Name:     name,
		Email:    email,
		Password: passwordHash,
		Role:     RoleCustomer,
	}

	return s.users.Create(ctx, user)
}

// ! Login
func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (LoginResponse, error) {
	user, err := s.users.GetByEmail(ctx, email)

	if err != nil {
		return LoginResponse{}, err
	}

	if !CheckPassword(password, user.Password) {
		return LoginResponse{}, ErrInvalidCredentials
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user)

	if err != nil {
		return LoginResponse{}, err
	}

	rawRefreshToken, err := GenerateRefreshToken()

	if err != nil {
		return LoginResponse{}, err
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	_, err = s.refreshTokens.Create(
		ctx,
		RefreshToken{
			UserID:    user.ID,
			TokenHash: HashRefreshToken(rawRefreshToken),
			ExpiresAt: expiresAt,
		},
	)

	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// ! Refresh
func (s *AuthService) Refresh(
	ctx context.Context,
	rawToken string,
) (LoginResponse, error) {
	tokenHash := HashRefreshToken(rawToken)

	refreshToken, err := s.refreshTokens.GetByTokenHash(
		ctx,
		tokenHash,
	)

	if err != nil {
		return LoginResponse{}, ErrRefreshTokenInvalid
	}

	if refreshToken.RevokedAt != nil {
		return LoginResponse{}, ErrRefreshTokenInvalid
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return LoginResponse{}, ErrRefreshTokenInvalid
	}

	user, err := s.users.GetByID(
		ctx,
		refreshToken.UserID,
	)

	if err != nil {
		return LoginResponse{}, ErrRefreshTokenInvalid
	}

	//* Rotate old refresh token.
	if err := s.refreshTokens.Revoke(
		ctx,
		refreshToken.ID,
	); err != nil {
		return LoginResponse{}, err
	}

	newRawToken, err := GenerateRefreshToken()

	if err != nil {
		return LoginResponse{}, err
	}

	newExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	_, err = s.refreshTokens.Create(
		ctx,
		RefreshToken{
			UserID:    user.ID,
			TokenHash: HashRefreshToken(newRawToken),
			ExpiresAt: newExpiresAt,
		},
	)

	if err != nil {
		return LoginResponse{}, err
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user)

	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: newRawToken,
		ExpiresAt:    newExpiresAt,
	}, nil
}

// ! Logout
func (s *AuthService) Logout(
	ctx context.Context,
	rawToken string,
) error {
	tokenHash := HashRefreshToken(rawToken)

	token, err := s.refreshTokens.GetByTokenHash(
		ctx,
		tokenHash,
	)

	if err != nil {
		return nil
	}

	return s.refreshTokens.Revoke(
		ctx,
		token.ID,
	)
}

//! Get User by id
func (s *AuthService) GetUserByID(
	ctx context.Context,
	id int64,
) (User, error){
	return s.users.GetByID(ctx, id)
}
