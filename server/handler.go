package server

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jiangklijna/web-shell/lib"
)

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
		path := strings.TrimPrefix(r.URL.Path, "/cmd/")
		if len(path) < 10 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		clientIP := r.RemoteAddr[0:strings.LastIndex(r.RemoteAddr, ":")]
		_, _, correctPath := lib.GenerateAll(username, password, clientIP, r.UserAgent())
		if path != correctPath {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// LoginHandler Login interface
func LoginHandler(username, password string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/json; charset=utf-8")
		clientIP := r.RemoteAddr[0:strings.LastIndex(r.RemoteAddr, ":")]
		secret := lib.GenerateSecret(clientIP, r.UserAgent())

		r.ParseForm()
		token := r.Form.Get("token")
		if token == "" {
			w.Write([]byte("{\"code\":1,\"msg\":\"invalid token!\",\"secret\":\"" + secret + "\"}"))
			return
		}

		time.Sleep(time.Duration(rand.Int63n(int64(time.Second / 2))))
		correctToken := lib.GenerateToken(username, password, secret)

		if token != correctToken {
			w.Write([]byte("{\"code\":1,\"msg\":\"Login incorrect!\",\"secret\":\"" + secret + "\"}"))
			return
		}
		path := lib.GeneratePath(secret, token)
		// login success
		w.Write([]byte("{\"code\":0,\"msg\":\"login success!\",\"path\":\"" + path + "\"}"))
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
			"%s Comleted %s %s in %v from %s",
			start.Format("2006/01/02 15:04:05"),
			r.Method,
			r.URL.Path,
			time.Since(start),
			r.RemoteAddr)
		println(str)
	})
}
