package main

const Version = "0.1"
const Server = "web-shell-" + Version

type Application struct {
	server *ShellServer
	paras  *Parameter
}

// New Application
func NewApp() *Application {
	return &Application{NewShellServer(), new(Parameter)}
}

// Init App
func (app *Application) Init() {
	app.paras.Init()
	app.server.Init(app.paras)
}

func main() {
	app := NewApp()
	app.Init()
	app.server.Run(app.paras)
}
