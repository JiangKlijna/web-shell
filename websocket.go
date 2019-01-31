package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"os/exec"
)

func execute(sh string, rw io.ReadWriter) error {
	cmd := exec.Command(sh)
	cmd.Stdin = rw
	cmd.Stdout = rw
	cmd.Stderr = rw
	err := cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}

type WebSocketIO websocket.Conn

func (io *WebSocketIO) Write(p []byte) (int, error) {
	ws := (*websocket.Conn)(io)
	return len(p), ws.WriteMessage(websocket.TextMessage, p)
}

func (io *WebSocketIO) Read(p []byte) (int, error) {
	//return os.Stdin.Read(p)
	ws := (*websocket.Conn)(io)
	_, data, err := ws.ReadMessage()
	if err != nil {
		return 0, err
	}
	fmt.Print(string(data))
	for i, b := range data {
		p[i] = b
	}
	return len(data), nil
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()
	wsio := (*WebSocketIO)(ws)
	eio := NewEncodingIO(wsio)
	err = execute("cmd", eio)
	if err != nil {
		log.Println(err)
		return
	}
}
