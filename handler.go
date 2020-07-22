package main

import (
	"crypto/md5"
	"crypto/sha512"

	"encoding/base64"
	"fmt"
	"hash"
	"math/rand"
	"net/http"
	"strings"
	"time"
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

// GenerateToken get token1 token2
// token1 = md5(md5(username+password+clientIP+userAgent)+username+password+clientIP+userAgent)
// token2 = sha512(token1^10)
func GenerateToken(username, password, clientIP, userAgent string) (string, string) {
	hex := func(h hash.Hash, val string) string {
		h.Write([]byte(val))
		return fmt.Sprintf("%x", h.Sum(nil))
	}
	token1 := hex(md5.New(), username+password+clientIP+userAgent)
	token1 = hex(md5.New(), token1+username+password+clientIP+userAgent)
	token2 := hex(sha512.New(), strings.Repeat(token1, 10))
	token2 = hex(sha512.New(), strings.Repeat(token2, 10))
	return token1, token2
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
		_, token2 := GenerateToken(username, password, clientIP, r.UserAgent())
		if path != token2 {
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
		r.ParseForm()
		token := r.Form.Get("token")
		if token == "" {
			w.Header().Set("Client-Ip", clientIP)
			w.Write([]byte("{\"code\":1,\"msg\":\"invalid token!\"}"))
			return
		}

		time.Sleep(time.Duration(rand.Int63n(int64(time.Second))))

		token1, token2 := GenerateToken(username, password, clientIP, r.UserAgent())

		if token == token1 {
			// login success && set cookie
			w.Write([]byte("{\"code\":0,\"msg\":\"login success!\",\"path\":\"" + token2 + "\"}"))
			return
		}
		w.Header().Set("Client-Ip", clientIP)
		w.Write([]byte("{\"code\":1,\"msg\":\"Login incorrect!\"}"))
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
