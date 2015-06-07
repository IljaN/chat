package main

import "fmt"

type Chat struct {
	Backend *Backend
	Rooms   map[string]Room
}

func (c *Chat) CreateRoom(r Room, locationFormat string) Room {
	rid := GenerateRoomId()
	r.Id = rid
	r.Location = fmt.Sprintf(locationFormat, rid)
	c.Rooms[rid] = r
	c.Backend.Room(rid)

	return c.Rooms[rid]
}
