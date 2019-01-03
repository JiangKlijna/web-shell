package main

import "net/http"

type ShellServer *http.ServeMux

func NewShellServer() ShellServer {
	return http.NewServeMux()
}