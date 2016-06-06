package server

import (
	"errors"
	"io"
	"log"
	"time"
)

const maxQueuedMsgCount = 1024

type client struct {
	id string
	rwc io.ReadWriteCloser
	msgs []string
	timer *time.Timer
}

func newClient(id string, t *time.Timer) *client {
	c := client{id: id, timer: t}
	return &c
}

func (c *client) setTimer(t *time.Timer) {
	if c.timer != nil {
		c.timer.Stop()
	}
	c.timer = t
}

func (c *client) register(rwc io.ReadWriteCloser) error {
	if c.rwc != nil {
		log.Printf("Не регистрируем, так как %sуже подключен", c.id)
		return errors.New("Повторная регистрация")
	}
	c.setTimer(nil)
	c.rwc = rwc
	return nil
}

func (c *client) deregister() {
	if c.rwc != nil {
		c.rwc.Close()
		c.rwc = nil
	}
}

func (c *client) registered() bool {
	return c.rwc != nil
}

func (c *client) enqueue(msg string) error {
	if len(c.msgs) >= maxQueuedMsgCount {
		return errors.New("Слишком много сообщений в очереди")
	}
	c.msgs = append(c.msgs, msg)
	return nil
}

func (c *client) sendQueued(other *client) error {
	if c.id == other.id || other.rwc == nil {
		return errors.New("Неверный клиент")
	}
	for _, m := range c.msgs {
		sendServerMsg(other.rwc, m)
	}
	c.msgs = nil
	log.Printf("Отправка сообщения из очереди от %s к %s", c.id, other.id)
	return nil
}

func (c *client) send(other *client, msg string) error {
	if c.id == other.id {
		return errors.New("Неверный клиент")
	}
	if other.rwc != nil {
		return sendServerMsg(other.rwc, msg)
	}
	return c.enqueue(msg)
}
