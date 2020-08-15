package main

import (
	"os/exec"
)

// MakeBuild Compile and Generate executable file
func MakeBuild() {
	// exist := fileExists(staticGenGoFile)

	// ps := append([]string{"build"}, goFiles...)
	// if exist {
	// 	ps = append(ps, staticGenGoFile)
	// }

	// if runtime.GOOS == "windows" {
	// 	ps = append(ps, "pty_windows.go")
	// } else {
	// 	ps = append(ps, "pty_notwin.go")
	// }
	res, err := exec.Command("go", "build").CombinedOutput()
	if err != nil {
		println("web-shell build error:", err.Error())
		println(string(res))
	} else {
		println("web-shell build successful")
	}
}
