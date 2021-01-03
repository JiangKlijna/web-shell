package main

import (
	"os/exec"
)

// MakeBuild Compile and Generate executable file
func MakeBuild() {
	res, err := exec.Command("go", "build").CombinedOutput()
	if err != nil {
		println("web-shell build failed:", err.Error(), "\n", string(res))
	} else {
		println("web-shell build success.")
	}
}
