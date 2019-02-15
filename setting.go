package main

import (
	"flag"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type Parameter struct {
	Port     string
	Username string
	Password string
	Command  string
	Encoding string
}

func (paras *Parameter) Init() {
	var (
		help, version bool
	)
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&version, "v", false, "show version and exit")
	flag.StringVar(&(paras.Port), "P", "2019", "listening port")
	flag.StringVar(&(paras.Username), "u", "admin", "username")
	flag.StringVar(&(paras.Password), "p", "admin", "password")
	flag.StringVar(&(paras.Command), "c", "", "command cmd or bash")
	flag.StringVar(&(paras.Encoding), "e", "utf8", "encoding")
	flag.Parse()
	if help {
		printUsage()
		flag.PrintDefaults()
		os.Exit(1)
	} else if version {
		printVersion()
		os.Exit(1)
	} else {
		paras.organize()
	}
}

// Organize command line parameters
func (paras *Parameter) organize() {
	_, err := strconv.Atoi(paras.Port)
	if err != nil {
		println("Port " + paras.Port + ":illegal")
		os.Exit(1)
	}
	paras.Command = strings.Trim(paras.Command, " ")
	if paras.Command == "" {
		paras.Command = defaultCommand()
	}
	err = InitEncodingIO(paras.Encoding)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func printUsage() {
	println(`Usage:
  web-shell [-P port] [-u username] [-p password] [-c command] [-e encoding]

Example:
  web-shell -P 2019 -u admin -p admin -c bash -e utf8

Options:`)
}

func printVersion() {
	println("web-shell version:", Version)
}

func defaultCommand() string {
	if runtime.GOOS == "windows" {
		return "cmd"
	} else {
		return "bash"
	}
}
