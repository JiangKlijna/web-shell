package main

// MakeDebug clean -> down -> build
func MakeDebug() {
	MakeClean()
	MakeDown()
	MakeBuild()
	invoke("./app")
}
