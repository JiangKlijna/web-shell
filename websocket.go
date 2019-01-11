package main

import (
	"github.com/axgle/mahonia"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"os/exec"
)

func execute(sh string, rw io.ReadWriter, down chan error) {
	cmd := exec.Command(sh)
	cmd.Stdin = rw
	cmd.Stdout = rw
	cmd.Stderr = rw

	err := cmd.Start()
	if err != nil {
		down <- err
		return
	}
	down <- cmd.Wait()
}

type WebSocketIO websocket.Conn

func (io *WebSocketIO) Write(p []byte) (n int, err error) {
	ws := (*websocket.Conn)(io)
	data := mahonia.NewDecoder("gbk").ConvertString(string(p))
	d := []byte(data)
	return len(d), ws.WriteMessage(websocket.TextMessage, d)
}

func (io *WebSocketIO) Read(p []byte) (n int, err error) {
	//return os.Stdin.Read(p)
	ws := (*websocket.Conn)(io)
	_, data, err := ws.ReadMessage()
	if err != nil {
		return 0, err
	}
	for i, b := range data {
		p[i] = b
	}
	return len(data), nil
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()
	done := make(chan error)
	wsio := (*WebSocketIO)(ws)
	execute("cmd", wsio, done)
}
