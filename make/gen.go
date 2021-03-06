package main

import (
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
)

func compressStatic(m *minify.M, filename string) ([]byte, error) {
	ext := last(strings.Split(filename, "."))
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if ext != "js" && ext != "css" && ext != "html" {
		return bs, nil
	}
	return m.Bytes(ext, bs)
}

// MakeGen generate staticGenGoFile
func MakeGen() {
	m := minify.New()
	m.AddFunc("css", css.Minify)
	m.AddFunc("html", html.Minify)
	m.AddFunc("js", js.Minify)

	genTime := strconv.FormatInt(time.Now().Unix(), 10)

	MakeDown()
	f, err := os.Create(staticGenGoFile)
	if err != nil {
		panic(err)
	}
	io.WriteString(f, `package server
import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"time"
)
var modtime = time.Unix(`+genTime+`, 0)
var errNotDir = errors.New("Not a folder")
// MemoryFile Read R file
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
	return nil, errNotDir
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
// FakeFileSystem Read R file system
type FakeFileSystem struct {
}
// Open open resources by R
func (ffs FakeFileSystem) Open(name string) (http.File, error) {
	if data, ok := R[name]; ok {
		if data != nil {
			return &MemoryFile{bytes.NewReader(data), int64(len(data)), name, false}, nil
		}
		return &MemoryFile{nil, 0, name, true}, nil
	}
	return nil, os.ErrNotExist
}
func init() {
	StaticHandler = http.FileServer(&FakeFileSystem{})
}
// R Static file resources
var R = map[string][]byte{`)
	getFiles(staticDir, func(sf *staticFile) {
		if sf == nil {
			return
		}
		if sf.isDir {
			io.WriteString(f, "\n\t\"/"+sf.path)
			io.WriteString(f, `":nil,`)
		} else {
			bs, err := compressStatic(m, sf.name)
			if err != nil {
				println(sf.name)
				panic(err)
			}
			io.WriteString(f, "\n\t\"/"+sf.path)
			io.WriteString(f, `":{`)
			for _, b := range bs {
				io.WriteString(f, strconv.Itoa(int(b)))
				io.WriteString(f, ",")
			}
			io.WriteString(f, "},")
		}
		println(sf.name, "generate successful")
	})
	io.WriteString(f, "\n\t\"/\":nil,")
	io.WriteString(f, "\n}")
	if err := f.Close(); err != nil {
		panic(err)
	}
}

// staticFile file info
type staticFile struct {
	isDir      bool
	name, path string
}

// getFiles func(string, func(*StaticFile))
func getFiles(dir string, callback func(*staticFile)) {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range fs {
		name := dir + "/" + f.Name()
		sf := &staticFile{f.IsDir(), name, name[len(staticDir)+1:]}
		callback(sf)
		if f.IsDir() {
			getFiles(name, callback)
		}
	}
}
