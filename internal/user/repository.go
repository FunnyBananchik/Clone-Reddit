package user

import (
	"errors"
	"redditclone/internal/utils"
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorUsername    = errors.New("username already exists")
	ErrorUnknownUser = errors.New("there is no user with that name or password")
)

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type account struct {
	user     User
	password []byte
}

type RepositoryUsers struct {
	users map[string]account
	mu    sync.RWMutex
}

func NewRepository() *RepositoryUsers {
	return &RepositoryUsers{users: make(map[string]account)}
}

func (r *RepositoryUsers) Create(username, password string) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	key := strings.ToLower(username)
	if _, ok := r.users[key]; ok {
		return User{}, ErrorUsername
	}
	newUser := User{Id: utils.GenerateID(), Username: username}
	r.users[key] = account{user: newUser, password: hash}
	return newUser, nil
}

func (r *RepositoryUsers) Login(username, password string) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := strings.ToLower(username)
	u, ok := r.users[key]
	if !ok {
		return User{}, ErrorUnknownUser
	}
	err := bcrypt.CompareHashAndPassword(u.password, []byte(password))
	if err != nil {
		return User{}, ErrorUnknownUser
	}
	return u.user, nil
}

func (r *RepositoryUsers) CheckExists(username string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.users[strings.ToLower(username)]
	return ok
}
