package main

import "os"

// MakeDebug clean -> down -> build
func MakeDebug() {
	MakeClean()
	MakeDown()
	MakeBuild()
	invoke(append([]string{"./web-shell"}, os.Args[2:]...)...)
}
