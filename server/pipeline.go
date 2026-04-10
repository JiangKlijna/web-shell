package server

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/gorilla/websocket"
	"github.com/jiangklijna/web-shell/lib"
	"github.com/runletapp/go-console"
)

// PipeLine Connect websocket and childprocess
type PipeLine struct {
	pty console.Console
	skt *websocket.Conn
}

// NewPipeLine Malloc PipeLine
func NewPipeLine(conn *websocket.Conn, command string) (*PipeLine, error) {
	proc, err := console.New(120, 60)
	if err != nil {
		return nil, err
	}
	err = proc.Start([]string{command})
	if err != nil {
		return nil, err
	}
	return &PipeLine{proc, conn}, nil
}

// ReadSktAndWritePty read skt and write pty
func (w *PipeLine) ReadSktAndWritePty(logChan chan string) {
	for {
		mt, payload, err := w.skt.ReadMessage()
		if err != nil && err != io.EOF {
			logChan <- fmt.Sprintf("ReadSktAndWritePty: websocket read failed: %v", err)
			return
		}
		if mt != websocket.TextMessage {
			logChan <- fmt.Sprintf("ReadSktAndWritePty: invalid message type %d", mt)
			return
		}
		var msg lib.Message
		err = json.Unmarshal(payload, &msg)
		if err != nil {
			logChan <- fmt.Sprintf("ReadSktAndWritePty: invalid message format: %v", err)
			return
		}
		switch msg.Type {
		case lib.TypeResize:
			var size []int
			err := json.Unmarshal(msg.Data, &size)
			if err != nil {
				logChan <- fmt.Sprintf("ReadSktAndWritePty: invalid resize data: %v", err)
				return
			}
			if len(size) != 2 || size[0] <= 0 || size[1] <= 0 {
				logChan <- fmt.Sprintf("ReadSktAndWritePty: invalid resize values: %v", size)
				return
			}
			err = w.pty.SetSize(size[0], size[1])
			if err != nil {
				logChan <- fmt.Sprintf("ReadSktAndWritePty: pty resize failed: %v", err)
				return
			}
		case lib.TypeData:
			var dat string
			err := json.Unmarshal(msg.Data, &dat)
			if err != nil {
				logChan <- fmt.Sprintf("ReadSktAndWritePty: invalid data format: %v", err)
				return
			}
			_, err = w.pty.Write([]byte(dat))
			if err != nil {
				logChan <- fmt.Sprintf("ReadSktAndWritePty: pty write failed: %v", err)
				return
			}
		default:
			logChan <- fmt.Sprintf("ReadSktAndWritePty: unknown message type %d", msg.Type)
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
			logChan <- fmt.Sprintf("ReadPtyAndWriteSkt: pty read failed: %v", err)
			return
		}
		err = w.skt.WriteMessage(websocket.TextMessage, buf[:n])
		if err != nil {
			logChan <- fmt.Sprintf("ReadPtyAndWriteSkt: websocket write failed: %v", err)
			return
		}
	}
}
