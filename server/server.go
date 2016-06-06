package server

import (
	"crypto/tls"
	"golang.org/x/net/websocket"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const registerTimeoutSec = 10

const wsReadTimeoutSec = 60 * 60 * 24

type Server struct {
	*roomTable
}

func NewServer(rs string) *Server {
	return &Server{
		roomTable: newRoomTable(time.Second * registerTimeoutSec, rs),
	}
}

func (s *Server) Run(p int, useTls bool) {
	http.Handle("/ws", websocket.Handler(s.wsHandler))
	http.HandleFunc("/", s.httpHandler)

	var e error

	pstr := ":" + strconv.Itoa(p)
	if useTls {
		config := &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			},
			PreferServerCipherSuites: true,
		}
		server := &http.Server{Addr: pstr, Handler: nil, TLSConfig: config }

		e = server.ListenAndServeTLS("/cert/cert.pem", "/cert/key.pem")
	} else {
		e = http.ListenAndServe(pstr, nil)
	}

	if e != nil {
		log.Fatal("Run: " + e.Error())
	}
}

func (s *Server) httpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, DELETE")

	p := strings.Split(r.URL.Path, "/")
	if len(p) != 3 {
		s.httpError("Неверный путь: " + r.URL.Path, w)
		return
	}
	rid, cid := p[1], p[2]

	switch r.Method {
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.httpError("Невозможно прочитать тело запроса: " + err.Error(), w)
			return
		}
		m := string(body)
		if m == "" {
			s.httpError("Пустое тело запроса", w)
			return
		}
		if err := s.roomTable.send(rid, cid, m); err != nil {
			s.httpError("Невозможно отправить сообщение: " + err.Error(), w)
			return
		}
	case "DELETE":
		s.roomTable.remove(rid, cid)
	default:

		return
	}

	io.WriteString(w, "OK\n")
}

func (s *Server) wsHandler(ws *websocket.Conn) {
	var rid, cid string

	registered := false

	var msg wsClientMsg
	loop:
	for {
		err := ws.SetReadDeadline(time.Now().Add(time.Duration(wsReadTimeoutSec) * time.Second))
		if err != nil {
			s.wsError("ws.SetReadDeadline error: " + err.Error(), ws)
			break
		}

		err = websocket.JSON.Receive(ws, &msg)
		if err != nil {
			if err.Error() != "EOF" {
				s.wsError("websocket.JSON.Receive error: " + err.Error(), ws)
			}
			break
		}

		switch msg.Cmd {
		case "register":
			if registered {
				s.wsError("Повторный запрос на регистрацию", ws)
				break loop
			}
			if msg.RoomID == "" || msg.ClientID == "" {
				s.wsError("Неверный запрос на регистрацию: не переданы'clientid' или 'roomid'", ws)
				break loop
			}
			if err = s.roomTable.register(msg.RoomID, msg.ClientID, ws); err != nil {
				s.wsError(err.Error(), ws)
				break loop
			}
			registered, rid, cid = true, msg.RoomID, msg.ClientID

			defer s.roomTable.deregister(rid, cid)
			break
		case "send":
			if !registered {
				s.wsError("Клиент не зарегистрирован", ws)
				break loop
			}
			if msg.Msg == "" {
				s.wsError("Неверный запрос на отправку: отсутствует 'msg'", ws)
				break loop
			}
			s.roomTable.send(rid, cid, msg.Msg)
			break
		default:
			s.wsError("Неверное сообщение: отсутствует 'cmd'", ws)
			break
		}
	}
	ws.Close()
}

func (s *Server) httpError(msg string, w http.ResponseWriter) {
	err := errors.New(msg)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (s *Server) wsError(msg string, ws *websocket.Conn) {
	sendServerErr(ws, msg)
}
