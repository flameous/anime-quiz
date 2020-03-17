package quiz

import (
	"encoding/json"
	"golang.org/x/net/websocket"
)

type connection interface {
	send(serverMessage) error
}

type WSConnection struct {
	Conn  *websocket.Conn
	extra struct{}
}

func (w *WSConnection) send(msg serverMessage) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = w.Conn.Write(b)
	if err != nil {
		return err
	}
	return nil
}
