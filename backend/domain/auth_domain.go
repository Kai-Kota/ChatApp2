package domain

import (
	"errors"

	"github.com/Kai-Kota/ChatApp2/backend/store"
	"golang.org/x/crypto/bcrypt"
)

type IAuthDomain interface {
	Signup(userName, password string) error
	Login(userName, password string) (string, error)
}

type AuthDomain struct {
	authStore store.IAuthStore
}

func NewAuthDomain(authStore store.IAuthStore) IAuthDomain {
	return &AuthDomain{
		authStore: authStore,
	}
}

func (d *AuthDomain) Signup(userName, password string) error {
	if userName == "" || password == "" {
		return errors.New("userName and password required")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return d.authStore.Signup(userName, string(hashed))
}

func (d *AuthDomain) Login(userName, password string) (string, error) {
	if userName == "" || password == "" {
		return "", errors.New("userName and password required")
	}
	return d.authStore.Login(userName, password)
}
