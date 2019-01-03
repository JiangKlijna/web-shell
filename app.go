package main

import "net/http"

const Version = "0.1"

type Application struct {
	mux   ShellServer
	paras *Parameter
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

func main() {
	app := NewApp()
	app.Init()
	app.Run()
}
