package server

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/jiangklijna/web-shell/lib"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	paseto "aidanwoods.dev/go-paseto"
	"github.com/gorilla/websocket"
)

var sessionKey = paseto.NewV4SymmetricKey()

// HTMLDirHandler FileServer
func HTMLDirHandler() http.Handler {
	return http.FileServer(http.Dir("html"))
}

// GetMethodHandler Only allow GET requests
func GetMethodHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// PostMethodHandler Only allow POST requests
func PostMethodHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// VerifyHandler Login verification
func VerifyHandler(username, password string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if username == "" && password == "" {
			// authentication disabled, permit all traffic
		} else {
			token := strings.TrimPrefix(r.URL.Path, "/cmd/")
			if len(token) < 10 {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			p := paseto.NewParser()
			if _, err := p.ParseV4Local(sessionKey, token, nil); err != nil {
				log.Printf("Invalid token: %v", err)
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// LoginHandler Login interface
func LoginHandler(username, password string) http.Handler {
	if username == "" || password == "" {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lib.HttpWriteJSON(w, 0, lib.LoginResult{
				Code: 0,
				Msg:  "Logged-in automatically (no authentication required)",
				Path: "noauth--login-not-required",
			})
		})
	}

	sha256User := lib.HashCalculation(sha256.New(), username)
	sha256Pass := lib.HashCalculation(sha256.New(), password)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const halfSecond = int64(time.Second / 2)
		time.Sleep(time.Duration(rand.Int63n(halfSecond)))

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			lib.HttpWriteJSON(w, http.StatusBadRequest, lib.LoginResult{Code: 1, Msg: "Invalid request body"})
			return
		}

		var req lib.LoginRequest
		if err := json.Unmarshal(body, &req); err != nil {
			lib.HttpWriteJSON(w, http.StatusBadRequest, lib.LoginResult{Code: 1, Msg: "Invalid JSON"})
			return
		}

		if req.Token != "" {
			p := paseto.NewParser()
			if _, err := p.ParseV4Local(sessionKey, req.Token, nil); err == nil {
				log.Println("resuming session from stored token")
				lib.HttpWriteJSON(w, 0, lib.LoginResult{Code: 0, Msg: "login success!", Path: req.Token})
				return
			}
		}

		if sha256User != req.Username || sha256Pass != req.Password {
			lib.HttpWriteJSON(w, 0, lib.LoginResult{Code: 1, Msg: "Login incorrect!"})
			return
		}

		tokenObj := paseto.NewToken()
		tokenObj.SetIssuedAt(time.Now())
		tokenObj.SetNotBefore(time.Now())
		tokenObj.SetExpiration(time.Now().Add(96 * time.Hour))
		tokenBytes := tokenObj.V4Encrypt(sessionKey, nil)

		lib.HttpWriteJSON(w, 0, lib.LoginResult{Code: 0, Msg: "login success!", Path: tokenBytes})
	})
}

// ConnectionHandler Make websocket and childprocess communicate
func ConnectionHandler(command string) http.Handler {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrader.Upgrade error:", err.Error())
			return
		}
		defer conn.Close()

		pl, err := NewPipeLine(conn, command)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}
		defer pl.pty.Close()

		logChan := make(chan string)
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

// ContentPathHandler content path prefix
func ContentPathHandler(contentpath string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, contentpath)
		r.URL.Path = p
		next.ServeHTTP(w, r)
	})
}

// LoggingHandler Log print
func LoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.Header().Add("Server", Server)
		next.ServeHTTP(w, r)
		str := fmt.Sprintf(
			"%s Completed %s %s in %v from %s",
			start.Format("2006/01/02 15:04:05"),
			r.Method,
			r.URL.Path,
			time.Since(start),
			r.RemoteAddr)
		println(str)
	})
}
