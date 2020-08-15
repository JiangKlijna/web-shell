package main

import "os"

// MakeRun down -> gen -> build
func MakeRun() {
	MakeGen()
	MakeBuild()
	invoke(append([]string{"./web-shell"}, os.Args[2:]...)...)
}
