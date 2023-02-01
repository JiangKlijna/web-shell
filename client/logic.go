package client

import (
	"crypto/md5"
	"errors"
	"github.com/jiangklijna/web-shell/lib"
	"log"

	"github.com/gorilla/websocket"
)

// LoginServer get websocket path
func LoginServer(https bool, username, password, host, port, contentpath string, get func(url string) (map[string]interface{}, error)) (string, error) {
	protocol := "http"
	if https {
		protocol = "https"
	}
	md5User := lib.HashCalculation(md5.New(), username)
	md5Pass := lib.HashCalculation(md5.New(), password)

	var LoginURL = protocol + "://" + host + ":" + port + contentpath + "/login"
	data, err := get(LoginURL + "?username=" + md5User + "&password=" + md5Pass)
	if err != nil {
		return "", err
	}
	if data["code"] != 0.0 {
		return "", errors.New(data["msg"].(string))
	}
	return data["path"].(string), nil
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
