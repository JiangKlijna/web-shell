package main

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"log"
	"strconv"

	"encoding/base64"
	"fmt"
	"hash"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// HTMLDirHandler FileServer
func HTMLDirHandler() http.Handler {
	return http.FileServer(http.Dir("html"))
}

// GetMethodHandler Only allow GET requests
func GetMethodHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GenerateToken Get secret token path
// secret = sha224(clientIP+userAgent+pid+Server).reverse()
// token = md5(secret+md5(username+secret+password)+secret)
// path = sha512(secret.reverse()^5+token.reverse()^5).reverse()
func GenerateToken(username, password, clientIP, userAgent string) (string, string, string) {
	hex := func(h hash.Hash, val string) string {
		h.Write([]byte(val))
		return fmt.Sprintf("%x", h.Sum(nil))
	}
	reverse := func(s string) string {
		runes := []rune(s)
		for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
			runes[from], runes[to] = runes[to], runes[from]
		}
		return string(runes)
	}
	pid := strconv.Itoa(os.Getpid())
	secret := reverse(hex(sha256.New224(), clientIP+userAgent+pid+Server))
	token := hex(md5.New(), secret+hex(md5.New(), username+secret+password)+secret)
	path := reverse(hex(sha512.New(), strings.Repeat(reverse(secret), 5)+strings.Repeat(reverse(token), 5)))
	return secret, token, path
}

// VerifyHandler Login verification
func VerifyHandler(username, password string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/cmd/")
		if path == "" {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}
		clientIP := r.RemoteAddr[0:strings.LastIndex(r.RemoteAddr, ":")]
		_, _, correctPath := GenerateToken(username, password, clientIP, r.UserAgent())
		if path != correctPath {
			http.Error(w, "403 page forbidden", http.StatusForbidden)
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
		secret, correctToken, path := GenerateToken(username, password, clientIP, r.UserAgent())

		r.ParseForm()
		token := r.Form.Get("token")
		if token == "" {
			w.Write([]byte("{\"code\":1,\"msg\":\"invalid token!\",\"secret\":\"" + secret + "\"}"))
			return
		}

		time.Sleep(time.Duration(rand.Int63n(int64(time.Second))))

		if token != correctToken {
			w.Write([]byte("{\"code\":1,\"msg\":\"Login incorrect!\",\"secret\":\"" + secret + "\"}"))
			return
		}
		// login success
		w.Write([]byte("{\"code\":0,\"msg\":\"login success!\",\"path\":\"" + path + "\"}"))
	})
}

// ConnectionHandler Make websocket and childprocess communicate
func ConnectionHandler(parms *Parameter) http.Handler {
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		pl, err := NewPipeLine(conn, parms.Command)
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

// AuthHandler @Deprecated http authentication
func AuthHandler(username, password string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Dotcoo User Login"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		auths := strings.SplitN(auth, " ", 2)
		if len(auths) != 2 {
			w.Write([]byte("Authorization Error!\n"))
			return
		}
		authMethod := auths[0]
		if authMethod != "Basic" {
			w.Write([]byte("AuthMethod Error!\n"))
			return
		}
		authB64 := auths[1]
		authstr, err := base64.StdEncoding.DecodeString(authB64)
		if err != nil {
			w.Write([]byte("Unauthorized!\n"))
			return
		}
		userPwd := strings.SplitN(string(authstr), ":", 2)
		if len(userPwd) != 2 {
			w.Write([]byte("Type Error!\n"))
			return
		}
		if username != userPwd[0] || password != userPwd[1] {
			w.Header().Set("WWW-Authenticate", `Basic realm="Dotcoo User Login"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
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
