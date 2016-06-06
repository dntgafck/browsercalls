package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const maxRoomCapacity = 2

type room struct {
	parent *roomTable
	id     string
	clients         map[string]*client
	registerTimeout time.Duration
	roomSrvUrl      string
}

func newRoom(p *roomTable, id string, to time.Duration, rs string) *room {
	return &room{parent: p, id: id, clients: make(map[string]*client), registerTimeout: to, roomSrvUrl: rs}
}

func (rm *room) client(clientID string) (*client, error) {
	if c, ok := rm.clients[clientID]; ok {
		return c, nil
	}
	if len(rm.clients) >= maxRoomCapacity {
		log.Printf("Комната %s полна, клиент %s не будет добавлен", rm.id, clientID)
		return nil, errors.New("Достигнута максимальная вместимость комнаты")
	}

	var timer *time.Timer
	if rm.parent != nil {
		timer = time.AfterFunc(rm.registerTimeout, func() {
			if c := rm.clients[clientID]; c != nil {
				rm.parent.removeIfUnregistered(rm.id, c)
			}
		})
	}
	rm.clients[clientID] = newClient(clientID, timer)

	log.Printf("Added client %s to room %s", clientID, rm.id)

	return rm.clients[clientID], nil
}

func (rm *room) register(clientID string, rwc io.ReadWriteCloser) error {
	c, err := rm.client(clientID)
	if err != nil {
		return err
	}
	if err = c.register(rwc); err != nil {
		return err
	}

	log.Printf("Client %s registered in room %s", clientID, rm.id)

	// Отправка сообщений из очереди от уже зарегистрированного клиента
	if len(rm.clients) > 1 {
		for _, otherClient := range rm.clients {
			otherClient.sendQueued(c)
		}
	}
	return nil
}

// отправляет сообщения для другого клиента в комнате или добавляет их в очередь
func (rm *room) send(srcClientID string, msg string) error {
	src, err := rm.client(srcClientID)
	if err != nil {
		return err
	}

	// отправить сообщение в очередь
	if len(rm.clients) == 1 {
		return rm.clients[srcClientID].enqueue(msg)
	}

	// отправить сообщение другому клиенту в комнате
	for _, oc := range rm.clients {
		if oc.id != srcClientID {
			return src.send(oc, msg)
		}
	}

	// Комната повреждениа.
	return errors.New(fmt.Sprintf("Комната  повреждена %+v", rm))
}

func (rm *room) remove(clientID string) {
	if c, ok := rm.clients[clientID]; ok {
		c.deregister()
		delete(rm.clients, clientID)
		log.Printf("Removed client %s from room %s", clientID, rm.id)

		// Send bye to the room Server.
		resp, err := http.Post(rm.roomSrvUrl+"/bye/"+rm.id+"/"+clientID, "text", nil)
		if err != nil {
			log.Printf("Failed to post BYE to room server %s: %v", rm.roomSrvUrl, err)
		}
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}
}

func (rm *room) empty() bool {
	return len(rm.clients) == 0
}

func (rm *room) wsCount() int {
	count := 0
	for _, c := range rm.clients {
		if c.registered() {
			count += 1
		}
	}
	return count
}
