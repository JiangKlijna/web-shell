package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
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

var go_file = []string{"app.go", "setting.go", "server.go", "handler.go"}

var static_file = map[string]string{
	"index.html":    "/",
	"index.js":      "/index.js",
	"index.css":     "/index.css",
	"xterm.min.js":  "/xterm.min.js",
	"xterm.min.css": "/xterm.min.css",
}

var xterm_file = map[string]string{
	"xterm.min.js":  "https://cdn.bootcss.com/xterm/3.9.1/xterm.min.js",
	"xterm.min.css": "https://cdn.bootcss.com/xterm/3.9.1/xterm.min.css",
}

// Download static resources
func down() {
	get := func(url string) ([]byte, error) {
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
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

// generate gen_go_file
func gen() {
	down()
	buf := bytes.Buffer{}
	buf.WriteString(`package main
import "net/http"
func init() {
	StaticHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if data, ok := R[r.RequestURI]; ok {
			w.Write(data)
		} else {
			http.Error(w, "File not Found", http.StatusNotFound)
		}
	})
}
var R = map[string][]byte{`)
	for filename, path := range static_file {
		bs, err := ioutil.ReadFile(static_dir + "/" + filename)
		if err != nil {
			panic(err)
		}
		buf.WriteString("\n\t\"" + path)
		buf.WriteString(`":  []byte{`)
		for _, b := range bs {
			buf.WriteString(strconv.Itoa(int(b)))
			buf.WriteString(",")
		}
		buf.WriteString("},")
		fmt.Println(filename + " generate successful")
	}
	buf.WriteString("\n}")
	err := ioutil.WriteFile(gen_go_file, buf.Bytes(), 0662)
	if err != nil {
		panic(err)
	}
}

// Compile and Generate executable file
func build() {
	(func(exist bool) {
		ps := append([]string{"build"}, go_file...)
		if exist {
			ps = append(ps, gen_go_file)
		}
		_, err := exec.Command("go", ps...).CombinedOutput()
		if err != nil {
			fmt.Println("web-shell build error")
		} else {
			fmt.Println("web-shell build successful")
		}
	})((func(path string) bool {
		if stat, err := os.Stat(path); err == nil {
			return !stat.IsDir()
		}
		return false
	})(gen_go_file))
}

func invoke(sh ...string) {
	//CombinedOutput
	cmd := exec.Command(sh[0], sh[1:]...)
	fmt.Println(sh)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Print(line)
	}
	fmt.Println()
	cmd.Wait()
}

// down -> gen -> build
func run() {
	gen()
	build()
}

// clean -> down -> build
func debug() {
	clean()
	down()
	build()
}

func clean() {
	exec.Command("go", "clean").CombinedOutput()
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
		case "debug":
			return debug
		case "down":
			return down
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
