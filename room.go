package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Room struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (r *Room) String() string {
	return r.Name
}

func GenerateRoomId() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%x", rand.Int())
}
