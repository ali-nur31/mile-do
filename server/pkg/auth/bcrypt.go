package auth

import "golang.org/x/crypto/bcrypt"

type BcryptPasswordManager struct {
}

func NewBcryptPasswordManager() *BcryptPasswordManager {
	return &BcryptPasswordManager{}
}

func (bpm *BcryptPasswordManager) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func (bpm *BcryptPasswordManager) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
