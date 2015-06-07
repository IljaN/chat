package user

import (
	"fmt"
	"time"
	"math/rand"
	"github.com/garyburd/redigo/redis"
	"golang.org/x/crypto/bcrypt"
	"log"
	"encoding/hex"
)

const bcryptCost int = 10



type Manager struct {
	conn redis.Conn
}

func NewManager(c redis.Conn) *Manager {
	return &Manager{c}
}

func (m *Manager) RegisterUser(name string, password string) {
	//hash := hashPassword(password)
	//u := User{generateId(), name, hash, ""}
	//m.persistUser(u);
}




func generateId() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%x", rand.Int())
}

func hashPassword(passwd string) (string) {
	hash_bytes, err := bcrypt.GenerateFromPassword([]byte(passwd), bcryptCost)

	if (err != nil) {
		log.Panicf("Error:", err)
	}


	return hex.EncodeToString(hash_bytes)
}

