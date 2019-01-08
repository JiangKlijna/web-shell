package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func HtmlDirHandler() http.Handler {
	return http.FileServer(http.Dir("html"))
}

func MimeHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var filext string
		if strings.HasSuffix(r.RequestURI, "/") {
			filext = ".html"
		} else {
			filext = filepath.Ext(r.RequestURI)
		}
		ctype := mime.TypeByExtension(filext)
		w.Header().Set("Content-Type", ctype)
		next.ServeHTTP(w, r)
	})
}

func GetMethodHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthHandler(username, password string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", Server)
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Dotcoo User Login"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		auths := strings.SplitN(auth, " ", 2)
		if len(auths) != 2 {
			io.WriteString(w, "Authorization Error!\n")
			return
		}
		authMethod := auths[0]
		if authMethod != "Basic" {
			io.WriteString(w, "AuthMethod Error!\n")
			return
		}
		authB64 := auths[1]
		authstr, err := base64.StdEncoding.DecodeString(authB64)
		if err != nil {
			io.WriteString(w, "Unauthorized!\n")
			return
		}
		userPwd := strings.SplitN(string(authstr), ":", 2)
		if len(userPwd) != 2 {
			io.WriteString(w, "Type Error!\n")
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

func LoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		str := fmt.Sprintf(
			"%s Comleted %s %s in %v from %s",
			start.Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			time.Since(start),
			r.RemoteAddr)
		go fmt.Println(str)
	})
}
