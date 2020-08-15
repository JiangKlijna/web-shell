package client

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"

	"github.com/gorilla/websocket"
)

// Version WebShell Client current version
const Version = "1.0"

// UserAgent Request header[User-Agent]
var UserAgent = fmt.Sprintf("web-shell-client/%s (%s; %s; %s)", Version, runtime.GOOS, runtime.GOARCH, runtime.Version())

// WebShellClient connect to WebShellServer
type WebShellClient struct {
	Client *http.Client
	Dialer *websocket.Dialer
}

// Init http client
func (c *WebShellClient) Init() {
	c.Client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	c.Dialer = &websocket.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
}

// Run WebShellClient
func (c *WebShellClient) Run(host, post, contentpath string) {
	path, err := LoginServer(host, post, contentpath, c.GetJSON)
	if err != nil {
		log.Println("Login to Server failed:", err.Error())
	}
	ConnectSocket(host, post, contentpath, path, UserAgent, c.GetWebsocket)
}

// GetRes http get request
func (c *WebShellClient) GetRes(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", UserAgent)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// GetJSON http get request and parse JSON
func (c *WebShellClient) GetJSON(url string) (map[string]interface{}, error) {
	res, err := c.GetRes(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.New("response status is " + strconv.Itoa(res.StatusCode))
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(bytes, &data)
	return data, err
}

// GetWebsocket get websocket connection
func (c *WebShellClient) GetWebsocket(url string) (*websocket.Conn, error) {
	h := make(http.Header)
	h["User-Agent"] = []string{UserAgent}
	skt, _, err := c.Dialer.Dial(url, h)
	return skt, err
}
