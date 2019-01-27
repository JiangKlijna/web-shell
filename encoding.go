package main

import (
	"github.com/axgle/mahonia"
	"io"
)

var charset *mahonia.Charset

func InitEncodingIO(encoding string) *mahonia.Charset {
	charset = mahonia.GetCharset(encoding)
	return charset
}

type EncodingIO struct {
	rw io.ReadWriter
}

func NewEncodingIO(rw io.ReadWriter) io.ReadWriter {
	if charset.Name == "UTF-8" {
		return rw
	}
	return &EncodingIO{rw}
}

func (eio *EncodingIO) Write(p []byte) (int, error) {
	return eio.rw.Write(p)
}

func (eio *EncodingIO) Read(p []byte) (int, error) {
	return eio.rw.Read(p)
}

