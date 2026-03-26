package main

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
)

const WinptyUrl = "https://github.com/rprichard/winpty/releases/download/0.4.3/winpty-0.4.3-msvc2015.zip"

// MakeWinpty download winpty.zip & parse zip & save winpty.dll
func MakeWinpty() {
	data, err := httpGet(WinptyUrl)
	if err != nil {
		println("download winpty failed.")
		panic(err)
	}
	reader := bytes.NewReader(data)
	zipReader, err := zip.NewReader(reader, int64(len(data)))
	if err != nil {
		println("parse winpty.zip failed.")
		panic(err)
	}
	saveFileByZip(zipReader, "x64/bin/winpty.dll", "winpty.dll")
	//saveFileByZip(zipReader, "x64/bin/winpty-agent.exe", "winpty-agent.exe")
}

func saveFileByZip(zipReader *zip.Reader, zipFileName, toFileName string) {
	zipFile, err := zipReader.Open(zipFileName)
	if err != nil {
		println("open ", zipFileName, " failed.")
		panic(err)
	}
	localF1, err := os.Create(toFileName)
	if err != nil {
		println("create file ", toFileName, " failed.")
		panic(err)
	}
	_, err = io.Copy(localF1, zipFile)
	if err != nil {
		println("save file ", toFileName, " failed.")
		panic(err)
	}
}
