package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
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

const winpty_dir = "winpty"
const winpty_zip = winpty_dir + ".zip"
const winpty_download = "https://github.com/rprichard/winpty/releases/download/0.4.3/winpty-0.4.3-msvc2015.zip"

var go_file = []string{"app.go", "setting.go", "handler.go", "websocket.go"}

var xterm_files = []string{
	"https://unpkg.com/xterm@4.0.0/lib/xterm.js",
	"https://unpkg.com/xterm@4.0.0/css/xterm.css",

	"https://unpkg.com/xterm-addon-fit@0.3.0/lib/xterm-addon-fit.js",
	"https://unpkg.com/xterm-addon-web-links@0.2.1/lib/xterm-addon-web-links.js",
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

// Download static resources
func down() {
	get := func(url string) ([]byte, error) {
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		if res.StatusCode < 200 || res.StatusCode >= 300 {
			return nil, errors.New("response status is " + strconv.Itoa(res.StatusCode))
		}
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

	// download pty and un zip
	if runtime.GOOS == "windows" {
		if fileExists(winpty_dir+"/winpty.dll") && fileExists(winpty_dir+"/winpty-agent.exe") {
			return
		}

		if fileExists(winpty_zip) {
			fmt.Println(winpty_zip + " already exist")
		} else {
			fmt.Println(winpty_zip + " is downloading")
			data, err := get(winpty_download)
			if err != nil {
				panic(err)
			}
			err = ioutil.WriteFile(winpty_zip, data, 0664)
			if err != nil {
				panic(err)
			}
			fmt.Println(winpty_zip + " download successful")
		}
		os.Mkdir(winpty_dir, 0664)
		reader, err := zip.OpenReader(winpty_zip)
		if err != nil {
			panic(err)
		}
		var winpty_dll string
		var winpty_agent string
		if runtime.GOARCH == "amd64" {
			winpty_dll = "x64/bin/winpty.dll"
			winpty_agent = "x64/bin/winpty-agent.exe"
		} else {
			winpty_dll = "ia32/bin/winpty.dll"
			winpty_agent = "ia32/bin/winpty-agent.exe"
		}

		// foreach zip file
		for _, file := range reader.File {
			println(file.Name)
			if file.Name == winpty_dll {
				dst, err := os.OpenFile(winpty_dir+"/winpty.dll", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
				if err != nil {
					panic(err)
				}
				src, err := file.Open()
				if err != nil {
					panic(err)
				}
				_, err = io.Copy(dst, src)
				if err != nil {
					panic(err)
				}
				src.Close()
				dst.Close()
			} else if file.Name == winpty_agent {
				f, err := os.OpenFile(winpty_dir+"/winpty-agent.exe", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
				if err != nil {
					panic(err)
				}
				src, err := file.Open()
				if err != nil {
					panic(err)
				}
				_, err = io.Copy(f, src)
				if err != nil {
					panic(err)
				}
				src.Close()
				f.Close()
			}
		}
		reader.Close()
		fmt.Println(winpty_zip + " unzip successful")
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
	"time"
)
var modtime = time.Now()
var notDir = errors.New("Not a folder")
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
	return nil, notDir
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
		return nil, os.ErrNotExist
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
	if runtime.GOOS == "windows" {
		ps = append(ps, "pty_windows.go")
	} else {
		ps = append(ps, "pty_notwin.go")
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
	err := cmd.Run()
	if err != nil {
		println(err.Error())
	}
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
	os.Remove(winpty_zip)
}

func main() {
	(func() func() {
		if len(os.Args) < 2 {
			return build
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
