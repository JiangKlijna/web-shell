package main

import (
	"fmt"
	"os"
	"os/exec"
)

// MakeRemark make description
const MakeRemark = `web-shell Project Build Tool
usage:
	make down		Download static resources
	make gen		Generate go file by static resources
	make build		Compile and generate executable files
	make run		Run web-shell
	make debug		Debug web-shell
	make clean		Clean tmp files
`

const staticGenGoFile = "static_gen.go"

const staticDir = "html"

var goFiles = []string{"app.go", "setting.go", "handler.go", "websocket.go"}

var staticFiles = []string{
	"https://cdnjs.cloudflare.com/ajax/libs/blueimp-md5/2.16.0/js/md5.min.js",
	"https://unpkg.com/xterm@4.7.0/lib/xterm.js",
	"https://unpkg.com/xterm@4.7.0/css/xterm.css",

	"https://unpkg.com/xterm-addon-fit@0.4.0/lib/xterm-addon-fit.js",
	"https://unpkg.com/xterm-addon-webgl@0.8.0/lib/xterm-addon-webgl.js",
	"https://unpkg.com/xterm-addon-web-links@0.4.0/lib/xterm-addon-web-links.js",
}

func fileExists(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return !stat.IsDir()
	}
	return false
}

func fileDir(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return stat.IsDir()
	}
	return false
}

func last(arr []string) string {
	return arr[len(arr)-1]
}

func invoke(sh ...string) {
	//CombinedOutput
	cmd := exec.Command(sh[0], sh[1:]...)
	println(fmt.Sprint(sh))
	//cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		println(err.Error())
	}
}

func main() {
	(func() func() {
		if len(os.Args) < 2 {
			return MakeBuild
		}
		switch os.Args[1] {
		case "run":
			return MakeRun
		case "debug":
			return MakeDebug
		case "down":
			return MakeDown
		case "gen":
			return MakeGen
		case "build":
			return MakeBuild
		case "clean":
			return MakeClean
		default:
			return func() {
				println(MakeRemark)
			}
		}
	})()()
}
