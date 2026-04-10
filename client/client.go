package client

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/jiangklijna/web-shell/lib"
)

// Version WebShell Client current version
const Version = "3.0"

// UserAgent Request header[User-Agent]
var UserAgent = fmt.Sprintf("web-shell-client/%s (%s; %s; %s)", Version, runtime.GOOS, runtime.GOARCH, runtime.Version())

// WebShellClient connect to WebShellServer
type WebShellClient struct {
	Client *http.Client
	Dialer *websocket.Dialer
}

// Init http client
func (c *WebShellClient) Init(https bool, crt, key, rootcrt string) {
	if https {
		tlsConfig := &tls.Config{}
		if crt != "" && key != "" && rootcrt != "" {
			cliCrt, err := tls.LoadX509KeyPair(crt, key)
			if err != nil {
				log.Fatalln("Load crt or key file failed:", err.Error())
			}
			tlsConfig.RootCAs = lib.ReadCertPool(rootcrt)
			tlsConfig.Certificates = []tls.Certificate{cliCrt}
		} else if crt != "" {
			tlsConfig.RootCAs = lib.ReadCertPool(crt)
		} else {
			tlsConfig.InsecureSkipVerify = true
		}
		c.Client = &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
		c.Dialer = &websocket.Dialer{TLSClientConfig: tlsConfig}
	} else {
		c.Client = &http.Client{}
		c.Dialer = &websocket.Dialer{}
	}
}

// Run WebShellClient
func (c *WebShellClient) Run(https bool, username, password, host, post, contentpath string) {
	path, err := LoginServer(https, username, password, host, post, contentpath, c.PostJSON)
	if err != nil {
		log.Println("Login to Server failed:", err.Error())
		return
	}
	ConnectSocket(https, host, post, contentpath, path, UserAgent, c.GetWebsocket)
}

// PostRes http post request
func (c *WebShellClient) PostRes(url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// PostJSON http post request and parse JSON
func (c *WebShellClient) PostJSON(url string, body []byte) ([]byte, error) {
	res, err := c.PostRes(url, body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.New("response status is " + strconv.Itoa(res.StatusCode))
	}
	return io.ReadAll(res.Body)
}

// GetWebsocket get websocket connection
func (c *WebShellClient) GetWebsocket(url string) (*websocket.Conn, error) {
	h := make(http.Header)
	h["User-Agent"] = []string{UserAgent}
	skt, _, err := c.Dialer.Dial(url, h)
	return skt, err
}
