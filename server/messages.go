package server

import (
	"encoding/json"
	"io"
)

// Сообщение от клиента
type wsClientMsg struct {
	Cmd      string `json:"cmd"`
	RoomID   string `json:"roomid"`
	ClientID string `json:"clientid"`
	Msg      string `json:"msg"`
}

// Сообщение для клиента
type wsServerMsg struct {
	Msg   string `json:"msg"`
	Error string `json:"error"`
}

func sendServerMsg(w io.Writer, msg string) error {
	m := wsServerMsg{
		Msg: msg,
	}
	return send(w, m)
}

func sendServerErr(w io.Writer, errMsg string) error {
	m := wsServerMsg{
		Error: errMsg,
	}
	return send(w, m)
}

func send(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	if err := enc.Encode(data); err != nil {
		return err
	}
	return nil
}
