package server

import (
	"fmt"
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/json; charset=utf-8")

		if username == "" || password == "" {
			// Authentication disabled, return a success regardless
			w.Write([]byte("{\"code\":0,\"msg\":\"Logged-in automatically (no authentication required)\",\"path\":\"noauth--login-not-required\"}"))
			return
		}

		const halfSecond = int64(time.Second / 2)
		time.Sleep(time.Duration(rand.Int63n(halfSecond)))

		if token := r.URL.Query().Get("token"); token != "" {
			p := paseto.NewParser()
			if _, err := p.ParseV4Local(sessionKey, token, nil); err == nil {
				log.Println("resuming session from stored token")
				w.Write([]byte("{\"code\":0,\"msg\":\"login success!\",\"path\":\"" + token + "\"}"))
				return
			}
		}

		sentUser, sentPass := r.URL.Query().Get("username"), r.URL.Query().Get("password")

		if username != sentUser || password != sentPass {
			w.Write([]byte("{\"code\":1,\"msg\":\"Login incorrect!\"}"))
			return
		}

		// login success
		token := paseto.NewToken()

		token.SetIssuedAt(time.Now())
		token.SetNotBefore(time.Now())
		token.SetExpiration(time.Now().Add(96 * time.Hour))
		tokenBytes := token.V4Encrypt(sessionKey, nil)

		w.Write([]byte("{\"code\":0,\"msg\":\"login success!\",\"path\":\"" + tokenBytes + "\"}"))
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
