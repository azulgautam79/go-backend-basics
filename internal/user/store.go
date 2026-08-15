package user

// import (
// 	"errors"
// 	"time"
// )

// var (
// 	ErrUserNotFound     = errors.New("user not found")
// 	ErrEmailAlreadyUsed = errors.New("email already exists")
// 	ErrInvalidSession   = errors.New("invalid session")
// )

// type Store struct {
// 	users         map[int64]User
// 	byEmail       map[string]int64
// 	refreshTokens map[string]RefreshToken
// 	nextID        int64
// }

// func NewStore() *Store {
// 	return &Store{
// 		users:         make(map[int64]User),
// 		byEmail:       make(map[string]int64),
// 		refreshTokens: make(map[string]RefreshToken),
// 		nextID:        1,
// 	}
// }

// // ! Create User
// func (s *Store) Create(user User) (User, error) {

// 	if _, exists := s.byEmail[user.Email]; exists {
// 		return User{}, ErrEmailAlreadyUsed
// 	}

// 	user.ID = s.nextID
// 	s.nextID++

// 	s.users[user.ID] = user
// 	s.byEmail[user.Email] = user.ID

// 	return user, nil
// }

// // ! Get User By Id
// func (s *Store) GetByID(id int64) (User, error) {

// 	user, exists := s.users[id]

// 	if !exists {
// 		return User{}, ErrUserNotFound
// 	}

// 	return user, nil
// }

// // ! Get User By Email
// func (s *Store) GetByEmail(email string) (User, error) {

// 	id, exists := s.byEmail[email]

// 	if !exists {
// 		return User{}, ErrUserNotFound
// 	}

// 	return s.GetByID(id)
// }

// // ! Create Refresh Tokens
// func (s *Store) CreateRefreshToken(userID int64, duration time.Duration) (RefreshToken, error) {
// 	token, err := GenerateRefreshToken()

// 	if err != nil {
// 		return RefreshToken{}, err
// 	}

// 	refreshToken := RefreshToken{
// 		Token:     token,
// 		UserID:    userID,
// 		ExpiresAt: time.Now().Add(duration),
// 		Revoked:   false,
// 	}

// 	s.refreshTokens[token] = refreshToken

// 	return refreshToken, nil
// }

// // ! Get Refresh Token
// func (s *Store) GetRefreshToken(token string) (RefreshToken, error) {
// 	refreshToken, exists := s.refreshTokens[token]

// 	if !exists {
// 		return RefreshToken{}, errors.New("refresh token not found")
// 	}

// 	if refreshToken.Revoked {
// 		return RefreshToken{}, errors.New("refresh token revoked")
// 	}

// 	if time.Now().After(refreshToken.ExpiresAt) {
// 		return RefreshToken{}, errors.New("refresh token expired")
// 	}

// 	return refreshToken, nil
// }

// // ! Revoke Refresh Token
// func (s *Store) RevokeRefreshToken(token string) error {
// 	refreshToken, exists := s.refreshTokens[token]

// 	if !exists {
// 		return errors.New("refresh token not found")
// 	}

// 	refreshToken.Revoked = true

// 	s.refreshTokens[token] = refreshToken

// 	return nil
// }
