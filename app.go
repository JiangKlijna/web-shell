package main

import "net/http"

type Application struct {
	mux   *http.ServeMux
	paras *Parameter
}

// Init App
func (app *Application) Init() {
	app.mux = NewServeMux()
	app.paras = NewParameter()
}

// Start App
func (app *Application) Run() {

}

func main() {
	app := &Application{}
	app.Init()
	app.Run()
}
