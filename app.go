package main

const Version = "0.1"
const Server = "web-shell-" + Version

var app *Application

type Application struct {
	server *ShellServer
	paras  *Parameter
}

// New Application
func NewApp() *Application {
	return &Application{NewShellServer(), &Parameter{}}
}

// Init App
func (app *Application) Init() {
	app.paras.Init()
	app.server.Init(app.paras)
}

func init() {
	app = NewApp()
}

func main() {
	app.Init()
	app.server.Run(app.paras)
}
