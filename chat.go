package main

import (
	"fmt"
	"log"
)

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

func (c *Chat) DissolveRoom(id string) bool {
	_, roomExists := c.Rooms[id]
	_, bcastExists := c.Backend.roomChannels[id]

	if roomExists != bcastExists {
		if roomExists {
			log.Fatalf("Room with id %v exists but has no broadcaster.")
		}

		if bcastExists {
			log.Fatalf("Broadcaster for room id %v exists but has no room.")
		}
	} else if !roomExists && !bcastExists {
		return false
	}

	delete(c.Rooms, id)
	delete(c.Backend.roomChannels, id)

	return true
}
