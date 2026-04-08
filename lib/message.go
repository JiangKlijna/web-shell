package lib

import "encoding/json"

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

// MessageClient Websocket Communication data format
type MessageClient struct {
	Type int         `json:"t"`
	Data interface{} `json:"d"`
}

// LoginRequest Login request body
type LoginRequest struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResult Login response body
type LoginResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Path string `json:"path"`
}
