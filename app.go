package main

import "net/http"

type Application struct {
	mux   *http.ServeMux
	paras *Parameter
}

func NewApp() *Application {
	return &Application{http.NewServeMux(), &Parameter{}}
}

// Init App
func (app *Application) Init() {
	
}

// Start App
func (app *Application) Run() {

}

func main() {
	app := NewApp()
	app.Init()
	app.Run()
}
