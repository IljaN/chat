package user

import (
	"encoding/hex"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"time"
)

const bcryptCost int = 10

type Manager struct {
	conn redis.Conn
}

func NewManager(c redis.Conn) *Manager {
	return &Manager{c}
}

func (m *Manager) register(name string, password string) {
	//hash := hashPassword(password)
	//u := User{generateId(), name, hash, ""}
	//m.persistUser(u);
}

func generateId() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%x", rand.Int())
}

func hashPassword(passwd string) string {
	hash_bytes, err := bcrypt.GenerateFromPassword([]byte(passwd), bcryptCost)

	if err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(hash_bytes)
}
