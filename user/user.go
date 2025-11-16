package user

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string
	Name     string
	Password []byte
}

type Repo interface {
	Register(name, password string) (*User, error)
	Login(name, password string) (*User, error)
}

type MemoryRepo struct {
	mu    sync.Mutex
	Users map[string]*User
}

func CreateMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		Users: make(map[string]*User),
	}
}

func (r *MemoryRepo) Register(name, password string) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.Users[name]; ok {
		return nil, errors.New("Данный пользователь уже зарегистрирован")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := User{
		ID:       uuid.NewString(),
		Name:     name,
		Password: hashedPassword,
	}
	r.Users[name] = &u
	return &u, nil
}

func (r *MemoryRepo) Login(name, password string) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.Users[name]
	if !ok {
		return nil, errors.New("Пользователя не мущестует")
	}

	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		return nil, errors.New("Неверный пароль")
	}
	return u, nil
}
