package main

import "net/http"

const Version = "0.1"

var app *Application

type Application struct {
	server ShellServer
	paras  *Parameter
}

// New Application
func NewApp() *Application {
	return &Application{ShellServer(http.NewServeMux()), &Parameter{}}
}

// Init App
func (app *Application) Init() {
	app.paras.Init()
}

// Start App
func (app *Application) Run() {

}

func init() {
	app = NewApp()
}

func main() {
	app.Init()
	app.Run()
}
