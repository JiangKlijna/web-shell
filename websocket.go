package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"unicode/utf8"
)

const (
	TypeErr = iota
	TypeData
	TypeResize
)

type Message struct {
	Type int             `json:"t"`
	Data json.RawMessage `json:"d"`
}

type PTY interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Close()
	SetSize(w, h uint16)
}

type PipeLine struct {
	pty  PTY
	conn *websocket.Conn
}

func NewPipeLine(pty PTY, conn *websocket.Conn) *PipeLine {
	return &PipeLine{pty, conn}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func PtyHandler(parms *Parameter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer conn.Close()
		pty, err := OpenPty(parms.Command)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}
		defer pty.Close()
		p := NewPipeLine(pty, conn)
		go p.WritePump()
		p.ReadPump()
	})
}

func (w *PipeLine) ReadPump() {
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
			var size []uint16
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

func (w *PipeLine) WritePump() {
	buf := make([]byte, 8192)
	reader := bufio.NewReader(w.pty)
	var buffer bytes.Buffer
	for {
		n, err := reader.Read(buf)
		if err != nil {
			log.Printf("Error Failed to read from pty master: %s", err)
			return
		}
		//read byte array as Unicode code points (rune in go)
		bufferBytes := buffer.Bytes()
		runeReader := bufio.NewReader(bytes.NewReader(append(bufferBytes[:], buf[:n]...)))
		buffer.Reset()
		i := 0
		for i < n {
			char, charLen, e := runeReader.ReadRune()
			if e != nil {
				log.Printf("Error Failed to read from pty master: %s", err)
				return
			}
			if char == utf8.RuneError {
				runeReader.UnreadRune()
				break
			}
			i += charLen
			buffer.WriteRune(char)
		}
		err = w.conn.WriteMessage(websocket.TextMessage, buffer.Bytes())
		if err != nil {
			log.Printf("Error Failed to send UTF8 char: %s", err)
			return
		}
		buffer.Reset()
		if i < n {
			buffer.Write(buf[i:n])
		}
	}
}
