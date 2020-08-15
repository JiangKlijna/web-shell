package client

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gorilla/websocket"
	"github.com/jiangklijna/web-shell/lib"
	"github.com/nsf/termbox-go"
)

// PipeLine Connect websocket and childprocess
type PipeLine struct {
	skt *websocket.Conn
}

// NewPipeLine Malloc PipeLine
func NewPipeLine(conn *websocket.Conn) (*PipeLine, error) {
	return &PipeLine{conn}, nil
}

// ReadSktAndWriteStdio read skt and write stdout
func (w *PipeLine) ReadSktAndWriteStdio(logChan chan string) {
	for {
		mt, payload, err := w.skt.ReadMessage()
		if err != nil && err != io.EOF {
			logChan <- fmt.Sprintf("Error ReadSktAndWriteTer websocket ReadMessage failed: %s", err)
			return
		}
		if mt != websocket.TextMessage {
			logChan <- fmt.Sprintf("Error ReadSktAndWriteTer Invalid message type %d", mt)
			return
		}
		os.Stdout.Write(payload)
	}
}

// ReadStdioAndWriteSkt read stdin and write skt
func (w *PipeLine) ReadStdioAndWriteSkt(logChan chan string) {
	err := termbox.Init()
	if err != nil {
		logChan <- fmt.Sprintf("Error ReadTerAndWriteSkt Init Termbox failed: %s", err)
		return
	}
	defer termbox.Close()
	for {
		var msg lib.MessageClient
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Ch == 0 {
				ev.Ch = rune(ev.Key)
			}
			msg = lib.MessageClient{Type: lib.TypeData, Data: string(ev.Ch)}
		case termbox.EventResize:
			msg = lib.MessageClient{Type: lib.TypeResize, Data: []termbox.Attribute{termbox.ColorDefault, termbox.ColorDefault}}
		case termbox.EventError:
			logChan <- fmt.Sprintf("Error ReadTerAndWriteSkt Termbox PollEvent failed: %s", err)
			return
		default:
			break
		}
		data, err := json.Marshal(msg)
		if err != nil {
			logChan <- fmt.Sprintf("Error ReadTerAndWriteSkt json.Marshal failed: %s", err)
			return
		}
		err = w.skt.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			logChan <- fmt.Sprintf("Error ReadTerAndWriteSkt skt write failed: %s", err)
			return
		}
	}
}
