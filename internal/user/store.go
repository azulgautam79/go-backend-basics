package user

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailAlreadyUsed = errors.New("email already exists")
	ErrInvalidSession   = errors.New("invalid session")
)

type Store struct {
	users    map[int64]User
	byEmail  map[string]int64
	sessions map[string]int64
	nextID   int64
}

func NewStore() *Store {
	return &Store{
		users:    make(map[int64]User),
		byEmail:  make(map[string]int64),
		sessions: make(map[string]int64),
		nextID:   1,
	}
}

//! Create User
func (s *Store) Create(user User) (User, error) {

	if _, exists := s.byEmail[user.Email]; exists {
		return User{}, ErrEmailAlreadyUsed
	}

	user.ID = s.nextID
	s.nextID++

	s.users[user.ID] = user
	s.byEmail[user.Email] = user.ID

	return user, nil
}

//! Get User By Id
func (s *Store) GetByID(id int64) (User, error) {

	user, exists := s.users[id]

	if !exists {
		return User{}, ErrUserNotFound
	}

	return user, nil
}

//! Get User By Email
func (s *Store) GetByEmail(email string) (User, error) {

	id, exists := s.byEmail[email]

	if !exists {
		return User{}, ErrUserNotFound
	}

	return s.GetByID(id)
}

// ! Create Session
func (s *Store) CreateSession(userID int64) (string, error) {
	token, err := GenerateToken()

	if err != nil {
		return "", err
	}

	s.sessions[token] = userID
	return token, nil
}

// ! Get User By Session
func (s *Store) GetUserBySession(token string) (User, error) {
	userID, exists := s.sessions[token]

	if !exists {
		return User{}, errors.New("invalid session")
	}

	return s.GetByID(userID)
}
