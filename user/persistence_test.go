package user

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"os"
	"os/exec"
	"testing"
)

var redisProcess *exec.Cmd
var redisConn redis.Conn

func setUp() {
	redisProcess = checkRedis()
	redisConn = createRedisCon()

}

func TestMain(m *testing.M) {
	setUp()
	os.Exit(m.Run())
	redisConn.Do("FLUSHALL")
}

func TestPersistAndLoadUser(t *testing.T) {

	redisConn.Do("FLUSHALL")

	s := getSut()
	u := getTestUser()

	s.persistUser(u)
	loaded, _ := s.loadById(u.Id)

	if loaded.Id != u.Id {
		t.Fail()
	}
	if loaded.Name != u.Name {
		t.Fail()
	}
	if loaded.passwordHash != u.passwordHash {
		t.Fail()
	}
	if loaded.authToken != u.authToken {
		t.Fail()
	}

	redisConn.Do("FLUSHALL")
}

func TestLoadByName(t *testing.T) {

	u := getTestUser()

	// Write test data to redis
	redisConn.Do("HMSET", redisUserPrefix+u.Id,
		"id", u.Id,
		"name", u.Name,
		"passwordHash", u.passwordHash,
		"authToken", u.authToken)

	redisConn.Do("SET", redisUserNameIdxPrefix+u.Name, u.Id)
	redisConn.Flush()

	s := getSut()

	loaded, err := s.loadByName("Max Mustermann")

	if err != nil {
		t.Log("Could not load by Name Max Mustermann")
		t.Fail()
	}

	if loaded.Id != u.Id {
		t.Fail()
	}
	if loaded.Name != u.Name {
		t.Fail()
	}
	if loaded.passwordHash != u.passwordHash {
		t.Fail()
	}
	if loaded.authToken != u.authToken {
		t.Fail()
	}

}

func getSut() Persistence {
	return Persistence{redis: redisConn}
}

func getTestUser() User {
	return User{
		Id:           "someId123",
		Name:         "Max Mustermann",
		passwordHash: "abcdefg",
		authToken:    "someToken"}
}

func checkRedis() *exec.Cmd {
	redis := exec.Command("redis-server", "--port", "1432")
	go redis.Start()

	return redis

}

func createRedisCon() redis.Conn {
	conn, err := redis.Dial("tcp", ":1432")

	if err != nil {
		log.Fatalf("Error (Could not connect to redis testserver):", err)
	}

	return conn

}
