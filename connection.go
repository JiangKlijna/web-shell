package main

import (
	"encoding/json"
	"fmt"
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

// PipeLine Connect websocket and childprocess
type PipeLine struct {
	pty console.Console
	skt *websocket.Conn
}

// NewPipeLine Malloc PipeLine
func NewPipeLine(pty console.Console, skt *websocket.Conn) *PipeLine {
	return &PipeLine{pty, skt}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// CommunicationHandler Make websocket and childprocess communicate
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
		defer proc.Close()
		err = proc.Start([]string{parms.Command})
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}

		logChan := make(chan string)
		pl := NewPipeLine(proc, conn)
		go pl.ReadSktAndWritePty(logChan)
		go pl.ReadPtyAndWriteSkt(logChan)

		errlog := <-logChan
		log.Println(errlog)
		go func() {
			<-logChan
			close(logChan)
		}()
	})
}

// ReadSktAndWritePty read skt and write pty
func (w *PipeLine) ReadSktAndWritePty(logChan chan string) {
	for {
		mt, payload, err := w.skt.ReadMessage()
		if err != nil && err != io.EOF {
			logChan <- fmt.Sprintf("Error ReadSktAndWritePty websocket ReadMessage failed: %s", err)
			return
		}
		if mt != websocket.TextMessage {
			logChan <- fmt.Sprintf("Error ReadSktAndWritePty Invalid message type %d", mt)
			return
		}
		var msg Message
		err = json.Unmarshal(payload, &msg)
		if err != nil {
			logChan <- fmt.Sprintf("Error ReadSktAndWritePty Invalid message %s", err)
			return
		}
		switch msg.Type {
		case TypeResize:
			var size []int
			err := json.Unmarshal(msg.Data, &size)
			if err != nil {
				logChan <- fmt.Sprintf("Error ReadSktAndWritePty Invalid resize message: %s", err)
				return
			}
			err = w.pty.SetSize(size[0], size[1])
			if err != nil {
				logChan <- fmt.Sprintf("Error ReadSktAndWritePty pty resize failed: %s", err)
				return
			}
		case TypeData:
			var dat string
			err := json.Unmarshal(msg.Data, &dat)
			if err != nil {
				logChan <- fmt.Sprintf("Error ReadSktAndWritePty Invalid data message %s", err)
				return
			}
			_, err = w.pty.Write([]byte(dat))
			if err != nil {
				logChan <- fmt.Sprintf("Error ReadSktAndWritePty pty write failed: %s", err)
				return
			}
		default:
			logChan <- fmt.Sprintf("Error ReadSktAndWritePty Invalid message type %d", mt)
			return
		}
	}
}

// ReadPtyAndWriteSkt read pty and write skt
func (w *PipeLine) ReadPtyAndWriteSkt(logChan chan string) {
	buf := make([]byte, 4096)
	for {
		n, err := w.pty.Read(buf)
		if err != nil {
			logChan <- fmt.Sprintf("Error ReadPtyAndWriteSkt pty read failed: %s", err)
			return
		}
		err = w.skt.WriteMessage(websocket.TextMessage, buf[:n])
		if err != nil {
			logChan <- fmt.Sprintf("Error ReadPtyAndWriteSkt skt write failed: %s", err)
			return
		}
	}
}
