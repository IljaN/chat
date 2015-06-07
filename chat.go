package main

type Chat struct {
	Backend *Backend
	Rooms   map[string]Room
}

func (c *Chat) CreateRoom(r Room) string {
	rid := GenerateRoomId()
	r.Id = rid
	c.Rooms[rid] = r
	c.Backend.Room(rid)

	return rid
}
