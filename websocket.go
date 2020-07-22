package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/runletapp/go-console"
)

// Message.Type
const (
	TypeErr = iota
	TypeData
	TypeResize
)

// Message Websocket Communication data format
type Message struct {
	Type int             `json:"t"`
	Data json.RawMessage `json:"d"`
}

// PipeLine Connect websocket and pty
type PipeLine struct {
	pty  console.Console
	conn *websocket.Conn
}

// NewPipeLine Malloc PipeLine
func NewPipeLine(pty console.Console, conn *websocket.Conn) *PipeLine {
	return &PipeLine{pty, conn}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// CommunicationHandler Make websocket and pty communicate
func CommunicationHandler(parms *Parameter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer conn.Close()
		proc, err := console.New(120, 60)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}
		err = proc.Start([]string{"wsl"})
		if err != nil {
			println(err.Error())
		}
		defer proc.Close()
		p := NewPipeLine(proc, conn)
		go p.ReadPtyAndWriteWebsocket()
		p.ReadWebsocketAndWritePty()
	})
}

// ReadWebsocketAndWritePty read websocket and write pty
func (w *PipeLine) ReadWebsocketAndWritePty() {
	for {
		mt, payload, err := w.conn.ReadMessage()
		if err != nil {
			if err != io.EOF {
				log.Printf("Error conn.ReadMessage failed: %s\n", err)
				return
			}
		}
		if mt != websocket.TextMessage {
			log.Printf("Error Invalid message type %d\n", mt)
			return
		}
		var msg Message
		err = json.Unmarshal(payload, &msg)
		if err != nil {
			log.Printf("Error Invalid message %s\n", err)
			return
		}
		switch msg.Type {
		case TypeResize:
			var size []int
			err := json.Unmarshal(msg.Data, &size)
			if err != nil {
				log.Printf("Error Invalid resize message: %s\n", err)
			} else {
				w.pty.SetSize(size[0], size[1])
			}
		case TypeData:
			var dat string
			err := json.Unmarshal(msg.Data, &dat)
			if err != nil {
				log.Printf("Error Invalid data message %s\n", err)
			} else {
				w.pty.Write([]byte(dat))
			}
		default:
			log.Printf("Error Invalid message type %d\n", mt)
			return
		}
	}
}

// ReadPtyAndWriteWebsocket read pty and write websocket
func (w *PipeLine) ReadPtyAndWriteWebsocket() {
	buf := make([]byte, 8192)
	// reader := bufio.NewReader(w.pty)
	for {
		n, err := w.pty.Read(buf)
		if err != nil {
			log.Printf("Error Failed to read from pty master: %s", err)
			return
		}
		err = w.conn.WriteMessage(websocket.TextMessage, buf[:n])
		if err != nil {
			log.Printf("Error Failed to send UTF8 char: %s", err)
			return
		}
	}
}
