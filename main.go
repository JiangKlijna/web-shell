package main

import (
	"github.com/jiangklijna/web-shell/cmd"
)

func main() {
	parms := new(cmd.Parameter)
	parms.Init()
	parms.Run()
}
