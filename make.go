package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

const xterm_version = "3.14.5"

var go_file = []string{"app.go", "setting.go", "handler.go", "websocket.go"}

var xterm_files = []string{
	"https://cdn.bootcss.com/xterm/" + xterm_version + "/xterm.min.js",
	"https://cdn.bootcss.com/xterm/" + xterm_version + "/xterm.min.css",
	"https://cdn.bootcss.com/xterm/" + xterm_version + "/addons/fit/fit.min.js",
	"https://cdn.bootcss.com/xterm/" + xterm_version + "/addons/webLinks/webLinks.min.js",
}

func fileExists(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return !stat.IsDir()
	}
	return false
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
	last := func(arr []string) string {
		return arr[len(arr)-1]
	}
	for _, url := range xterm_files {
		name := last(strings.Split(url, "/"))
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
import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"syscall"
	"time"
)
var modtime = time.Now()
type MemoryFile struct {
	*bytes.Reader
	size  int64
	name  string
	isDir bool
}
func (m *MemoryFile) Close() error {
	return nil
}
func (m *MemoryFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, errors.New("no dir")
}
func (m *MemoryFile) Stat() (os.FileInfo, error) {
	return m, nil
}
func (m *MemoryFile) Name() string {
	return m.name
}
func (m *MemoryFile) Size() int64 {
	return m.size
}
func (m *MemoryFile) Mode() os.FileMode {
	return os.ModePerm
}
func (m *MemoryFile) ModTime() time.Time {
	return modtime
}
func (m *MemoryFile) IsDir() bool {
	return m.isDir
}
func (m *MemoryFile) Sys() interface{} {
	return nil
}
type FakeFileSystem struct {
}
func (ffs FakeFileSystem) Open(name string) (http.File, error) {
	if data, ok := R[name]; ok {
		if data != nil {
			return &MemoryFile{bytes.NewReader(data), int64(len(data)), name, false}, nil
		} else {
			return &MemoryFile{nil, 0, name, true}, nil
		}
	} else {
		return nil, syscall.ERROR_PATH_NOT_FOUND
	}
}
func init() {
	StaticHandler = http.FileServer(&FakeFileSystem{})
}
var R = map[string][]byte{`)
	type StaticFile struct {
		isDir      bool
		name, path string
	}
	var getFiles func(string, func(*StaticFile))
	getFiles = func(dir string, callback func(*StaticFile)) {
		fs, err := ioutil.ReadDir(dir)
		if err != nil {
			panic(err)
		}
		for _, f := range fs {
			name := dir + "/" + f.Name()
			sf := &StaticFile{f.IsDir(), name, name[len(static_dir)+1:]}
			callback(sf)
			if f.IsDir() {
				getFiles(name, callback)
			}
		}
	}
	getFiles(static_dir, func(sf *StaticFile) {
		if sf == nil {
			return
		}
		if sf.isDir {
			buf.WriteString("\n\t\"/" + sf.path)
			buf.WriteString(`":nil,`)
		} else {
			bs, err := ioutil.ReadFile(sf.name)
			if err != nil {
				panic(err)
			}
			buf.WriteString("\n\t\"/" + sf.path)
			buf.WriteString(`":{`)
			for _, b := range bs {
				buf.WriteString(strconv.Itoa(int(b)))
				buf.WriteString(",")
			}
			buf.WriteString("},")
		}
		fmt.Println(sf.name + " generate successful")
	})
	buf.WriteString("\n\t\"/\":nil,")
	buf.WriteString("\n}")
	err := ioutil.WriteFile(gen_go_file, buf.Bytes(), 0662)
	if err != nil {
		panic(err)
	}
}

// Compile and Generate executable file
func build() {
	exist := fileExists(gen_go_file)

	ps := append([]string{"build"}, go_file...)
	if exist {
		ps = append(ps, gen_go_file)
	}
	res, err := exec.Command("go", ps...).CombinedOutput()
	if err != nil {
		fmt.Println("web-shell build error : " + err.Error() + "\n" + string(res))
	} else {
		fmt.Println("web-shell build successful")
	}
}

func invoke(sh ...string) {
	//CombinedOutput
	cmd := exec.Command(sh[0], sh[1:]...)
	fmt.Println(sh)
	//cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// down -> gen -> build
func run() {
	gen()
	build()
	invoke("app")
}

// clean -> down -> build
func debug() {
	clean()
	down()
	build()
	invoke("app")
}

func clean() {
	invoke("go", "clean")
	err := os.Remove(gen_go_file)
	if err != nil {
		panic(err)
	}
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
