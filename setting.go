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
	paras.Port = strconv.Itoa(*flag.Int("P", 2019, "listening port"))
	paras.Username = *flag.String("u", "admin", "username")
	paras.Password = *flag.String("p", "admin", "password")
	paras.Command = *flag.String("c", "", "command cmd or bash")
	paras.Encoding = *flag.String("e", "utf8", "encoding")
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
	paras.Command = strings.Trim(paras.Command, " ")
	if paras.Command == "" {
		paras.Command = defaultCommand()
	}
	err := InitEncodingIO(paras.Encoding)
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
