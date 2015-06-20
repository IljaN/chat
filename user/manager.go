package user

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"time"
)

const bcryptCost int = 10

type Manager struct {
	persistence   *Persistence
	authenticator *Authenticator
}

func NewManager(p *Persistence, a *Authenticator) *Manager {
	return &Manager{p, a}
}

func (m *Manager) Register(name string, password string) {
	hash := hashPassword(password)

	hashPassword(password)
	var u = User{
		Id:           generateId(),
		Name:         name,
		passwordHash: hash,
		authToken:    ""}

	m.persistence.persist(u)
}

func (m *Manager) Login(username, password string) {
	u, err := m.persistence.loadByName(username)

	if err != nil {
		log.Fatal(err)
	}

	hashedPwBytes := []byte(u.passwordHash)
	pwBytes := []byte(password)

	err = bcrypt.CompareHashAndPassword(hashedPwBytes, pwBytes)

	if err != nil {

	}

}

func (m *Manager) Authenticated(token string) (isAuthenticated bool, err error) {
	err = m.authenticator.Validate(token)

	if err != nil {
		return false, err
	}

	return true, err

}

func generateId() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%x", rand.Int())
}

func hashPassword(p string) string {
	var pwBytes = []byte(p)

	hashedPassword, err := bcrypt.GenerateFromPassword(pwBytes, bcryptCost)
	if err != nil {
		panic(err)
	}

	return string(hashedPassword)
}
