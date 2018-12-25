package main

import (
	"fmt"
	"os"
	"os/exec"
)

const MakeRemark = `web-shell Project Build Tool
	make run
	make gen
	make build
	make clean`

const gen_go_file = "static_gen.go"

const static_dir = "html"

var static_file = []string{"index.js", "index.css", "xterm.min.js", "xterm.min.css",}

var xterm_file = map[string]string{
	"xterm.min.js":  "https://cdn.bootcss.com/xterm/3.9.1/xterm.min.js",
	"xterm.min.css": "https://cdn.bootcss.com/xterm/3.9.1/xterm.min.css",
}

func gen() {

}

func build() {
	gen()
}

func run() {
	build()
}

func clean() {
	cmd := exec.Command("go", "clean")
	cmd.CombinedOutput()
	os.Remove(gen_go_file)
}

func main() {
	(func() func() {
		if len(os.Args) < 2 {
			fmt.Println(MakeRemark)
			os.Exit(1)
		}
		switch os.Args[1] {
		case "run":
			return run
		case "gen":
			return gen
		case "build":
			return build
		case "clean":
			return clean
		default:
			return func() {
				fmt.Println(MakeRemark)
				os.Exit(1)
			}
		}
	})()()
}
