package main

import (
	"fmt"
	"os"
)

const MakeRemark = `web-shell Project Build Tool
	make gen
	make build
	make clean`

func gen() {

}

func build() {
	gen()
}

func clean() {

}

func main() {
	(func() func() {
		if len(os.Args) < 2 {
			fmt.Println(MakeRemark)
			os.Exit(1)
		}
		switch os.Args[1] {
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
