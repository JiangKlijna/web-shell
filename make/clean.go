package main

import "os"

// MakeClean go clean && delete gen_go_file
func MakeClean() {
	invoke("go", "clean")
	os.Remove(staticGenGoFile)
}
