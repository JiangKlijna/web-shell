package main

import (
	"errors"
	"github.com/axgle/mahonia"
	"io"
)

var charset *mahonia.Charset

func InitEncodingIO(encoding string) error {
	charset = mahonia.GetCharset(encoding)
	if charset == nil {
		return errors.New(encoding + " is not support!")
	}
	return nil
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
	data := charset.NewDecoder().ConvertString(string(p))
	return eio.rw.Write([]byte(data))
}

func (eio *EncodingIO) Read(p []byte) (int, error) {
	data := charset.NewEncoder().ConvertString(string(p))
	return eio.rw.Read([]byte(data))
}
