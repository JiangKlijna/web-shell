package main

import (
	"flag"
	"fmt"
	"os"
)

type Parameter struct {
	Port     int
	Username string
	Password string
}

func (paras *Parameter) Init() {
	var (
		help, version bool
	)
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&version, "v", false, "show version and exit")
	paras.Port = *flag.Int("P", 2019, "listening port")
	paras.Username = *flag.String("u", "admin", "username")
	paras.Password = *flag.String("p", "admin", "password")
	flag.Parse()
	flag.Usage = usage
	if help {
		printVersion()
		flag.Usage()
		flag.PrintDefaults()
	} else if version {
		printVersion()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: web-shell [-P port] [-u username] [-p password]
Example: web-shell -P 2019 -u admin -p admin

Options:
`)
}

func printVersion() {
	os.Stderr.WriteString("web-shell version: ")
	os.Stderr.WriteString(Version)
	os.Stderr.WriteString("\n")
}
