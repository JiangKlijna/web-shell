package main

import (
	"encoding/base64"
	"io"
	"net/http"
	"strings"
)

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
		if username == userPwd[0] && password == userPwd[1] {
			next.ServeHTTP(w, r)
		}
	})
}
