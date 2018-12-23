package main

import (
	"fmt"
	"os"
)

const MakeRemark = `web-shell Project Build Tool
	make gen
	make build
	make clean
	make publish`

func gen() {

}

func build() {
	gen()
}

func clean() {

}

func publish() {

}

func main() {
	(func() func() {
		if len(os.Args) < 2 {
			fmt.Println(MakeRemark)
			os.Exit(1)
		}
		switch os.Args[1] {
		case "build":
			return build
		case "clean":
			return clean
		case "publish":
			return publish
		default:
			return func() {
				fmt.Println(MakeRemark)
				os.Exit(1)
			}
		}
	})()()
}
