package main

type EncodingIO struct {
}

func (io *EncodingIO) Write(p []byte) (int, error) {
	return len(p), nil
}

func (io *EncodingIO) Read(p []byte) (int, error) {
	return len(p), nil
}
