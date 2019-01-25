package main

import "bufio"

type EncodingIO struct {
	rw bufio.ReadWriter
}

func (io *EncodingIO) Write(p []byte) (int, error) {
	return io.rw.Write(p)
}

func (io *EncodingIO) Read(p []byte) (int, error) {
	return io.rw.Read(p)
}
