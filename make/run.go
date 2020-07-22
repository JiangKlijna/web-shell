package main

// MakeRun down -> gen -> build
func MakeRun() {
	MakeGen()
	MakeBuild()
	invoke("./app")
}
