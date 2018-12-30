package main

type Application struct {
	paras *Parameter
}

// Init App
func (app *Application) Init() {
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
