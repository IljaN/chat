package user

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

const redisUserPrefix = "user:"
const redisUserNameIdxPrefix = "idx_user_name:"

type Persistence struct {
	RedisConn redis.Conn
}

func (p *Persistence) loadByName(name string) (User, error) {
	r, err := redis.String(p.RedisConn.Do("GET", redisUserNameIdxPrefix+name))

	if err != nil {
		return User{}, err
	}

	u, err := p.loadById(r)

	if err != nil {
		log.Printf(
			"WARNING: Found byName index %v for user %v with id %v but no entry in the user hashtable",
			redisUserNameIdxPrefix+name,
			name,
			r)
	}

	return u, err

}

func (p *Persistence) loadById(id string) (User, error) {
	umap, err := redis.StringMap(p.RedisConn.Do("HGETALL", redisUserPrefix+id))

	if err != nil {
		return User{}, err
	}

	u := User{umap["id"], umap["name"], umap["passwordHash"], umap["authToken"]}

	return u, err
}

func (p *Persistence) persist(u User) {
	_, err := p.RedisConn.Do("HMSET", redisUserPrefix+u.Id,
		"id", u.Id,
		"name", u.Name,
		"passwordHash", u.passwordHash,
		"authToken", u.authToken)

	if err != nil {
		log.Panicf("Error:", err)
	}

	_, err = p.RedisConn.Do("SET", redisUserNameIdxPrefix+u.Name, u.Id)

	if err != nil {
		log.Panicf("Error:", err)
	}

	p.RedisConn.Flush()
}
