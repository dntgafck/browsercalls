package server

import (
	"io"
	"log"
	"sync"
	"time"
)

type roomTable struct {
	lock            sync.Mutex
	rooms           map[string]*room
	registerTimeout time.Duration
	roomSrvUrl      string
}

func newRoomTable(to time.Duration, rs string) *roomTable {
	return &roomTable{rooms: make(map[string]*room), registerTimeout: to, roomSrvUrl: rs}
}

func (rt *roomTable) room(id string) *room {
	rt.lock.Lock()
	defer rt.lock.Unlock()

	return rt.roomLocked(id)
}

// создает комнату без лока.
func (rt *roomTable) roomLocked(id string) *room {
	if r, ok := rt.rooms[id]; ok {
		return r
	}
	rt.rooms[id] = newRoom(rt, id, rt.registerTimeout, rt.roomSrvUrl)
	log.Printf("Created room %s", id)

	return rt.rooms[id]
}

// удаляет клиента, если комната окажется пустой - удаляет комнату
func (rt *roomTable) remove(rid string, cid string) {
	rt.lock.Lock()
	defer rt.lock.Unlock()

	rt.removeLocked(rid, cid)
}

// удаляет клиента без лока
func (rt *roomTable) removeLocked(rid string, cid string) {
	if r := rt.rooms[rid]; r != nil {
		r.remove(cid)
		if r.empty() {
			delete(rt.rooms, rid)
			log.Printf("Removed room %s", rid)
		}
	}
}

// отправляет сообщения в комнату, если комната не существует - он ее создает
func (rt *roomTable) send(rid string, srcID string, msg string) error {
	rt.lock.Lock()
	defer rt.lock.Unlock()

	r := rt.roomLocked(rid)
	return r.send(srcID, msg)
}

// отправляет запрос на регистрацию в комнату
func (rt *roomTable) register(rid string, cid string, rwc io.ReadWriteCloser) error {
	rt.lock.Lock()
	defer rt.lock.Unlock()

	r := rt.roomLocked(rid)
	return r.register(cid, rwc)
}

func (rt *roomTable) deregister(rid string, cid string) {
	rt.lock.Lock()
	defer rt.lock.Unlock()

	if r := rt.rooms[rid]; r != nil {
		if c := r.clients[cid]; c != nil {
			if c.registered() {
				c.deregister()

				c.setTimer(time.AfterFunc(rt.registerTimeout, func() {
					rt.removeIfUnregistered(rid, c)
				}))

				log.Printf("Deregistered client %s from room %s", c.id, rid)
				return
			}
		}
	}
}

func (rt *roomTable) removeIfUnregistered(rid string, c *client) {
	log.Printf("Removing client %s from room %s due to timeout", c.id, rid)

	rt.lock.Lock()
	defer rt.lock.Unlock()

	if r := rt.rooms[rid]; r != nil {
		if c == r.clients[c.id] {
			if !c.registered() {
				rt.removeLocked(rid, c.id)
				return
			}
		}
	}
}

func (rt *roomTable) wsCount() int {
	rt.lock.Lock()
	defer rt.lock.Unlock()

	count := 0
	for _, r := range rt.rooms {
		count = count + r.wsCount()
	}
	return count
}
