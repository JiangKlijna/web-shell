package main

import (
	"flag"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// Parameter Command line parameters
type Parameter struct {
	Port     string
	Username string
	Password string
	Command  string
}

// Init Parameter
func (parms *Parameter) Init() {
	var (
		help, version bool
	)
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&version, "v", false, "show version and exit")
	flag.StringVar(&(parms.Port), "P", "2019", "listening port")
	flag.StringVar(&(parms.Username), "u", "admin", "username")
	flag.StringVar(&(parms.Password), "p", "admin", "password")
	flag.StringVar(&(parms.Command), "c", "", "command cmd or bash")
	flag.Parse()
	if help {
		printUsage()
		flag.PrintDefaults()
		os.Exit(1)
	} else if version {
		printVersion()
		os.Exit(1)
	} else {
		parms.organize()
	}
}

// Organize command line parameters
func (parms *Parameter) organize() {
	_, err := strconv.Atoi(parms.Port)
	if err != nil {
		println("Port " + parms.Port + ":illegal")
		os.Exit(1)
	}
	parms.Command = strings.Trim(parms.Command, " ")
	if parms.Command == "" {
		parms.Command = defaultCommand()
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
	}
	return "bash"
}
