package main

import (
	"flag"
	"log"
	"github.com/osvaldshpengler/browsercalls/server"
)

var tls = flag.Bool("tls", false, "Используется HTTPS")
var port = flag.Int("port", 80, "Порт сервера")
var roomSrv = flag.String("room-server", "http://localhost", "Адрес размещения приложения")

func main() {
	flag.Parse()

	log.Printf("Запуск сервера: tls = %t, port = %d, room-server=%s", *tls, *port, *roomSrv)

	c := server.NewServer(*roomSrv)
	c.Run(*port, *tls)
}
