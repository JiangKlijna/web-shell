package main

type Application struct {
}

// Init App
func (app *Application) Init() {

}

// Start App
func (app *Application) Run() {

}

func main() {
	app := &Application{}
	app.Init()
	app.Run()
}
