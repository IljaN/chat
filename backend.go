package chat

import "github.com/dustin/go-broadcast"

type Backend struct {
	roomChannels map[string]broadcast.Broadcaster
}

func NewBackend() *Backend {
	return &Backend{roomChannels: make(map[string]broadcast.Broadcaster)}

}

func (b *Backend) GetRoomChannels() []string {
	keys := make([]string, len(b.roomChannels))
	i := 0
	for k := range b.roomChannels {
		keys[i] = k
		i += 1
	}
	return keys
}

func (b *Backend) OpenListener(roomid string) chan interface{} {
	listener := make(chan interface{})
	b.Room(roomid).Register(listener)
	return listener
}

func (b *Backend) CloseListener(roomid string, listener chan interface{}) {
	b.Room(roomid).Unregister(listener)
	close(listener)
}

func (b *Backend) DeleteBroadcast(roomid string) {
	bcast, ok := b.roomChannels[roomid]
	if ok {
		bcast.Close()
		delete(b.roomChannels, roomid)
	}
}

func (b *Backend) Room(roomid string) broadcast.Broadcaster {

	bcast, ok := b.roomChannels[roomid]
	if !ok {
		bcast = broadcast.NewBroadcaster(10)
		b.roomChannels[roomid] = bcast
	}
	return bcast
}
