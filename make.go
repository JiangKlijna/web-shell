package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

const MakeRemark = `web-shell Project Build Tool
usage:
	make down		Download static resources
	make gen		Generate go file by static resources
	make build		Compile and generate executable files
	make run		Run web-shell
	make debug		Debug web-shell
	make clean		Clean tmp files
`

const gen_go_file = "static_gen.go"

const static_dir = "html"

var static_file = []string{"index.js", "index.css", "xterm.min.js", "xterm.min.css",}

var xterm_file = map[string]string{
	"xterm.min.js":  "https://cdn.bootcss.com/xterm/3.9.1/xterm.min.js",
	"xterm.min.css": "https://cdn.bootcss.com/xterm/3.9.1/xterm.min.css",
}

// Download static resources
func down() {
	get := func(url string) ([]byte, error) {
		res, err := http.Get(url)
		defer res.Body.Close()
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(res.Body)
	}
	fileExists := func(path string) bool {
		if stat, err := os.Stat(path); err == nil {
			return !stat.IsDir()
		}
		return false
	}
	for name, url := range xterm_file {
		filename := static_dir + "/" + name
		exist := fileExists(filename)
		if exist {
			fmt.Println(filename + " already exist")
			continue
		}
		data, err := get(url)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(filename, data, 0664)
		if err != nil {
			panic(err)
		}
		fmt.Println(filename + " download successful")
	}
}

func gen() {
	down()
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
		case "down":
			return down
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
