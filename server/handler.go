package server

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jiangklijna/web-shell/lib"

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
		token := strings.TrimPrefix(r.URL.Path, "/cmd/")
		if len(token) < 10 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		p := paseto.NewParser()
		parsedToken, err := p.ParseV4Local(sessionKey, token, nil)
		if err != nil {
			log.Printf("Invalid token: %v", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		uaHash, err := parsedToken.GetString("ua")
		if err != nil || uaHash != lib.HashCalculation(sha256.New(), r.Header.Get("User-Agent")) {
			log.Printf("User-Agent mismatch")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoginHandler Login interface
func LoginHandler(username, password string) http.Handler {
	loginRateLimit := lib.NewRateLimiter(10 * time.Millisecond)
	authRateLimit := lib.NewRateLimiter(time.Second)

	sha256User := lib.HashCalculation(sha256.New(), username)
	sha256Pass := lib.HashCalculation(sha256.New(), password)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !loginRateLimit() {
			lib.HttpWriteJSON(w, http.StatusTooManyRequests, lib.LoginResult{Code: 1, Msg: "Too many requests"})
			return
		}

		body, err := io.ReadAll(r.Body)
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
			parsedToken, err := p.ParseV4Local(sessionKey, req.Token, nil)
			if err == nil {
				uaHash, err := parsedToken.GetString("ua")
				if err == nil && uaHash == lib.HashCalculation(sha256.New(), r.Header.Get("User-Agent")) {
					log.Println("resuming session from stored token")
					lib.HttpWriteJSON(w, 0, lib.LoginResult{Code: 0, Msg: "login success!", Path: req.Token})
					return
				}
			}
		}

		if !authRateLimit() {
			lib.HttpWriteJSON(w, http.StatusTooManyRequests, lib.LoginResult{Code: 1, Msg: "Too many login attempts"})
			return
		}

		if sha256User != req.Username || sha256Pass != req.Password {
			lib.HttpWriteJSON(w, 0, lib.LoginResult{Code: 1, Msg: "Login incorrect!"})
			return
		}

		tokenObj := paseto.NewToken()
		tokenObj.SetIssuedAt(time.Now())
		tokenObj.SetNotBefore(time.Now())
		tokenObj.SetExpiration(time.Now().Add(24 * time.Hour))
		tokenObj.Set("ua", lib.HashCalculation(sha256.New(), r.Header.Get("User-Agent")))
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
			origin := r.Header.Get("Origin")
			if origin == "" {
				return true
			}

			originURL, err := url.Parse(origin)
			if err != nil {
				return false
			}

			return originURL.Host == r.Host
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
			log.Println("NewPipeLine error:", err.Error())
			conn.WriteMessage(websocket.TextMessage, []byte("NewPipeLine failed"))
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
		w.Header().Add("X-Powered-By", XPoweredBy)
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
