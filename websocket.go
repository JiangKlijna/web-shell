package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"os/exec"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

type WebSocketIO struct {
	conn *websocket.Conn
}

func NewWebSocketIO(w http.ResponseWriter, r *http.Request) (*WebSocketIO, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return &WebSocketIO{conn}, nil
}

func (io *WebSocketIO) execute(sh string) error {
	cmd := exec.Command(sh)
	cmd.Stdin = io
	cmd.Stdout = io
	cmd.Stderr = io
	return cmd.Run()
}

func (io *WebSocketIO) Write(p []byte) (int, error) {
	return len(p), io.conn.WriteMessage(websocket.BinaryMessage, p)
}

func (io *WebSocketIO) Read(p []byte) (int, error) {
	_, data, err := io.conn.ReadMessage()
	if err != nil {
		return 0, err
	}
	copy(p, data)
	return len(data), nil
}

func WebsocketHandler(parms *Parameter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io, err := NewWebSocketIO(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer io.conn.Close()
		err = io.execute(parms.Command)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
