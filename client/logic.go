package client

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"github.com/jiangklijna/web-shell/lib"
	"log"

	"github.com/gorilla/websocket"
)

// LoginServer get websocket path
func LoginServer(https bool, username, password, host, port, contentpath string, post func(url string, body []byte) ([]byte, error)) (string, error) {
	protocol := "http"
	if https {
		protocol = "https"
	}
	sha256User := lib.HashCalculation(sha256.New(), username)
	sha256Pass := lib.HashCalculation(sha256.New(), password)

	LoginURL := protocol + "://" + host + ":" + port + contentpath + "/login"
	reqBody := lib.LoginRequest{
		Username: sha256User,
		Password: sha256Pass,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	respBytes, err := post(LoginURL, bodyBytes)
	if err != nil {
		return "", err
	}
	var result lib.LoginResult
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", err
	}
	if result.Code != 0 {
		return "", errors.New(result.Msg)
	}
	return result.Path, nil
}

// ConnectSocket c
func ConnectSocket(https bool, host, port, contentpath, path, UserAgent string, conn func(url string) (*websocket.Conn, error)) {
	protocol := "ws"
	if https {
		protocol = "wss"
	}
	skt, err := conn(protocol + "://" + host + ":" + port + contentpath + "/cmd/" + path)
	if err != nil {
		log.Println("Connect to WebSocket failed:", err.Error())
		return
	}
	pl, _ := NewPipeLine(skt)

	logChan := make(chan string)
	go pl.ReadSktAndWriteStdio(logChan)
	go pl.ReadStdioAndWriteSkt(logChan)

	errlog := <-logChan
	log.Println(errlog)
	go func() {
		<-logChan
		close(logChan)
	}()
}
